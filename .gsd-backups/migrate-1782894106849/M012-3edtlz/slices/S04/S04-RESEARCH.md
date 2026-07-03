# S04 — Research

**Date:** 2026-05-20

## Summary

M012 resolves the M011 tokenizer-parity question with a precise answer: the current pure-Go `sugarme/tokenizer` path is not equivalent to Hugging Face, but exact parity is achievable through Go bindings to Hugging Face's Rust `tokenizers` implementation. The S01 baseline established authoritative HF token IDs and attention masks. S02 proved the current Go tokenizer fails all five fixed probes. S03 proved `github.com/daulet/tokenizers` plus `libtokenizers.a` matches all five probes exactly.

The milestone should not proceed to ONNX performance benchmarking yet. Tokenizer correctness is solved only in isolation; it is not integrated into the fd API runtime. The remaining blocker is build and packaging: `daulet/tokenizers` requires CGO and a native/static `libtokenizers.a` at link time. Importing that dependency into default API code without build tags or native artifact packaging would risk breaking normal TEI/default builds and CI.

The next milestone should therefore be a native packaging and build-tag integration milestone, not a throughput benchmark. Once a tagged opt-in build can link the HF tokenizer binding reproducibly in Docker/CI without affecting default TEI builds, then the ONNX embedder can use the parity-correct tokenizer and rerun TEI-vs-ONNX cosine. Only after cosine equivalence passes should ONNX latency/throughput benchmarking resume.

## Recommendation

Recommended next milestone: **HF tokenizer native packaging and opt-in build integration**.

Scope should include:

1. Add a documented `hf_tokenizers` or `onnx_hf_tokenizer` build-tag path for the ONNX tokenizer implementation.
2. Define how `libtokenizers.a` is obtained, checksummed, cached, and supplied in Docker/CI/local builds.
3. Preserve default builds without native tokenizer requirements.
4. Add token parity tests under the tagged build path using the S01 baseline.
5. Only after tagged token parity passes, rerun Go ONNX vs TEI cosine with isolated Redis cache namespace.

Do **not** run ONNX throughput benchmarks from the current default build. The default ONNX code still uses `sugarme/tokenizer`, which S02 proved non-equivalent.

## Implementation Landscape

### Key Files

- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt` — HF tokenization source of truth.
- `benchmark-results/fd-tokenizer-go-current-m012-s02.txt` — current `sugarme/tokenizer` mismatch evidence.
- `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt` — passing HF Rust tokenizers binding evidence.
- `tools/compare_tokenizers.py` — reusable harness for baseline/current/binding comparisons.
- `api/embed/onnx.go` — still uses `sugarme/tokenizer`; should not be treated as semantically correct until tagged integration replaces tokenization.
- `api/Dockerfile`, `.github/workflows/go-quality.yml`, and `api/go.mod` — likely future files for native packaging/build-tag integration.

### Build Order

1. Design native artifact handling for `libtokenizers.a` with checksum and architecture metadata.
2. Add build-tagged tokenizer implementation so default builds remain unaffected.
3. Update CI/Docker/local docs for tagged opt-in ONNX builds.
4. Add tagged token parity tests against S01 baseline.
5. Rerun isolated-cache ONNX cosine only after token parity passes in the actual API build.
6. Run performance benchmarks only after cosine equivalence passes.

### Verification Approach

Final M012 verification should confirm:

- Go tests pass in the default build.
- Lint passes in the default build.
- Default TEI health remains ok.
- Tokenizer artifacts parse and contain no raw probe text.
- `go-current` remains expected-fail until runtime integration changes.
- `go-hf-binding` remains pass when supplied with `libtokenizers.a`.
- GitNexus scope is understood.

## Constraints

- TEI remains default.
- ONNX remains opt-in.
- Default Go/API builds must not require `libtokenizers.a` until a build-tag/package strategy exists.
- No raw probe text in artifacts.
- No ONNX performance benchmark until token parity is integrated and cosine equivalence passes.

## Common Pitfalls

- **Mistaking isolated parity for runtime integration** — S03 proves candidate correctness, not packaged API readiness.
- **Breaking default builds** — native tokenizer imports can force link-time requirements even when ONNX is not selected at runtime.
- **Skipping cosine after tokenizer fix** — token parity is necessary but not sufficient; rerun TEI-vs-ONNX cosine before benchmarking speed.
- **Shared Redis cache masking** — continue using isolated cache namespaces for backend comparison.

## Open Risks

- Native library distribution may add CI/Docker complexity.
- Static library architecture/version drift could make reproducible builds harder.
- Build tags may complicate developer ergonomics if not documented well.

## Sources

- HF baseline: `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`
- Current Go mismatch: `benchmark-results/fd-tokenizer-go-current-m012-s02.txt`
- HF binding parity pass: `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt`
- S03 summaries: `.gsd/milestones/M012-3edtlz/slices/S03/S03-SUMMARY.md`
