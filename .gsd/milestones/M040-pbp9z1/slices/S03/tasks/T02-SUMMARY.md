---
id: T02
parent: S03
milestone: M040-pbp9z1
key_files:
  - benchmark-results/fd-legal-model-quick-gate-m040-s03.md
  - benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md
key_decisions:
  - Treat HTTP 200 `/health` without runtime metadata as baseline unavailable and defer candidate evaluation fail-closed.
  - Keep candidate scope capped at `BAAI/bge-m3` and `intfloat/multilingual-e5-large`; do not broaden the model search.
duration: 
verification_result: passed
completed_at: 2026-05-22T08:13:38.713Z
blocker_discovered: false
---

# T02: Published the canonical sanitized S03 legal model quick-gate artifact with two bounded candidates and a fail-closed defer verdict.

**Published the canonical sanitized S03 legal model quick-gate artifact with two bounded candidates and a fail-closed defer verdict.**

## What Happened

Generated `benchmark-results/fd-legal-model-quick-gate-m040-s03.md` using the T01 evaluator against the intended same-host API surface instead of changing only the request `model`. The local API and Docker services were checked first; `/health` returned HTTP 200 but did not include the contract-required runtime metadata block, so the evaluator correctly stopped before smoke embeddings or legal retrieval metrics and recorded a baseline-unavailable deferral. The candidate shortlist remained capped at two plausible multilingual/Russian-capable models (`BAAI/bge-m3` and `intfloat/multilingual-e5-large`), each represented as a separate endpoint/runtime label with isolated cache namespaces (`m040-s03-candidate-bge-m3` and `m040-s03-candidate-multilingual-e5-large`). The artifact records corpus hash, bounded doc/query counts, candidate list, endpoint labels, expected dimensions, cache namespaces, stop reasons, redaction status, and the S04 recommendation to defer candidate replacement rather than expand into an open-ended bakeoff. The T01 dry-run artifact was retained but clearly marked non-canonical so S04 consumes only the canonical T02 artifact.

## Verification

Verified the canonical artifact with `tools/verify_legal_model_quick_gate_artifact.py --max-candidates 2`, compiled both quick-gate scripts with `py_compile`, and ran additional content checks for candidate count, redaction, isolated namespaces, and deferred outcome. The artifact contains no raw legal corpus text or secret material per verifier-backed redaction checks.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker compose ps && curl -fsS --max-time 5 http://127.0.0.1:8000/health && curl -fsS --max-time 5 http://127.0.0.1:30080/health` | 0 | ✅ pass - services were inspectable; API /health lacked runtime metadata and TEI health returned no useful JSON | 184ms |
| 2 | `python3 tools/evaluate_legal_model_quick_gate.py --corpus tests/44-FZ-2026-articles.jsonl --output benchmark-results/fd-legal-model-quick-gate-m040-s03.md --baseline-api-url http://127.0.0.1:8000 --baseline-model deepvk/USER-bge-m3 --baseline-runtime-label current-same-host-api-tei --baseline-dimensions 1024 --baseline-cache-namespace m040-s03-baseline-deepvk --candidate-api-url http://127.0.0.1:18001 --candidate-model BAAI/bge-m3 --candidate-runtime-label candidate-bge-m3-separate-endpoint-required --candidate-dimensions 1024 --candidate-cache-namespace m040-s03-candidate-bge-m3 --candidate-api-url http://127.0.0.1:18002 --candidate-model intfloat/multilingual-e5-large --candidate-runtime-label candidate-multilingual-e5-large-separate-endpoint-required --candidate-dimensions 1024 --candidate-cache-namespace m040-s03-candidate-multilingual-e5-large --max-docs 128 --max-title-queries 32 --max-self-queries 32 --batch-size 16 --timeout-seconds 10` | 0 | ✅ pass - canonical fail-closed artifact generated with two deferred candidates | 238ms |
| 3 | `python3 tools/verify_legal_model_quick_gate_artifact.py --artifact benchmark-results/fd-legal-model-quick-gate-m040-s03.md --max-candidates 2` | 0 | ✅ pass - artifact valid | 57ms |
| 4 | `python3 -m py_compile tools/evaluate_legal_model_quick_gate.py tools/verify_legal_model_quick_gate_artifact.py` | 0 | ✅ pass - quick-gate scripts compile | 52ms |
| 5 | `inline artifact content checks for candidate_count=2, raw_text_logged=false, isolated cache namespaces, and defer_candidate outcomes` | 0 | ✅ pass - bounded scope and redaction-sensitive fields confirmed | 40ms |

## Deviations

Live legal retrieval metrics were not run because the baseline `/health` response lacked the contract-required `runtime` block; per the task failure modes, the artifact blocks/defer live comparison instead of continuing with uninspectable runtime metadata.

## Known Issues

The currently running same-host API returns `/health` without runtime metadata, so a future live comparison must first restore or redeploy the runtime block required by `docs/same-host-embedding-service-contract.md`.

## Files Created/Modified

- `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`
- `benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md`
