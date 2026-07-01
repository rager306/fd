---
id: S02
parent: M033-aym8ih
milestone: M033-aym8ih
provides:
  - Actionable ONNX Runtime wheel provisioning behavior for future hosted proof.
requires:
  []
affects:
  []
key_files:
  - tools/provision_onnx_artifacts.py
  - docs/onnx-artifacts/PROVISIONING.md
  - benchmark-results/fd-onnx-runtime-wheel-provisioning-m033-s02.txt
key_decisions:
  - D031: provision ONNX Runtime wheel/zip sources by extracting the configured member with size/sha verification while preserving direct-file fallback.
patterns_established:
  - Wheel extraction must read only configured member paths and must verify destination size/sha before use.
observability_surfaces:
  - Provisioning dry-run includes manifest-derived ONNX Runtime expected sha.
  - Outcome artifact records positive/negative/fallback probes.
drill_down_paths:
  - .gsd/milestones/M033-aym8ih/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M033-aym8ih/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M033-aym8ih/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T07:24:41.006Z
blocker_discovered: false
---

# S02: Wheel provisioning documentation and closure

**Documented and verified ONNX Runtime wheel provisioning support.**

## What Happened

S02 documented the wheel extraction behavior, recorded outcome and D031, and ran final verification. The provisioning helper now has documented, tested behavior for `.whl`/`.zip` runtime sources and direct `.so` fallback, without weakening artifact safety or changing runtime defaults. GitNexus HIGH scope is expected because the central provisioning flow changed and was covered by targeted probes and full guardrails.

## Verification

S02 verification passed: compile/dry-run/verifiers, synthetic wheel probes, Go tests, lint, actionlint, tagged tests, Docker build, leak checks, tracked binary hygiene, GitNexus detect, background process check, and port check all passed.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Final zip-member regular-file check needed one correction to tolerate normal zip entries without POSIX file type bits while rejecting symlinks/directories/non-regular POSIX entries. Verification was rerun after the correction.

## Known Limitations

The real PyPI wheel was not provisioned in hosted CI. Exact ONNX model source remains blocked. ONNX remains opt-in experimental.

## Follow-ups

Next gate remains exact ONNX model binary source: either exact-binary immutable hosting/mirroring or a full reproducible-export workflow plus renewed quality/performance/package proof. Hosted workflow proof still requires explicit push approval.

## Files Created/Modified

- `tools/provision_onnx_artifacts.py` — Adds safe ONNX Runtime wheel/zip member extraction and manifest-derived runtime expected sha handling.
- `docs/onnx-artifacts/PROVISIONING.md` — Documents wheel extraction behavior and limits.
- `benchmark-results/fd-onnx-runtime-wheel-provisioning-m033-s02.txt` — Outcome artifact for M033.
- `.gsd/DECISIONS.md` — Decision D031 added by GSD.
