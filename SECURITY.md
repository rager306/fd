# Security Policy

## Supported versions

| Version | Supported |
|---------|-----------|
| 1.0.x (current release line) | ✅ |
| older `master` snapshots | ❌ |

`fd` is developed on `master`; fixes land there and are tagged as releases. Only the latest release line receives security fixes.

## Reporting a vulnerability

**Preferred: GitHub private vulnerability reporting.** Open <https://github.com/rager306/fd/security/advisories/new> — the report stays private until a fix is published. If you cannot use it, email <alexey.podolinniy@gmail.com> with `fd security` in the subject.

Include:

- affected version, tag, or commit SHA;
- the exact request or command that triggers the issue, and the configuration in use (environment variables with secret values removed);
- what an attacker gains — for example: reading cached embeddings belonging to another user, bypassing `FD_API_KEY`, exhausting memory, or degrading TEI availability;
- any evidence you already have (logs, traces, response bodies).

Do not include real API keys, real customer text, or production embeddings in a report.

## What to expect

- Acknowledgement within **5 business days** (best effort: this project has a single maintainer).
- An assessment and severity call, then a fix or a documented decision not to fix.
- Coordinated disclosure: please give us a chance to ship a fix before publishing. We will credit you in `CHANGELOG.md` and the release notes if you want the credit.
- Target timelines: critical issues as soon as practical; lower severity in the next regular release.

## Scope

In scope:

- the `fd` API service in `api/` — authentication, request validation, cache-key isolation, error handling, lifecycle gates;
- the container/compose configuration in this repository;
- anything that lets one same-host consumer read or corrupt another's embeddings.

Out of scope:

- deployments operated by other people — test against your own instance only;
- third-party components: [HuggingFace TEI](https://github.com/huggingface/text-embeddings-inference), [Redis](https://github.com/redis/redis), and the Go toolchain report their own vulnerabilities upstream;
- denial-of-service testing against shared or production hosts. A local `docker compose` stack is fine;
- findings that require the operator to deliberately disable a guard (for example setting `REDIS_CACHE_NO_EXPIRE` or running without `FD_API_KEY`) unless the guard is documented as safe by default and is not.

## Hardening already in place

- `FD_API_KEY` is required for protected endpoints; unset means protected routes fail closed with `401`.
- API keys are compared in constant time over equal-length SHA-256 digests.
- Request bodies are capped (10 MB / 128 inputs / 1.2M characters), and `/v1/embeddings` runs behind a lifecycle capacity gate.
- Cache keys are namespaced by `EMBEDDING_CACHE_VERSION` and the model/tokenizer/chunking namespace components, so cached vectors cannot silently cross models.
- `SetTrustedProxies(nil)` is set: forwarded headers are not trusted by default.
- CI runs `go test`, `golangci-lint` (18 linters, includes `gosec`), and `govulncheck` on every push and pull request to `master`.

## Operational notes for operators

- Never commit or log secrets; `FD_API_KEY` belongs in `.env` (gitignored) or the shell environment.
- Bind the service to the same host or a trusted network segment: `fd` is designed for same-host consumers, not public internet exposure.
- Redis is bound to `127.0.0.1` by default in `docker-compose.override.yaml`. Keep it that way unless the network is trusted.
