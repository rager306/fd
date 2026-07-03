---
id: T03
parent: S01
milestone: M007-wvm7b3
key_files:
  - README.md
  - .golangci.yml
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T12:09:28.271Z
blocker_discovered: false
---

# T03: Designed minimal GitHub Actions workflow mirroring local Go test/lint commands.

**Designed minimal GitHub Actions workflow mirroring local Go test/lint commands.**

## What Happened

Recorded the workflow design. The workflow will be `.github/workflows/go-quality.yml`, named `Go Quality`, triggered on `push` and `pull_request` with path filters for `api/**`, `.golangci.yml`, `README.md`, and the workflow file itself. Permissions will be `contents: read`. One `quality` job will run on `ubuntu-latest`, use `actions/checkout@v4`, `actions/setup-go@v5` with Go 1.25.x and cache enabled for `api/go.sum`, run `go test ./... -short` in `api/`, and run the pinned GolangCI-Lint command from M006 with `--config ../.golangci.yml ./...`.

## Verification

Design satisfies docs and local command parity constraints.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `workflow design from docs + README/.golangci.yml inspection` | 0 | ✅ pass: design recorded | 0ms |

## Deviations

Remote CI run verification is intentionally excluded because pushing is an external action requiring explicit confirmation.

## Known Issues

Workflow will be locally validated only until the user explicitly approves push.

## Files Created/Modified

- `README.md`
- `.golangci.yml`
