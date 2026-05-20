---
id: T01
parent: S02
milestone: M017-j10hmp
key_files:
  - benchmark-results/fd-onnx-512-outcome-m017-s02.txt
key_decisions:
  - 512-token ONNX is necessary but not sufficient for strict legal equivalence.
  - Recommended next path is hybrid: 512-token ONNX baseline plus deterministic chunking for fragments above 512 tokens.
duration: 
verification_result: passed
completed_at: 2026-05-20T07:29:53.136Z
blocker_discovered: false
---

# T01: Wrote the 512-token ONNX outcome assessment.

**Wrote the 512-token ONNX outcome assessment.**

## What Happened

Created the M017 S02 outcome assessment. It compares M015 128-token failure, M016 Python 512 diagnostic, and M017 tagged Go 512 legal gate. The assessment concludes that 512-token ONNX dramatically improves ranking parity but still fails strict cosine equivalence due long fragments above 512 tokens.

## Verification

Outcome artifact hygiene check passed with required metrics and no raw legal text leaks.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python outcome artifact hygiene check` | 0 | ✅ pass — m017_s02_outcome_hygiene=pass; raw_legal_text_leaks=0 | 0ms |

## Deviations

None.

## Known Issues

Outcome artifact is a decision artifact; no chunking implementation was added in M017.

## Files Created/Modified

- `benchmark-results/fd-onnx-512-outcome-m017-s02.txt`
