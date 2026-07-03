---
id: T01
parent: S01
milestone: M025-9bvjxa
key_files:
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
key_decisions:
  - Hosted CI is blocked until the ONNX model has an immutable external source URL/cache entry.
  - Native tokenizer `latest` URL remains acceptable only with mandatory checksum verification, but production should pin or mirror it.
duration: 
verification_result: passed
completed_at: 2026-05-20T11:37:34.353Z
blocker_discovered: false
---

# T01: Documented the ONNX artifact provisioning and cache contract.

**Documented the ONNX artifact provisioning and cache contract.**

## What Happened

Wrote the artifact provisioning contract and linked it from the ONNX artifact README. The contract lists required artifacts, destination paths, verification requirements, cache layout, source selection recommendations, failure diagnostics, and current blockers. It explicitly preserves TEI/default behavior and refuses fake CI readiness while artifact sources are missing.

## Verification

New paths were confirmed unused before writing, and the contract states blockers without committing binaries or secrets.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `test ! -e docs/onnx-artifacts/PROVISIONING.md && test ! -e tools/provision_onnx_artifacts.py` | 0 | ✅ pass — m025_new_paths_ok | 0ms |
| 2 | `write docs/onnx-artifacts/PROVISIONING.md; edit docs/onnx-artifacts/README.md` | 0 | ✅ pass — provisioning contract written and linked | 0ms |

## Deviations

None.

## Known Issues

ONNX model external source is not yet defined; ONNX Runtime pinned source/hash is not yet encoded in a tracked manifest.

## Files Created/Modified

- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/README.md`
