# S02: Select safe TEI startup mitigation

**Goal:** Choose the smallest safe TEI startup mitigation that does not change fd runtime semantics or embedding identity.
**Demo:** A selected mitigation is encoded as compose/docs changes or a documented no change decision with rationale.

## Must-Haves

- Candidate options compared with tradeoffs.
- If a safe config change exists, compose/docs are updated.
- If no safe change exists, blocker/limitation is documented with operator guidance.
- fd remains TEI-only with runtime identity unchanged.

## Proof Level

- This slice proves: code/config review plus smoke verification

## Integration Closure

Provides candidate config for controlled proof in S03.

## Verification

- Documents startup expectation and failure modes for operators.

## Tasks

- [x] **T01: Inventoried TEI `/data` cache and confirmed USER-bge-m3 Candle files are present.** `est:45m`
  Read-only inspect `/data` inside the running TEI container and host mounted cache paths if discoverable. Record required Candle/tokenizer/config/safetensors files, ONNX absence, file sizes, and whether `HF_HUB_OFFLINE=1` is likely safe to attempt. Do not print secrets or file contents.
  - Files: `documents/tei-startup-mitigation-m045.md`
  - Verify: Artifact lists cache files needed for offline Candle startup and states whether the cache appears complete.

- [x] **T02: Selected `HF_HUB_OFFLINE=1` as the S03 mitigation candidate.** `est:30m`
  Compare `HF_HUB_OFFLINE=1`, local model path, no-change documentation, and rejected ONNX artifact option. Select the S03 candidate or record a blocker. No runtime restart.
  - Files: `documents/tei-startup-mitigation-m045.md`
  - Verify: Artifact has a clear selected candidate, rejected options, risk, rollback, and success criteria for S03.

- [x] **T03: Prepared compose/docs candidate for TEI offline cache startup without restarting runtime.** `est:45m`
  If T02 selects a config change, update compose/docs so future TEI starts include the candidate environment. Do not restart the running container. If no change is safe, write an explicit no-change limitation instead.
  - Files: `docker-compose.yaml`, `docs/same-host-embedding-service-contract.md`, `documents/tei-startup-mitigation-m045.md`
  - Verify: `docker compose config tei` reflects candidate config if changed; current running container remains unchanged and healthy.

- [x] **T04: Verified non-destructive S02 state with fd and TEI smoke checks.** `est:30m`
  Run fd health/ready/embedding smoke, direct TEI smoke, and compose config check after any file changes. Run Go gates if code/config/docs changes warrant them. Confirm no restart/recreate occurred.
  - Files: `documents/tei-startup-mitigation-m045.md`
  - Verify: Smoke checks pass; compose candidate is visible; current runtime identity remains TEI 1024; no restart was performed.

## Files Likely Touched

- documents/tei-startup-mitigation-m045.md
- docker-compose.yaml
- docs/same-host-embedding-service-contract.md
