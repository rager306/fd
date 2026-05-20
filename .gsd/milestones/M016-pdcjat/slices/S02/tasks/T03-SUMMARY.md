---
id: T03
parent: S02
milestone: M016-pdcjat
key_files:
  - benchmark-results/fd-legal-divergence-profile-m016-s01.txt
  - benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt
key_decisions:
  - Root-cause verdict: confirmed/narrowed to ONNX max sequence length 128 truncation for the severe M015 legal divergence outliers.
  - Remediation path: S03 should prioritize a 512-token ONNX export/runtime gate, plus chunking for cases beyond 512; do not pursue tokenizer replacement or model alternatives first.
duration: 
verification_result: passed
completed_at: 2026-05-20T05:41:30.949Z
blocker_discovered: false
---

# T03: Recorded the root-cause verdict: ONNX 128-token truncation caused the severe legal divergence outliers.

**Recorded the root-cause verdict: ONNX 128-token truncation caused the severe legal divergence outliers.**

## What Happened

Interpreted S01/S02 evidence. S01 showed every M015 worst case exceeds 128 tokens. S02 showed local ONNX at 512 nearly restores TEI cosine for cases that fit within 512, while the two cases still truncated at 512 remain below the strict 0.999 threshold. This rejects tokenizer mismatch as the primary current explanation and confirms the current ONNX 128-token runtime/export path as the main cause of severe outliers.

## Verification

Root-cause verdict is directly supported by S01 and S02 artifacts.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `S01 profile artifact review` | 0 | ✅ pass — 17/17 worst cases truncated at 128; 2/17 at 512 | 0ms |
| 2 | `S02 sequence diagnostic artifact review` | 0 | ✅ pass — 128 mean cosine 0.9204953; 512 mean cosine 0.99885631 | 0ms |

## Deviations

None.

## Known Issues

Two worst cases still truncate at 512 and only reach ~0.990 cosine against TEI, so 512 alone may not be sufficient for all legal corpus cases. Full corpus rerun is required after remediation.

## Files Created/Modified

- `benchmark-results/fd-legal-divergence-profile-m016-s01.txt`
- `benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt`
