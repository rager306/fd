---
id: T01
parent: S01
milestone: M029-4nh2ca
key_files:
  - tools/provision_onnx_artifacts.py
key_decisions:
  - Remote artifact URLs are HTTPS-only by default, private-network hosts are blocked by default, and an explicit allowed-host policy is available through CLI/env.
  - Redirects are disabled for artifact downloads to avoid redirect-to-private bypass.
duration: 
verification_result: passed
completed_at: 2026-05-21T04:31:30.145Z
blocker_discovered: false
---

# T01: Hardened remote ONNX artifact URL handling with policy and byte limits.

**Hardened remote ONNX artifact URL handling with policy and byte limits.**

## What Happened

Hardened remote artifact source handling in the provisioning helper. Remote URLs now require HTTPS, can be constrained by explicit allowed hosts, block private/loopback/link-local/reserved/multicast/unspecified DNS resolutions by default, disable redirects, sanitize displayed URLs by stripping query strings, and enforce `Content-Length` plus streaming byte caps. Local deterministic probes verified localhost/http blocking, disallowed host blocking, and oversized remote Content-Length failure.

## Verification

Python compile and local security probes passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py` | 0 | ✅ pass | 0ms |
| 2 | `gsd_exec M029 provisioning hardening local security probes` | 0 | ✅ pass — block_http_localhost, block_disallowed_host, remote_content_length_cap | 118ms |

## Deviations

Remote download cap was verified through a deterministic monkeypatched opener to avoid real network access.

## Known Issues

The manual workflow has not yet been updated to expose an explicit allowed-host input; default helper behavior is safer but workflow operator host allowlists remain a future workflow UX improvement if needed.

## Files Created/Modified

- `tools/provision_onnx_artifacts.py`
