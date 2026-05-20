---
id: M024-b8pfpl
title: "ONNX 1024 packaged performance benchmark"
status: complete
completed_at: 2026-05-20T11:31:47.624Z
key_decisions:
  - D022: packaged ONNX Docker 1024 is locally performance-viable after legal quality pass, but remains opt-in experimental until artifact provisioning/CI, operational diagnostics, and rollout gates pass.
key_files:
  - benchmark-results/fd-benchmark-m024-onnx-docker1024.txt
  - benchmark-results/fd-onnx-docker-performance-outcome-m024-s02.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - Container-specific restart commands are required for valid L2 restart benchmark metrics.
  - Packaged ONNX preserves throughput and cold latency viability, but warm single-request latency and batch L1 p95 differ from local non-container ONNX.
  - Benchmark metadata should eventually include Docker image ID/digest or in-container library hash for packaged runtimes.
---

# M024-b8pfpl: ONNX 1024 packaged performance benchmark

**M024 proved the packaged ONNX Docker image remains locally performance-viable while preserving TEI/default behavior.**

## What Happened

M024 validated packaged ONNX Docker performance. The benchmark targeted the packaged ONNX image on port 18000 with isolated cache namespace `m024-onnx-docker-benchmark` and a restart command that restarted `fd-onnx-m024-bench` for L2 persistence measurement. The benchmark completed with best cold latency 7.6ms, warm latency mean 2.03ms, max throughput about 864 req/s, Redis L2 restart 3.36ms, batch L1 p95 8.91ms, batch L2 p95 5.47ms, and chunk reuse warm p95 5.24ms. The outcome artifact compares this against M014 TEI and M019 local ONNX evidence. D022 records the result as packaged performance viability only; ONNX remains experimental until artifact provisioning/CI and operational rollout gates pass. Final verification passed across actionlint, scripts, CI-safe verifier, default tests/lint/Docker, tagged tests, artifact hygiene, binary hygiene, cleanup, and GitNexus scope.

## Success Criteria Results

- Benchmark artifact: PASS.
- Config snapshot: PASS.
- Performance viability: PASS.
- Outcome/decision: PASS.
- Guardrails: PASS.

## Definition of Done Results

- Packaged ONNX benchmark run: met.
- Sanitized effective config captured: met.
- Outcome comparison against prior baselines: met.
- Decision scoping performance pass: met.
- Default guardrails: met.
- Cleanup and binary hygiene: met.
- ONNX remains opt-in experimental: met.

## Requirement Outcomes

- Packaged ONNX performance: validated.
- Production/default switch: still blocked.
- Hosted artifact provisioning/CI: still blocked.
- Operational rollout: still blocked.
- Default guardrails: validated.

## Deviations

None. Metadata caveat: host-side benchmark cannot hash the in-container ONNX Runtime shared library path; recorded in task and outcome artifacts.

## Follow-ups

Next recommended milestone: external artifact provisioning/cache design for hosted CI and deployments. Alternative: operational diagnostics and rollout/rollback gate for missing or mismatched ONNX artifacts.
