---
id: T01
parent: S02
milestone: M010-84qfzu
key_files:
  - .gsd/milestones/M010-84qfzu/slices/S02/tasks/T01-SUMMARY.md
key_decisions:
  - Use `/v1/embeddings` as the baseline comparator API because it returns float arrays directly.
  - Use stable probe labels and character counts in artifacts; do not print raw probe texts.
  - Default comparator env vars should mirror benchmark defaults: `http://localhost:8000`, `deepvk/USER-bge-m3`, and `1024` dimensions.
duration: 
verification_result: mixed
completed_at: 2026-05-19T18:34:47.947Z
blocker_discovered: false
---

# T01: Defined the S02 comparator contract: TEI API baseline, fixed Russian/legal probes, sanitized labels/counts only, and shape/norm/hash/cosine output fields.

**Defined the S02 comparator contract: TEI API baseline, fixed Russian/legal probes, sanitized labels/counts only, and shape/norm/hash/cosine output fields.**

## What Happened

Defined the dense comparator contract for S02. The comparator should call the current TEI-backed `/v1/embeddings` endpoint because it returns JSON float arrays directly and matches the production OpenAI-compatible API surface. The baseline probe set should be fixed, non-sensitive, Russian/legal-style text embedded as script constants, but artifacts must print only stable labels and character counts. Required fields are API URL, model, expected dimensions, probe label, character count, response dimensions, finite-value check, L2 norm, normalized-within-tolerance result, deterministic vector hash, and pairwise cosine similarities. Runtime configuration should use safe env defaults aligned with `benchmark.py`: `COMPARE_API_URL=http://localhost:8000`, `COMPARE_MODEL=deepvk/USER-bge-m3`, and `COMPARE_DIMENSIONS=1024`.

## Verification

Reviewed existing benchmark configuration and API handlers. Confirmed `/v1/embeddings` returns float embeddings with model and dimensions, and the comparator design excludes raw probe texts from artifact output.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Read `benchmark.py` config/env/sanitization pattern.` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Read `api/handlers/embeddings.go` response shape for `/v1/embeddings`.` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Read `api/handlers/batch.go` and chose `/v1/embeddings` over batch because JSON float arrays are simpler for the baseline comparator.` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Read `.gitignore` to confirm generated runtime/model artifacts should not be staged.` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None. Contract was defined from existing API/benchmark patterns before implementation.

## Known Issues

The exact cosine thresholds for ONNX-vs-TEI equivalence are not set in T01; T02/T03 will first establish baseline shape/norm/hash/cosine output so S03 can choose threshold based on observed FP32 export behavior.

## Files Created/Modified

- `.gsd/milestones/M010-84qfzu/slices/S02/tasks/T01-SUMMARY.md`
