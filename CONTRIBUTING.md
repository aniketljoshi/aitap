# Contributing to aitap

Thanks for wanting to help build aitap.

This project is still small enough that a thoughtful issue, a crisp bug report, or a docs fix can
change the product quickly. You do not need to land a giant feature to make a meaningful
contribution.

Before participating, please read the [Code of Conduct](CODE_OF_CONDUCT.md).

## Good First Contributions

These are all valuable:

- Tighten request or response parsing for a provider
- Improve streaming or token accounting tests
- Fix edge cases in export or redaction
- Polish TUI readability and navigation
- Improve docs, examples, or onboarding
- Add missing provider coverage where the proxy shape is already understood

## Before You Start

For small fixes, feel free to open a PR directly.

For larger changes, please open an issue first so we can align on scope. This helps avoid duplicate
work and makes reviews faster.

Useful starting points:

- [Support Guide](SUPPORT.md)
- [Security Policy](SECURITY.md)
- [GitHub issue templates](https://github.com/aniketljoshi/aitap/issues/new/choose)

## Local Setup

### Prerequisites

- Go 1.22 or newer
- Git

### Clone and build

```bash
git clone https://github.com/aniketljoshi/aitap.git
cd aitap
go build -o bin/aitap ./cmd/aitap
```

### Run tests

```bash
go test ./...
```

### Run locally

```bash
go run ./cmd/aitap
```

If you already use `make`, the repo includes:

```bash
make build
make test
make run
```

## Development Workflow

1. Fork the repo and create a focused branch from `main`.
2. Make the smallest change that fully solves the problem.
3. Add or update tests when behavior changes.
4. Run `go test ./...`.
5. Open a pull request using the repo template.

Recommended branch name patterns:

- `fix/sse-parser-edge-case`
- `feat/openrouter-pricing`
- `docs/readme-quickstart`
- `refactor/proxy-error-handling`

## Project Conventions

### Code style

- Follow standard Go formatting with `gofmt`.
- Prefer small, explicit functions over clever abstractions.
- Keep provider-specific logic easy to trace.
- Add comments only when behavior is not obvious from the code.

### Testing

Please include tests for:

- New request or response parsing behavior
- Streaming edge cases
- Provider detection changes
- Redaction or export behavior

If your change is intentionally hard to test, explain why in the PR.

### Commits

Conventional Commits are preferred:

- `feat: add openrouter forward proxy example`
- `fix: handle empty SSE chunk gracefully`
- `docs: refresh contributing guide`
- `test: cover redaction fallback`

## Areas That Benefit From Extra Care

### Provider and pricing changes

If you update provider support:

- Keep `internal/provider/detect.go` in sync
- Update request or response parsing as needed
- Add tests near the changed parser or provider logic
- Refresh provider details in [README.md](README.md)

### Proxy behavior changes

If you touch `cmd/aitap/proxy.go`, please test both:

- Forward proxy mode
- HTTP proxy mode

Document any known limitation clearly if the behavior is provider-specific.

### TUI changes

The terminal UI should stay readable on smaller terminal sizes and under active streaming load.
If you adjust layout or key handling, mention the UX impact in the PR.

## Pull Request Checklist

Before opening a PR, make sure you can say yes to most of these:

- I read the relevant docs for the area I changed
- I kept the change focused
- I added or updated tests when behavior changed
- I ran `go test ./...`
- I updated docs if the user-facing workflow changed
- I used clear commit and PR descriptions

## Review Expectations

Reviews focus on correctness, clarity, and maintenance cost.

That usually means reviewers will look for:

- Behavioral regressions
- Missing tests
- Unclear naming or control flow
- Docs drift
- Hidden complexity in proxy and parser logic

Small, focused PRs get merged much faster than broad refactors.

## Reporting Bugs and Requesting Features

Please use the issue templates so reports arrive with the context needed to reproduce or evaluate
them:

- Bug report
- Feature request
- Help or setup question

If you are reporting a security issue, do not open a public issue. Follow the instructions in
[SECURITY.md](SECURITY.md).

## Questions

If you are unsure where to start, open a help issue and describe:

- What you are trying to do
- What you already checked
- Where you got stuck

That is enough to start a productive conversation.
