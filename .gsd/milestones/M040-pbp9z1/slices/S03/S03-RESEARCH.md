# S03 — Research

**Date:** 2026-05-22

## Summary

S03 should be a bounded legal-domain quick gate, not a model bakeoff. The milestone decision D040 and requirement R008 constrain the slice to 1-2 plausible candidates and require legal-domain evidence before any model can challenge `deepvk/USER-bge-m3`. Because `fd`'s `/v1/embeddings` request `model` field is compatibility metadata only, each candidate must run as a separately configured service/runtime, not as per-request model selection.

The existing legal evaluator is built for TEI-vs-ONNX parity on the same model, not for absolute candidate model selection. It can still be reused as the legal corpus runner and sanitized artifact renderer, but a candidate gate should compare candidate retrieval metrics against the current `deepvk/USER-bge-m3` baseline and record a conservative outcome: keep current, reject candidate, or defer candidate for a dedicated model bakeoff. The quick gate should stop as soon as candidates are unavailable, dimension/API incompatible, slower without quality upside, or legally worse.

Active requirement: R008 is owned by this slice and is the controlling constraint. S03 also supports R001 because any model replacement requires Russian legal corpus evidence. R006 will consume S03 only as one row in the final evidence envelope; S03 must not hijack S04 into open-ended experiments.

## Recommendation

Use a two-phase quick gate:

1. **Candidate shortlist gate:** select at most two candidates that are plausible for Russian/legal embeddings and operationally available on this host through the existing TEI/OpenAI-compatible `/v1/embeddings` path. If no candidate can be started without broad provisioning, write a deferral/reject artifact and stop.
2. **Legal evidence gate:** run the existing `tests/44-FZ-2026-articles.jsonl` corpus through current baseline and each candidate with sanitized outputs. Accept only candidates that meet or exceed baseline legal retrieval metrics within a small bounded run and do not introduce operational incompatibility. Otherwise recommend keeping `deepvk/USER-bge-m3`.

Prefer adding a small candidate-evaluation wrapper/tool only if needed to avoid misusing TEI-vs-ONNX parity semantics. Do not modify the main API to support per-request model selection; S01/D041 explicitly made response model and `/health.runtime.model` authoritative and deployment-scoped.

## Implementation Landscape

### Key Files

- `tools/evaluate_legal_retrieval.py` — Existing sanitized legal retrieval evaluator. It loads `tests/44-FZ-2026-articles.jsonl`, excludes raw text from artifacts, hashes texts, validates embedding dimensions, calculates top-1 agreement/overlap/recall/cross-backend cosine, and writes Markdown. It is parity-oriented but reusable for corpus parsing, request batching, metrics, and redaction patterns.
- `tests/44-FZ-2026-articles.jsonl` — Existing Russian legal corpus used by M039. Prior profile: 94 articles, 1,655 non-invalid candidate documents, 76 title queries; quick gate should use bounded `--max-docs`, `--max-title-queries`, `--max-self-queries` values already supported by the evaluator.
- `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt` — Baseline legal evidence for current `deepvk/USER-bge-m3` TEI vs packaged ONNX parity. It records corpus hash, no raw text, thresholds, selected docs/queries, and PASS.
- `docs/same-host-embedding-service-contract.md` — Contract constraints: no per-request model selection, response model and `/health.runtime.model` are authoritative, alternative model evaluation is explicitly bounded/non-goal outside S03.
- `docker-compose.yaml` / `docker-compose.override.yaml` — Default TEI setup with `MODEL_ID=${MODEL_ID:-deepvk/USER-bge-m3}`. Candidate runs should use separate service/container/port or explicit restart/reconfiguration with isolated cache namespace.
- `api/main.go` — Deployment-time `MODEL_ID`, `EMBEDDING_BACKEND`, health metadata, and cache namespace wiring. Reinforces that model is process configuration, not request routing.
- `benchmark.py` — Optional speed/cache sanity check for a candidate that passes legal quality. Use with candidate-specific namespace and label; not a substitute for legal gate.
- Candidate result target: `benchmark-results/fd-legal-model-quick-gate-m040-s03.md` or one artifact per candidate, with raw text excluded and final keep/defer/reject recommendation.

### Build Order

1. **Confirm exact gate shape before provisioning.** Decide maximum candidate count (1-2), pass/defer/reject outcomes, metrics to compare, and stop conditions. This is the most important S03 control because open-ended experimentation is explicitly out of scope.
2. **Select candidates from availability, not curiosity.** Prefer candidates that can be served by existing TEI on the current host and produce expected dimensions through `/v1/embeddings`. Record if a model is skipped because weights are unavailable, too large, license/format is unclear, or TEI cannot serve it quickly.
3. **Adapt/reuse legal evaluator.** If the current parity script is sufficient, run baseline as `tei` and candidate as `onnx` labels but name labels truthfully (for example `baseline-deepvk` and `candidate-...`) and set candidate dimensions. If cross-backend cosine is not meaningful across different models, add a small wrapper or script mode that compares retrieval metrics (recall/MRR/top-k) rather than vector identity.
4. **Run bounded corpus evidence.** Use the existing corpus, `max_docs`/query limits, batch size, and timeouts. Save sanitized Markdown with corpus hash, model IDs, dimensions, runtime labels, selected counts, thresholds, metrics, and verdict. No raw legal text.
5. **Optional speed sanity only for survivors.** If a candidate clearly passes legal metrics, run a small same-host benchmark/cache check to see whether it challenges the evidence envelope. If legal quality is worse or unavailable, skip speed work and recommend keep current/defer.
6. **Document outcome for S04.** Produce a concise S03 summary: keep current, defer specific candidate, or reject candidate, with exact evidence and caveats.

### Verification Approach

- Dry-run or parse path validates `tests/44-FZ-2026-articles.jsonl` and artifact redaction shape before live candidate calls.
- `/health` for baseline and each candidate reports expected backend/model/dimensions/cache namespace; smoke `/v1/embeddings` returns matching vector dimensions.
- Legal gate artifact(s) contain `raw_text_logged: false` or equivalent statement, corpus hash, selected document/query counts, candidate model IDs, runtime labels, dimensions, metrics, verdict, and caveat.
- If modifying `tools/evaluate_legal_retrieval.py` or adding a wrapper, run Python compile/check and the tool's dry-run mode; if Go code is changed unexpectedly, run `cd api && go test ./... -short`.
- Leak check over S03 artifacts: no raw legal text, secrets, signed URLs, private keys, or token values.
- Bound check: no more than two candidates appear in artifacts unless the user explicitly expands scope.

## Don't Hand-Roll

| Problem | Existing Solution | Why Use It |
|---------|------------------|------------|
| Legal corpus parsing and redacted artifact shape | `tools/evaluate_legal_retrieval.py` | Already hashes raw texts, excludes raw legal corpus text, records corpus/config/thresholds, and handles batching/errors. |
| Same-host API contract checks | `/health` + `/v1/embeddings` smoke from S01 contract | Avoids inventing a second readiness model and respects deployment-scoped model semantics. |
| Cache isolation | Existing `EMBEDDING_CACHE_VERSION`/namespace env and `/health.runtime.cache_namespace` | Prevents candidate vectors contaminating baseline or ONNX namespaces. |

## Constraints

- R008 caps this at 1-2 plausible candidates and legal-domain evidence. Stop rather than expanding into a model search.
- The API does not route by request `model`; candidate model testing requires separate runtime configuration/container/port or explicit service restart.
- Different models may have different embedding dimensions; the evaluator currently enforces one `--dimensions` value per run. Candidate comparisons must configure dimensions correctly and avoid cross-model cosine thresholds that assume identical model spaces.
- Current `tools/evaluate_legal_retrieval.py` pass criteria are for TEI-vs-ONNX parity of the same model. For different models, top-1 agreement and cross-backend cosine against baseline are not necessarily valid replacement metrics; retrieval quality against known legal query/doc relationships is the safer comparator.
- Raw legal text must never be recorded in artifacts.
- Redis namespaces must be isolated per candidate to avoid reusing baseline vectors.

## Common Pitfalls

- **Turning S03 into a bakeoff** — Do not search broadly once one or two plausible candidates are selected; write defer/reject if setup becomes expensive.
- **Misusing cross-model cosine** — Cross-backend cosine near 1.0 only makes sense for the same model/runtime parity. Different embedding models live in different vector spaces, so use retrieval metrics and known-item behavior instead.
- **Trusting request `model`** — S01/D041 says request `model` is OpenAI metadata only. Starting a candidate requires runtime config changes.
- **Dimension mismatch** — Candidate vectors may not be 1024-dimensional. Configure evaluator dimensions per candidate and record them.
- **Cache contamination** — Candidate and baseline caches must use disjoint namespaces.

## Open Risks

- Candidate model availability on the current host may be poor or TEI may fail to load alternatives within reasonable time. That is an acceptable deferral outcome for S03.
- The existing legal corpus has synthetic known-item/title/self-document metrics and is not a human-labeled relevance benchmark. A candidate should not replace current model based solely on this quick gate; at most it can justify a later dedicated bakeoff.
- If `tools/evaluate_legal_retrieval.py` needs a new candidate-comparison mode, that is a code change touching evaluation semantics and should get tests/dry-run verification.

## Skills Discovered

| Technology | Skill | Status |
|------------|-------|--------|
| Embedding models / sentence transformers | sentence-transformers | available |
| API contract / model semantics | api-design | available |
| Bounded decision interrogation | grill-me | available |
| Documentation artifact writing | write-docs | available |
| TEI / Hugging Face serving | none installed in available skills | none found/installed |

## Sources

- Existing legal evaluator behavior from `tools/evaluate_legal_retrieval.py`.
- M039 legal PASS artifact from `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt`.
- Same-host contract constraints from `docs/same-host-embedding-service-contract.md` and D041.
- Project memory: MEM041 (same-host legal/speed boundary), MEM044 (bounded alternative gate), MEM040 (avoid open-ended ONNX/model experimentation), MEM015 (prior alternative-model research track only after legal divergence context).
