---
id: S02
parent: M039-aexhf5
milestone: M039-aexhf5
provides:
  - Fresh packaged Go ONNX smoke, legal, and performance evidence for image `fd-api:onnx1024-m039`.
requires:
  []
affects:
  []
key_files:
  - benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt
  - benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt
  - benchmark-results/fd-onnx-docker-target-runtime-acceptance-m039-s02.txt
key_decisions: []
patterns_established:
  - Packaged ONNX containers must include `ONNX_RUNTIME_SHA256` to report runtime library verification in `/health`.
observability_surfaces:
  - Legal effective configuration, benchmark sanitized configuration, packaged acceptance matrix, image id, runtime SHA verification, cache namespaces.
drill_down_paths:
  - .gsd/milestones/M039-aexhf5/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M039-aexhf5/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M039-aexhf5/slices/S02/tasks/T03-SUMMARY.md
  - .gsd/milestones/M039-aexhf5/slices/S02/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T11:32:32.703Z
blocker_discovered: false
---

# S02: Packaged ONNX closure

**Completed packaged Docker ONNX legal and performance proof for the current artifact setup.**

## What Happened

S02 reran legal retrieval and performance gates through the freshly built packaged Go ONNX Docker image. Legal retrieval passed against TEI/default Go API and packaged Go ONNX API with isolated namespace. Benchmark.py ran against the packaged Go ONNX endpoint with sanitized config and isolated namespace. The acceptance matrix summarizes passed, skipped, and still-open gates. Final guardrails passed.

## Verification

S02 verification passed: packaged legal gate, packaged performance benchmark, outcome checks, final guardrails, GitNexus detect, Docker cleanup, background process check, and port check all passed.

## Requirements Advanced

- implicit target-runtime validation requirement — Target-runtime validation now includes packaged Docker Go ONNX smoke, legal, and performance gates.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Redis L2 restart subcheck was skipped because the packaged container was managed outside benchmark.py with empty restart command. An initial benchmark artifact check expected an absent `Benchmark completed` marker and was corrected to metric-based validation without rerunning the benchmark.

## Known Limitations

Hosted workflow proof, exact ONNX source proof, Redis L2 restart proof, and production/default promotion remain open.

## Follow-ups

Next gate should either wire a controllable Docker restart command into benchmark.py for Redis L2 restart proof, or resolve exact ONNX binary hosting/reproducible-export proof before hosted workflow evidence.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt` — Packaged legal retrieval target-runtime artifact.
- `benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt` — Packaged performance benchmark artifact.
- `benchmark-results/fd-onnx-docker-target-runtime-acceptance-m039-s02.txt` — Packaged acceptance coverage matrix artifact.
