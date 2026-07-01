---
id: S04
parent: M040-pbp9z1
milestone: M040-pbp9z1
provides:
  - Machine-checkable TEI-vs-ONNX same-host runtime recommendation.
  - Discoverable final operating contract link from the same-host service contract.
  - Fail-closed semantic validation of evidence envelope, redaction, and readiness-gate semantics.
requires:
  - slice: S01
    provides: Same-host local HTTP service contract and readiness semantics.
  - slice: S02
    provides: Packaged ONNX Docker restart/cache/preflight/legal/audit evidence.
  - slice: S03
    provides: Bounded alternative-model legal-domain quick-gate with `defer_candidate` result.
affects:
  - Milestone M040-pbp9z1 validation and completion.
key_files:
  - tools/verify_m040_s04_recommendation.py
  - benchmark-results/fd-runtime-recommendation-m040-s04.md
  - docs/same-host-embedding-service-contract.md
key_decisions:
  - TEI remains current/default until an explicit operator opt-in to ONNX.
  - Packaged ONNX is recommended only for same-host performance deployments that satisfy the S01 operating contract, artifact/tokenizer/runtime preflight, isolated Redis cache namespace, and live inference smoke check.
  - Alternative model replacement remains `defer_candidate` and fail-closed based on S03.
  - Hosted/remote CI proof is not a same-host runtime readiness gate.
patterns_established:
  - Final recommendation artifacts should be paired with a semantic verifier that fails closed on missing evidence links, weakened caveats, unsafe fallback language, and redaction violations.
observability_surfaces:
  - `tools/verify_m040_s04_recommendation.py` provides localized semantic failures for documentation drift, missing evidence, prohibited leakage patterns, and readiness-gate regressions.
  - `benchmark-results/fd-runtime-recommendation-m040-s04.md` exposes decision inputs, caveats, operator checks, non-goals, and source artifact links for future agents.
drill_down_paths:
  - .gsd/milestones/M040-pbp9z1/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M040-pbp9z1/slices/S04/tasks/T02-SUMMARY.md
  - .gsd/milestones/M040-pbp9z1/slices/S04/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-22T08:33:53.093Z
blocker_discovered: false
---

# S04: Runtime recommendation and operating contract

**Published and machine-verified the final TEI-vs-ONNX same-host runtime recommendation and linked operating contract for fd.**

## What Happened

S04 assembled the upstream S01/S02/S03 evidence into a final recommendation artifact and a fail-closed semantic verifier. The verifier requires the final artifact to cite the same-host contract and packaged ONNX evidence, preserve the exact recommendation stance, include cache namespace isolation and no-silent-fallback semantics, defer alternative model replacement, keep hosted/remote CI proof out of the readiness gate, and reject prohibited raw-text or secret-looking evidence. The final artifact recommends continuing to use `deepvk/USER-bge-m3`, keeping TEI as the current/default posture until an operator explicitly switches, and preferring packaged ONNX only for same-host performance deployments that satisfy S01 `/health.runtime` metadata, artifact/tokenizer/runtime preflight, optional `ONNX_RUNTIME_SHA256` runtime-library integrity, isolated `EMBEDDING_CACHE_VERSION`, and smoke `POST /v1/embeddings`. It consumes S02 Docker restart/cache/legal/preflight/audit evidence and S03's `defer_candidate` result without reopening those decisions. The same-host service contract now links to the final recommendation for discoverability.

## Verification

Fresh closeout verification ran through gsd_exec c52073f9-7ea0-4b13-9efa-99d54193c6f0 and passed: `python3 -m py_compile tools/verify_m040_s04_recommendation.py`; `python3 tools/verify_m040_s04_recommendation.py --self-test`; full verifier validation of `benchmark-results/fd-runtime-recommendation-m040-s04.md` against all required S02/S03 evidence inputs; `rg -n fd-runtime-recommendation-m040-s04 docs/same-host-embedding-service-contract.md`; and `npx gitnexus detect-changes --repo fd`, which reported no unexpected current changes. Milestone status also showed all three S04 tasks complete before slice closure.

## Requirements Advanced

- R001 — Final recommendation preserves deepvk/USER-bge-m3 and avoids overclaiming unproven legal-domain replacements.
- R002 — Final artifact and verifier require isolated Redis cache namespace semantics for ONNX/TEI comparisons and operation.
- R003 — Operating contract documents runtime/env/cache knobs including explicit ONNX switch, artifact/tokenizer/runtime preflight, optional runtime-library hash, and `EMBEDDING_CACHE_VERSION`.
- R004 — Verifier scans final and evidence artifacts for prohibited secret/raw-text patterns and the final artifact contains a redaction section.
- R005 — Final artifact consumes and links from the S01 same-host HTTP service contract rather than duplicating or weakening it.
- R008 — Final stance consumes S03 and keeps alternative model replacement deferred/fail-closed.
- R009 — Final stance and verifier require no request-level fallback and smoke readiness beyond `/health`.

## Requirements Validated

- R006 — Final recommendation artifact passed semantic validation against S02/S03 source evidence in gsd_exec c52073f9-7ea0-4b13-9efa-99d54193c6f0.
- R007 — Final verifier and artifact passed hosted/remote CI readiness-gate rejection semantics in gsd_exec c52073f9-7ea0-4b13-9efa-99d54193c6f0.

## New Requirements Surfaced

- None.

## Requirements Invalidated or Re-scoped

- None. — 

## Operational Readiness

None.

## Deviations

GitNexus MCP tools were not directly exposed in this harness; the available repo-scoped CLI equivalent `npx gitnexus detect-changes --repo fd` was used for the change-scope audit. T03 made no source edits.

## Known Limitations

This slice is a documentation/evidence assembly slice and does not rerun runtime services or Docker benchmarks. ONNX runtime-library integrity remains a required operator check when `ONNX_RUNTIME_SHA256` is configured; S02 evidence recorded `runtime_library_verified=false` because that value was not set. Alternative model replacement remains deferred/fail-closed.

## Follow-ups

Proceed to milestone validation/completion for M040-pbp9z1. No execution follow-up is required for S04.

## Files Created/Modified

- `tools/verify_m040_s04_recommendation.py` — Semantic verifier and self-tests for the final M040 S04 recommendation artifact.
- `benchmark-results/fd-runtime-recommendation-m040-s04.md` — Final TEI-vs-ONNX same-host runtime recommendation, evidence envelope, operating contract, caveats, non-goals, source artifacts, and redaction posture.
- `docs/same-host-embedding-service-contract.md` — Added discoverability link to the final M040 S04 runtime recommendation.
