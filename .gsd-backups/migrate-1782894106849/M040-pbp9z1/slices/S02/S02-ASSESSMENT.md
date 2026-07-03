---
sliceId: S02
uatType: artifact-driven
verdict: PASS
date: 2026-05-22T08:03:13Z
---

# UAT Result — S02

## Checks

| Check | Mode | Result | Notes |
|-------|------|--------|-------|
| Run the S02 artifact verifier with benchmark, preflight, legal, and audit inputs. | artifact | PASS | `gsd_exec` run `9160b76f-f79e-4a40-b427-f6965e1f95df` executed `python3 tools/verify_m040_s02_artifacts.py --benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt --preflight benchmark-results/fd-m040-s02-onnx-docker-preflight.txt --legal benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt --audit benchmark-results/fd-m040-s02-proof-audit.txt`; output: `M040 S02 artifact verification: PASS`. |
| Confirm verifier output is `M040 S02 artifact verification: PASS`. | artifact | PASS | The verifier output in `9160b76f-f79e-4a40-b427-f6965e1f95df` exactly included `M040 S02 artifact verification: PASS` and exited 0. |
| Confirm the preflight artifact records ONNX backend, model `deepvk/USER-bge-m3`, 1024 dimensions, and the S02 cache namespace. | artifact | PASS | Semantic checks found the preflight artifact non-empty (`43158` bytes) and confirmed ONNX backend evidence, model `deepvk/USER-bge-m3`, `1024`, and cache namespace `m040-s02-onnx-restart`. |
| Confirm the benchmark artifact contains measured restart/cache sections rather than skipped restart proof sections. | artifact | PASS | Semantic checks found the benchmark artifact non-empty (`9266` bytes), restart terms present, cache terms plus namespace present, and no `restart proof skipped` / `skipped restart proof` / restart-skip wording. |
| Confirm the legal guard artifact ends in PASS for the packaged ONNX proof run. | artifact | PASS | Semantic checks found the legal artifact non-empty (`11200` bytes) and PASS evidence at legal closeout. |
| Confirm the audit artifact reports PASS, no obvious secret/raw probe leakage, API proof container removal, and port 18000 cleanup. | artifact | PASS | Semantic checks found the audit artifact non-empty (`1307` bytes), audit verdict PASS, `leaks=[]` across prohibited private-key/bearer/secret/signed-URL patterns, API container absent/removal evidence, Redis preservation evidence, and port `18000` cleanup evidence. |
| Confirm Docker state has no running/stale `fd-m040-s02-onnx-api` container and that only the isolated Redis proof container may remain alive. | runtime | PASS | Live Docker query in `9160b76f-f79e-4a40-b427-f6965e1f95df` reported `api container query: <none>` and `redis container query: fd-m040-s02-redis Up 19 minutes 127.0.0.1:16379->6379/tcp`; this matches the expected cleanup/preservation state. |
| Confirm `127.0.0.1:18000` does not accept connections after cleanup. | runtime | PASS | Live socket check in `9160b76f-f79e-4a40-b427-f6965e1f95df` reported `PASS: 127.0.0.1:18000 does not accept connections (ConnectionRefusedError: [Errno 111] Connection refused)`. |

## Overall Verdict

PASS — all S02 artifact-driven and live cleanup checks passed, including verifier output, artifact semantics, Docker container cleanup, and port 18000 closure.

## Notes

- Primary evidence: `.gsd/exec/9160b76f-f79e-4a40-b427-f6965e1f95df.stdout` and empty stderr for the passing final verification run.
- An earlier exploratory `gsd_exec` run (`397eff94-d276-4215-823e-d107681cd2cf`) used an overly broad leak heuristic that flagged the audit phrase `bearer tokens` in the prohibited-pattern checklist; the final verification refined the check to match actual bearer token/private key/secret/signed URL material and found `leaks=[]` while the project verifier also passed.
- No human follow-up is required for this artifact-driven UAT.
