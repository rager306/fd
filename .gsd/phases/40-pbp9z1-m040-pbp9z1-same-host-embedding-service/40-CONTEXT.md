# M040-pbp9z1: Same-host embedding service readiness

**Gathered:** 2026-05-21
**Status:** Ready for planning

## Project Description

`fd` is a Go embedding API service for Russian/legal-domain workloads. The goal of this milestone is not to continue ONNX experimentation for its own sake; it is to prepare `fd` as a same-host local embedding service that neighboring services can depend on for excellent legal-domain quality and optimal speed on this host.

## Why This Milestone

M038 proved the current ONNX artifact through actual local Go endpoints. M039 proved the dedicated packaged ONNX Docker runtime through smoke, legal, and performance gates. The remaining product question is now broader and more practical: which runtime and operating contract should serve neighboring local HTTP clients on this host?

M040 turns the accumulated evidence into a same-host service readiness proof and runtime recommendation.

## User-Visible Outcome

### When this milestone is complete, the user can:

- Inspect a same-host embedding service contract that tells local services how to call `fd`, what runtime they are hitting, and how to interpret health/readiness.
- Inspect measured evidence for TEI vs ONNX as same-host runtimes, including quality, speed, restart/cache behavior, and operational caveats.
- See a practical runtime recommendation for local operation on this host.

### Entry point / environment

- Entry point: local HTTP API: `/v1/embeddings`, batch embeddings endpoint, `/health`.
- Environment: local same-host service on the current Ubuntu/KVM host, Docker where packaged proof is needed.
- Live dependencies involved: Redis on localhost, TEI/default service on localhost where baseline comparison is needed, packaged ONNX Docker container where ONNX lifecycle proof is needed.

## Completion Class

- Contract complete means: service contract artifact exists and maps endpoints, runtime/env expectations, no-silent-fallback behavior, timeout/retry guidance, and health metadata.
- Integration complete means: benchmark/retrieval tools exercise real TEI and/or ONNX HTTP endpoints with isolated namespaces.
- Operational complete means: packaged ONNX lifecycle/restart/cache behavior is proven or a blocker is recorded truthfully.

## Final Integrated Acceptance

To call this milestone complete, we must prove:

- A neighboring local HTTP client has a documented and verified contract for calling the service and interpreting readiness.
- The packaged ONNX Docker runtime can be restarted through a controlled script/hook and benchmark evidence records Redis L2 restart behavior rather than silently skipping it.
- The final runtime recommendation is based on legal quality, local speed, restart/cache behavior, health/preflight clarity, and operational simplicity.
- Alternative model checks remain bounded and legal-domain-specific; they cannot replace `deepvk/USER-bge-m3` without legal proof.

## Architectural Decisions

### Same-host service, not ONNX experimentation

**Decision:** Scope M040 around preparing `fd` as a local same-host embedding service for neighboring HTTP clients.

**Rationale:** The user clarified that the main goal is excellent quality and optimal speed for services on this host, not experiments beyond that boundary.

**Alternatives Considered:**
- Continue closing ONNX gates as independent experiments — rejected because ONNX is a candidate runtime, not the goal.
- Broaden into remote/hosted deployment proof — rejected because hosted GitHub Actions proof is not required for this project.

### Evidence-envelope runtime recommendation

**Decision:** Recommend TEI vs ONNX using an evidence envelope: legal-domain quality, same-host speed, restart/cache behavior, health/preflight clarity, and operational simplicity.

**Rationale:** Raw speed alone is not enough for a local infrastructure dependency. Quality regressions, mixed-vector cache contamination, or opaque lifecycle behavior would invalidate a runtime even if it is fast.

**Alternatives Considered:**
- Speed-first selection — rejected because legal quality and operability are load-bearing.
- Conservative TEI default forever — rejected because ONNX has strong measured evidence and should be recommended if it clears the envelope.

### Scripted packaged restart proof

**Decision:** Prove packaged restart/cache behavior with a small local Docker restart script/command wired into `BENCHMARK_API_RESTART_COMMAND` or an equivalent benchmark-compatible hook.

**Rationale:** M039 benchmark truthfully skipped Redis L2 restart because the container was managed externally. M040 should close that gap without rewriting the benchmark unnecessarily.

**Alternatives Considered:**
- First-class Docker semantics inside `benchmark.py` — stronger integration but broader code change.
- Manual restart proof outside benchmark artifacts — less invasive but weaker comparability.

### Bounded alternative model gate

**Decision:** Include only a bounded quick gate for 1-2 plausible alternative embedding model candidates using legal-domain evidence.

**Rationale:** The user allows testing other models only at the boundary, if useful for quality/speed, and explicitly does not want open-ended experiments.

**Alternatives Considered:**
- Full model bakeoff — deferred unless the quick gate shows a reason.
- No candidate check at all — rejected because the user allowed a small legal-domain candidate check.

## Error Handling Strategy

- Fail fast on invalid runtime, artifact, tokenizer, ONNX Runtime, or sequence-length configuration.
- `/health` is the primary safe readiness surface and must expose backend, model, dimensions, artifact/tokenizer/runtime verification, provider, cache namespace, sequence length, and production/default flag without secrets.
- No silent per-request fallback between TEI and ONNX, tokenizers, or models within one service run.
- Redis namespace contamination is a correctness risk; TEI/ONNX comparisons must isolate namespaces or explicitly record cache flush side effects.
- Local clients should use bounded timeouts and retry only transient service/transport failures; bad input, model mismatch, and config failures should not be blindly retried.
- Rollback is an operational restart/reconfiguration procedure, not hidden request-level behavior.

## Risks and Unknowns

- Redis L2 restart proof may require a small script or benchmark-harness adjustment — this is the highest remaining operational gap from M039.
- Candidate model availability may be poor or operationally expensive — the quick gate must stop before it becomes open-ended research.
- Benchmark artifacts may still reflect host-state side effects such as Redis `FLUSHALL` — M040 must record these effects clearly.
- ONNX package runtime requires `ONNX_RUNTIME_SHA256` for `/health` to report `runtime_library_verified=true` — missing it would weaken readiness proof.

## Existing Codebase / Prior Art

- `benchmark.py` — existing benchmark harness with sanitized effective config snapshot and `BENCHMARK_API_RESTART_COMMAND` restart hook.
- `benchmark.py:restart_api_for_l2_check` — existing restart function called from benchmark main path.
- `Dockerfile.onnx` — dedicated opt-in packaged ONNX image definition.
- `tools/build_onnx_image.sh` — builds the packaged ONNX image from local verified artifacts.
- `tools/evaluate_legal_retrieval.py` — legal-domain TEI-vs-ONNX parity evaluator.
- `api/handlers/health.go` — health/readiness surface for runtime metadata.
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` — ONNX artifact/runtime/source contract.
- `benchmark-results/fd-onnx-docker-target-runtime-acceptance-m039-s02.txt` — latest packaged ONNX evidence.

## Relevant Requirements

- R001 — M040 preserves Russian/legal-domain quality through no-regression gates and candidate model limits.
- R002 — M040 advances long-lived Redis L2 reuse by proving restart behavior.
- R003 — M040 documents runtime/cache env contract for same-host operation.
- R004 — M040 keeps benchmark artifacts comparable through sanitized config.
- R005 — M040 establishes the same-host HTTP embedding service contract.
- R006 — M040 produces an evidence-based runtime recommendation.
- R007 — M040 excludes hosted GitHub Actions proof as a readiness gate.
- R008 — M040 bounds alternative model checks to legal-domain evidence.
- R009 — M040 prevents silent runtime/model fallback from becoming part of the contract.

## Scope

### In Scope

- Same-host local HTTP client contract for `/v1/embeddings`, batch embeddings, and `/health`.
- TEI-vs-ONNX runtime recommendation based on measured evidence.
- Packaged Docker ONNX restart/cache proof using Redis L2 and isolated namespaces.
- Legal-domain no-regression evidence for current runtime and bounded candidates.
- A final recommendation artifact with evidence table and operating contract.

### Out of Scope / Non-Goals

- Hosted GitHub Actions proof as a readiness gate.
- Remote workflow dispatch, push, upload, or artifact mirroring.
- ONNX experimentation outside same-host service readiness.
- Rust rewrite, provider/INT8/NUMA/OpenVINO experiments without direct same-host evidence need.
- Embedded/library integration surface.
- Automatic per-request fallback between TEI and ONNX.
- Production/default switch as an action in this milestone.

## Technical Constraints

- TEI remains production/default unless an explicit later decision changes it.
- ONNX remains opt-in and must use `onnx hf_tokenizers` tags.
- Packaged ONNX containers must include `ONNX_RUNTIME_SHA256` to verify runtime library health metadata.
- Raw legal text, raw probe text, secrets, and signed artifact URLs must not be recorded in artifacts.
- Redis namespaces must be isolated for TEI/ONNX and candidate comparisons.
- Alternative model checks must be bounded and legal-domain-specific.

## Integration Points

- Local HTTP clients — consume `/v1/embeddings`, batch endpoint, and `/health`.
- Redis — provides L2 embedding cache and restart persistence proof.
- TEI/default API — baseline runtime and quality reference.
- Packaged ONNX Docker API — candidate optimized runtime proof.
- Legal corpus evaluator — verifies no-regression on Russian/legal domain.
- Benchmark harness — produces latency/throughput/cache/restart evidence.

## Testing Requirements

- Static: service contract and recommendation artifacts exist and contain required sections.
- Command: Go tests and lint pass if code changes; Python tools compile if scripts change.
- Behavioral: `/health`, `/v1/embeddings`, benchmark restart proof, and legal evaluator run against real local endpoints.
- Safety: artifact leak checks prove no raw legal/probe text, secrets, or signed URLs were recorded.
- Cleanup: no leftover M040 containers/processes and port 18000 clean after runtime proof.
- Graph: GitNexus impact/detect used around code changes and before closure.

## Acceptance Criteria

- S01 defines same-host contract with endpoint, env, health, retry/timeout, and no-silent-fallback expectations.
- S02 proves packaged restart/cache lifecycle or records a real blocker; Redis L2 restart is no longer silently skipped.
- S03 evaluates 1-2 candidate models only as a bounded legal-domain quick gate and does not widen scope.
- S04 produces a runtime recommendation artifact covering quality, speed, restart/cache, health/preflight, operational simplicity, caveats, and explicit non-goals.

## Open Questions

- Candidate model shortlist is intentionally deferred to S03 so it can be evidence-scoped and stopped if it grows beyond M040.
- Exact runtime recommendation is intentionally deferred to S04 after S01-S03 evidence.
