---
id: T03
parent: S01
milestone: M011-33b7wf
key_files:
  - .gsd/milestones/M011-33b7wf/slices/S01/S01-RESEARCH.md
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
key_decisions:
  - S02 must validate manifest/file identity before any ONNX Runtime load.
  - Explicit ONNX requests should fail fast on missing/checksum mismatch instead of silently falling back to TEI for benchmark evidence.
  - The tracked manifest is a local prototype contract, not a production artifact distribution mechanism.
duration: 
verification_result: passed
completed_at: 2026-05-19T18:58:40.510Z
blocker_discovered: false
---

# T03: Saved the S01 artifact validation contract for future opt-in ONNX runtime work.

**Saved the S01 artifact validation contract for future opt-in ONNX runtime work.**

## What Happened

Documented the S01 artifact validation contract in `S01-RESEARCH.md`. The research explains that `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` is the tracked source of truth for local prototype validation while the ONNX binary stays ignored under `.gsd/runtime/`. It specifies validation order, missing-file behavior, checksum-mismatch behavior, metadata checks, no silent fallback for explicit ONNX benchmark runs, and the no-production-runtime-change boundary.

## Verification

Verified S01 research exists and includes missing artifact, checksum mismatch, and production_default=false language. Manifest JSON parses successfully.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `test -f .gsd/milestones/M011-33b7wf/slices/S01/S01-RESEARCH.md && grep -q "missing artifact" .gsd/milestones/M011-33b7wf/slices/S01/S01-RESEARCH.md && grep -q "checksum mismatch" .gsd/milestones/M011-33b7wf/slices/S01/S01-RESEARCH.md && grep -q "production_default=false" .gsd/milestones/M011-33b7wf/slices/S01/S01-RESEARCH.md && python3 -m json.tool docs/onnx-artifacts/user-bge-m3-dense-fp32.json >/tmp/fd-onnx-manifest.json` | 0 | ✅ pass | 0ms |

## Deviations

None. S01 research was written through the GSD artifact tool.

## Known Issues

Production artifact distribution remains unresolved and is intentionally deferred to later slices or future milestones.

## Files Created/Modified

- `.gsd/milestones/M011-33b7wf/slices/S01/S01-RESEARCH.md`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
