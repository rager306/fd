# Project

## What This Is

`fd` is a Go embedding API service for Russian/legal-domain embedding workloads. It currently exposes local HTTP embedding endpoints, uses Redis caching, and has both the measured TEI/default runtime and an opt-in ONNX runtime path for `deepvk/USER-bge-m3`.

## Core Value

The service must provide high-quality Russian/legal-domain embeddings to neighboring services on the same host with predictable speed, cache behavior, and operational diagnostics.

## Project Shape

- **Complexity:** complex
- **Why:** The project crosses runtime selection, legal-domain quality gates, Redis cache correctness, Docker lifecycle, local HTTP service contracts, and benchmark comparability.

## Current State

TEI remains the production/default runtime. ONNX 1024 has passed local Go runtime and packaged Docker smoke/legal/performance evidence, but ONNX is still opt-in. M040 reframes the work around same-host embedding service readiness rather than ONNX experimentation.

## Architecture / Key Patterns

- Go API under `api/` exposes embedding handlers and health metadata.
- Redis provides L2 embedding cache; namespace isolation is required for TEI/ONNX comparisons.
- `benchmark.py` records sanitized effective configuration and supports `BENCHMARK_API_RESTART_COMMAND` for restart checks.
- ONNX runtime is explicit opt-in behind `onnx` and `hf_tokenizers` build tags and requires verified artifacts.
- Dedicated ONNX Docker packaging uses `Dockerfile.onnx` and `tools/build_onnx_image.sh`.

## Capability Contract

See `.gsd/REQUIREMENTS.md` for the explicit capability contract, requirement status, and coverage mapping.

## Milestone Sequence

- [x] M038: Go ONNX target runtime acceptance proof — Real Go endpoint smoke, legal, and performance evidence exists for the current ONNX artifact.
- [x] M039: Packaged Go ONNX target runtime rerun — Dedicated packaged ONNX Docker smoke, legal, and performance evidence exists.
- [ ] M040: Same-host embedding service readiness — Produce a same-host service contract, lifecycle/cache proof, bounded model quick gate, and runtime recommendation.

Actual active milestone ID for the next planned work is `M040-pbp9z1`.
