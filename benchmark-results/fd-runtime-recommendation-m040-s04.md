# M040 S04 Runtime Recommendation

## Recommendation

Use `deepvk/USER-bge-m3` as the embedding model. Prefer packaged ONNX as the same-host performance runtime only when an operator explicitly deploys it under the S01 same-host HTTP contract: `/health.runtime` metadata, artifact and tokenizer verification, optional `ONNX_RUNTIME_SHA256` runtime-library integrity, isolated `EMBEDDING_CACHE_VERSION`, and a smoke `POST /v1/embeddings` readiness check.

TEI remains the current/default posture until an operator explicitly switches to ONNX with `EMBEDDING_BACKEND=onnx` and the required ONNX environment. The service must run exactly one backend per process lifetime: no request-level fallback, no silent TEI fallback inside an ONNX process, and no per-request model selection.

Alternative model replacement is not accepted by this recommendation. S03 records `defer_candidate` / fail-closed outcomes for the bounded candidates, so BAAI/bge-m3 and intfloat/multilingual-e5-large are not recommended replacements for `deepvk/USER-bge-m3`.

## Decision Inputs

- S01 same-host contract requires HTTP consumers to treat `/health` as liveness/configuration metadata and to use smoke embedding for end-to-end inference readiness.
- S02 packaged ONNX evidence used `deepvk/USER-bge-m3`, ONNX artifact `user-bge-m3-dense-fp32`, dimensions `1024`, runtime label `onnx-docker-m040-s02`, and cache namespace seed `m040-s02-onnx-restart`.
- S02 preflight observed `/health.runtime.backend=onnx`, `/health.runtime.model=deepvk/USER-bge-m3`, `/health.runtime.artifact_id=user-bge-m3-dense-fp32`, `artifact_verified=true`, `tokenizer_verified=true`, and `runtime_library_verified=false` because `ONNX_RUNTIME_SHA256` was not set.
- S02 benchmark evidence showed API-only restart with Redis L2 reuse: after API restart the benchmark recorded `redis_delta/l2_after_api_restart` with Redis hits and no misses, plus cached batch L2 reuse after restart.
- S02 legal retrieval parity evidence recorded `PASS` for runtime label `docker-onnx-go-api-m040-s02`, cache namespace `m040-s02-onnx-restart`, model `deepvk/USER-bge-m3`, `raw_text_logged=false`, and no regression against the TEI comparison envelope.
- S02 proof audit recorded legal gate `PASS`, leak audit `PASS`, no prohibited patterns found, cleanup of the S02 proof ONNX API container, port 18000 clear, and no blockers.
- S03 bounded candidate gate ended with `defer_candidate`; BAAI/bge-m3 and intfloat/multilingual-e5-large remained `failed_closed`, and cross-model cosine/parity was explicitly not an acceptance metric.

## Same-Host Operating Contract

1. Keep TEI as the production/default runtime unless the operator intentionally restarts the service with ONNX configuration.
2. For ONNX, set `EMBEDDING_BACKEND=onnx` and supply the packaged ONNX manifest, tokenizer path, and ONNX Runtime library path before startup.
3. Treat `GET /health` as liveness and runtime metadata only. `/health` does not perform live inference and does not prove vector correctness; readiness requires a smoke embedding request.
4. Verify `/health.runtime` reports the intended backend, model, artifact, dimensions, verification booleans, provider, and `cache_namespace` before directing clients to the instance.
5. Run a smoke `POST /v1/embeddings` with short non-sensitive input and verify the response shape, model, dimensions, and embedding vector presence.
6. Isolate Redis cache namespace whenever comparing or switching TEI and ONNX. Use a disjoint `EMBEDDING_CACHE_VERSION` and verify `runtime.cache_namespace`; otherwise ONNX can serve cached TEI vectors and produce false equivalence.
7. Preserve the no-silent-fallback rule: startup should fail closed when ONNX preflight fails, and requests should not fall back to a different backend or model.
8. Hosted CI or GitHub Actions are not a readiness gate for this same-host runtime contract; same-host readiness is established by local runtime metadata plus smoke embedding.

## Required Operator Checks

Before switching from TEI to ONNX, perform these checks and record their results in the deployment notes:

- Configure `EMBEDDING_BACKEND=onnx`.
- Configure `ONNX_ARTIFACT_MANIFEST` for the packaged `user-bge-m3-dense-fp32` manifest and confirm the manifest is readable.
- Configure `ONNX_TOKENIZER_PATH` and confirm tokenizer verification succeeds for `deepvk/USER-bge-m3`.
- Configure `ONNX_RUNTIME_LIBRARY` and confirm the runtime library loads.
- If an integrity boundary is required, set `ONNX_RUNTIME_SHA256`; without it, `runtime_library_verified=false` is expected and must be documented as a weaker integrity posture rather than a runtime failure.
- Confirm `/health.runtime.artifact_verified=true` and `/health.runtime.tokenizer_verified=true`.
- Confirm `/health.runtime.model=deepvk/USER-bge-m3`, `/health.runtime.backend=onnx`, `/health.runtime.dimensions=1024`, and the expected provider.
- Confirm Redis cache namespace isolation with `EMBEDDING_CACHE_VERSION`; for S02 the evidence namespace seed was `m040-s02-onnx-restart`.
- Run a smoke embedding request through `POST /v1/embeddings` and verify response shape rather than relying on `/health` alone.
- Confirm rollback remains a service restart back to TEI and does not rely on request-level fallback.

## Evidence Links

- `tools/verify_m040_s02_artifacts.py`
- `tools/verify_legal_model_quick_gate_artifact.py`
- `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt`
- `benchmark-results/fd-m040-s02-onnx-docker-preflight.txt`
- `benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt`
- `benchmark-results/fd-m040-s02-proof-audit.txt`
- `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`
- `docs/same-host-embedding-service-contract.md`

## Caveats

- The S02 legal quality evidence is bounded to the recorded same-host TEI-vs-ONNX parity gate. It supports no-regression for the measured Russian legal retrieval envelope; it does not prove universal legal quality or replacement suitability for other corpora.
- S02 proved the packaged ONNX path under local Docker/API conditions and recorded evidence values, not future performance thresholds. Operators should treat latency and throughput numbers as environment-specific diagnostics.
- `runtime_library_verified=false` in S02 because `ONNX_RUNTIME_SHA256` was not set. This does not invalidate the run, but it means the runtime library integrity boundary was weaker than a SHA256-pinned deployment.
- `/health` alone does not prove live inference readiness, vector correctness, or cache behavior. A smoke embedding request is required before routing consumers.
- Cache isolation is mandatory for meaningful comparison; a shared Redis namespace can hide backend differences by returning old cached vectors.
- ONNX must fail closed with no fallback and no silent fallback to TEI inside an ONNX process.
- Hosted CI is not a readiness gate for same-host operation; local metadata and smoke embedding are the readiness evidence.

## Non-Goals

- Accepting BAAI/bge-m3, intfloat/multilingual-e5-large, or any other alternative model as a replacement is a non-goal. S03 leaves those candidates deferred and fail-closed.
- Establishing a general model registry, open-ended model bake-off, or cross-model cosine acceptance metric is out of scope.
- Requiring hosted CI or GitHub Actions as same-host runtime readiness evidence is a non-goal; hosted CI is not a readiness gate.
- Adding request-level TEI-to-ONNX or ONNX-to-TEI fallback is a non-goal because it would undermine reproducibility and cache isolation.
- Replacing the S01 HTTP contract with a library or embedded integration contract is out of scope.
- Proving INT8, NUMA, OpenVINO, GPU, or non-packaged ONNX variants is out of scope for this recommendation.

## Redaction

This recommendation intentionally keeps raw benchmark text, raw legal corpus text, smoke payload text, secrets, tokens, signed URLs, and private material excluded. It references only sanitized evidence markers, model/runtime identifiers, booleans, dimensions, cache namespace labels, verdicts, and source artifact paths. Source artifacts used by the verifier also state `raw_text_logged=false` or equivalent not logged redaction status where applicable.
