# S02: Runtime hardening notes

**Goal:** Document Redis and TEI operational hardening notes from runtime validation and record the ONNX/Candle decision.
**Demo:** After this, runtime hardening notes are explicit and not buried in milestone summaries.

## Must-Haves

- Redis localhost binding rationale is documented.
- Redis `vm.overcommit_memory` warning is documented as host-level deployment note.
- TEI Candle fallback and ONNX artifact choice is recorded as measured follow-up, not required immediate fix.
- Architecture/runtime decision is recorded via GSD decision tooling.

## Proof Level

- This slice proves: docs plus GSD decision record

## Integration Closure

Operational notes become discoverable in README and decision history before deployment or future optimization work.

## Verification

- Future operators know expected runtime warnings, host-level actions, and TEI backend tradeoffs.

## Tasks

- [x] **T01: Inspect hardening evidence** `est:small`
  Inspect current README/compose TEI and Redis notes plus M003/M004 evidence to define exact operational hardening documentation.
  - Files: `README.md`, `docker-compose.yaml`, `docker-compose.override.yaml`, `.gsd/milestones/M003-xx4yc3/M003-xx4yc3-ASSESSMENT.md`, `.gsd/milestones/M004-9886ht/M004-9886ht-SUMMARY.md`
  - Verify: Findings recorded.

- [x] **T02: Document runtime hardening notes** `est:small`
  Update README with runtime hardening notes for Redis localhost binding, Redis overcommit warning, LOG_LEVEL debug cache events, and TEI ONNX/Candle status.
  - Files: `README.md`
  - Verify: README contains required Redis/TEI/logging notes.

- [x] **T03: Record TEI backend decision** `est:small`
  Record a GSD decision that current runtime stays on measured TEI Candle fallback and ONNX export is future measured optimization, not default requirement.
  - Files: `.gsd/DECISIONS.md`
  - Verify: GSD decision saved.

## Files Likely Touched

- README.md
- docker-compose.yaml
- docker-compose.override.yaml
- .gsd/milestones/M003-xx4yc3/M003-xx4yc3-ASSESSMENT.md
- .gsd/milestones/M004-9886ht/M004-9886ht-SUMMARY.md
- .gsd/DECISIONS.md
