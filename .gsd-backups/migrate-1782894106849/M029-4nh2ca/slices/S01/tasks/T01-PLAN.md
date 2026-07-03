---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Harden remote artifact URL downloads

Add URL policy and bounded streaming to `tools/provision_onnx_artifacts.py`: HTTPS-only by default for remote URLs, block private/localhost/link-local addresses, support explicit allowed hosts, redact query strings in display/errors, and enforce max download bytes while streaming.

## Inputs

- `M028 security review findings`
- `tools/provision_onnx_artifacts.py`

## Expected Output

- `Updated provisioning helper URL policy`

## Verification

Local probes cover blocked localhost/private URL and oversized download.

## Observability Impact

Safer actionable URL failure messages.
