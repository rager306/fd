---
id: S01
parent: M045-d0e5xq
milestone: M045-d0e5xq
provides:
  - Candidate mitigation list for S02.
  - Current known-good runtime baseline before any startup config changes.
requires:
  []
affects:
  []
key_files:
  - documents/tei-startup-recon-m045.md
key_decisions:
  - D049: TEI startup probing is a separate operations milestone with non-destructive discovery before controlled proof.
patterns_established:
  - Treat TEI internal ORT/ONNX probing as external runtime behavior, not fd ONNX backend scope.
  - Do not run destructive TEI restarts without capture and rollback plan.
observability_surfaces:
  - documents/tei-startup-recon-m045.md records current runtime state, startup warning signatures, smoke evidence, and candidate mitigations.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-14T11:32:01.479Z
blocker_discovered: false
---

# S01: Non destructive TEI startup recon

**Captured non-destructive TEI startup/runtime recon and identified offline-cache as the leading mitigation candidate.**

## What Happened

S01 gathered read-only evidence for the running TEI-only deployment. The recon captured Docker inspect state, compose config, safe env subset, container health, fd health/readiness, fd and direct TEI embedding smoke, and recent TEI startup logs. Upstream TEI source/docs research found no obvious force-Candle or disable-ORT CLI flag, and showed ORT/ONNX probing can happen before Candle fallback. The leading candidate for S02 is `HF_HUB_OFFLINE=1` with a complete existing `/data` cache; local model path is secondary; adding ONNX files is rejected for current fd scope.

## Verification

UAT PASS via gsd_uat_exec evidence: artifact check `93a29b67-73ac-4644-9fdc-c5054ea171c2`, fd health runtime check `307f2271-c9d4-4315-b4ea-89b6c52d5e18`, direct TEI smoke `558c59fc-491f-4117-91b4-98dbdcc3f537`. No restart or recreate was performed.

## Requirements Advanced

- R028 — Advanced by documenting current startup behavior and identifying candidate mitigations.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

S01 did not prove startup-time improvement because controlled restart proof is S03 scope. `HF_HUB_OFFLINE=1` can break startup if required cached files are missing, so S02 must inventory cache before config changes.

## Follow-ups

S02 should inventory `/data` cached model files, then decide whether to prepare `HF_HUB_OFFLINE=1` as a compose candidate for controlled proof.

## Files Created/Modified

- `documents/tei-startup-recon-m045.md` — New S01 recon artifact with runtime evidence, log timeline, upstream source findings, candidate mitigations, and S02 recommendation.
