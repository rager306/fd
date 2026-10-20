# Contributing to fd

Thanks for taking the time to contribute. This document covers the fastest paths to a merged change and the gates your change has to pass.

## Ground rules

- Be respectful; the project follows the [Code of Conduct](./CODE_OF_CONDUCT.md).
- Report security issues privately per [SECURITY.md](./SECURITY.md), not in a public issue.
- Keep one logical change per pull request. A fix rides with the change that introduced it when that is possible.
- Do not mix a feature with unrelated lint fixes or refactors: it hides both from review.
- Never commit secrets. `FD_API_KEY`, tokens, and customer text stay out of the repository, tests, and logs.

## Development environment

Prerequisites: Go 1.26+ (the module `api/` declares `go 1.26`), Docker with Compose v2, `golangci-lint` v2 (or let `go run` fetch it), and `git`.

```bash
git clone https://github.com/rager306/fd.git
cd fd
cp .env.example .env          # set FD_API_KEY to a local random value
docker compose up -d          # fd_api + fd_tei + fd_redis
```

The API listens on `http://localhost:8000`. Protected routes need `Authorization: Bearer <FD_API_KEY>`; `/live`, `/ready`, `/health`, `/v1/healthcheck`, `/openapi.json`, and `/docs` stay public. TEI loads a BERT model on CPU at first start and can take minutes under load; the API's warmup recovery loop handles that — no manual restart needed.

## The gates

Run these before you open or update a pull request. CI runs the same ones.

```bash
cd api
go test ./... -short                     # the CI suite
go test ./... -race -count=1             # full suite, includes race-sensitive cache/queue tests
go vet ./...
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

Black-box and integration checks, when your change touches the HTTP contract or lifecycle:

```bash
# from the repository root, with the stack running
FD_API_KEY=$(docker exec fd_api printenv FD_API_KEY)
python3 tools/verify_fd_v2_contract.py
cd tests/integration && FD_INTEGRATION_API_KEY=<same value> go test -v .
```

## What belongs in a pull request

- A clear description: what changed, why, and what you verified. Link the issue it closes (`Closes #123`).
- Tests for behavior changes. Bug fixes need a regression test that fails without the fix.
- A `CHANGELOG.md` entry under `[Unreleased]` when the change is user-visible: new endpoint, new environment variable, changed default, changed error shape, or changed cache-key format.
- Documentation updates when you add or change configuration: `README.md`, `.env.example`, and the same-host contract in `docs/same-host-embedding-service-contract.md`.
- No test for a statically defined value and no test asserting a magic string; prefer behavioral assertions.
- Keep comment noise down: comments explain why, not what the next line does.

## Commit and PR style

- Conventional Commits: `feat(cache): ...`, `fix(api): ...`, `docs: ...`, `chore(deps): ...`, `perf(cache): ...`, `test(queue): ...`.
- Subject in the imperative mood, under ~72 characters; details in the body.
- Do not add `Co-authored-by` trailers or set a git identity for someone else's work.
- Rebase or merge `master` before requesting review if the branch is behind; no merge commits inside a feature branch unless the change needs them.

## Changing configuration or the contract

Configuration is the public surface of this service. When you add or rename an environment variable:

1. Read the value with the helpers in `api/internal/envutil` (strict parsing; a misconfigured value must fail startup rather than silently take a fallback).
2. Document it in the `README.md` configuration tables and add it to `.env.example`.
3. If it changes cache-key identity, say so explicitly and bump `EMBEDDING_CACHE_VERSION` or document the required flush.

`docs/same-host-embedding-service-contract.md` is the canonical consumer contract. New consumer-facing behavior belongs there first.

## Review

- A pull request is ready when CI is green, the description states what was verified, and review comments are addressed or explicitly declined with a reason.
- Expect the maintainer to ask for the verification command you ran — "tests pass" without the command is not enough.
- Large mechanical changes (renames, formatting) should be their own pull request.

## Cutting a release

1. Move the `[Unreleased]` entries in `CHANGELOG.md` under a new version heading with today's date, and update the badge/version references in `README.md` if the release changes them.
2. Tag the commit on `master` and push the tag: `git tag v1.2.3 && git push origin v1.2.3`.
3. The `Release notes` workflow reacts to the tag: it walks the Conventional Commit range with `git-cliff` (configuration in `cliff.toml`) and creates the GitHub release, or updates it if one already exists.
4. Check the generated notes against `CHANGELOG.md` before announcing. `CHANGELOG.md` is authoritative for release history; the workflow only fills the release page.

To preview the notes for the newest tag without publishing anything, run the workflow manually (`gh workflow run release.yml`) — the dispatch path prints the generated notes and stops.

## Working with coding agents

If you use an AI coding agent in this repository, follow `AGENTS.md` at the repository root. It defines the required pre-change caller analysis, the commit gate, and where to file defects in the agent tooling itself.

## License

By contributing you agree that your contribution is licensed under the [MIT License](./LICENSE).
