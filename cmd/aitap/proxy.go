package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/aniketljoshi/aitap/internal/model"
	"github.com/aniketljoshi/aitap/internal/parser"
	"github.com/aniketljoshi/aitap/internal/provider"
)

// providerUpstream maps a provider prefix path to its real HTTPS upstream.
// Users set their SDK base URL to http://localhost:9119/openai/v1 etc.
var providerUpstream = map[string]providerRoute{
	"/openai/":     {host: "api.openai.com", provider: model.ProviderOpenAI},
	"/anthropic/":  {host: "api.anthropic.com", provider: model.ProviderAnthropic},
	"/google/":     {host: "generativelanguage.googleapis.com", provider: model.ProviderGoogle},
	"/openrouter/": {host: "openrouter.ai", provider: model.ProviderOpenRouter},
	"/ollama/":     {host: "localhost:11434", provider: model.ProviderOllama},
}

type providerRoute struct {
	host     string
	provider model.Provider
}

// Known LLM API hosts for HTTP_PROXY mode detection.
var llmHosts = []string{
	"api.openai.com",
	"api.anthropic.com",
	"generativelanguage.googleapis.com",
	"openrouter.ai",
	"localhost:11434",
	"127.0.0.1:11434",
}

func startProxy(port int, callChan chan<- *model.Call, filterProvider string) error {
	proxy := &llmProxy{
		callChan:       callChan,
		filterProvider: filterProvider,
	}

	addr := fmt.Sprintf(":%d", port)
	log.Printf("aitap proxy listening on %s", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: proxy,
	}
	return server.ListenAndServe()
}

type llmProxy struct {
	callChan       chan<- *model.Call
	filterProvider string
}

func (p *llmProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// --- Mode 1: Forward proxy (CONNECT for HTTPS) ---
	if r.Method == http.MethodConnect {
		p.tunnelConnect(w, r)
		return
	}

	// --- Mode 2: Forward proxy mode via path prefix ---
	// e.g. POST http://localhost:9119/openai/v1/chat/completions
	for prefix, route := range providerUpstream {
		if strings.HasPrefix(r.URL.Path, prefix) {
			p.handleForwardProxy(w, r, prefix, route)
			return
		}
	}

	// --- Mode 3: HTTP_PROXY mode (client sends absolute URL) ---
	if r.URL.IsAbs() && isLLMHost(r.URL.Host) {
		p.handleHTTPProxy(w, r)
		return
	}

	// Non-LLM traffic in proxy mode — just forward
	if r.URL.IsAbs() {
		p.forwardPassthrough(w, r)
		return
	}

	// Direct request to aitap with unknown path — show help
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "aitap is running.\n\n")
	fmt.Fprintf(w, "Forward proxy mode:\n")
	fmt.Fprintf(w, "  OPENAI_BASE_URL=http://localhost:%d/openai/v1\n", 9119)
	fmt.Fprintf(w, "  ANTHROPIC_BASE_URL=http://localhost:%d/anthropic\n", 9119)
	fmt.Fprintf(w, "\nHTTP proxy mode:\n")
	fmt.Fprintf(w, "  HTTP_PROXY=http://localhost:%d\n", 9119)
}

// handleForwardProxy handles requests via path-prefix routing.
// The user's SDK points base_url at http://localhost:9119/openai/v1
// and we strip the prefix and forward to the real API over HTTPS.
func (p *llmProxy) handleForwardProxy(w http.ResponseWriter, r *http.Request, prefix string, route providerRoute) {
	// Apply provider filter
	if p.filterProvider != "" && string(route.provider) != p.filterProvider {
		p.forwardToUpstream(w, r, route, prefix)
		return
	}

	startTime := time.Now()

	// Read request body
	var reqBody []byte
	if r.Body != nil {
		reqBody, _ = io.ReadAll(r.Body)
		r.Body.Close()
	}

	// Build call record
	call := &model.Call{
		Timestamp: startTime,
		StartTime: startTime,
		Provider:  route.provider,
		Endpoint:  strings.TrimPrefix(r.URL.Path, prefix[:len(prefix)-1]),
	}
	parser.ParseRequest(route.provider, reqBody, call)

	// Determine upstream URL
	scheme := "https"
	if route.provider == model.ProviderOllama {
		scheme = "http"
	}
	upstreamPath := strings.TrimPrefix(r.URL.Path, prefix[:len(prefix)-1])
	upstreamURL := fmt.Sprintf("%s://%s%s", scheme, route.host, upstreamPath)
	if r.URL.RawQuery != "" {
		upstreamURL += "?" + r.URL.RawQuery
	}

	// Build upstream request
	upReq, err := http.NewRequestWithContext(r.Context(), r.Method, upstreamURL, bytes.NewReader(reqBody))
	if err != nil {
		http.Error(w, "Failed to create upstream request: "+err.Error(), http.StatusBadGateway)
		return
	}

	// Copy headers (skip hop-by-hop)
	copyHeaders(r.Header, upReq.Header)

	// Set correct Host header
	upReq.Host = route.host

	client := &http.Client{
		Timeout: 5 * time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		},
	}
	resp, err := client.Do(upReq)
	if err != nil {
		http.Error(w, "Upstream error: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	call.StatusCode = resp.StatusCode

	// Check if this is a streaming response (SSE)
	isSSE := strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream")

	if isSSE {
		p.handleSSEResponse(w, resp, call, startTime)
	} else {
		p.handleBufferedResponse(w, resp, call, startTime)
	}
}

// handleSSEResponse streams SSE chunks to the client while capturing data for the TUI.
func (p *llmProxy) handleSSEResponse(w http.ResponseWriter, resp *http.Response, call *model.Call, startTime time.Time) {
	// Copy response headers to client
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)

	flusher, canFlush := w.(http.Flusher)

	// Read SSE events and forward to client while accumulating data
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 256*1024), 1024*1024)

	var sseChunks []string

	for scanner.Scan() {
		line := scanner.Text()

		// Forward every line to the client immediately
		fmt.Fprintf(w, "%s\n", line)
		if canFlush {
			flusher.Flush()
		}

		// Collect data lines for parsing
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data != "[DONE]" {
				sseChunks = append(sseChunks, data)
			}
		}
	}

	// Final newline
	fmt.Fprintf(w, "\n")
	if canFlush {
		flusher.Flush()
	}

	call.Latency = time.Since(startTime)
	call.IsStreaming = true
	call.Completed = true

	// Parse the accumulated SSE chunks
	parser.ParseSSEChunks(call.Provider, sseChunks, call)
	call.EstimatedCost = provider.EstimateCost(call.Model, call.InputTokens, call.OutputTokens)

	p.callChan <- call
}

// handleBufferedResponse reads the full response body, parses it, and forwards to client.
func (p *llmProxy) handleBufferedResponse(w http.ResponseWriter, resp *http.Response, call *model.Call, startTime time.Time) {
	respBody, _ := io.ReadAll(resp.Body)

	call.Latency = time.Since(startTime)
	call.ResponseBody = string(respBody)
	call.Completed = true

	parser.ParseResponse(call.Provider, respBody, call)

	// Forward response to client
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)

	p.callChan <- call
}

// handleHTTPProxy handles traditional HTTP_PROXY mode (absolute URL requests).
func (p *llmProxy) handleHTTPProxy(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	prov := provider.DetectProvider(r.URL.Host)

	// Apply filter
	if p.filterProvider != "" && string(prov) != p.filterProvider {
		p.forwardPassthrough(w, r)
		return
	}

	// Read request body
	var reqBody []byte
	if r.Body != nil {
		reqBody, _ = io.ReadAll(r.Body)
		r.Body.Close()
	}

	call := &model.Call{
		Timestamp: startTime,
		StartTime: startTime,
		Provider:  prov,
		Endpoint:  r.URL.Path,
	}
	parser.ParseRequest(prov, reqBody, call)

	// Forward to actual API
	url := r.URL.String()
	upReq, err := http.NewRequestWithContext(r.Context(), r.Method, url, bytes.NewReader(reqBody))
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusBadGateway)
		return
	}
	copyHeaders(r.Header, upReq.Header)

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(upReq)
	if err != nil {
		http.Error(w, "Upstream error: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	call.StatusCode = resp.StatusCode

	isSSE := strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream")
	if isSSE {
		p.handleSSEResponse(w, resp, call, startTime)
	} else {
		p.handleBufferedResponse(w, resp, call, startTime)
	}
}

// forwardToUpstream forwards a request to the upstream without capturing (filtered out).
func (p *llmProxy) forwardToUpstream(w http.ResponseWriter, r *http.Request, route providerRoute, prefix string) {
	scheme := "https"
	if route.provider == model.ProviderOllama {
		scheme = "http"
	}
	upstreamPath := strings.TrimPrefix(r.URL.Path, prefix[:len(prefix)-1])
	upstreamURL := fmt.Sprintf("%s://%s%s", scheme, route.host, upstreamPath)
	if r.URL.RawQuery != "" {
		upstreamURL += "?" + r.URL.RawQuery
	}

	var body io.Reader
	if r.Body != nil {
		body = r.Body
		defer r.Body.Close()
	}

	upReq, err := http.NewRequestWithContext(r.Context(), r.Method, upstreamURL, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	copyHeaders(r.Header, upReq.Header)
	upReq.Host = route.host

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(upReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// forwardPassthrough forwards non-LLM traffic without inspection.
func (p *llmProxy) forwardPassthrough(w http.ResponseWriter, r *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// tunnelConnect handles HTTPS CONNECT tunneling (passthrough, no MITM).
func (p *llmProxy) tunnelConnect(w http.ResponseWriter, r *http.Request) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	targetConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		targetConn.Close()
		return
	}

	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	go transfer(targetConn, clientConn)
	go transfer(clientConn, targetConn)
}

func transfer(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()
	io.Copy(dst, src)
}

// copyHeaders copies HTTP headers, skipping hop-by-hop headers.
func copyHeaders(src, dst http.Header) {
	hopByHop := map[string]bool{
		"Connection":          true,
		"Proxy-Connection":    true,
		"Proxy-Authorization": true,
		"Proxy-Authenticate":  true,
		"Te":                  true,
		"Trailer":             true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
	}
	for k, v := range src {
		if hopByHop[k] {
			continue
		}
		for _, vv := range v {
			dst.Add(k, vv)
		}
	}
}

func isLLMHost(host string) bool {
	h := strings.ToLower(host)
	for _, llm := range llmHosts {
		if strings.Contains(h, llm) {
			return true
		}
	}
	return false
}
