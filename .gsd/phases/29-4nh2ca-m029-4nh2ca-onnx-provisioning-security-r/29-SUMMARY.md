---
id: M029-4nh2ca
title: "ONNX provisioning security remediation"
status: complete
completed_at: 2026-05-21T04:38:25.962Z
key_decisions:
  - D027: M029 remediates M028 MEDIUM provisioning risks but does not complete ONNX rollout readiness.
key_files:
  - tools/provision_onnx_artifacts.py
  - docs/onnx-artifacts/PROVISIONING.md
  - benchmark-results/fd-onnx-provisioning-security-remediation-m029-s02.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - Checksum verification protects artifact integrity after download, but URL policy and byte caps are required to protect CI runner availability before verification.
  - Disabling redirects is a simple way to avoid redirect-to-private bypass in artifact provisioning.
  - Archive extraction should validate metadata before copying bytes, even when extraction writes to a fixed destination.
---

# M029-4nh2ca: ONNX provisioning security remediation

**M029 hardened ONNX artifact provisioning against arbitrary URL/unbounded download and oversized archive risks.**

## What Happened

M029 remediated the two medium security findings from M028. The provisioning helper now requires HTTPS for remote artifact URLs, disables redirects, blocks private/loopback/link-local/reserved/multicast/unspecified resolved addresses by default, supports optional allowed-host policy, and enforces Content-Length plus streaming byte caps. Native tokenizer archive extraction now requires a regular file member and checks member size before copy using manifest expected size or a hard cap. Documentation and outcome artifacts record the new policy, evidence, and remaining LOW findings. TEI remains default and ONNX remains opt-in experimental.

## Success Criteria Results

- MEDIUM-1 remediation: PASS.
- MEDIUM-2 remediation: PASS.
- Guardrails: PASS.
- Production safety: PASS (no switch).

## Definition of Done Results

- URL safety policy: met.
- Private/localhost blocking: met.
- Download byte limits: met.
- Archive member type/size caps: met.
- Tests/probes/guardrails: met.
- Docs/outcome/decision: met.
- Production/default switch: not performed.

## Requirement Outcomes

- ONNX provisioning security: advanced and validated for M028 MEDIUM findings.
- Hosted workflow trust: improved but still requires immutable sources and actual run.
- M028 LOW findings: remain open.
- Production rollout: still blocked.

## Deviations

None.

## Follow-ups

Next recommended gates: remediate M028 LOW findings for path output/path roots, or select immutable artifact sources and prepare hosted workflow proof. Remote push/workflow dispatch still requires explicit user approval.
