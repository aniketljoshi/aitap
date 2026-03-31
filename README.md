<div align="center">

# 🔍 aitap

### mitmproxy for LLMs, but pleasant.

A single-binary terminal UI that intercepts and inspects LLM API calls in real-time.<br/>
See every prompt, response, token count, latency, and cost — **without changing your code.**

[![Go Version](https://img.shields.io/github/go-mod/go-version/aniketjoshi/aitap?style=flat-square&color=00ADD8)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/aniketjoshi/aitap?style=flat-square&color=orange)](https://github.com/aniketjoshi/aitap/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/aniketjoshi/aitap?style=flat-square)](https://goreportcard.com/report/github.com/aniketjoshi/aitap)

<br/>

```
 aitap  :9119  |  4 calls  |  8.2k tokens  |  $0.057

  ~ 1 │ anthropic  │ claude-sonnet-4-...    │  1.2k›890 │  $0.012 │  2.3s
    2 │ openai     │ gpt-4o                 │   340›210 │  $0.004 │  1.1s
  ~ 3 │ anthropic  │ claude-sonnet-4-...    │ 4.1k›1.5k │  $0.041 │  4.7s
  ~ 4 │ ollama     │ llama3:8b              │   890›450 │    free │  3.2s
```

</div>

---

## Table of Contents

- [Quick Start](#quick-start)
- [What You See](#what-you-see)
- [Supported Providers](#supported-providers)
- [Features](#features)
- [CLI Reference](#cli)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

## Quick Start

### Install

```bash
# Homebrew (macOS / Linux)
brew install aniketjoshi/tap/aitap

# Go install
go install github.com/aniketjoshi/aitap/cmd/aitap@latest

# Or build from source
git clone https://github.com/aniketjoshi/aitap.git
cd aitap && make build
```

### Forward proxy mode (recommended)

Point your SDK's base URL at aitap. No certificates, no MITM — just plain HTTP locally.

```bash
# Terminal 1 — start aitap
aitap

# Terminal 2 — configure your SDK and run your app
export OPENAI_BASE_URL=http://localhost:9119/openai/v1       # OpenAI
export ANTHROPIC_BASE_URL=http://localhost:9119/anthropic     # Anthropic
export OLLAMA_HOST=http://localhost:9119/ollama               # Ollama
export GOOGLE_API_BASE=http://localhost:9119/google           # Google

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

## What You See

```
 aitap  :9119  |  4 calls  |  8.2k tokens  |  $0.057

  ~ 1 │ anthropic  │ claude-sonnet-4-...    │  1.2k›890 │  $0.012 │  2.3s
    2 │ openai     │ gpt-4o                 │   340›210 │  $0.004 │  1.1s
  ~ 3 │ anthropic  │ claude-sonnet-4-...    │ 4.1k›1.5k │  $0.041 │  4.7s
  ~ 4 │ ollama     │ llama3:8b              │   890›450 │    free │  3.2s

  ── Request ──
  system: You are a helpful coding assistant.
  user: Summarize the latest PR changes and suggest...

  ── Response ──
  assistant: Based on the diff, here are the key changes:...

  status=200  in=890  out=450  cost=free  latency=3.2s

  j/k navigate  enter expand  q quit
```

## Supported Providers

| Provider | Chat | Streaming (SSE) | Cost Tracking |
| :--- | :---: | :---: | :---: |
| **OpenAI** | ✅ | ✅ | ✅ |
| **Anthropic** | ✅ | ✅ | ✅ |
| **Google** | ✅ | ✅ | ✅ |
| **Ollama** | ✅ | ✅ | free |
| **OpenRouter** | ✅ | ✅ | ✅ |

## Features

| | |
| :--- | :--- |
| 🔀 **Two proxy modes** | Forward proxy (base URL rewrite, zero certs) or HTTP proxy (env var) |
| 📦 **Single binary** | No runtime, no dependencies, no cloud account |
| 🖥️ **Live TUI** | See calls as they happen, expand any to inspect request/response |
| 🌊 **Streaming support** | Parses SSE chunks in real-time, accumulates tokens and text |
| 💰 **Cost tracking** | Per-call and session totals with current model pricing |
| 🔒 **Secret redaction** | API keys, bearer tokens, AWS keys automatically masked in exports |
| 💾 **Session export** | Save to JSONL for sharing, debugging, or post-analysis |
| 🎯 **Provider filtering** | Focus on just one provider at a time |

## CLI

```
aitap                         # start on :9119
aitap --port 8080             # custom port
aitap --export session.jsonl  # auto-export on exit
aitap --redact                # mask secrets in export
aitap --filter anthropic      # only show Anthropic calls
aitap --version               # show version
```

## How It Works

```
   Your App                aitap (:9119)              LLM API
  ┌─────────┐   HTTP     ┌──────────────┐   HTTPS   ┌──────────┐
  │  SDK /   │──────────▶│  intercept   │──────────▶│  OpenAI  │
  │  Script  │◀──────────│  + display   │◀──────────│  Claude   │
  └─────────┘            └──────────────┘            └──────────┘
                               │
                          ┌────▼────┐
                          │  Live   │
                          │  TUI    │
                          └─────────┘
```

**Forward proxy mode** (recommended): Your SDK sends requests to `http://localhost:9119/openai/v1/chat/completions`. aitap strips the provider prefix, forwards to `https://api.openai.com/v1/chat/completions`, and streams the response back while capturing metadata. No certificates needed because your SDK talks HTTP to localhost, and aitap talks HTTPS to the upstream API.

**HTTP proxy mode**: Set `HTTP_PROXY` env var. Works for plain HTTP traffic like local Ollama. For HTTPS APIs via this mode, traffic passes through without inspection (CONNECT tunnel).

> [!NOTE]
> **aitap never stores or transmits your data.** Everything stays in memory unless you explicitly export with `--export`.

## Architecture

```
aitap/
├── cmd/aitap/           # Entry point and proxy server
│   ├── main.go          # CLI flags, startup, shutdown
│   └── proxy.go         # HTTP proxy (forward + HTTP_PROXY modes)
└── internal/
    ├── model/           # Data types (Call, Session, Provider)
    ├── parser/          # Request/response parsing per provider
    │   ├── parse.go     # Non-streaming parsers
    │   └── sse.go       # SSE streaming parsers
    ├── provider/        # Provider detection and pricing
    ├── redact/          # Secret redaction
    ├── export/          # JSONL export
    └── tui/             # Bubble Tea terminal UI
```

## Why Not X?

| Tool | Difference |
| :--- | :--- |
| **Langfuse / Phoenix** | Cloud platforms with dashboards. aitap is a local dev tool — no account, no API key. |
| **llm-interceptor** | Requires mitmproxy + certificate setup. aitap is one binary, zero config. |
| **llm.log** | SQLite-focused logging. aitap is TUI-first with live streaming inspection. |
| **LiteLLM** | Full proxy/gateway platform. aitap is read-only, zero config, zero overhead. |

## Contributing

We welcome contributions of all kinds! Please read our [Contributing Guide](CONTRIBUTING.md) to get started.

See also:

- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security Policy](SECURITY.md)

## Development

```bash
git clone https://github.com/aniketjoshi/aitap.git
cd aitap
make build    # build binary
make test     # run tests
make run      # build and run
```

## License

[MIT](LICENSE) — made with ☕ by [Aniket Joshi](https://github.com/aniketjoshi)

---

<div align="center">

**[Report Bug](https://github.com/aniketjoshi/aitap/issues/new?template=bug_report.yml)** · **[Request Feature](https://github.com/aniketjoshi/aitap/issues/new?template=feature_request.yml)** · **[Discussions](https://github.com/aniketjoshi/aitap/discussions)**

If aitap helps your workflow, consider giving it a ⭐

</div>
