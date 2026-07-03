# S04 Research: Runtime recommendation and operating contract

## Summary

S04 is a documentation/evidence-integration slice, not a new runtime experiment. S01 already produced the same-host HTTP contract and safe `/health.runtime` semantics; S02 proved packaged ONNX Docker restart/cache behavior with legal no-regression evidence; S03 bounded alternative-model exploration and failed closed. The remaining work is to publish a final TEI-vs-ONNX recommendation and operating contract that a fresh local-service operator can act on.

The evidence points to recommending **the current `deepvk/USER-bge-m3` model on the packaged ONNX backend as the preferred same-host performance runtime when deployed with the S01 contract and isolated cache namespace**, while keeping the important caveat that TEI remains the current production/default unless the operator explicitly switches runtimes. Alternative model replacement should remain deferred. ONNX clears legal parity and local speed/restart/cache evidence; the main caveat is operational rollout discipline, especially runtime metadata, `ONNX_RUNTIME_SHA256`, and cache namespace isolation.

## Active Requirements S04 Owns or Supports

- **R006 (primary owner):** final TEI-vs-ONNX recommendation must use the evidence envelope: legal quality, same-host performance, restart/cache behavior, health/preflight clarity, and operational simplicity.
- **R007 (primary owner):** do not make hosted GitHub Actions proof, remote dispatch, push/upload, or artifact mirroring a readiness gate.
- **R001 (support):** recommendation must preserve Russian/legal quality; S02 legal retrieval is the current proof, S03 prevents unproven model replacement.
- **R002 (support):** operating contract must preserve long-lived Redis L2 reuse; S02 proved L2 survives API restart.
- **R003 (support):** operating contract must restate env/runtime/cache knobs with safe defaults and validation boundaries.
- **R004 (support):** final artifact/verifier should preserve sanitized config/evidence comparability and leak boundaries.
- **R005 (support):** final recommendation should consume, not replace, the same-host `/health`, `/v1/embeddings`, and `/embeddings/batch` contract.
- **R008 (support, already validated):** alternative candidates remain bounded and deferred fail-closed.
- **R009 (support):** no silent runtime/model/tokenizer fallback; rollback is explicit restart/reconfiguration.

## Skills Applied

The requested skills were not available as a callable tool in this harness, so I read their installed `SKILL.md` files and applied these rules:

- **api-design:** callers outlive assumptions; document endpoint/status/retry/error semantics honestly; do not imply `200` health equals live inference readiness.
- **design-an-interface:** surface distinct shapes before choosing. For S04 the viable shapes are: (1) narrative doc only, (2) benchmark artifact only, or (3) machine-checkable recommendation artifact plus links to the contract. Recommend shape (3): an auditable recommendation artifact with a verifier, optionally linked from docs/README.
- **grill-me:** codebase before question. The load-bearing decisions are already resolved by D038-D041 and S01-S03; no user interview is needed unless the executor wants to change the recommendation away from evidence.
- **observability:** log/record decisions, not activity. For S04, the recommendation artifact should expose the decision inputs and caveats in a future-agent-readable table instead of prose-only claims.
- **write-docs:** write for a fresh reader with a named post-read action: “choose and operate the local embedding runtime safely.” Avoid burying non-goals and caveats.

## Skills Discovered

Installed skills already cover the direct cross-cutting work: `api-design`, `observability`, `write-docs`, and GitNexus skills. I ran `npx skills find` for direct technologies:

- `go http service`: found generic Go service skills, but S04 should not need Go implementation changes.
- `docker redis`: found Redis/Docker skills, but S02 already closed runtime proof and S04 only consumes it.
- `onnx runtime`: no skills found.
- `text embeddings`: found generic embedding-provider skills, but they are not specific to this local fd/TEI/ONNX evidence envelope.

No new skills were installed because the slice is evidence synthesis, not novel implementation in those technologies.

## Evidence Inventory for the Planner

### S01 contract inputs

- `docs/same-host-embedding-service-contract.md` is canonical for neighboring local HTTP clients.
- `/health` is metadata, not a live inference probe. Full readiness still needs a smoke `POST /v1/embeddings`.
- `/v1/embeddings` request `model` is compatibility metadata only; response `model` and `/health.runtime.model` are authoritative.
- TEI/default health includes only safe runtime fields: backend, model, dimensions, production_default, cache_namespace.
- ONNX health adds artifact/tokenizer/runtime verification fields, provider, sequence length, artifact_id.
- Cache namespace isolation is mandatory when comparing/switching TEI and ONNX.
- No silent fallback: backend is startup-only; rollback/switch is restart + cache namespace/flush discipline.

### S02 packaged ONNX proof inputs

Key files:

- `tools/run_m040_s02_docker_restart_proof.sh` — reproducible packaged ONNX restart/cache proof runner.
- `tools/verify_m040_s02_artifacts.py` — semantic verifier for S02 artifacts.
- `benchmark-results/fd-m040-s02-onnx-docker-preflight.txt` — build/start/health/smoke evidence.
- `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt` — benchmark and restart/cache evidence.
- `benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt` — legal TEI-vs-ONNX parity gate.
- `benchmark-results/fd-m040-s02-proof-audit.txt` — leak/cleanup audit.

Important facts:

- ONNX health: backend `onnx`, model `deepvk/USER-bge-m3`, dimensions `1024`, `artifact_verified=true`, `tokenizer_verified=true`, provider `CPUExecutionProvider`, cache namespace includes `m040-s02-onnx-restart`.
- Caveat: `runtime_library_verified=false` in the proof health because `ONNX_RUNTIME_SHA256` was not set in that run. The operating contract should recommend setting it for stronger preflight integrity, not treat it as a quality failure.
- Packaged ONNX benchmark summary: best cold latency `11.9ms`, warm latency mean `1.79ms`, max throughput `~799 req/s` at 4 concurrent, Redis L2 restart `4.73ms after API restart`, batch L2 p95 `3.89ms after API restart`, chunk reuse warm p95 `7.96ms`.
- Legal retrieval gate: `PASS`; cross-backend cosine minimum `0.99989883`, ONNX recall ratio `1.0`, mean overlap@5 `0.997701`, top1 agreement `1.0`, ONNX/TEI MRR effectively equal (`0.86256` vs `0.862557`). Caveat: unlabeled corpus proves TEI-vs-ONNX parity and synthetic known-item behavior, not absolute human relevance.
- Cleanup audit: no prohibited leaks; ONNX API proof container removed; port 18000 clear; Redis proof container intentionally left running for local cache reuse.

### S03 bounded candidate gate inputs

Key files:

- `tools/evaluate_legal_model_quick_gate.py`
- `tools/verify_legal_model_quick_gate_artifact.py`
- `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`

Important facts:

- Candidate count is exactly 2: `BAAI/bge-m3` and `intfloat/multilingual-e5-large`.
- Canonical verdict: `defer_candidate`.
- Stop reason: current same-host baseline `/health` lacked the contract-required runtime object at runtime, so live candidate comparison was unavailable.
- Cross-model cosine/top-1 parity is explicitly not an acceptance metric for replacement; retrieval metrics and operational compatibility are required.
- S04 should keep `deepvk/USER-bge-m3` as the model for recommendation and defer model replacement until future endpoints expose `/health.runtime` and pass legal retrieval metrics.

### Prior performance context for TEI vs ONNX

Useful historical benchmark summaries in `benchmark-results/`:

- `fd-benchmark-m014-tei-baseline.txt`: TEI best cold `59.0ms`, warm mean `2.25ms`, max throughput `~750 req/s`, Redis L2 restart `2.82ms`.
- `fd-benchmark-m014-onnx-hf-tokenizer.txt`: ONNX best cold `10.2ms`, warm mean `1.63ms`, max throughput `~891 req/s`, Redis L2 restart `2.70ms`.
- `fd-benchmark-m039-docker-onnx-target-runtime.txt`: packaged ONNX best cold `10.3ms`, warm mean `1.52ms`, max throughput `~937 req/s`, but Redis L2 restart was skipped.
- `fd-benchmark-m040-s02-onnx-docker-restart.txt`: packaged ONNX closes the restart gap with Redis L2 restart measured.

Use M040/S02 as the primary ONNX evidence. Use older TEI baseline only as context unless the executor chooses to rerun TEI.

## Implementation Landscape

Likely files to create/modify:

- **Create** `benchmark-results/fd-runtime-recommendation-m040-s04.md` (or equivalent canonical final artifact). It should be the auditable evidence envelope and operating contract.
- **Create** `tools/verify_m040_s04_recommendation.py` to validate the final artifact has required sections, recommendation stance, evidence citations, non-goals, redaction/no-secret boundaries, and S02/S03 semantic inputs. Follow the style of `tools/verify_m040_s02_artifacts.py` and `tools/verify_legal_model_quick_gate_artifact.py`.
- **Optionally modify** `README.md` and/or `docs/same-host-embedding-service-contract.md` to link the final recommendation artifact. If modifying existing docs, keep the canonical endpoint contract in S01 doc; do not duplicate all recommendation content.
- **Do not modify runtime code** unless execution discovers a hard acceptance gap. S04 should consume prior runtime evidence, not rewire health or embeddings behavior.

Suggested final artifact sections:

1. `# M040 S04 Runtime Recommendation and Operating Contract`
2. `## Recommendation` — explicit stance. Suggested wording: recommend packaged ONNX as the preferred same-host performance runtime for `deepvk/USER-bge-m3` when deployed with S01 contract, isolated Redis namespace, smoke readiness, and ONNX artifact/tokenizer/runtime preflight; TEI remains current production/default until explicit operator switch.
3. `## Evidence Envelope` — table/JSON-like block for legal quality, speed, restart/cache, health/preflight, operational simplicity, candidate models.
4. `## Operating Contract` — endpoints/readiness, timeout/retry, env vars, cache namespace, restart/rollback, no-silent-fallback.
5. `## Caveats and Required Operator Checks` — runtime library SHA, health metadata, smoke request, cache isolation, raw text redaction, old/stale service warning.
6. `## Non-Goals` — no hosted CI gate, no remote deployment proof, no request-level fallback, no per-request model selection, no open-ended model bakeoff.
7. `## Source Artifacts` — list S01/S02/S03 artifacts and what each contributes.
8. `## Redaction` — affirm no raw legal/probe text, secrets, signed URLs.

Verifier expectations should prefer semantics over timing thresholds. Avoid brittle exact latency matching except for checking that the artifact includes the M040/S02 summary values or source artifact references.

## Natural Seams / Task Candidates

1. **Evidence extraction + recommendation shape**
   - Read S01/S02/S03 artifacts and produce a compact evidence table.
   - Decide exact recommendation language and caveat stance.
   - No code changes required.

2. **Final artifact authoring**
   - Write `benchmark-results/fd-runtime-recommendation-m040-s04.md` with the sections above.
   - Ensure a fresh operator can answer: which runtime, under what preconditions, how to verify readiness, how to rollback, what not to do.

3. **Artifact verifier**
   - Add `tools/verify_m040_s04_recommendation.py`.
   - Validate required headings, required source references, explicit recommendation, `defer_candidate`, S02 legal/restart evidence, no hosted CI requirement, no raw text/secrets.
   - Include `--self-test` with a passing and failing fixture if time permits, mirroring S03 verifier style.

4. **Discoverability/update docs if needed**
   - Optionally add a small link from README or `docs/same-host-embedding-service-contract.md` related documents.
   - Do not duplicate S01 endpoint semantics; the final recommendation should point back to the contract.

5. **Closeout verification and GSD completion**
   - Run verifier, Python compile, leak checks, and any doc-link checks.
   - Use `gsd_task_complete`/slice completion in executor phases.

## First Proof / Highest-Risk Unblocker

Build the verifier early, before polishing prose. The highest risk is a plausible-sounding recommendation that accidentally violates an envelope constraint: recommending an unproven alternative model, making hosted CI a blocker, omitting cache namespace isolation, or claiming `/health` is live inference readiness.

First proof should be something like:

```bash
python3 tools/verify_m040_s04_recommendation.py \
  --artifact benchmark-results/fd-runtime-recommendation-m040-s04.md \
  --s02-benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt \
  --s02-legal benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt \
  --s03-gate benchmark-results/fd-legal-model-quick-gate-m040-s03.md
```

The verifier should fail until the artifact contains the correct recommendation stance and evidence/caveats.

## Verification Plan

Minimum commands/checks for executors:

```bash
python3 -m py_compile tools/verify_m040_s04_recommendation.py
python3 tools/verify_m040_s04_recommendation.py --self-test
python3 tools/verify_m040_s04_recommendation.py \
  --artifact benchmark-results/fd-runtime-recommendation-m040-s04.md \
  --s02-benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt \
  --s02-preflight benchmark-results/fd-m040-s02-onnx-docker-preflight.txt \
  --s02-legal benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt \
  --s02-audit benchmark-results/fd-m040-s02-proof-audit.txt \
  --s03-gate benchmark-results/fd-legal-model-quick-gate-m040-s03.md
```

Optional if docs/README are modified:

```bash
rg -n "fd-runtime-recommendation-m040-s04|Runtime Recommendation|same-host" README.md docs/same-host-embedding-service-contract.md
```

Leak/redaction check should include the final artifact and existing M040 artifacts. Reuse prohibited patterns from S02/S03 verifiers: private keys, bearer/API tokens, signed URLs, raw benchmark probes, raw legal text markers. The final artifact should cite counts/hashes/metrics only.

If code is modified despite this research recommending docs/tooling only, run:

```bash
(cd api && go test ./... -short)
```

Also follow project policy: run GitNexus impact before editing existing symbols and `gitnexus_detect_changes()` before closure if those tools are available to the executor.

## Risks and Watch-outs

- **Current live `http://127.0.0.1:8000/health` may be stale.** A quick check during research returned `status=ok` but no runtime keys; this matches S03’s fail-closed stop reason. Do not use the currently running baseline service as evidence for S04 unless it is rebuilt/restarted with S01 changes and verified.
- **ONNX proof has `runtime_library_verified=false`.** This does not invalidate quality/performance, but the final operating contract should instruct operators to set `ONNX_RUNTIME_SHA256` when they want runtime-library integrity surfaced as true.
- **Cache contamination remains the correctness footgun.** Any TEI-vs-ONNX switch/comparison must use isolated `EMBEDDING_CACHE_VERSION` or deliberate Redis flush with recorded side effects.
- **Do not overclaim legal quality.** S02 proves TEI-vs-ONNX parity on the sanitized 44-FZ corpus, not absolute human relevance. Phrase quality as no-regression/parity.
- **Alternative model replacement is not available.** S03 deferred candidates fail-closed; S04 should not recommend BAAI or E5.
- **Do not reintroduce hosted GitHub Actions as a blocker.** R007 and MEM039 explicitly exclude this from M040 readiness.
- **Avoid path/secret leaks in public docs.** The final artifact may cite file paths because it is an evidence artifact; any README/contract doc changes should be concise and not expose raw text, tokens, signed URLs, or local absolute runtime paths.

## Recommended S04 Outcome

Publish a machine-checkable final artifact that says:

- Use `deepvk/USER-bge-m3`.
- Prefer packaged ONNX for same-host performance once explicitly deployed with S01 health metadata, artifact/tokenizer verification, optional runtime-library SHA verification, isolated Redis namespace, and smoke readiness.
- Keep TEI as the current/default fallback posture until an operator intentionally switches; do not implement request-level fallback.
- Defer alternative model replacement.
- Treat hosted/remote CI proof as out of scope.

This completes R006/R007 and ties off R001-R005/R008/R009 without more runtime work.