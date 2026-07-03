# S03: HF tokenizer binding feasibility and parity

**Goal:** Test whether a Go binding to Hugging Face Rust tokenizers can exactly match the S01 baseline, then either integrate that path or close with blocker evidence.
**Demo:** After this, Go tokenizer parity either passes through a tested HF Rust tokenizers binding or is blocked with concrete linking/correctness evidence.

## Must-Haves

- Candidate HF Rust tokenizers Go binding is tested or concrete setup blocker is recorded.
- If candidate passes all S01 probes, runtime integration is implemented behind opt-in ONNX path and tests pass.
- If candidate fails or cannot be linked, blocker evidence is persisted and runtime code remains unchanged.
- No raw probe text appears in artifacts.
- TEI remains default.

## Proof Level

- This slice proves: Isolated dependency feasibility plus token-level comparator proof; Go tests if runtime code changes.

## Integration Closure

If parity passes and packaging is tractable, S03 can update the opt-in ONNX tokenizer path; otherwise it preserves the ONNX blocker and gives S04 a clear decision.

## Verification

- Records dependency/linker setup, candidate token outputs, exact mismatch/pass evidence, and production packaging caveats.

## Tasks

- [x] **T01: Probe HF Rust tokenizers binding feasibility** `est:medium`
  Run an isolated feasibility probe for `github.com/daulet/tokenizers` or equivalent HF Rust tokenizers binding: determine module import, native library requirements, and whether prebuilt linux-amd64 assets can be used locally.
  - Verify: A temp-module or isolated command either encodes one probe or records exact linker/setup blocker.

- [x] **T02: Compare HF binding candidate against baseline** `est:medium`
  Compare the candidate binding output against S01 baseline for all fixed probes and persist a sanitized candidate artifact.
  - Files: `tools/compare_tokenizers.py`, `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt`
  - Verify: Candidate comparison exits 0 only if parity passes; otherwise writes mismatch/setup artifact.

- [x] **T03: Integrate parity path or record blocker** `est:medium`
  If candidate parity passes, integrate the binding into `api/embed/onnx.go` and add token parity tests. If it does not pass or cannot be packaged, do not integrate; write blocker summary instead.
  - Files: `api/embed/onnx.go`, `api/embed/onnx_test.go`, `api/go.mod`, `api/go.sum`
  - Verify: If code changes: Go tests pass. If blocker: blocker evidence names exact failure and runtime code remains unchanged.

- [x] **T04: Verify S03 outcome** `est:small`
  Run S03 verification: parser/leak checks, Go tests/lint if code changed, GitNexus detect_changes, and default TEI health if runtime code changed.
  - Verify: All applicable verification gates pass and S04 decision input is clear.

## Files Likely Touched

- tools/compare_tokenizers.py
- benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt
- api/embed/onnx.go
- api/embed/onnx_test.go
- api/go.mod
- api/go.sum
