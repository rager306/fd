---
id: M039-aexhf5
title: "Packaged Go ONNX target runtime rerun"
status: complete
completed_at: 2026-05-21T11:33:13.296Z
key_decisions: []
key_files:
  - benchmark-results/fd-onnx-docker-smoke-m039-s01.txt
  - benchmark-results/fd-onnx-docker-smoke-rerun-m039-s01.txt
  - benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt
  - benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt
  - benchmark-results/fd-onnx-docker-target-runtime-acceptance-m039-s02.txt
lessons_learned:
  - Packaged ONNX containers need `ONNX_RUNTIME_SHA256` in the runtime environment for `/health` to report `runtime_library_verified=true`.
  - Benchmark.py artifacts can lack a generic completion marker; final checks should validate actual metrics/config markers instead.
---

# M039-aexhf5: Packaged Go ONNX target runtime rerun

**Produced packaged Docker Go ONNX smoke, legal, and performance evidence for the current artifact setup.**

## What Happened

M039 moved from M038's local Go target-runtime proof to packaged Docker Go ONNX proof. It verified local artifacts, built a fresh dedicated image `fd-api:onnx1024-m039`, smoke-tested the packaged endpoint twice, reran Russian/legal retrieval through actual packaged Go ONNX endpoint against TEI/default Go API, and reran performance benchmark against the packaged endpoint. The milestone records the runtime SHA verification requirement, benchmark side effects, skipped Redis L2 restart proof, and remaining production blockers. TEI remains production/default and ONNX remains opt-in experimental.

## Success Criteria Results

- PASS — Packaged Docker Go ONNX evidence exists.
- PASS — Evidence is through actual packaged endpoint.
- PASS — Redis namespaces and benchmark side effects are explicit.
- PASS — ONNX remains opt-in experimental and TEI remains production/default.

## Definition of Done Results

- Done — Packaged ONNX image `fd-api:onnx1024-m039` built from current artifacts.
- Done — Packaged Go ONNX smoke passed with artifact/tokenizer/runtime library verification.
- Done — User-requested packaged smoke rerun passed.
- Done — Packaged legal retrieval gate passed through actual packaged endpoint.
- Done — Packaged performance benchmark passed through actual packaged endpoint.
- Done — Outcome records image id, artifact hashes, namespaces, metrics, skipped subchecks, and non-actions without raw text/secrets.
- Done — Final guardrails passed.

## Requirement Outcomes

The implicit target-runtime validation requirement is advanced with packaged Docker Go evidence, but immutable source, hosted workflow, Redis L2 restart, rollout, and production decision requirements remain open.

## Deviations

First smoke revealed `runtime_library_verified=false` without `ONNX_RUNTIME_SHA256`; accepted smoke/rerun/legal/perf runs set it explicitly. Redis L2 restart proof remains skipped because benchmark.py was run against externally managed Docker containers.

## Follow-ups

Next recommended gates: add a controllable Docker restart harness for Redis L2 proof, or resolve exact ONNX binary source/reproducible-export proof before hosted workflow evidence.
