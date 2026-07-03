---
id: T02
parent: S04
milestone: M040-pbp9z1
key_files:
  - benchmark-results/fd-runtime-recommendation-m040-s04.md
  - docs/same-host-embedding-service-contract.md
key_decisions:
  - The recommendation treats S02 values as evidence markers rather than future thresholds.
  - TEI remains current/default until explicit operator opt-in to ONNX.
  - Alternative model replacement remains deferred/fail-closed based on S03.
duration: 
verification_result: passed
completed_at: 2026-05-22T08:31:06.660Z
blocker_discovered: false
---

# T02: Authored the M040 S04 TEI-vs-ONNX same-host runtime recommendation and linked it from the same-host contract.

**Authored the M040 S04 TEI-vs-ONNX same-host runtime recommendation and linked it from the same-host contract.**

## What Happened

Created `benchmark-results/fd-runtime-recommendation-m040-s04.md` as the final machine-checkable operating contract and recommendation artifact. It states the required stance: keep `deepvk/USER-bge-m3`, prefer packaged ONNX only for explicit same-host opt-in under the S01 contract, keep TEI as the current/default posture until an operator switches, require artifact/tokenizer/runtime preflight and smoke embedding readiness, isolate Redis cache namespace, prohibit request-level fallback, and keep alternative model replacement deferred/fail-closed. The artifact summarizes S02 benchmark/preflight/legal/audit evidence as sanitized evidence values and records the `runtime_library_verified=false` caveat because `ONNX_RUNTIME_SHA256` was not set. Added a concise related-document link from `docs/same-host-embedding-service-contract.md` to the final recommendation without duplicating the contract.

## Verification

Ran the required semantic verifier against the final artifact and all S02/S03 inputs. The first run failed because the Redaction section used wording outside the verifier’s accepted marker set; I changed the sentence to explicitly say raw/private material is `excluded` and source payloads are `not logged`, then reran successfully. Also verified the same-host contract contains the discoverability link to the new artifact.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 tools/verify_m040_s04_recommendation.py --artifact benchmark-results/fd-runtime-recommendation-m040-s04.md --s02-benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt --s02-preflight benchmark-results/fd-m040-s02-onnx-docker-preflight.txt --s02-legal benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt --s02-audit benchmark-results/fd-m040-s02-proof-audit.txt --s03-gate benchmark-results/fd-legal-model-quick-gate-m040-s03.md` | 1 | ❌ fail (redaction wording marker missing; fixed before final pass) | 69ms |
| 2 | `python3 tools/verify_m040_s04_recommendation.py --artifact benchmark-results/fd-runtime-recommendation-m040-s04.md --s02-benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt --s02-preflight benchmark-results/fd-m040-s02-onnx-docker-preflight.txt --s02-legal benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt --s02-audit benchmark-results/fd-m040-s02-proof-audit.txt --s03-gate benchmark-results/fd-legal-model-quick-gate-m040-s03.md` | 0 | ✅ pass | 60ms |
| 3 | `python3 - <<'PY'
from pathlib import Path
contract = Path('docs/same-host-embedding-service-contract.md').read_text(encoding='utf-8')
needle = '../benchmark-results/fd-runtime-recommendation-m040-s04.md'
if needle not in contract:
    raise SystemExit(f'missing link: {needle}')
print('same-host contract link verification: PASS')
PY` | 0 | ✅ pass | 80ms |

## Deviations

None. The verifier-required section name `Decision Inputs` was used to satisfy the existing semantic verifier while still covering the task plan’s evidence envelope content.

## Known Issues

None.

## Files Created/Modified

- `benchmark-results/fd-runtime-recommendation-m040-s04.md`
- `docs/same-host-embedding-service-contract.md`
