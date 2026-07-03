---
id: T02
parent: S02
milestone: M045-d0e5xq
key_files:
  - documents/tei-startup-mitigation-m045.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T11:37:43.751Z
blocker_discovered: false
---

# T02: Selected `HF_HUB_OFFLINE=1` as the S03 mitigation candidate.

**Selected `HF_HUB_OFFLINE=1` as the S03 mitigation candidate.**

## What Happened

Compared candidates from S01 and cache inventory. Selected `HF_HUB_OFFLINE=1` because the required cached safetensors/tokenizer files are present and Hugging Face Hub docs say offline mode avoids HTTP calls and uses cached files only. Rejected adding ONNX artifacts because it conflicts with M042 TEI-only posture; kept local model path and Candle-only image as fallback options.

## Verification

`documents/tei-startup-mitigation-m045.md` contains selected candidate, rationale, rejected options, risks, rollback, and S03 success criteria.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `review mitigation document for selected candidate and rejected alternatives` | 0 | ✅ pass: mitigation candidate selected | 120000ms |

## Deviations

None.

## Known Issues

Offline mode may fail if TEI needs an uncached metadata file; S03 proof must capture failure and rollback.

## Files Created/Modified

- `documents/tei-startup-mitigation-m045.md`
