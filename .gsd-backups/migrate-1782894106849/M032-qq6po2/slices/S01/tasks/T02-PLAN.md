---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Verify failure modes

Run negative probes by feeding temporary tampered copies of metadata/manifest to prove checksum/toolchain/source mismatch failures are explicit and sanitized.

## Inputs

- `tools/verify_onnx_export_contract.py`

## Expected Output

- `Task summary evidence`

## Verification

Tampered artifact sha, model revision, and package version probes fail with expected labels and safe paths.

## Observability Impact

Proves diagnostic failure modes are actionable.
