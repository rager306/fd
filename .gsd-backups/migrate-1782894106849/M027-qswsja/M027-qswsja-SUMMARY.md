---
id: M027-qswsja
title: "ONNX artifact preflight diagnostics"
status: complete
completed_at: 2026-05-20T12:45:16.622Z
key_decisions:
  - D025: M027 authorizes preflight hardening only; no production/default promotion and no runtime provider enumeration claim.
key_files:
  - api/embed/onnx_manifest.go
  - api/embed/onnx_manifest_test.go
  - api/main.go
  - api/main_test.go
  - api/handlers/health.go
  - api/handlers/health_test.go
  - docs/onnx-artifacts/OPERATIONS.md
  - benchmark-results/fd-onnx-preflight-diagnostics-outcome-m027-s02.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - Manifest source-file metadata is a useful seam for validating tokenizer artifacts at startup.
  - Runtime library hashing should remain explicit/opt-in until artifact source and caching contracts are immutable.
  - Provider diagnostics should be honest about what is validated: configured support, not enumerated runtime availability.
---

# M027-qswsja: ONNX artifact preflight diagnostics

**M027 added tokenizer, runtime-library, and provider preflight diagnostics for opt-in ONNX startup.**

## What Happened

M027 implemented the next ONNX operational startup diagnostic gate. The ONNX manifest validation result now exposes tokenizer JSON metadata, and startup config validates tokenizer JSON existence, size, and checksum when that metadata is present. Startup can optionally validate the ONNX Runtime shared library sha through `ONNX_RUNTIME_SHA256`, and rejects unsupported configured providers before initialization. Safe `/health.runtime` metadata now includes provider, tokenizer verification, and runtime library verification flags while excluding path-like sensitive fields. Operations docs, outcome artifacts, and D025 document that this is hardening only: TEI remains default and ONNX remains opt-in experimental.

## Success Criteria Results

- Tokenizer preflight: PASS.
- Optional runtime sha preflight: PASS.
- Provider validation: PASS.
- Health metadata: PASS.
- Guardrails: PASS.
- Production safety: PASS (no switch).

## Definition of Done Results

- Tokenizer JSON checksum preflight: met.
- Optional ONNX Runtime sha preflight: met.
- Provider config validation: met.
- Safe health metadata: met.
- Default TEI unaffected: met.
- Docs/outcome/decision: met.
- Production promotion: not performed.

## Requirement Outcomes

- ONNX operational preflight: advanced and validated.
- Default TEI behavior: validated.
- Security review: still pending.
- Hosted CI/artifact provisioning: still pending.
- Production rollout: still blocked.

## Deviations

GitNexus pre-commit scope was high because M027 intentionally touched manifest validation, startup config, health metadata, tests, and docs. Direct impact analysis and full guardrails covered this scope.

## Follow-ups

Next recommended gate is security review for artifact path handling, startup error messages, signed URL handling, and logging redaction; after that, runtime library source manifest/provider enumeration or hosted artifact provisioning can proceed.
