## What changed

<!-- What behavior changed, and why. Link the issue: Closes #123 -->

## How it was verified

<!-- The exact commands you ran and their result. "Tests pass" is not enough. -->

```bash
cd api && go test ./... -short
```

## Checklist

- [ ] One logical change; a fix rides with the change that introduced it.
- [ ] `go test ./... -short` passes; `-race` run for cache, queue, or coalescer changes.
- [ ] `go vet ./...` and `golangci-lint` (v2.12.2, `../.golangci.yml`) are clean.
- [ ] `govulncheck ./...` reports no reachable vulnerabilities.
- [ ] Behavior change has a test; a bug fix has a regression test that fails without it.
- [ ] User-visible change is in `CHANGELOG.md` under `[Unreleased]`.
- [ ] New or changed configuration is documented in `README.md` and `.env.example`.
- [ ] No secrets, customer text, or embeddings in code, tests, fixtures, or logs.
- [ ] Callers of changed symbols were checked (`go test` covers it where tests exist).
