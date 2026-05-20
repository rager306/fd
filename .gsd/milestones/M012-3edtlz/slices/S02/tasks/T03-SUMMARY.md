---
id: T03
parent: S02
milestone: M012-3edtlz
key_files:
  - .gsd/milestones/M012-3edtlz/slices/S02/S02-RESEARCH.md
  - benchmark-results/fd-tokenizer-go-current-m012-s02.txt
key_decisions:
  - Do not patch `sugarme/tokenizer` blindly; current evidence points to implementation-level divergence for XLM-R/Unigram/Metaspace behavior.
  - S03 should first test `github.com/daulet/tokenizers` or related Hugging Face Rust tokenizers bindings in isolation.
  - If Rust-backed Go bindings pass parity, then consider integration; otherwise record packaging/correctness blocker and evaluate sidecar/defer.
duration: 
verification_result: passed
completed_at: 2026-05-20T02:09:59.230Z
blocker_discovered: false
---

# T03: Evaluated tokenizer alternatives and recommended testing Hugging Face Rust tokenizers Go bindings before any runtime integration.

**Evaluated tokenizer alternatives and recommended testing Hugging Face Rust tokenizers Go bindings before any runtime integration.**

## What Happened

Researched the tokenizer mismatch and candidate paths. The current tokenizer JSON uses Unigram plus Sequence normalization, Metaspace pre-tokenization, TemplateProcessing, and Metaspace decoding. A special-token toggle does not fix the mismatch. The strongest candidate is a Go binding to Hugging Face's Rust `tokenizers`, such as `github.com/daulet/tokenizers` or the related `crowdsecurity/go-tokenizers`, because it loads `tokenizer.json` and should mirror HF behavior if linked correctly. The research recommends S03 test that binding in isolation before modifying runtime code.

## Verification

Research artifact exists and recommends the S03 path with evidence from the current mismatch artifact and external library docs.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `search_and_read Go Hugging Face tokenizers tokenizer.json XLMRoberta SentencePiece library` | 0 | ✅ pass — found HF tokenizer docs and Go binding candidates | 0ms |
| 2 | `fetch_page README for daulet/tokenizers and crowdsecurity/go-tokenizers` | 0 | ✅ pass — both expose Go bindings for HF Tokenizers and tokenizer.json loading, with native library requirements | 0ms |
| 3 | `python3 inspect tokenizer.json model/pretokenizer/postprocessor` | 0 | ✅ pass — tokenizer JSON uses Unigram, Sequence normalizer, Metaspace pre-tokenizer, TemplateProcessing | 0ms |
| 4 | `gsd_summary_save S02 RESEARCH` | 0 | ✅ pass — wrote .gsd/milestones/M012-3edtlz/slices/S02/S02-RESEARCH.md | 0ms |

## Deviations

Did not add candidate dependencies in S02; kept this slice to comparison evidence and research so S03 can make the dependency decision with a focused feasibility test.

## Known Issues

Native HF tokenizers bindings likely require CGO and `libtokenizers` static/native library packaging. That may be acceptable for correctness proof but needs explicit operational planning before production use.

## Files Created/Modified

- `.gsd/milestones/M012-3edtlz/slices/S02/S02-RESEARCH.md`
- `benchmark-results/fd-tokenizer-go-current-m012-s02.txt`
