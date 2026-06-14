---
id: S02
parent: M041-4tw0w7
milestone: M041-4tw0w7
provides:
  - Lifecycle primitives and probes for downstream S03 observability work.
  - Lifecycle gate and capacity overload behavior for `/v1/embeddings`.
  - Graceful shutdown orchestration for deployment/runtime validation.
requires:
  []
affects:
  []
key_files:
  - api/lifecycle/state.go
  - api/lifecycle/warmup.go
  - api/lifecycle/shutdown.go
  - api/middleware/lifecycle.go
  - api/handlers/probes.go
  - api/main.go
  - api/fd_v2_lifecycle_integration_test.go
key_decisions:
  - D044: `FD_MAX_IN_FLIGHT` provides default-off overload behavior for F-2 without a worker-pool redesign.
patterns_established:
  - Lifecycle state is the single source for warmup readiness, shutdown, drain, and last error.
  - Lifecycle-gated embedding routes should run validation first, then readiness/shutdown/capacity gating, then inference.
  - M043 Go gates remain mandatory for Go slices: `go test ./...`, golangci-lint v2.12.2, and govulncheck.
observability_surfaces:
  - Structured logs for warmup start/done/fail and shutdown signal/complete/force-close.
  - /live and /ready probes with canonical lifecycle error envelopes.
  - Lifecycle last error in `State` for warmup/runtime failure diagnosis.
drill_down_paths:
  - .gsd/milestones/M041-4tw0w7/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M041-4tw0w7/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M041-4tw0w7/slices/S02/tasks/T03-SUMMARY.md
  - .gsd/milestones/M041-4tw0w7/slices/S02/tasks/T04-SUMMARY.md
  - .gsd/milestones/M041-4tw0w7/slices/S02/tasks/T05-SUMMARY.md
  - .gsd/milestones/M041-4tw0w7/slices/S02/tasks/T06-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-06-14T06:20:18.881Z
blocker_discovered: false
---

# S02: Lifecycle warmup readiness and graceful shutdown

**Implemented fd lifecycle readiness: model pre-warm, /live and /ready probes, /v1/embeddings lifecycle gating, graceful SIGTERM/SIGINT drain, and executable F-1/F-2/F-5 integration coverage.**

## What Happened

S02 built a complete lifecycle layer around fd's embedding runtime. T01 introduced `lifecycle.State` with readiness, shutdown, in-flight tracking, drain waiting, last error, context helpers, and the default singleton. T02 added `lifecycle.PreWarm` and asynchronous startup warmup in `main.go`, leaving readiness false until a dummy embedding succeeds and recording/logging warmup errors. T03 added `/live` and `/ready`: liveness is process-only, readiness is lifecycle-backed and returns `model_not_loaded` with `Retry-After: 5` before warmup. T04 added `/v1/embeddings` lifecycle middleware so unready/shutting-down requests fail before inference and accepted requests are tracked. T05 added lifecycle-owned SIGTERM/SIGINT graceful shutdown with a shared 30s server/in-flight drain deadline and force-close error path. T06 added executable integration-style tests for startup sequence, F-1, F-2, and F-5; to make F-2 reachable, it added `FD_MAX_IN_FLIGHT` as a default-off max in-flight gate that emits `model_overloaded`.

## Verification

All task-level evidence passed. Final T06 evidence includes `go test . -run TestFdV2Lifecycle -v`, `go test ./...`, golangci-lint v2.12.2 with repo config, and govulncheck with 0 reachable vulnerabilities. Evidence files are under benchmark-results/m041-s02-t01-*.txt through m041-s02-t06-*.txt, including benchmark-results/m041-s02-t06-lifecycle-integration.txt.

## Requirements Advanced

- R011 — Completed and verified lifecycle readiness, overload, and shutdown behavior.

## Requirements Validated

- R011 — S02 implements pre-warm, /live, /ready, model_not_loaded/model_overloaded 503+Retry-After behavior, SIGTERM/SIGINT drain, and validates F-1/F-2/F-5 via `benchmark-results/m041-s02-t06-lifecycle-integration.txt` plus full Go gates.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

T06 used `api/fd_v2_lifecycle_integration_test.go` instead of root-level `tests/integration/fd_v2_lifecycle_test.go` because root-level `go test ./tests/integration/...` is not valid with the current Go module layout. F-2 required a minimal production overload path (`FD_MAX_IN_FLIGHT`) because `model_overloaded` previously existed only as an error code.

## Known Limitations

`FD_MAX_IN_FLIGHT` defaults to 0 (unlimited) to preserve current production behavior; operators must set a positive value to enable active overload rejection. `/embeddings/batch` is not lifecycle-gated in S02 because the planned scope targeted `/v1/embeddings`.

## Follow-ups

S03 should add request observability/metrics around lifecycle outcomes. Future work can decide whether legacy `/embeddings/batch` should also use the lifecycle gate.

## Files Created/Modified

- `api/lifecycle/state.go` — Lifecycle state, readiness, shutdown, in-flight tracking, capacity tracking, last error, context helpers.
- `api/lifecycle/warmup.go` — PreWarm primitive for startup dummy inference.
- `api/lifecycle/shutdown.go` — Signal-driven graceful shutdown orchestration.
- `api/middleware/lifecycle.go` — Readiness/shutdown/capacity gate for `/v1/embeddings`.
- `api/handlers/probes.go` — /live and /ready endpoint handlers.
- `api/main.go` — Wired lifecycle state, warmup, probes, gate, capacity env, and graceful shutdown.
- `api/fd_v2_lifecycle_integration_test.go` — Executable lifecycle scenario coverage for startup, F-1, F-2, and F-5.
