---
id: S03
parent: M025-9bvjxa
milestone: M025-9bvjxa
provides:
  - Hosted CI skeleton for future full ONNX proof.
requires:
  []
affects:
  - Future hosted ONNX CI proof once artifact sources exist
key_files:
  - .github/workflows/onnx-packaging.yml
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
key_decisions:
  - Manual workflow only; no push/PR trigger until immutable artifact sources exist.
  - Explicit artifact inputs required; signed/secret-bearing URLs need future masked-secret wiring.
patterns_established:
  - Do not make full ONNX CI required until artifacts are externally provisioned.
  - Manual workflow dispatch is the bridge between local proof and hosted CI proof.
  - Keep normal Go Quality artifact-free and TEI/default-oriented.
observability_surfaces:
  - Manual workflow phase logs, provisioning JSON, verifier output, tagged test phases, image build phase.
drill_down_paths:
  - .gsd/milestones/M025-9bvjxa/slices/S03/tasks/T01-SUMMARY.md
  - .gsd/milestones/M025-9bvjxa/slices/S03/tasks/T02-SUMMARY.md
  - .gsd/milestones/M025-9bvjxa/slices/S03/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T11:48:46.281Z
blocker_discovered: false
---

# S03: Full hosted ONNX CI skeleton

**S03 added a safe manual hosted ONNX packaging workflow skeleton and verified guardrails.**

## What Happened

S03 added and verified the hosted ONNX CI skeleton. The workflow is manual-only, requires explicit artifact sources, provisions artifacts, verifies them, runs tagged tests, and builds the opt-in ONNX Docker image. Documentation warns against plain signed URL inputs and explains the workflow remains non-required until immutable artifact sources exist. Final verification passed across actionlint, provisioning checks, default tests/lint/Docker, tagged checks, binary hygiene, cleanup, and GitNexus scope.

## Verification

All S3 verification passed.

## Requirements Advanced

- onnx-hosted-ci-skeleton — Added safe manual hosted workflow skeleton for ONNX provisioning/build proof.

## Requirements Validated

- m025-s03-actionlint — actionlint passed for Go Quality and ONNX packaging workflows.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

Workflow has not run end-to-end in GitHub because artifact URLs are not configured and no push was made.

## Follow-ups

After immutable artifact URLs/cache are available, manually dispatch `.github/workflows/onnx-packaging.yml`, then add hosted legal/performance jobs or required CI status based on evidence.

## Files Created/Modified

- `.github/workflows/onnx-packaging.yml` — Manual ONNX packaging workflow skeleton.
- `docs/onnx-artifacts/README.md` — Docs mention manual workflow and URL safety guidance.
- `docs/onnx-artifacts/PROVISIONING.md` — Provisioning docs describe manual workflow and safety constraints.
