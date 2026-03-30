package parser

import (
	"encoding/json"

	"github.com/aniketjoshi/aitap/internal/model"
	"github.com/aniketjoshi/aitap/internal/provider"
)

// ParseRequest extracts model and messages from a request body.
func ParseRequest(p model.Provider, body []byte, call *model.Call) {
	switch p {
	case model.ProviderOpenAI, model.ProviderOpenRouter:
		parseOpenAIRequest(body, call)
	case model.ProviderAnthropic:
		parseAnthropicRequest(body, call)
	case model.ProviderOllama:
		parseOllamaRequest(body, call)
	case model.ProviderGoogle:
		parseGoogleRequest(body, call)
	}
}

// ParseResponse extracts tokens, response text, and calculates cost.
func ParseResponse(p model.Provider, body []byte, call *model.Call) {
	switch p {
	case model.ProviderOpenAI, model.ProviderOpenRouter:
		parseOpenAIResponse(body, call)
	case model.ProviderAnthropic:
		parseAnthropicResponse(body, call)
	case model.ProviderOllama:
		parseOllamaResponse(body, call)
	case model.ProviderGoogle:
		parseGoogleResponse(body, call)
	}

	call.EstimatedCost = provider.EstimateCost(call.Model, call.InputTokens, call.OutputTokens)
}

// --- OpenAI format ---

func parseOpenAIRequest(body []byte, call *model.Call) {
	var req struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		Tools []json.RawMessage `json:"tools"`
		Stream bool `json:"stream"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return
	}
	call.Model = req.Model
	call.IsStreaming = req.Stream
	call.Tools = len(req.Tools)
	for _, m := range req.Messages {
		call.Messages = append(call.Messages, model.Message{
			Role:    m.Role,
			Content: truncate(m.Content, 500),
		})
	}
}

func parseOpenAIResponse(body []byte, call *model.Call) {
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return
	}
	call.InputTokens = resp.Usage.PromptTokens
	call.OutputTokens = resp.Usage.CompletionTokens
	if len(resp.Choices) > 0 {
		call.ResponseText = truncate(resp.Choices[0].Message.Content, 500)
	}
}

// --- Anthropic format ---

func parseAnthropicRequest(body []byte, call *model.Call) {
	var req struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content json.RawMessage `json:"content"`
		} `json:"messages"`
		System string            `json:"system"`
		Tools  []json.RawMessage `json:"tools"`
		Stream bool              `json:"stream"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return
	}
	call.Model = req.Model
	call.IsStreaming = req.Stream
	call.SystemPrompt = truncate(req.System, 300)
	call.Tools = len(req.Tools)
	for _, m := range req.Messages {
		content := extractContentString(m.Content)
		call.Messages = append(call.Messages, model.Message{
			Role:    m.Role,
			Content: truncate(content, 500),
		})
	}
}

func parseAnthropicResponse(body []byte, call *model.Call) {
	var resp struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return
	}
	call.InputTokens = resp.Usage.InputTokens
	call.OutputTokens = resp.Usage.OutputTokens
	if len(resp.Content) > 0 {
		call.ResponseText = truncate(resp.Content[0].Text, 500)
	}
}

// --- Ollama format ---

func parseOllamaRequest(body []byte, call *model.Call) {
	var req struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		Stream *bool `json:"stream"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return
	}
	call.Model = req.Model
	if req.Stream != nil {
		call.IsStreaming = *req.Stream
	}
	for _, m := range req.Messages {
		call.Messages = append(call.Messages, model.Message{
			Role:    m.Role,
			Content: truncate(m.Content, 500),
		})
	}
}

func parseOllamaResponse(body []byte, call *model.Call) {
	var resp struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		PromptEvalCount int `json:"prompt_eval_count"`
		EvalCount       int `json:"eval_count"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return
	}
	call.InputTokens = resp.PromptEvalCount
	call.OutputTokens = resp.EvalCount
	call.ResponseText = truncate(resp.Message.Content, 500)
}

// --- Google format (simplified) ---

func parseGoogleRequest(body []byte, call *model.Call) {
	var req struct {
		Contents []struct {
			Role  string `json:"role"`
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"contents"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return
	}
	call.Model = "gemini" // extracted from URL path typically
	for _, c := range req.Contents {
		text := ""
		for _, p := range c.Parts {
			text += p.Text
		}
		call.Messages = append(call.Messages, model.Message{
			Role:    c.Role,
			Content: truncate(text, 500),
		})
	}
}

func parseGoogleResponse(body []byte, call *model.Call) {
	var resp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
		} `json:"usageMetadata"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return
	}
	call.InputTokens = resp.UsageMetadata.PromptTokenCount
	call.OutputTokens = resp.UsageMetadata.CandidatesTokenCount
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		call.ResponseText = truncate(resp.Candidates[0].Content.Parts[0].Text, 500)
	}
}

// --- Helpers ---

func extractContentString(raw json.RawMessage) string {
	// Anthropic content can be string or array of content blocks
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &blocks); err == nil {
		for _, b := range blocks {
			if b.Type == "text" {
				return b.Text
			}
		}
	}
	return string(raw)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
