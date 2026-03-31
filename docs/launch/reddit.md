# Reddit launch posts

---

## r/LocalLLaMA

### Title

I built a terminal UI that shows you exactly what your LLM agent is sending and receiving — one binary, zero config

### Body

I kept running into the same problem: my agent would behave unexpectedly and I had no quick way
to see the actual API traffic without digging through logs or adding print statements everywhere.

So I built **aitap** — a single Go binary that sits between your app and the LLM API and shows
every call in a terminal UI.

**How to use it:**

```bash
# Install
go install github.com/aniketljoshi/aitap/cmd/aitap@latest

# Start
aitap

# Point your SDK at it
export OPENAI_BASE_URL=http://localhost:9119/openai/v1
# or
export OLLAMA_HOST=http://localhost:9119/ollama

# Run your app as usual
python my_agent.py
```

**What you see:**

- Every call: provider, model, input/output tokens, cost, latency
- Expand any call to see the full prompt and response
- Works with streaming (SSE) — tokens are parsed as they arrive
- Session totals so you know your total spend

**Works with:** OpenAI, Anthropic, Google, Ollama, OpenRouter

**What makes it different from Langfuse/LangSmith:**

- No cloud account, no API key, no dashboard
- One binary, runs locally
- Starts in 1 second, not 10 minutes of setup
- Your data never leaves your machine

Try `aitap --demo` to see it with sample traffic without calling any API.

GitHub: https://github.com/aniketljoshi/aitap

[GIF/screenshot here]

Would love feedback — especially from people running local models with Ollama. That's the
zero-friction path since it's already HTTP.

---

## r/ChatGPTCoding

### Title

Built a "Wireshark for LLM calls" — see exactly what your AI coding agent sends to OpenAI/Claude/Gemini

### Body

If you use Claude Code, Codex CLI, aider, or any coding agent — you've probably wondered
what's actually being sent to the API.

I built **aitap**: a terminal tool that intercepts LLM API calls and shows them in real-time.

```
aitap  :9119  |  4 calls  |  8.2k tokens  |  $0.057

  ~   1 | anthropic | claude-sonnet-4...   |   1.2k>890 | $0.012 |  2.3s
      2 | openai    | gpt-4o               |    340>210 | $0.004 |  1.1s
  ~   3 | google    | gemini-2.5-pro       |  4.1k>1.5k | $0.041 |  4.7s
      4 | ollama    | llama3:8b            |    890>450 |   free |  3.2s
```

- Single binary, no setup, no cloud account
- Shows prompts, responses, tokens, cost, latency
- Works with streaming responses
- Supports OpenAI, Anthropic, Google, Ollama, OpenRouter

Point your SDK's base URL at `http://localhost:9119/openai/v1` and you can watch every call
your agent makes.

GitHub: https://github.com/aniketljoshi/aitap

[GIF/screenshot here]

---

## r/golang

### Title

aitap: a Bubble Tea TUI for inspecting LLM API traffic in real-time (forward proxy + SSE parsing)

### Body

I built a developer tool in Go that acts as a local proxy for LLM API calls. It parses
request/response bodies (including SSE streaming) for OpenAI, Anthropic, Google, and Ollama
and displays them in a Bubble Tea terminal UI.

**Architecture highlights:**

- Forward proxy mode: SDK points base URL at localhost, aitap forwards to upstream over HTTPS.
  No certificates, no MITM. Clean separation.
- SSE streaming: buffered line scanner parses `data:` prefixed chunks per-provider (OpenAI
  content deltas, Anthropic message_start/content_block_delta, Ollama NDJSON, Google cumulative
  usageMetadata)
- Provider-aware: detects provider from path prefix or hostname, routes to correct parser
- Zero external runtime deps: single binary, pure Go build
- Bubble Tea + Lipgloss for the TUI with provider color coding

**What I'd love feedback on:**

- The proxy design in `cmd/aitap/proxy.go` — handling both forward proxy (path prefix) and
  HTTP_PROXY (absolute URL) modes in the same handler
- SSE parsing approach in `internal/parser/sse.go` — each provider has different streaming
  event shapes
- Whether the `estimateTokens()` fallback (chars/4) is good enough when providers don't include
  usage in stream chunks

GitHub: https://github.com/aniketljoshi/aitap

`go install github.com/aniketljoshi/aitap/cmd/aitap@latest`

[GIF here]
