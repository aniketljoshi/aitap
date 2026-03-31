# Security Policy

Security reports are welcome and appreciated.

If you believe you have found a vulnerability in aitap, please do **not** open a public GitHub
issue or pull request. Report it privately so there is time to investigate and ship a fix safely.

## Supported Versions

| Version | Status |
| --- | --- |
| `main` | Supported |
| Latest release | Supported |
| Older releases | Best effort only |

## How To Report

Please email [aniket@aitap.dev](mailto:aniket@aitap.dev) with:

- A short summary of the issue
- The affected version, commit, or branch if known
- Step-by-step reproduction details
- Expected impact
- Any proof-of-concept material or logs that help validate the report

If the issue involves secrets, authentication headers, or private prompts, please redact or
minimize sensitive data wherever possible.

## What Happens Next

The maintainer will aim to:

- Acknowledge receipt within 72 hours
- Triage and confirm severity within 7 days
- Keep you updated on remediation progress
- Credit you for the report if you want public acknowledgment

## Scope

Examples of security-sensitive areas in this repo include:

- Proxy forwarding behavior
- Request and response parsing
- Secret handling and redaction
- Exported session data
- Accidental leakage of headers, API keys, or prompts

## Responsible Disclosure

Please avoid:

- Opening a public issue before a fix is ready
- Sharing exploit details broadly before maintainers have had time to respond
- Accessing or retaining data that is not yours

Good-faith research that helps protect users is appreciated.
