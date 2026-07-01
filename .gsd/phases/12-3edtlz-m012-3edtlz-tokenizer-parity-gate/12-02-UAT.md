# S02: Go tokenizer candidate comparison — UAT

**Milestone:** M012-3edtlz
**Written:** 2026-05-20T02:12:24.315Z

# S02 UAT — Go tokenizer candidate comparison

## Checks

- [x] Current Go tokenizer path compared against S01 HF baseline.
- [x] Mismatch artifact exists at `benchmark-results/fd-tokenizer-go-current-m012-s02.txt`.
- [x] All five fixed probes fail parity, proving current path is not sufficient.
- [x] Artifact records mismatch details without raw probe text.
- [x] Alternative path researched: Go bindings to Hugging Face Rust tokenizers.
- [x] S03 recommendation is explicit.
- [x] Parser/leak checks passed.
- [x] GitNexus scope reviewed; medium risk is limited to comparator tool flows.

## UAT Result

Pass. S03 should test HF Rust tokenizer bindings in isolation before any runtime integration.

