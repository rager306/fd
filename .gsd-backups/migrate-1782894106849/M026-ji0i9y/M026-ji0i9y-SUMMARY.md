---
id: M026-ji0i9y
title: "ONNX operational diagnostics implementation"
status: complete
completed_at: 2026-05-20T12:12:08.559Z
key_decisions:
  - D024: M026 implements first ONNX startup diagnostics gate; production remains blocked until remaining diagnostics, provisioning, security, and rollout gates pass.
key_files:
  - api/handlers/health.go
  - api/handlers/health_test.go
  - api/main.go
  - api/main_test.go
  - api/embed/onnx_manifest.go
  - api/embed/onnx_manifest_test.go
  - docs/onnx-artifacts/OPERATIONS.md
  - benchmark-results/fd-onnx-operational-diagnostics-outcome-m026-s02.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - Health metadata can be added safely with a handler constructor while preserving legacy default health shape.
  - Manifest runtime contract should be parsed into code and enforced during config load, not only documented in artifacts.
  - GitNexus pre-commit breadth can be high for startup/manifest implementation; final post-commit reindex/detect is necessary to confirm stable scope.
---

# M026-ji0i9y: ONNX operational diagnostics implementation

**M026 implemented safe ONNX startup diagnostics and health metadata while preserving TEI/default behavior.**

## What Happened

M026 implemented the first code-level ONNX operational diagnostics gate. The health handler now supports safe runtime metadata while preserving default TEI health shape. ONNX manifest validation now carries `validated_max_sequence_length`; startup config rejects `ONNX_MAX_SEQUENCE_LENGTH` values above that manifest contract. Main logs safe ONNX preflight metadata and Redis cache namespace, and wires ONNX runtime metadata into `/health` only when ONNX is active. Tests cover default health, ONNX metadata safety, manifest contract parsing, sequence mismatch, and runtime health status. Operations docs and the outcome artifact now distinguish implemented diagnostics from remaining gaps. ONNX remains opt-in experimental.

## Success Criteria Results

- Startup diagnostics: PASS.
- Health metadata: PASS.
- Default guardrails: PASS.
- Docs/outcome: PASS.
- Production safety: PASS (no switch).

## Definition of Done Results

- ONNX startup diagnostics implemented: met for manifest/sequence/health metadata scope.
- Health metadata safe and opt-in: met.
- Tests cover TEI/default and ONNX diagnostics: met.
- Docs updated: met.
- Default guardrails: met.
- ONNX production promotion: not performed.

## Requirement Outcomes

- ONNX diagnostics: advanced and partially validated.
- Default TEI behavior: validated.
- Production rollout: still blocked.
- Remaining diagnostics/security/provisioning gates: deferred explicitly.

## Deviations

GitNexus pre-commit scope was high because the implementation intentionally touched startup config, manifest validation, health handler, tests, and docs. Direct impact analysis was performed before edits and all guardrails passed.

## Follow-ups

Next gates: implement tokenizer JSON checksum preflight and ONNX Runtime sha/provider diagnostics, then perform security review for artifact path/URL/logging handling before any rollout proof.
