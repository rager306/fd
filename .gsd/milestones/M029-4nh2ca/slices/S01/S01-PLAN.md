# S01: Provisioning URL and archive hardening

**Goal:** Remediate M028 MEDIUM findings in the provisioning helper with deterministic tests/probes.
**Demo:** After this, provisioning remote URLs and archives have bounded, policy-checked behavior.

## Must-Haves

- Private/localhost URL hosts are blocked by default.
- Optional allowed-host policy exists for HTTPS artifact hosts.
- Downloads enforce max bytes while streaming.
- Archive members must be regular files and within expected/hard size caps before copy.
- Tests/probes pass without internet.

## Proof Level

- This slice proves: Python compile plus local security probes for URL policy, byte caps, and archive caps.

## Integration Closure

Makes manual/hosted ONNX artifact provisioning safer before using it as rollout evidence.

## Verification

- Failure messages should be actionable and avoid printing signed URL query strings.

## Tasks

- [x] **T01: Harden remote artifact URL downloads** `est:medium`
  Add URL policy and bounded streaming to `tools/provision_onnx_artifacts.py`: HTTPS-only by default for remote URLs, block private/localhost/link-local addresses, support explicit allowed hosts, redact query strings in display/errors, and enforce max download bytes while streaming.
  - Files: `tools/provision_onnx_artifacts.py`
  - Verify: Local probes cover blocked localhost/private URL and oversized download.

- [x] **T02: Harden archive member extraction** `est:small`
  Add archive member safeguards for native tokenizer extraction: require regular file, validate declared member size against expected size/hard cap before copying, and enforce copy byte limit.
  - Files: `tools/provision_onnx_artifacts.py`
  - Verify: Local probes cover oversized/non-regular archive member rejection.

- [x] **T03: Verify provisioning hardening** `est:medium`
  Run S01 verification: Python compile, deterministic security probes, provisioning dry-run and missing-source behavior, verifier allow-missing, GitNexus scope, binary hygiene.
  - Verify: All S01 checks pass.

## Files Likely Touched

- tools/provision_onnx_artifacts.py
