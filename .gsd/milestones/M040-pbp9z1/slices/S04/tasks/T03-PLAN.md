---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T03: Run closeout verification and change-scope audit

Why: S04 is the final assembly slice, so completion must rest on fresh evidence after the last file change, not on earlier verifier runs. Expected executor skills: verify-before-complete. Do: rerun Python compile, verifier self-test, full final-artifact validation, and a discoverability check for the link from the service contract; run `gitnexus_detect_changes()` before closure to confirm only the recommendation artifact, verifier, and contract link changed. If the verifier finds missing semantics or redaction issues, return to T01/T02 rather than weakening checks. Failure Modes (Q5): if GitNexus reports stale index, run the requested analyze command before rerunning change detection; if verification fails, do not mark the task or slice complete; if unexpected files changed, inspect and either revert or document why they are required. Load Profile (Q6): small local checks only; no runtime services, Docker, or network. Negative Tests (Q7): rely on the T01 self-test negative fixtures and the full verifier's redaction/evidence gate; no additional runtime negative test is required because this slice does not modify runtime code. Done when all verification commands exit 0 and GitNexus change detection matches the planned scope.

## Inputs

- `tools/verify_m040_s04_recommendation.py`
- `benchmark-results/fd-runtime-recommendation-m040-s04.md`
- `docs/same-host-embedding-service-contract.md`
- `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt`
- `benchmark-results/fd-m040-s02-onnx-docker-preflight.txt`
- `benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt`
- `benchmark-results/fd-m040-s02-proof-audit.txt`
- `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`

## Expected Output

- `tools/verify_m040_s04_recommendation.py`
- `benchmark-results/fd-runtime-recommendation-m040-s04.md`
- `docs/same-host-embedding-service-contract.md`

## Verification

python3 -m py_compile tools/verify_m040_s04_recommendation.py
python3 tools/verify_m040_s04_recommendation.py --self-test
python3 tools/verify_m040_s04_recommendation.py --artifact benchmark-results/fd-runtime-recommendation-m040-s04.md --s02-benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt --s02-preflight benchmark-results/fd-m040-s02-onnx-docker-preflight.txt --s02-legal benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt --s02-audit benchmark-results/fd-m040-s02-proof-audit.txt --s03-gate benchmark-results/fd-legal-model-quick-gate-m040-s03.md
rg -n fd-runtime-recommendation-m040-s04 docs/same-host-embedding-service-contract.md

## Observability Impact

Confirms the verifier remains the authoritative closeout signal and that the final artifact is discoverable from the service contract.
