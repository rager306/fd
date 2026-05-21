---
id: T02
parent: S01
milestone: M029-4nh2ca
key_files:
  - tools/provision_onnx_artifacts.py
key_decisions:
  - Native tokenizer archive members must be regular files and must fit the manifest expected size or hard cap before bytes are copied.
duration: 
verification_result: passed
completed_at: 2026-05-21T04:31:44.485Z
blocker_discovered: false
---

# T02: Hardened native tokenizer archive extraction with type and size caps.

**Hardened native tokenizer archive extraction with type and size caps.**

## What Happened

Hardened archive extraction for native tokenizer provisioning. The helper now checks the selected archive member is a regular file, rejects oversized members before copying, and streams extraction through the same bounded copy helper. Local probes covered oversized member rejection and non-regular member rejection.

## Verification

Python compile and local security probes passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py` | 0 | ✅ pass | 0ms |
| 2 | `gsd_exec M029 provisioning hardening local security probes` | 0 | ✅ pass — archive_member_expected_size_cap, archive_member_regular_file_required | 118ms |

## Deviations

None.

## Known Issues

Zip archives are not supported; existing tar behavior remains the path for native tokenizer archives.

## Files Created/Modified

- `tools/provision_onnx_artifacts.py`
