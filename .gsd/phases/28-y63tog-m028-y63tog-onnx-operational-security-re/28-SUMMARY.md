---
id: M028-y63tog
title: "ONNX operational security review"
status: complete
completed_at: 2026-05-21T04:19:09.848Z
key_decisions:
  - D026: M028 MEDIUM findings block treating hosted/manual ONNX packaging workflow as trusted rollout evidence until URL/download/archive protections are remediated.
key_files:
  - .gsd/milestones/M028-y63tog/slices/S01/S01-RESEARCH.md
  - .gsd/DECISIONS.md
lessons_learned:
  - Artifact hashes protect integrity after materialization, but they do not protect CI runners from arbitrary URL fetches or oversized downloads before verification.
  - Security review should remain separate from remediation to preserve a clean audit trail.
  - Health metadata is currently safe from path disclosure; startup/log errors are the remaining disclosure surface.
---

# M028-y63tog: ONNX operational security review

**M028 produced a read-only security review for ONNX operational surfaces and identified remediation blockers before trusted hosted packaging.**

## What Happened

M028 completed a read-only security review of ONNX operational artifact and startup surfaces. The review mapped Go startup env/path handling, manifest path resolution, health metadata, manual workflow inputs, provisioning helper downloads and archive extraction, verifier path handling, Docker packaging copies, and operations docs. It found two Medium issues around arbitrary outbound manual workflow URLs/unbounded downloads and unbounded archive member copying before verification, plus two Low issues around local path disclosure and repo-external manifest paths. It found no shell injection through workflow inputs, no health endpoint path disclosure, no artifact integrity bypass after download, no tar path traversal in the current fixed-destination extraction path, and no production-default bypass. No runtime code was modified.

## Success Criteria Results

- Code-cited review: PASS.
- Findings/non-findings: PASS.
- Read-only scope: PASS.
- Decision recorded: PASS.

## Definition of Done Results

- Security review maps actual attack surfaces: met.
- Findings include file:line, exploit, severity, reachability, impact, remediation: met.
- Non-findings and out-of-scope recorded: met.
- No code remediation performed: met.
- Report and decision persisted: met.
- Production/default promotion: not performed.

## Requirement Outcomes

- ONNX security review gate: validated.
- Hosted packaging workflow trust: blocked until MEDIUM findings are remediated.
- ONNX production/default switch: still blocked.

## Deviations

None.

## Follow-ups

Recommended next milestone: remediate MEDIUM-1 and MEDIUM-2 by adding approved URL source policy/private-address blocking/byte caps and archive member size/type caps before relying on hosted ONNX packaging workflow evidence.
