# S02: Opt in build tag boundary

**Goal:** Create an opt-in build-tag/package boundary for HF tokenizer binding that does not affect default TEI builds.
**Demo:** After this, default builds stay clean and the project has an explicit opt-in tokenizer build path or a packaging blocker.

## Must-Haves

- Default Go tests/lint pass without native library requirements.
- Opt-in tagged path compiles or fails with documented native artifact blocker.
- Dependency changes are isolated from default runtime.
- Tokenizer parity harness can use the tagged path if implemented.
- GitNexus scope reviewed.

## Proof Level

- This slice proves: Default tests/lint plus tagged build/probe or blocker evidence.

## Integration Closure

Determines whether runtime integration can proceed without breaking default builds.

## Verification

- Adds clear build errors/checks when native artifact is missing and verifies default build isolation.

## Tasks

- [x] **T01: Design opt-in build tag boundary** `est:small`
  Design the Go package/build-tag boundary for native HF tokenizers: file names, build tags, interface shape, default fallback behavior, and dependency isolation.
  - Files: `api/embed/onnx.go`, `api/embed/onnx_test.go`, `api/go.mod`
  - Verify: Task summary states exact files/build tags and default-build safety rule.

- [x] **T02: Implement tagged native tokenizer probe** `est:medium`
  Implement a minimal build-tagged tokenizer package/probe that imports `github.com/daulet/tokenizers` only under the opt-in tag and can encode fixed probes using the native artifact path. Do not replace runtime ONNX tokenizer yet.
  - Files: `api/embed/hf_tokenizer_native.go`, `api/embed/hf_tokenizer_native_test.go`, `api/go.mod`, `api/go.sum`
  - Verify: Default `go test ./... -short` passes; tagged test compiles/runs when CGO_LDFLAGS points at `.gsd/runtime/tokenizers/linux-amd64`.

- [x] **T03: Verify tagged tokenizer parity** `est:medium`
  Run parity test through the tagged package against the M012 baseline and persist a build-tag comparison artifact if possible. If tagged build cannot run, record exact blocker.
  - Files: `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt`
  - Verify: Tagged parity command exits 0 and writes PASS artifact, or writes blocker evidence.

- [x] **T04: Verify build tag isolation** `est:small`
  Run S02 verification: default tests/lint without native flags, tagged tests with native flags if implemented, artifact/leak checks, GitNexus detect_changes.
  - Files: `api/embed/hf_tokenizer_native.go`, `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt`
  - Verify: Default build and applicable tagged build gates pass; no raw text/native binaries tracked.

## Files Likely Touched

- api/embed/onnx.go
- api/embed/onnx_test.go
- api/go.mod
- api/embed/hf_tokenizer_native.go
- api/embed/hf_tokenizer_native_test.go
- api/go.sum
- benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt
