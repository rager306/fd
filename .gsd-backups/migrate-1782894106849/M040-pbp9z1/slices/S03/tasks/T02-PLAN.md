---
estimated_steps: 16
estimated_files: 1
skills_used: []
---

# T02: Run bounded candidate gate and publish S03 artifact

Expected executor skills for task-plan frontmatter: sentence-transformers, api-design, write-docs, verify-before-complete.

Why: S04 needs one concise evidence row saying whether alternative models should affect the same-host runtime recommendation. This task applies the tooling from T01 to at most two candidates and writes the canonical sanitized S03 artifact.

Do:
1. Use `docs/same-host-embedding-service-contract.md` and D041 semantics: do not test candidates by changing only the request `model`; each candidate must be a separately configured runtime/endpoint or an explicit truthful deferral.
2. Candidate shortlist is capped at two plausible Russian/legal candidates. Prefer candidates that can be served through the existing OpenAI-compatible `/v1/embeddings` path on this host with clear dimensions, such as `BAAI/bge-m3` first and `intfloat/multilingual-e5-large` or another locally available Russian-capable multilingual model second only if it is readily provisioned. Do not broaden the search if these fail.
3. Isolate cache namespaces for baseline and each candidate using candidate-specific labels (for example `m040-s03-baseline-deepvk`, `m040-s03-candidate-bge-m3`) or flush deliberately; record the chosen namespace and avoid reusing S02/ONNX namespaces.
4. For each candidate, first record `/health` model/backend/dimensions/cache namespace and a smoke embedding dimension check. If the model cannot be loaded, endpoint is unavailable, dimensions mismatch, license/format is unclear, or setup requires broad provisioning, stop that candidate and record `defer_candidate` or `reject_candidate` with exact evidence instead of continuing.
5. If a candidate reaches live legal evaluation, run the bounded corpus with small explicit limits (`--max-docs` no more than 256 and query limits no more than 64 each unless already justified by T01 defaults). Compare retrieval recall@k/MRR against the current `deepvk/USER-bge-m3` baseline; do not require cross-model cosine agreement.
6. Write `benchmark-results/fd-legal-model-quick-gate-m040-s03.md` as the canonical S03 result. Include corpus hash, selected doc/query counts, candidate list, endpoint/runtime labels, dimensions, cache namespaces, metrics or stop reasons, `raw_text_logged: false`, caveats about the synthetic/unlabeled corpus, and final recommendation for S04: keep current, reject candidate, or defer candidate.
7. Remove or clearly mark any dry-run artifact from T01 as non-canonical if it remains in `benchmark-results`; the canonical artifact is the only S03 output S04 should consume.

Done when: The canonical artifact validates, considers no more than two candidates, contains no raw legal text or secret material, and gives S04 a bounded recommendation that cannot expand into an open-ended model bakeoff.

Threat Surface (Q3): local candidate endpoints may return untrusted/malformed JSON or secret-bearing diagnostics; artifacts must sanitize errors and never log raw legal text, tokens, signed URLs, private keys, or environment values.
Requirement Impact (Q4): closes R008, supports R001 legal no-regression and R006 evidence envelope. Decisions preserved: D037/D038/D040/D041. Re-test bound count, legal metrics/stop reasons, redaction, and cache namespace isolation.
Failure Modes (Q5): baseline unavailable -> artifact must block/defer live comparison; candidate unavailable -> defer candidate with setup phase; candidate slower without quality upside -> keep current; candidate legal metrics worse -> reject candidate; verifier failure -> do not mark slice complete.
Load Profile (Q6): shared resources are TEI/API containers, Redis L2 cache, CPU/RAM, and model downloads. The quick gate intentionally keeps the corpus bounded; it is not a throughput benchmark except optional speed sanity for legal survivors.
Negative Tests (Q7): verify candidate count remains <=2, artifact redaction passes, no raw corpus sample appears, no request-model selector semantics are claimed, and a candidate without health/smoke success cannot be reported as accepted.

## Inputs

- `tools/evaluate_legal_model_quick_gate.py`
- `tools/verify_legal_model_quick_gate_artifact.py`
- `tests/44-FZ-2026-articles.jsonl`
- `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt`
- `docs/same-host-embedding-service-contract.md`
- `docker-compose.yaml`
- `docker-compose.override.yaml`
- `benchmark.py`

## Expected Output

- `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`

## Verification

python3 tools/verify_legal_model_quick_gate_artifact.py --artifact benchmark-results/fd-legal-model-quick-gate-m040-s03.md --max-candidates 2
python3 -m py_compile tools/evaluate_legal_model_quick_gate.py tools/verify_legal_model_quick_gate_artifact.py

## Observability Impact

Produces the S03 inspection surface for future agents and S04: one sanitized markdown artifact with candidate phases, health/smoke metadata, legal metric results or stop reasons, cache namespace evidence, final bounded verdict, and verifier-backed redaction status.
