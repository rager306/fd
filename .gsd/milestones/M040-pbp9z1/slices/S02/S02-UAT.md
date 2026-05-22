# S02: Docker restart and Redis L2 proof — UAT

**Milestone:** M040-pbp9z1
**Written:** 2026-05-22T08:01:53.580Z

# S02 UAT: Docker restart and Redis L2 proof

## UAT Type
Operational proof / artifact verification for same-host packaged ONNX runtime readiness.

## Preconditions
- Work from `/root/fd`.
- Docker is available on the host.
- The S02 proof artifacts exist under `benchmark-results/`:
  - `fd-m040-s02-onnx-docker-preflight.txt`
  - `fd-benchmark-m040-s02-onnx-docker-restart.txt`
  - `fd-legal-retrieval-m040-s02-onnx-docker-restart.txt`
  - `fd-m040-s02-proof-audit.txt`
- The verifier exists at `tools/verify_m040_s02_artifacts.py`.

## Numbered Steps
1. Run the S02 artifact verifier with benchmark, preflight, legal, and audit inputs.
2. Confirm verifier output is `M040 S02 artifact verification: PASS`.
3. Confirm the preflight artifact records ONNX backend, model `deepvk/USER-bge-m3`, 1024 dimensions, and the S02 cache namespace.
4. Confirm the benchmark artifact contains measured restart/cache sections rather than skipped restart proof sections.
5. Confirm the legal guard artifact ends in PASS for the packaged ONNX proof run.
6. Confirm the audit artifact reports PASS, no obvious secret/raw probe leakage, API proof container removal, and port 18000 cleanup.
7. Confirm Docker state has no running/stale `fd-m040-s02-onnx-api` container and that only the isolated Redis proof container may remain alive.
8. Confirm `127.0.0.1:18000` does not accept connections after cleanup.

## Expected Outcomes
- Packaged ONNX proof artifacts validate successfully.
- Runtime metadata is observable and safe to consume downstream.
- Redis L2 behavior is proven across an API-only restart, not hidden by a Redis restart or broad cache contamination.
- Legal no-regression guard passes for the proof run.
- Cleanup evidence is explicit and leaves the proof port clear.

## Edge Cases
- If Docker or ONNX prerequisites are unavailable, the proof runner must produce exact blocker evidence rather than silently skipping restart/cache proof.
- If benchmark restart sections are skipped, S02 UAT fails.
- If the legal artifact is missing or non-PASS, S02 UAT fails.
- If raw probe/legal text, obvious secrets, bearer tokens, private key blocks, or signed URL query material appear in artifacts, S02 UAT fails.
- If the API proof container remains present or port 18000 is still bound, S02 UAT fails.

## Not Proven By This UAT
- Production deployment sizing, multi-host networking, or external Redis exposure safety.
- Long-duration soak behavior under concurrent neighbor-service traffic.
- A final TEI-vs-ONNX recommendation; S04 combines this evidence with S03 model quality gating.
- Alternative model quality beyond the current `deepvk/USER-bge-m3` packaged ONNX proof.
