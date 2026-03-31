# Hacker News — Show HN post

## Title

Show HN: aitap -- a local traffic inspector for LLMs with a terminal UI

## URL

https://github.com/aniketljoshi/aitap

## Text (for the first comment — post this immediately after submitting)

---

Hi HN, I built aitap because I kept hitting the same problem: my agent would do something
unexpected and I had no quick way to see what was actually being sent to the LLM.

Existing options are either full observability platforms that need a cloud account and API keys
(Langfuse, LangSmith, Phoenix), or MITM proxies that require certificate setup (llm-interceptor).
I wanted something I could start in one command and get immediate visibility.

aitap is a single Go binary that sits between your app and the LLM provider. It shows every
request and response in a terminal UI — model, tokens, cost, latency, and the full
prompt/response when you expand a call.

**How it works:**

It has two modes. The recommended one is forward proxy: you point your SDK's base URL at
`http://localhost:9119/openai/v1` (or `/anthropic`, `/google`, `/ollama`). aitap strips the
prefix, forwards to the real API over HTTPS, and streams the response back. No TLS interception,
no certificates.

**What it shows:**

- Live call list with provider, model, token counts, estimated cost, latency
- Streaming-aware — SSE chunks are parsed as they arrive
- Expandable detail pane with full request/response
- Session totals (total spend, total tokens)
- JSONL export for sharing or post-analysis

**What it doesn't do:**

- No cloud. No dashboard. No account.
- No SDK wrapping — works at the HTTP layer
- No data leaves your machine unless you explicitly export
- No eval, no scoring, no replay — just inspection

**Stack:** Go, Bubble Tea (TUI), single binary, zero runtime dependencies.

Try it: `go install github.com/aniketljoshi/aitap/cmd/aitap@latest`

Or `aitap --demo` to see sample traffic without making any API calls.

Happy to answer questions about the architecture, proxy design, or SSE parsing approach.

---

## Timing notes

- Best posting times for Show HN: Tuesday-Thursday, 8-10am ET
- Reply to every comment in the first 2 hours
- Be honest about limitations (no HTTPS MITM yet, pricing may drift)
