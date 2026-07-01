# S02: Lifecycle warmup readiness and graceful shutdown — UAT

**Milestone:** M041-4tw0w7
**Written:** 2026-06-14T06:20:18.882Z

# UAT — M041 S02 Lifecycle warmup readiness and graceful shutdown

## Verdict
PASS

## Checks

- [x] Startup readiness sequence: `/live` returns 200 while `/ready` returns 503 `model_not_loaded` before warmup, then `/ready` returns 200 after warmup.
  - Evidence: `benchmark-results/m041-s02-t06-lifecycle-integration.txt` (`TestFdV2LifecycleStartupSequence`).
- [x] F-1 startup race: `/v1/embeddings` returns 503 `model_not_loaded` + `Retry-After: 5` before warmup and 200 after readiness flips.
  - Evidence: `benchmark-results/m041-s02-t06-lifecycle-integration.txt` (`TestFdV2LifecycleF1ModelNotLoadedThenReady`).
- [x] F-2 overload: when `FD_MAX_IN_FLIGHT` capacity is exhausted, `/v1/embeddings` returns 503 `model_overloaded` + `Retry-After: 5`, then returns 200 after capacity recovers.
  - Evidence: `benchmark-results/m041-s02-t06-lifecycle-integration.txt` (`TestFdV2LifecycleF2ModelOverloadedThenRecovers`).
- [x] F-5 shutdown: after shutdown begins, new `/v1/embeddings` requests return 503 `shutting_down` + `Retry-After: 30`, while the in-flight request drains normally.
  - Evidence: `benchmark-results/m041-s02-t06-lifecycle-integration.txt` (`TestFdV2LifecycleF5ShutdownRejectsNewAndDrainsInflight`).
- [x] Static quality gates: full `go test ./...`, golangci-lint v2.12.2, and govulncheck pass with 0 reachable vulnerabilities.
  - Evidence: `benchmark-results/m041-s02-t06-go-test.txt`, `benchmark-results/m041-s02-t06-lint.txt`, `benchmark-results/m041-s02-t06-govulncheck.txt`.

## Notes
The original plan named `tests/integration/fd_v2_lifecycle_test.go`, but the repository's Go module root is `api/`, so executable lifecycle integration tests live in `api/fd_v2_lifecycle_integration_test.go` and run with `cd api && go test . -run TestFdV2Lifecycle -v`.

