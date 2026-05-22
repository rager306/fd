# S04: Runtime recommendation and operating contract — UAT

**Milestone:** M040-pbp9z1
**Written:** 2026-05-22T08:33:53.094Z

## UAT Type
Documentation and semantic validation UAT for the final same-host embedding runtime recommendation.

## Preconditions
- Work from `/root/fd`.
- S01, S02, and S03 artifacts exist, including the same-host contract, S02 packaged ONNX evidence files, and S03 candidate-gate artifact.
- No Docker service or live embedding runtime is required for this UAT.

## Steps
1. Compile the verifier: `python3 -m py_compile tools/verify_m040_s04_recommendation.py`.
2. Run verifier self-tests: `python3 tools/verify_m040_s04_recommendation.py --self-test`.
3. Validate the final recommendation artifact against evidence inputs:
   ```bash
   python3 tools/verify_m040_s04_recommendation.py \
     --artifact benchmark-results/fd-runtime-recommendation-m040-s04.md \
     --s02-benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt \
     --s02-preflight benchmark-results/fd-m040-s02-onnx-docker-preflight.txt \
     --s02-legal benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt \
     --s02-audit benchmark-results/fd-m040-s02-proof-audit.txt \
     --s03-gate benchmark-results/fd-legal-model-quick-gate-m040-s03.md
   ```
4. Confirm discoverability from the same-host contract: `rg -n fd-runtime-recommendation-m040-s04 docs/same-host-embedding-service-contract.md`.
5. Confirm change-scope audit: `npx gitnexus detect-changes --repo fd`.

## Expected Outcomes
- Steps 1-3 exit 0 and the full verifier prints `PASS benchmark-results/fd-runtime-recommendation-m040-s04.md`.
- Step 4 finds the discoverability link in `docs/same-host-embedding-service-contract.md`.
- Step 5 reports no unexpected change scope for closeout.
- The recommendation remains: TEI is current/default; packaged ONNX is preferred only after explicit operator switch and required same-host preflight/smoke/cache isolation checks; alternative models remain deferred/fail-closed.

## Edge Cases
- If cache namespace isolation language is removed, the verifier must fail.
- If the artifact accepts an alternative model replacement instead of `defer_candidate`, the verifier must fail.
- If hosted or remote CI is described as the readiness gate, the verifier must fail.
- If `/health` alone is treated as live inference readiness without smoke `POST /v1/embeddings`, the verifier must fail.
- If raw legal/probe text, secrets, bearer tokens, signed URLs, or similar prohibited patterns appear in scanned artifacts, the verifier must fail.

## Not Proven By This UAT
- It does not rerun Docker, TEI, ONNX, Redis, or legal retrieval benchmarks; those were proven in upstream slices.
- It does not prove ONNX runtime-library integrity when `ONNX_RUNTIME_SHA256` is unset; the final artifact documents this S02 caveat.
- It does not validate hosted/remote deployment behavior; hosted CI is intentionally outside the same-host readiness gate.
