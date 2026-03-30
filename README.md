# aitap

**mitmproxy for LLMs, but pleasant.**

A single-binary terminal UI that intercepts and inspects LLM API calls in real-time. See every prompt, response, token count, latency, and cost — without changing your code.

```
brew install aniketjoshi/tap/aitap    # macOS
go install github.com/aniketjoshi/aitap/cmd/aitap@latest
```

## Quick start

### Forward proxy mode (recommended)

Point your SDK's base URL at aitap. No certificates, no MITM — just plain HTTP locally.

```bash
# Start aitap
aitap

# In another terminal, configure your SDK:

# OpenAI
export OPENAI_BASE_URL=http://localhost:9119/openai/v1

# Anthropic
export ANTHROPIC_BASE_URL=http://localhost:9119/anthropic

# Ollama
export OLLAMA_HOST=http://localhost:9119/ollama

# Google
export GOOGLE_API_BASE=http://localhost:9119/google

# Then run your app as usual
python my_agent.py
```

### HTTP proxy mode

Works for HTTP traffic (e.g., local Ollama).

```bash
aitap

# In another terminal
export HTTP_PROXY=http://127.0.0.1:9119
python my_agent.py
```

## What you see

```
 aitap  :9119  |  4 calls  |  8.2k tokens  |  $0.057

  ~ 1 | anthropic  | claude-sonnet-4-...    |  1.2k>890 |  $0.012 |  2.3s
    2 | openai     | gpt-4o                 |   340>210 |  $0.004 |  1.1s
  ~ 3 | anthropic  | claude-sonnet-4-...    | 4.1k>1.5k |  $0.041 |  4.7s
  ~ 4 | ollama     | llama3:8b              |   890>450 |    free |  3.2s

  -- Request --
  system: You are a helpful coding assistant.
  user: Summarize the latest PR changes and suggest...

  -- Response --
  assistant: Based on the diff, here are the key changes:...

  status=200  in=890  out=450  cost=free  latency=3.2s

  j/k navigate  enter expand  q quit
```

## Supported providers

| Provider   | Chat | Streaming (SSE) | Cost tracking |
|------------|------|-----------------|---------------|
| OpenAI     | Yes  | Yes             | Yes           |
| Anthropic  | Yes  | Yes             | Yes           |
| Google     | Yes  | Yes             | Yes           |
| Ollama     | Yes  | Yes             | free          |
| OpenRouter | Yes  | Yes             | Yes           |

## Features

- **Two proxy modes** — forward proxy (base URL rewrite, zero certs) or HTTP proxy (env var)
- **Single binary** — no runtime, no dependencies, no cloud account
- **Live TUI** — see calls as they happen, expand any to inspect request/response
- **Streaming support** — parses SSE chunks in real-time, accumulates tokens and text
- **Cost tracking** — per-call and session totals with current model pricing
- **Secret redaction** — API keys, bearer tokens, AWS keys automatically masked in exports
- **Session export** — save to JSONL for sharing, debugging, or post-analysis
- **Provider filtering** — focus on just one provider at a time

## CLI

```
aitap                         # start on :9119
aitap --port 8080             # custom port
aitap --export session.jsonl  # auto-export on exit
aitap --redact                # mask secrets in export
aitap --filter anthropic      # only show Anthropic calls
aitap --version               # show version
```

## How it works

**Forward proxy mode** (recommended): Your SDK sends requests to `http://localhost:9119/openai/v1/chat/completions`. aitap strips the provider prefix, forwards to `https://api.openai.com/v1/chat/completions`, and streams the response back while capturing metadata. No certificates needed because your SDK talks HTTP to localhost, and aitap talks HTTPS to the upstream API.

**HTTP proxy mode**: Set `HTTP_PROXY` env var. Works for plain HTTP traffic like local Ollama. For HTTPS APIs via this mode, traffic passes through without inspection (CONNECT tunnel).

aitap parses request/response bodies for known LLM API formats, extracts metadata (model, tokens, timing, cost), and displays everything in a navigable terminal UI.

**aitap never stores or transmits your data.** Everything stays in memory unless you explicitly export with `--export`.

## Why not X?

| Tool | Difference |
|------|-----------|
| Langfuse / Phoenix | Cloud platforms with dashboards. aitap is a local dev tool — no account, no API key. |
| llm-interceptor | Requires mitmproxy + certificate setup. aitap is one binary, zero config. |
| llm.log | SQLite-focused logging. aitap is TUI-first with live streaming inspection. |
| LiteLLM | Full proxy/gateway platform. aitap is read-only, zero config, zero overhead. |

## Development

```bash
git clone https://github.com/aniketjoshi/aitap.git
cd aitap
make build    # build binary
make test     # run tests
make run      # build and run
```

## License

MIT
