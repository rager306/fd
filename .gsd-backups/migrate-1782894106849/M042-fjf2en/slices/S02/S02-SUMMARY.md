---
id: S02
parent: M042-fjf2en
milestone: M042-fjf2en
provides:
  - TEI-only current runtime/build/docs posture.
  - Validated R027.
  - Deferred R021 with rationale.
requires:
  - slice: S01
    provides: TEI RCA and ONNX deferral rationale.
affects:
  []
key_files:
  - api/main.go
  - api/main_test.go
  - api/go.mod
  - api/go.sum
  - README.md
  - docs/same-host-embedding-service-contract.md
  - docs/fd-v2.md
  - docker-compose.yaml
  - benchmark-results/m042-s02-tei-only-check.txt
key_decisions:
  - D047: ONNX removed from current product/runtime path, future research only.
  - D048: fd ONNX removal and TEI internal ONNX probing are separate layers.
patterns_established:
  - Use TEI-only active runtime posture; preserve ONNX history as artifacts, not operator path.
  - Reject non-TEI EMBEDDING_BACKEND at startup.
observability_surfaces:
  - Runtime health remains TEI-only; TEI-only check artifact captures active posture and dependency graph.
drill_down_paths:
  - .gsd/milestones/M042-fjf2en/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M042-fjf2en/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M042-fjf2en/slices/S02/tasks/T03-SUMMARY.md
  - .gsd/milestones/M042-fjf2en/slices/S02/tasks/T04-SUMMARY.md
  - .gsd/milestones/M042-fjf2en/slices/S02/tasks/T05-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-06-14T10:54:30.038Z
blocker_discovered: false
---

# S02: TEI active-path cleanup and safe mitigation

**Made fd TEI-only in active startup/build/docs/CI surfaces, removed active ONNX code paths, and passed final gates.**

## What Happened

S02 implemented the user-approved TEI-first rescope. T01 inventoried ONNX active surfaces and defined the remove/keep boundary. T02 removed ONNX as an accepted runtime backend from `api/main.go`; fd now constructs TEI only and rejects non-TEI `EMBEDDING_BACKEND`. T03 removed active ONNX build-tagged embedder/manifest/tokenizer files, `Dockerfile.onnx`, and ONNX/tokenizer module dependencies. T04 updated README, same-host contract, fd-v2 docs, compose, and CI to remove active ONNX operator/build instructions while preserving historical research artifacts. T05 ran final TEI-only checks and mandatory Go/static/security gates, validated R027, and deferred R021.

## Verification

Fresh final gates passed after all changes: `benchmark-results/m042-s02-go-test.txt` shows `go test ./...` passed; `benchmark-results/m042-s02-lint.txt` shows golangci-lint v2.12.2 passed with 0 issues; `benchmark-results/m042-s02-govulncheck.txt` shows 0 reachable vulnerabilities. `benchmark-results/m042-s02-tei-only-check.txt` verifies ONNX active files/workflow removed, no ONNX/runtime tokenizer deps, and compose no longer uses `--dtype`. Structured UAT result saved with PASS.

## Requirements Advanced

- R027 — Implemented and validated TEI-only active posture.
- R021 — Deferred async chunking because S02 scope changed to TEI-only cleanup.

## Requirements Validated

- R027 — `benchmark-results/m042-s02-tei-only-check.txt` plus final gates validate active ONNX removal.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Original S02 async chunking plan was replaced after S01 RCA and user direction. R021 is deferred rather than validated. ONNX is not implemented; it is removed from the active path and preserved only as historical/future research.

## Known Limitations

The external HuggingFace TEI binary may still internally probe ORT/ONNX and delay startup; fd no longer exposes ONNX as an active runtime, but TEI startup stabilization remains a separate external-runtime task.

## Follow-ups

Validate the milestone or close M042 depending on whether skipped S03 is accepted as satisfying the revised scope. If further optimization is needed, plan TEI startup stabilization or safe TEI request-shaping as a new slice/milestone.

## Files Created/Modified

None.
