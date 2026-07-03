---
id: T03
parent: S01
milestone: M045-d0e5xq
key_files:
  - documents/tei-startup-recon-m045.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T11:30:47.899Z
blocker_discovered: false
---

# T03: Identified safe TEI startup mitigation candidates from source/docs without destructive commands.

**Identified safe TEI startup mitigation candidates from source/docs without destructive commands.**

## What Happened

Inspected upstream TEI router/backend source and Hugging Face Hub environment docs. Found CLI args such as `--dtype`, `--pooling`, and `--dense-path`, but no documented force-Candle or disable-ORT flag. Backend source shows ORT/ONNX probing can occur before Candle fallback. Identified `HF_HUB_OFFLINE=1` with a complete cache as the strongest candidate, local model path as a secondary candidate, and ONNX artifact addition as rejected for current fd scope.

## Verification

`documents/tei-startup-recon-m045.md` lists inspected sources, candidate mitigations, expected effects, risks, and S02 recommendation.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `fetch/read upstream TEI source and Hugging Face Hub environment variable docs; no container restart or docker run` | 0 | ✅ pass: mitigation candidates documented | 120000ms |

## Deviations

Did not run TEI binary `--help` inside a new container because prior `docker run --rm ... --help` was blocked/destructive. Used upstream source/docs instead.

## Known Issues

Need controlled proof before applying `HF_HUB_OFFLINE=1` permanently because incomplete cache could break startup.

## Files Created/Modified

- `documents/tei-startup-recon-m045.md`
