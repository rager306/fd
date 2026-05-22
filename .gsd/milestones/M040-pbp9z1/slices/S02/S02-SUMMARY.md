---
id: S02
parent: M040-pbp9z1
milestone: M040-pbp9z1
provides:
  - Packaged Docker ONNX restart/cache benchmark evidence with Redis L2 restart behavior and sanitized configuration for S04.
  - A reusable S02 proof runner and artifact verifier for same-host ONNX runtime readiness checks.
  - Legal no-regression and leak/cleanup audit evidence for the packaged ONNX proof path.
requires:
  - slice: S01
    provides: Same-host service contract, runtime metadata expectations, timeout/retry/no-silent-fallback framing, and current API surfaces.
affects:
  - S04
key_files:
  - tools/run_m040_s02_docker_restart_proof.sh
  - tools/verify_m040_s02_artifacts.py
  - benchmark-results/fd-m040-s02-onnx-docker-preflight.txt
  - benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt
  - benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt
  - benchmark-results/fd-m040-s02-proof-audit.txt
key_decisions:
  - Use a fixed S02 cache namespace `m040-s02-onnx-restart` and fixed proof container/network names so the restart command restarts only the API container while Redis remains alive.
  - Expose proof API and Redis only on localhost host bindings for the local proof path.
  - Treat legal retrieval and cleanup proof as first-class verifier inputs, not narrative-only evidence.
  - Preserve the isolated proof Redis container for local development/cache reuse while removing only the S02 ONNX API proof container.
patterns_established:
  - Use artifact-semantic verification for runtime proofs: assert health/runtime metadata, benchmark semantics, legal guard result, redaction boundary, and cleanup state rather than relying only on command exit codes.
  - Keep cache namespace isolation explicit for backend/runtime comparisons to avoid stale cross-run embedding cache contamination.
  - Target API-only restart commands during Redis L2 proof so restart latency/cache reuse evidence is not invalidated by Redis restarts.
observability_surfaces:
  - `/health.runtime.*` metadata for backend/model/dimensions/cache namespace/runtime details.
  - Smoke `/v1/embeddings` response model and 1024-dimensional vector shape.
  - Benchmark sanitized effective configuration and measured restart/cache sections.
  - Legal retrieval guard PASS artifact.
  - Proof audit artifact covering leak audit, container cleanup, Redis preservation, and port 18000 status.
drill_down_paths:
  - .gsd/milestones/M040-pbp9z1/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M040-pbp9z1/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M040-pbp9z1/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-22T08:01:53.580Z
blocker_discovered: false
---

# S02: Docker restart and Redis L2 proof

**Proved the packaged ONNX Docker API can start, smoke, restart only the API process, reuse Redis L2 across restart, pass the legal guard, and clean up with audited artifacts.**

## What Happened

S02 converted the previously skipped packaged ONNX restart/cache path into a reproducible proof. T01 added `tools/run_m040_s02_docker_restart_proof.sh` and `tools/verify_m040_s02_artifacts.py`, fixing the proof shape around a fixed S02 cache namespace, localhost-only bindings, and artifact-semantic verification. T02 ran the packaged Docker ONNX proof: the image/build path was exercised, the API started on the proof port with Redis external/alive, `/health` and smoke embedding evidence showed ONNX backend, model `deepvk/USER-bge-m3`, 1024 dimensions, and namespace/runtime metadata, then `benchmark.py` ran with `BENCHMARK_API_RESTART_COMMAND` targeting only the ONNX API container so restart sections were measured rather than skipped. T03 added and ran the legal retrieval guard, produced the proof audit, extended the verifier to require legal and audit closeout semantics, removed the ONNX API proof container, left the isolated Redis container alive intentionally for local cache reuse, and confirmed port 18000 was clear.

## Verification

Fresh closeout verification passed via `gsd_exec` run `74de24ec-8b60-4b3f-a8ce-557c489d5efc`: `python3 tools/verify_m040_s02_artifacts.py --benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt --preflight benchmark-results/fd-m040-s02-onnx-docker-preflight.txt --legal benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt --audit benchmark-results/fd-m040-s02-proof-audit.txt` returned `M040 S02 artifact verification: PASS`; `python3 -m py_compile tools/verify_m040_s02_artifacts.py` passed; `bash -n tools/run_m040_s02_docker_restart_proof.sh` passed; all four required artifacts were non-empty; Docker cleanup showed no `fd-m040-s02-onnx-api` container and Redis `fd-m040-s02-redis` still running as intended; a socket check confirmed `127.0.0.1:18000` refused connections after cleanup. A second artifact summary run `1daa69dd-2b09-417f-88b0-e6a22c417328` confirmed the preflight and benchmark artifacts include ONNX backend, model `deepvk/USER-bge-m3`, 1024 dimensions, cache namespace evidence, restart evidence, legal PASS, and audit PASS without needing to expose raw probe/legal text.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

- None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

The runner/verifier needed narrow fixes during T02/T03 before the full proof could complete: shell helpers now emit captured `/health` JSON and parse piped JSON correctly; smoke checks validate the observable 1024-dimensional embedding shape; verifier reads `benchmark.py` sanitized environment from `environment.values`; and verifier now accepts and enforces `--legal` plus `--audit` inputs. The host Python lacked the benchmark Redis dependency in the managed system environment, so Redis was installed into `.gsd/runtime/python-packages` and the proof ran with local `PYTHONPATH`.

## Known Limitations

The isolated proof Redis container `fd-m040-s02-redis` remains running intentionally for local development/cache reuse; the API proof container was removed and port 18000 is clear. S02 proves same-host packaged ONNX restart/cache behavior for the current model and proof configuration, not final production operations or alternative model selection.

## Follow-ups

S04 should consume S02's packaged ONNX restart/cache benchmark evidence together with S03's bounded legal model quick-gate to make the final TEI-vs-ONNX runtime recommendation and operating contract.

## Files Created/Modified

- `tools/run_m040_s02_docker_restart_proof.sh` — Reproducible proof runner for packaged Docker ONNX API restart/cache benchmark with Redis alive and S02 cache namespace.
- `tools/verify_m040_s02_artifacts.py` — Artifact verifier enforcing benchmark/preflight/legal/audit semantics, redaction boundaries, and cleanup proof.
- `benchmark-results/fd-m040-s02-onnx-docker-preflight.txt` — Preflight/build/start/health/smoke artifact for the packaged ONNX proof run.
- `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt` — Benchmark artifact with measured API restart and Redis L2 reuse sections.
- `benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt` — Legal retrieval no-regression guard result for the packaged ONNX proof run.
- `benchmark-results/fd-m040-s02-proof-audit.txt` — Leak audit and cleanup artifact proving API container removal and port 18000 cleanup while Redis remains alive intentionally.
