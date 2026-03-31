# X / Twitter launch posts

---

## Main launch tweet

I built a local traffic inspector for LLMs.

One binary. Zero config. Beautiful TUI.

See every prompt, response, token count, and cost from OpenAI, Claude, Gemini, and Ollama — while your agent is running.

aitap --demo to try it instantly.

github.com/aniketljoshi/aitap

[attach demo.gif]

---

## Thread (reply to main tweet)

### Tweet 2

How it works:

Point your SDK at http://localhost:9119/openai/v1

aitap forwards to the real API over HTTPS and shows every call in your terminal.

No certificates. No cloud account. No SDK wrapping.

### Tweet 3

What you see:

- Provider, model, tokens in/out, cost, latency
- Full prompt and response (expand any call)
- Streaming-aware — SSE chunks parsed in real-time
- Session totals so you know your total spend

### Tweet 4

Built with:

- Go (single binary)
- Bubble Tea (terminal UI)
- Lipgloss (styling)
- Zero runtime dependencies

Supports: OpenAI, Anthropic, Google, Ollama, OpenRouter

### Tweet 5

Why I built it:

Every observability tool for LLMs wants you to:
- Create an account
- Add an API key
- Install an SDK
- Configure a dashboard

I just wanted to see what my agent was sending.

So I built a tool that starts in 1 second and shows me.

### Tweet 6

Try it:

go install github.com/aniketljoshi/aitap/cmd/aitap@latest

Or just:

aitap --demo

to see sample traffic without making any API calls.

Star it if this is useful: github.com/aniketljoshi/aitap

---

## Shorter standalone tweet (for later repost)

"My agent is doing something weird" is a daily feeling.

aitap sits between your app and the LLM, shows every call in a terminal UI.

One binary. Zero setup. Works with OpenAI, Claude, Gemini, Ollama.

github.com/aniketljoshi/aitap

[attach demo-short.gif]

---

## Engagement tweet (for after initial traction)

People keep asking what LLM observability tool to use.

If you just need to see what's happening right now:

aitap

- No dashboard
- No account
- No API key
- Just a terminal and the truth

[attach screenshot of TUI with 4-5 calls]
