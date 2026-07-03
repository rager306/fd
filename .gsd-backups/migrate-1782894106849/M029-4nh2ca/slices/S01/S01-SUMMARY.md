---
id: S01
parent: M029-4nh2ca
milestone: M029-4nh2ca
provides:
  - Safer provisioning helper suitable for later hosted workflow proof after documentation/closure.
requires:
  []
affects:
  - S02 documentation and final closure
  - Manual hosted ONNX packaging workflow trust gate
key_files:
  - tools/provision_onnx_artifacts.py
key_decisions:
  - Remote URLs require HTTPS and are private-network blocked by default.
  - Artifact URL redirects are disabled.
  - Archive members are checked for regular-file type and expected/hard size before copying.
patterns_established:
  - Use deterministic monkeypatched probes instead of internet-dependent tests for artifact downloader security behavior.
  - Hash verification protects integrity after materialization; bounded streaming protects runner availability before verification.
observability_surfaces:
  - Actionable provisioning failure messages with query-redacted URL display.
drill_down_paths:
  - .gsd/milestones/M029-4nh2ca/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M029-4nh2ca/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M029-4nh2ca/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T04:32:56.554Z
blocker_discovered: false
---

# S01: Provisioning URL and archive hardening

**S01 remediated M028 medium-risk URL/download/archive provisioning findings.**

## What Happened

S01 hardened the ONNX provisioning helper against M028 MEDIUM findings. It added HTTPS-only remote URL policy, optional allowed-host enforcement, DNS private/loopback/link-local/reserved address blocking, disabled redirects, Content-Length and streaming byte caps, and archive member regular-file/size caps before extraction. Local deterministic probes and provisioning guardrails passed.

## Verification

All S01 checks passed.

## Requirements Advanced

- onnx-provisioning-security — Remediated M028 MEDIUM-1 arbitrary URL/unbounded download risk and MEDIUM-2 unbounded archive member risk.

## Requirements Validated

- m029-s01-guardrails — All S01 security probes and guardrails passed.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

Allowed-host policy is available via helper CLI/env but not yet surfaced as a workflow_dispatch input; M028 LOW path output/path-root findings remain future work.

## Follow-ups

S02 should document M028 MEDIUM remediation and remaining LOW findings, then run final checks and commit.

## Files Created/Modified

- `tools/provision_onnx_artifacts.py` — URL policy, bounded downloads, no-redirect behavior, and archive member size/type caps.
