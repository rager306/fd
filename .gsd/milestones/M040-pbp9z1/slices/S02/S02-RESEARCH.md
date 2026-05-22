# S02 — Research

**Date:** 2026-05-22

## Summary

S02 should close the exact gap left by M039: the packaged ONNX Docker benchmark already proved smoke, legal parity, and fast cache behavior, but its Redis L2 restart sections were skipped because `BENCHMARK_API_RESTART_COMMAND` was empty. The codebase already has the right seam: `benchmark.py` reads `BENCHMARK_API_RESTART_COMMAND`, runs it through `restart_api_for_l2_check()`, waits for `/health`, then measures single-request and batch Redis L2 behavior after API restart.

The recommended implementation is therefore a thin, benchmark-compatible Docker restart proof rather than a benchmark rewrite. Use the existing packaged ONNX image/build path, run the ONNX API on the same-host proof port, set isolated cache namespace values, set `ONNX_RUNTIME_SHA256` where available so `/health.runtime.runtime_library_verified=true`, and invoke `benchmark.py` with a small restart command that restarts only the packaged ONNX API container/process while leaving Redis alive.

Active requirements supported: R001 requires no quality regression for optimized runtime evidence; S02 should reuse the M039 legal gate or rerun it if runtime/container settings change. R002 requires long-lived Redis L2 reuse, which is the core proof for this slice. R004 requires sanitized effective configuration in artifacts; `benchmark.py` already records allowlisted env, Docker compose hashes/images, Redis metadata, runtime artifact metadata, and excludes raw benchmark texts.

## Recommendation

Take the existing hook-first approach from D039: provide a small local Docker restart command through `BENCHMARK_API_RESTART_COMMAND` and run the existing `benchmark.py` against the packaged ONNX API. Avoid adding first-class Docker lifecycle logic to `benchmark.py` unless the hook is insufficient, because the benchmark already records the required sections and comparability metadata.

The first proof should be operational, not code-heavy: confirm the packaged ONNX API can start, `/health` reports ONNX backend and verified runtime metadata, a smoke embedding works, Redis is not restarted, and Section 5/6 of `benchmark.py` no longer say `skipped`. If the Docker restart cannot be controlled in this environment, the executor should write a truthful blocker artifact with exact command/output rather than silently accepting another skipped benchmark.

## Implementation Landscape

### Key Files

- `benchmark.py` — Primary benchmark harness. Relevant seams: `API_RESTART_COMMAND = os.getenv("BENCHMARK_API_RESTART_COMMAND", "docker compose restart api")`, `restart_api_for_l2_check()`, Section 5 single-request Redis L2 restart proof, Section 6 batch Redis L2 restart proof, and sanitized `effective_config_snapshot()`.
- `tools/build_onnx_image.sh` — Builds the dedicated packaged ONNX image from local verified ONNX/model/tokenizer/runtime artifacts into `.gsd/runtime/docker/onnx1024-context` and tags it as `fd-api:onnx-1024` by default.
- `Dockerfile.onnx` — Packaged ONNX API image. Sets `EMBEDDING_BACKEND=onnx`, `ONNX_ARTIFACT_MANIFEST`, `ONNX_MODEL_PATH`, `ONNX_TOKENIZER_PATH`, `ONNX_RUNTIME_LIBRARY`, and `ONNX_MAX_SEQUENCE_LENGTH=1024` inside the container.
- `docs/same-host-embedding-service-contract.md` — S01 contract. Requires same-host clients and operators to inspect `/health.runtime.backend`, `runtime.model`, `runtime.dimensions`, `runtime.runtime_library_verified` for ONNX, and `runtime.cache_namespace`; also documents that Redis namespace isolation or flush is required when switching TEI/ONNX.
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` — ONNX artifact/runtime/source contract consumed by packaged proof and benchmark metadata.
- `benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt` — Prior M039 evidence. Important because it shows `BENCHMARK_API_RESTART_COMMAND` was empty and both Redis L2 restart sections were skipped.
- `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt` — Prior packaged ONNX legal parity PASS for `deepvk/USER-bge-m3`; reuse as baseline unless runtime/container inputs change.
- `docker-compose.yaml` / `docker-compose.override.yaml` — Existing TEI/API/Redis Compose setup. Redis is exposed to host only through the local override; do not expose Redis broadly.
- `README.md` — Documents benchmark behavior, Redis FLUSHALL side effect, and Redis cache persistence/tuning warnings.

### Build Order

1. **Preflight current packaged runtime inputs.** Confirm `tools/build_onnx_image.sh` prerequisites exist and build or reuse the `fd-api:onnx-1024` image. Capture if the ONNX Runtime library or tokenizer artifacts are missing; that is a real blocker.
2. **Start packaged ONNX API with Redis kept external/alive.** Use a distinct container name and host port (M039 used `http://localhost:18000`) and set an S02-specific cache namespace such as `EMBEDDING_CACHE_VERSION=m040-s02-onnx-restart`. Keep Redis on `127.0.0.1:6379`; do not restart Redis as part of the API restart command.
3. **Health/smoke before benchmark.** Check `/health` for `backend=onnx`, `model=deepvk/USER-bge-m3`, `dimensions=1024`, expected `cache_namespace`, and `runtime_library_verified=true` when `ONNX_RUNTIME_SHA256` is configured. Then run a smoke `/v1/embeddings` request because `/health` is not an inference probe.
4. **Run `benchmark.py` with the restart hook.** Set `BENCHMARK_API_URL=http://localhost:18000`, `BENCHMARK_RUNTIME_LABEL=docker-onnx-go-api-m040-s02`, `BENCHMARK_API_RESTART_COMMAND='<docker restart/wait command>'`, ONNX manifest/runtime env values, and isolated cache namespace env. Redirect output to a new `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt`.
5. **Validate artifact semantics.** The output must show `api_restart_command_configured: true`, Section 5 not skipped, Section 6 batch L2 not skipped, summary `Redis L2 restart:` with latency, and sanitized config with `input_texts_logged: false`.
6. **Quality guard if runtime inputs changed.** If image/runtime/model/tokenizer differ from M039 accepted inputs, rerun `tools/evaluate_legal_retrieval.py` against TEI and packaged ONNX using an isolated ONNX namespace and save `benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt`.

### Verification Approach

- `tools/build_onnx_image.sh` succeeds or records exact missing artifact blocker.
- Packaged ONNX `/health` returns safe ONNX runtime metadata with expected backend/model/dimensions/cache namespace and no secrets/paths beyond the contractually allowed fields.
- Smoke embedding against `http://localhost:18000/v1/embeddings` returns 1024 dimensions and response model `deepvk/USER-bge-m3`.
- `uv run --python 3.13 --with requests --with redis python benchmark.py` (or equivalent existing invocation) produces a benchmark artifact where Redis L2 restart and batch L2 restart are measured, not skipped.
- Artifact leak check: no raw benchmark/legal probe text, private keys, bearer tokens, signed URL query parameters, or secrets in `benchmark-results/*m040-s02*`.
- Cleanup: no leftover M040 proof container/processes, and port 18000 is clear after proof.
- If code changes are needed, run `cd api && go test ./... -short`; if only scripts/docs/artifacts change, compile/check touched Python scripts and use GitNexus detect before closeout.

## Constraints

- Redis must survive the API restart. Restarting Docker Compose as a stack or restarting Redis invalidates the L2 persistence proof.
- Use isolated Redis namespace values for ONNX proof. Memory MEM008 warns TEI/ONNX cache sharing can make false equivalence by returning the other runtime's cached vectors.
- `benchmark.py` intentionally calls `FLUSHALL`; do not run this proof against shared/production Redis.
- `/health` is not a live inference probe; S01 requires a smoke embedding request for full readiness.
- Packaged ONNX containers should include `ONNX_RUNTIME_SHA256`; without it `/health.runtime.runtime_library_verified` is weaker and S04 must record that caveat.
- Raw legal/probe texts and secrets must not be written to public artifacts.

## Common Pitfalls

- **Silent skipped proof** — If `BENCHMARK_API_RESTART_COMMAND` is unset/empty or fails, the benchmark will continue with skipped L2 sections. Treat skipped restart as S02 failure/blocker, not success.
- **Restarting Redis with the API** — This destroys the L2 survival signal. The restart command should target only the API container/process.
- **Cache contamination** — Reusing default namespaces across TEI and ONNX can return stale vectors. Set `EMBEDDING_CACHE_VERSION` and model/revision/tokenizer namespace components where relevant.
- **Over-trusting `/health`** — `/health` can be ok while inference is broken; always include smoke embedding.
- **Leaky artifacts** — Benchmark metadata is sanitized, but any ad-hoc command output must be checked before saving.

## Open Risks

- The current host may not have Docker permissions, the packaged image context, ONNX runtime library, or port 18000 available. If so, record exact blocker evidence.
- `BENCHMARK_API_RESTART_COMMAND` uses `shell=True`; keep the command simple and non-secret, and do not interpolate untrusted values.
- If the packaged proof uses `docker run` outside Compose, `benchmark.py`'s `docker compose config/images` metadata may not fully describe the running container. Capture explicit `docker ps`/image digest evidence in the benchmark wrapper artifact if needed.

## Skills Discovered

| Technology | Skill | Status |
|------------|-------|--------|
| Docker / Docker Compose | none installed in available skills | none found/installed |
| Redis | none installed in available skills | none found/installed |
| Go API / HTTP contract | api-design | available |
| Observability / artifact safety | observability | available |

## Sources

- Existing restart hook and benchmark sections from `benchmark.py`.
- Prior skipped restart evidence from `benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt`.
- Prior packaged legal PASS from `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt`.
- Same-host contract and cache isolation rules from `docs/same-host-embedding-service-contract.md`.
- Project memory: MEM043 (scripted Docker restart proof), MEM008 (Redis namespace contamination), MEM024 (packaged ONNX performance viability), MEM003 (long-lived Redis L2 policy).
