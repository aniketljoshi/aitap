# Contributing to aitap

First off, thank you for considering contributing to aitap! 🎉

Every contribution matters — whether it's fixing a typo, adding a provider, improving the TUI, or reporting a bug. This guide will help you get started.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Project Structure](#project-structure)
- [Adding a New Provider](#adding-a-new-provider)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [aniket@aitap.dev](mailto:aniket@aitap.dev).

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git

### Setup

```bash
git clone https://github.com/aniketjoshi/aitap.git
cd aitap
make build
make test
```

### Running locally

```bash
# Build and run
make run

# Or build then run separately
make build
./bin/aitap
```

## How to contribute

### Reporting bugs

Found a bug? [Open an issue](https://github.com/aniketjoshi/aitap/issues/new?template=bug_report.yml) with:

- What you expected to happen
- What actually happened
- Steps to reproduce
- Your OS, Go version, and aitap version (`aitap --version`)

### Suggesting features

Have an idea? [Open a feature request](https://github.com/aniketjoshi/aitap/issues/new?template=feature_request.yml). Before creating one, please check if a similar request already exists.

### Submitting code

1. **Fork** the repo and create a branch from `main`
2. **Name your branch** descriptively: `fix/streaming-parse-error`, `feat/azure-openai-support`, `docs/update-readme`
3. **Write tests** for any new functionality
4. **Run the full check** before submitting:

```bash
make test       # run tests
make build      # verify build
```

5. **Open a Pull Request** using our [PR template](.github/PULL_REQUEST_TEMPLATE.md)

## Development Workflow

```bash
# 1. Fork and clone
git clone https://github.com/<your-username>/aitap.git
cd aitap

# 2. Create a feature branch
git checkout -b feat/my-awesome-feature

# 3. Make your changes, write tests

# 4. Run checks
make test
make build

# 5. Commit with a descriptive message
git commit -m "feat: add Azure OpenAI provider support"

# 6. Push and open a PR
git push origin feat/my-awesome-feature
```

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Keep functions small and focused
- Add comments for non-obvious logic
- Use descriptive variable names
- Run `gofmt -s -w .` before committing

## Project Structure

```
aitap/
  cmd/aitap/         # Entry point and proxy server
    main.go          # CLI flags, startup, shutdown
    proxy.go         # HTTP proxy (forward + HTTP_PROXY modes)
  internal/
    model/           # Data types (Call, Session, Provider)
    parser/          # Request/response parsing per provider
      parse.go       # Non-streaming parsers
      sse.go         # SSE streaming parsers
    provider/        # Provider detection and pricing
    redact/          # Secret redaction
    export/          # JSONL export
    tui/             # Bubble Tea terminal UI
```

## Adding a New Provider

1. Add the provider constant in `internal/model/call.go`
2. Add host detection in `internal/provider/detect.go`
3. Add pricing in `internal/provider/detect.go`
4. Add request/response parser in `internal/parser/parse.go`
5. Add SSE parser in `internal/parser/sse.go`
6. Add the upstream route in `cmd/aitap/proxy.go`
7. Add tests for all the above
8. Update the README provider table

## Adding Pricing for New Models

Edit `internal/provider/detect.go` and add to the `Pricing` map. Use per-1M-token pricing from the provider's official page.

## Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/). This makes changelogs and release notes easier to generate.

| Prefix | Purpose |
| :--- | :--- |
| `feat:` | New feature |
| `fix:` | Bug fix |
| `docs:` | Documentation only |
| `test:` | Adding or updating tests |
| `refactor:` | Code change that neither fixes a bug nor adds a feature |
| `chore:` | Build process, CI, or tooling changes |

**Examples:**

```
feat: add Azure OpenAI provider support
fix: handle empty SSE chunks in streaming parser
docs: update README with Google provider setup
test: add coverage for redaction edge cases
```

## Pull Request Process

1. Ensure your branch is up to date with `main`
2. Fill out the [PR template](.github/PULL_REQUEST_TEMPLATE.md) completely
3. Ensure all tests pass (`make test`)
4. Keep PRs focused — one feature or fix per PR
5. Update documentation if behavior changes
6. Be patient — reviews may take a few days

> [!TIP]
> Small, focused PRs are reviewed faster than large ones.

## Questions?

Open a [discussion](https://github.com/aniketjoshi/aitap/discussions) or reach out in an issue. No question is too small.

---

Thank you for helping make aitap better!
