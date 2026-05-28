# Project

## What This Is

`fd` is a Go embedding API service for Russian/legal-domain embedding workloads. It exposes local HTTP embedding endpoints for same-host consumers, uses Redis caching, and supports both the measured TEI/default runtime and an explicit opt-in ONNX runtime path for `deepvk/USER-bge-m3`.

## Core Value

The service must provide high-quality Russian/legal-domain embeddings to neighboring same-host services with predictable speed, cache behavior, runtime identity, and operational diagnostics.

## Project Shape

- **Complexity:** complex
- **Why:** The project crosses runtime selection, legal-domain quality gates, Redis cache correctness, Docker lifecycle, local HTTP service contracts, and benchmark comparability.

## Current State

M040 is complete. `fd` now has a same-host local HTTP service contract, safe `/health.runtime` metadata, packaged ONNX restart/cache proof, a bounded alternative-model quick-gate artifact, and a machine-verified TEI-vs-ONNX runtime recommendation.

TEI remains the production/default runtime. ONNX 1024 has passed local Go runtime and packaged Docker smoke/legal/performance/restart-cache evidence, but remains explicit opt-in. Packaged ONNX is recommended only for same-host performance deployments that satisfy the operating contract: artifact/tokenizer/runtime preflight, cache namespace isolation, no silent request-level fallback, and a live `/v1/embeddings` smoke check.

`deepvk/USER-bge-m3` remains the model baseline. Alternative model replacement is deferred fail-closed until candidate and baseline services expose contract-required runtime metadata and legal-domain retrieval metrics prove a candidate can challenge the current model.

There is no active milestone. Future work should start from a new GSD milestone or quick task depending on scope.

## Communication

Project communication should be in Russian by default: dialogue, informational output, questions, discussion, intermediate status updates, and final summaries. Use another language only when the user explicitly asks for it.

## Architecture / Key Patterns

- Go API under `api/` exposes embedding handlers and health metadata.
- `/v1/embeddings` request `model` is OpenAI-compatibility metadata; clients must treat the response `model` and `/health.runtime.model` as authoritative runtime identity.
- `/health` exposes safe operational metadata but is not a live inference readiness probe; readiness still requires a smoke embedding request.
- Redis provides L2 embedding cache; namespace isolation is mandatory for TEI/ONNX or model comparisons.
- `benchmark.py` records sanitized effective configuration and supports `BENCHMARK_API_RESTART_COMMAND` for restart checks.
- ONNX runtime is explicit opt-in behind `onnx` and `hf_tokenizers` build tags and requires verified artifacts.
- Dedicated ONNX Docker packaging uses `Dockerfile.onnx` and `tools/build_onnx_image.sh`.
- Semantic artifact verifiers are used for high-stakes recommendation artifacts to prevent missing evidence, unsafe fallback language, or redaction regressions.

## Capability Contract

See `.gsd/REQUIREMENTS.md` for the explicit capability contract, requirement status, and coverage mapping. After M040 cleanup, all currently tracked requirements are validated by completed work; new runtime/model/service ambitions should be captured as new requirements before execution.

## Milestone Sequence

- [x] M038: Go ONNX target runtime acceptance proof — Real Go endpoint smoke, legal, and performance evidence exists for the current ONNX artifact.
- [x] M039: Packaged Go ONNX target runtime rerun — Dedicated packaged ONNX Docker smoke, legal, and performance evidence exists.
- [x] M040: Same-host embedding service readiness — Same-host contract, lifecycle/cache proof, bounded model quick gate, and runtime recommendation are complete.

## Likely Next Work

No next milestone is active. Plausible future work includes:

- rerun the bounded legal model gate once baseline and candidate endpoints expose contract-required `/health.runtime` metadata;
- perform an operator rollout rehearsal for explicit ONNX opt-in using the M040 operating contract;
- prepare a release/push plan for the local branch, which is ahead of `origin/master`.
