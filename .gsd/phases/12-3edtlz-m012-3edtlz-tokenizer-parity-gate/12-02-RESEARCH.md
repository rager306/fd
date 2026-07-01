# S02 — Research

**Date:** 2026-05-20

## Summary

The current pure-Go tokenizer path is not viable as-is for `deepvk/USER-bge-m3` ONNX equivalence. The durable comparison artifact `benchmark-results/fd-tokenizer-go-current-m012-s02.txt` shows that `github.com/sugarme/tokenizer/pretrained.FromFile + EncodeSingle(addSpecialTokens=true)` fails parity for all five fixed Russian/legal probes. Token counts differ on every probe, first input-ID mismatches appear as early as index 2, and attention-mask lengths differ because Go produces more tokens than the Hugging Face baseline.

A simple configuration toggle does not explain the mismatch. Running `EncodeSingle` without special tokens removes BOS/EOS but still leaves the core Russian tokenization divergent. The tokenizer JSON is an XLM-R-style Hugging Face tokenizers configuration: model type `Unigram`, normalizer `Sequence`, pre-tokenizer `Metaspace` with `add_prefix_space=true`, post-processor `TemplateProcessing`, decoder `Metaspace`, and five added tokens. This points to implementation-level divergence in the Go tokenizer stack rather than a missing special-token flag.

The best next path is to test bindings to Hugging Face's Rust `tokenizers` implementation, not to patch `sugarme/tokenizer` blind. `daulet/tokenizers` and the related `crowdsecurity/go-tokenizers` expose Go bindings for Hugging Face Tokenizers and can load a `tokenizer.json`, but require native `libtokenizers` static library handling via CGO/linker flags or prebuilt binaries. That operational cost is real, but exact tokenizer parity is more important than staying pure-Go if ONNX correctness is the gate.

## Recommendation

For S03, run an isolated feasibility test of `github.com/daulet/tokenizers` against the S01 baseline. Do not wire it into `api/embed/onnx.go` until the isolated comparison passes all five probes.

Decision tree for S03:

1. If `daulet/tokenizers` can be built/linked locally and matches the HF baseline exactly, switch the ONNX tokenizer path to that binding behind the opt-in ONNX backend and add token parity tests.
2. If the binding works but introduces unacceptable packaging complexity, record the blocker and decide between sidecar and defer.
3. If no Go binding can be made to match quickly, keep ONNX blocked and recommend either a Python/Rust tokenizer sidecar or deferring Go ONNX integration.

Do not spend S03 trying to hand-roll XLM-R/Unigram/Metaspace normalization. The risk of a subtly wrong tokenizer is too high, and M011 already showed that small preprocessing differences produce materially different embeddings.

## Implementation Landscape

### Key Files

- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt` — authoritative Hugging Face baseline.
- `benchmark-results/fd-tokenizer-go-current-m012-s02.txt` — current Go tokenizer mismatch evidence.
- `tools/compare_tokenizers.py` — comparison harness; can be extended for additional Go binding candidates.
- `api/embed/onnx.go` — current ONNX embedder uses `sugarme/tokenizer`; only modify in S03 if parity is proven.
- `api/go.mod` / `api/go.sum` — only add a new tokenizer binding if isolated proof passes.

### Build Order

1. Keep S02 as evidence/research; do not alter runtime code here.
2. In S03, create an isolated candidate comparison for `daulet/tokenizers` using the same S01 baseline.
3. If candidate comparison passes, then integrate the dependency and replace the tokenization path.
4. If it fails or cannot be linked reproducibly, close S03 with blocker evidence.

### Verification Approach

- Current mismatch evidence: `uv run --python 3.13 --with transformers --with torch --with sentencepiece python tools/compare_tokenizers.py --mode go-current` should write `benchmark-results/fd-tokenizer-go-current-m012-s02.txt` and exit `2` until parity is fixed.
- S03 candidate success requires the same comparison shape to exit `0` for all five probes.
- Raw probe text leakage checks remain mandatory.
- Go tests/lint remain mandatory if runtime code or dependencies change.

## Constraints

- TEI remains default.
- ONNX remains opt-in.
- No throughput benchmark until tokenizer parity and cosine equivalence pass.
- Tokenizer artifacts may include token IDs and hashes, but not raw probe text.
- Native bindings add CGO/static-library packaging work; this must be called out before production use.

## Common Pitfalls

- **Special-token false positive** — removing BOS/EOS changes length but does not fix the divergent core tokenization.
- **Hand-rolled tokenizer drift** — XLM-R uses Unigram plus normalizer/metaspace/post-processing details; a nearly-correct implementation can still break embeddings.
- **Ignoring linker artifacts** — Rust-backed Go bindings may solve correctness but require `libtokenizers` distribution and CI/runtime packaging.
- **Benchmarking failed parity** — current `go-current` mismatch means ONNX speed claims remain invalid.

## Open Risks

- `daulet/tokenizers` may require Rust toolchain or prebuilt `libtokenizers` that is not currently part of the project images/CI.
- The binding may match token IDs but introduce deployment complexity comparable to a sidecar.
- A Python tokenizer sidecar may be operationally simpler than native static-library packaging, but it changes the architecture and performance profile.

## Sources

- Current mismatch artifact: `benchmark-results/fd-tokenizer-go-current-m012-s02.txt`
- HF tokenizer baseline: `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`
- Local tokenizer JSON inspection: `model.type=Unigram`, `pre_tokenizer=Metaspace`, `post_processor=TemplateProcessing`
- Hugging Face tokenizer algorithms documentation: https://huggingface.co/docs/transformers/en/tokenizer_summary
- XLM-RoBERTa tokenizer implementation reference: https://github.com/huggingface/transformers/blob/main/src/transformers/models/xlm_roberta/tokenization_xlm_roberta.py
- `daulet/tokenizers` README: https://raw.githubusercontent.com/daulet/tokenizers/main/README.md
- `crowdsecurity/go-tokenizers` README: https://raw.githubusercontent.com/crowdsecurity/go-tokenizers/main/README.md
