# S03: Parity correct ONNX tokenizer integration

**Goal:** Integrate the parity-correct native HF tokenizer into the opt-in ONNX backend under the `hf_tokenizers` build tag, then rerun semantic equivalence if the tagged runtime starts.
**Demo:** After this, the ONNX tokenizer path either uses the parity-correct binding under opt-in conditions or remains blocked with exact evidence.

## Must-Haves

- Tagged ONNX embedder uses parity-correct native tokenizer path.
- Default untagged build remains green.
- Tagged tokenizer parity tests pass.
- If tagged ONNX API starts, isolated-cache TEI-vs-ONNX cosine is rerun.
- If tagged runtime cannot start, blocker evidence names exact native/ONNX issue.
- No ONNX throughput benchmark is run unless cosine passes.

## Proof Level

- This slice proves: Tagged Go tests plus optional local tagged ONNX API run and isolated-cache cosine artifact.

## Integration Closure

If successful, enables meaningful TEI-vs-ONNX cosine comparison with isolated Redis cache. If blocked, records exact build/runtime blocker.

## Verification

- Adds tokenizer implementation metadata and evidence for tagged ONNX runtime behavior.

## Tasks

- [x] **T01: Design tagged ONNX tokenizer integration** `est:small`
  Run impact analysis and design the smallest ONNX tokenizer abstraction needed to swap tokenizers under build tags without changing handlers/cache/API contract.
  - Files: `api/embed/onnx.go`, `api/embed/hf_tokenizer_native.go`
  - Verify: Impact analysis recorded; design names changed symbols and default-build behavior.

- [x] **T02: Implement tagged ONNX tokenizer path** `est:medium`
  Implement tokenizer abstraction and tagged native implementation so untagged builds keep current behavior and tagged builds use HF native tokenizer for ONNX input IDs/masks.
  - Files: `api/embed/onnx.go`, `api/embed/onnx_tokenizer_default.go`, `api/embed/onnx_tokenizer_hf.go`, `api/embed/hf_tokenizer_native.go`
  - Verify: Default tests pass; tagged tests pass with `CGO_LDFLAGS`.

- [x] **T03: Run tagged ONNX cosine comparison** `est:medium`
  Start tagged ONNX API locally using validated native tokenizer library and isolated Redis cache namespace, then rerun TEI-vs-ONNX cosine comparison. If startup fails, capture blocker evidence.
  - Files: `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt`
  - Verify: Cosine artifact passes threshold or records startup/runtime blocker; Redis namespace isolated.

- [x] **T04: Verify tagged ONNX integration outcome** `est:small`
  Run S03 verification: default tests/lint, tagged tests, artifact/leak checks, GitNexus detect_changes, and cleanup of local tagged server if started.
  - Verify: All applicable gates pass; no background tagged server remains.

## Files Likely Touched

- api/embed/onnx.go
- api/embed/hf_tokenizer_native.go
- api/embed/onnx_tokenizer_default.go
- api/embed/onnx_tokenizer_hf.go
- benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt
