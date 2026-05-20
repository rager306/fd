---
id: T02
parent: S01
milestone: M016-pdcjat
key_files:
  - tools/profile_legal_divergence.py
key_decisions:
  - The profiler uses local Hugging Face tokenizer diagnostics only; it does not claim to know TEI internal truncation behavior.
  - Profiler artifacts include IDs/hashes/token counts/truncation flags only, not raw legal text.
duration: 
verification_result: passed
completed_at: 2026-05-20T05:34:28.854Z
blocker_discovered: false
---

# T02: Implemented the sanitized worst-case legal divergence profiler.

**Implemented the sanitized worst-case legal divergence profiler.**

## What Happened

Implemented `tools/profile_legal_divergence.py`. The tool parses M015 metrics, extracts worst cross-backend cosine rows, resolves them against the 44-ФЗ corpus using stable fallback IDs, verifies text hashes, computes local HF tokenizer token counts and truncation diagnostics for configurable sequence lengths, and writes sanitized markdown output.

## Verification

`python3 -m py_compile tools/profile_legal_divergence.py` passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/profile_legal_divergence.py` | 0 | ✅ pass | 0ms |

## Deviations

None.

## Known Issues

Runtime diagnostics in S02 are still needed to compare actual TEI vs ONNX behavior. S01 only profiles tokenization/truncation risk for worst IDs.

## Files Created/Modified

- `tools/profile_legal_divergence.py`
