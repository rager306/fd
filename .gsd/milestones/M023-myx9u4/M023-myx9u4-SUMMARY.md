---
id: M023-myx9u4
title: "ONNX 1024 packaged legal quality gate"
status: complete
completed_at: 2026-05-20T11:02:46.476Z
key_decisions:
  - D021: packaged ONNX Docker 1024 passes selected Russian/legal quality, but ONNX remains opt-in experimental until packaged performance, artifact provisioning/CI, and operational rollout gates pass.
key_files:
  - benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt
  - benchmark-results/fd-onnx-docker-legal-outcome-m023-s02.txt
  - .github/workflows/go-quality.yml
  - .gsd/DECISIONS.md
lessons_learned:
  - `Dockerfile.onnx` must be allowed by binary hygiene checks; do not use a naive `.onnx$` regex over all tracked paths.
  - Packaged legal quality can match local M018 quality when using the staged Docker image and isolated Redis namespace.
  - Continue separating quality pass, performance pass, artifact provisioning, and production promotion decisions.
---

# M023-myx9u4: ONNX 1024 packaged legal quality gate

**M023 proved the packaged ONNX Docker image passes the selected Russian/legal quality gate while keeping ONNX experimental.**

## What Happened

M023 validated the dedicated ONNX 1024 Docker image against the Russian/legal retrieval gate. The existing TEI default API served as baseline on port 8000 and the packaged ONNX image ran on port 18000 with isolated cache namespace `m023-onnx-docker-legal`. The evaluator passed with minimum cross-backend cosine 0.99989883, top-1 agreement 1.0, mean overlap@5 0.997701, and ONNX recall ratio 1.0. Both primary and outcome artifacts exclude raw legal text. During hygiene verification, the workflow binary regex was corrected to allow the tracked source file `Dockerfile.onnx` while still blocking real binary artifacts. D021 records that this quality pass does not authorize production/default promotion. Final guardrails passed: actionlint, verifier, default tests/lint/Docker, tagged tests, binary hygiene, cleanup, and GitNexus scope.

## Success Criteria Results

- Packaged legal evaluator: PASS.
- Raw text hygiene: PASS.
- Redis isolation: PASS.
- Default guardrails: PASS.
- Runtime cleanup: PASS.

## Definition of Done Results

- Packaged legal gate run: met.
- Isolated Redis namespace: met (`m023-onnx-docker-legal`).
- Sanitized artifacts: met, raw leak checks 0.
- Default guardrails: met.
- Cleanup: met.
- ONNX remains opt-in experimental: met via D021.

## Requirement Outcomes

- Packaged ONNX legal quality: validated.
- Default TEI runtime: preserved.
- Binary hygiene: fixed and validated.
- Production promotion: still blocked pending packaged performance, hosted artifact provisioning/CI, and rollout diagnostics.

## Deviations

A CI binary-hygiene false positive was discovered and fixed: tracked source file `Dockerfile.onnx` matched the previous `.onnx$` binary regex. The corrected workflow exempts `Dockerfile.onnx` while still blocking real ONNX/native/runtime binaries.

## Follow-ups

Next recommended gate is packaged ONNX performance benchmark against TEI baseline using the M022 image and sanitized config snapshot. Another valid next gate is external artifact provisioning/cache design for hosted full ONNX CI.
