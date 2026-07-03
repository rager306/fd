---
id: M020-mvkq4d
title: "ONNX 1024 artifact contract"
status: complete
completed_at: 2026-05-20T10:03:47.898Z
key_decisions:
  - D018: The existing dynamic-axis ONNX manifest is the experimental ONNX 1024 runtime contract; production remains blocked on Docker/CI packaging and artifact provisioning.
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - .gsd/DECISIONS.md
lessons_learned:
  - Separate export provenance from validated runtime contract to avoid implying the ONNX binary was re-exported.
  - A manifest can document local quality/performance validation while still explicitly blocking production readiness.
  - Tracked metadata should link evidence artifacts and remaining gates so future agents do not repeat prior measurements unnecessarily.
---

# M020-mvkq4d: ONNX 1024 artifact contract

**M020 made the ONNX 1024 runtime contract explicit in tracked metadata while preserving experimental status and binary hygiene.**

## What Happened

M020 documented the ONNX 1024 runtime contract after quality and local performance gates passed. S01 updated the tracked ONNX manifest to preserve export provenance (`sequence_length=128`) while adding `runtime.validated_max_sequence_length=1024`, validated runtime env, M018 legal quality PASS evidence, M019 local performance evidence, sequence-length failure contract, and future production gates. S02 recorded D018 and ran fresh validation: JSON contract checks, evidence links, tracked binary hygiene, Go tests, lint, tagged HF tokenizer tests, runtime cleanup, and GitNexus scope all passed. TEI remains production/default and ONNX remains experimental.

## Success Criteria Results

- Tracked contract: PASS.
- Evidence links: PASS.
- Binary hygiene: PASS.
- Production safety: PASS.
- Next gate explicit: PASS.

## Definition of Done Results

- 1024 runtime contract tracked: met.
- Export sequence length 128 preserved: met.
- M018/M019 evidence linked: met.
- production_default remains false: met.
- No binaries tracked: met.
- Fresh tests/lint/tagged checks passed: met.
- Next Docker/CI packaging gate explicit: met.

## Requirement Outcomes

- ONNX 1024 artifact contract: validated.
- ONNX production readiness: remains blocked by Docker/CI packaging and artifact provisioning.
- Binary hygiene: validated with zero tracked `.onnx` or `libtokenizers.a` files.
- TEI default: preserved.

## Deviations

None.

## Follow-ups

Create the next GSD milestone for Docker/CI packaging and artifact provisioning: supply ONNX binary and native tokenizer without committing binaries, verify checksums, run packaged legal quality and performance gates, and add startup/failure observability for missing or mismatched artifacts.
