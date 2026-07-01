---
id: M025-9bvjxa
title: "ONNX artifact provisioning and hosted CI contract"
status: complete
completed_at: 2026-05-20T11:49:31.223Z
key_decisions:
  - D023: ONNX rollout remains staged and opt-in with TEI as default rollback path until diagnostics, artifact provisioning/CI, security review, and rollout proof pass.
key_files:
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/OPERATIONS.md
  - tools/provision_onnx_artifacts.py
  - .github/workflows/onnx-packaging.yml
  - docs/onnx-artifacts/README.md
  - .gsd/DECISIONS.md
lessons_learned:
  - Do not use fake artifact defaults for CI; missing external sources must be explicit blockers.
  - A manual workflow is safer than required CI until large artifacts have immutable sources.
  - Operational readiness should be documented separately from quality/performance pass evidence.
---

# M025-9bvjxa: ONNX artifact provisioning and hosted CI contract

**M025 established artifact provisioning, operational rollout, and manual hosted ONNX CI contracts without promoting ONNX or committing binaries.**

## What Happened

M025 converted the locally proven ONNX 1024 Docker path into a truthful provisioning, operations, and hosted-CI contract. S01 added `docs/onnx-artifacts/PROVISIONING.md` and `tools/provision_onnx_artifacts.py`, defining required artifacts, destinations, checksum policy, dry-run behavior, and missing-source blockers. S02 added `docs/onnx-artifacts/OPERATIONS.md` plus D023, defining startup diagnostics, failure messages, health expectations, staged rollout, and rollback to TEI. S03 added `.github/workflows/onnx-packaging.yml`, a manual-only hosted workflow skeleton that provisions explicit artifact sources, verifies checksums, runs tagged tests, and builds the opt-in ONNX image without becoming required CI. Verification passed across actionlint, provisioning dry-run, verifier, default tests/lint/Docker, tagged checks, binary hygiene, cleanup, and GitNexus scope.

## Success Criteria Results

- Provisioning: PASS.
- Operations: PASS.
- Manual CI skeleton: PASS.
- Guardrails: PASS.

## Definition of Done Results

- Provisioning contract: met.
- Provisioning helper: met.
- Operational diagnostics/rollout contract: met.
- Manual hosted ONNX workflow skeleton: met.
- No binaries tracked: met.
- Default guardrails: met.
- ONNX remains opt-in experimental: met.

## Requirement Outcomes

- External artifact provisioning/cache design: validated.
- Operational diagnostics/rollout contract: validated as documentation/policy.
- Hosted ONNX CI skeleton: validated syntactically and safely manual.
- Production promotion: remains blocked pending artifact source selection, workflow run, operational implementation, security review, and rollout proof.

## Deviations

None. Full hosted ONNX workflow was not executed because immutable artifact sources are not configured and no push/remote action was requested.

## Follow-ups

Next gates: choose an immutable artifact store/source for ONNX model and runtime assets, then manually dispatch the ONNX packaging workflow after push. Separately, implement startup/health operational diagnostics in code if production rollout remains a goal.
