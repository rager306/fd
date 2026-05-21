# S01: Target runtime validation contract

**Goal:** Document target runtime validation boundary and gates.
**Demo:** After this, the repo states exactly which Go/API/package gates are required before Python-generated/provisioned ONNX evidence can count for production acceptance.

## Must-Haves

- Python proof boundary is explicit.
- Go target runtime gates are explicit.
- Future Rust runtime gate rule is explicit.
- Redis namespace isolation is required for comparisons.
- No production/default promotion.

## Proof Level

- This slice proves: Manifest/docs/outcome checks; no runtime regeneration or external action.

## Integration Closure

Aligns M036 reproducible-export contract with Go production runtime and potential Rust future runtime.

## Verification

- Future agents can distinguish export helper checks, Go runtime smoke tests, packaged runtime gates, and future Rust gates.

## Tasks

- [x] **T01: Inspect target runtime surfaces** `est:small`
  Inspect existing Go ONNX runtime tests, API integration tests, legal evaluator, and benchmark harness to ground the target-runtime acceptance contract in actual project surfaces.
  - Files: `api/embed/onnx_test.go`, `api/embed/hf_tokenizer_native_test.go`, `api/main_test.go`, `api/handlers/embeddings_integration_test.go`, `tools/evaluate_legal_retrieval.py`, `benchmark.py`
  - Verify: Summarize Go/API/package validation surfaces and Python-helper boundary.

- [x] **T02: Persist target runtime validation contract** `est:small`
  Update ONNX manifest and provisioning docs with target-runtime validation contract: Python helper boundary, required Go API/package gates for any new or regenerated artifact, and equivalent gate rule for any future Rust backend.
  - Files: `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`, `docs/onnx-artifacts/PROVISIONING.md`
  - Verify: JSON/docs checks confirm Python boundary, Go gates, Rust gate rule, Redis namespace isolation, and no promotion claim.

- [x] **T03: Record target runtime outcome** `est:small`
  Update README and write outcome artifact summarizing target-runtime validation policy and remaining blockers.
  - Files: `docs/onnx-artifacts/README.md`, `benchmark-results/fd-onnx-target-runtime-validation-contract-m037-s01.txt`
  - Verify: Outcome/README checks pass; no raw text, secrets, signed URLs, or promotion overclaims.

- [x] **T04: Verify target runtime contract** `est:small`
  Run S01 verification: manifest JSON, contract marker checks, provisioning/export verifier, actionlint, and GitNexus detect.
  - Verify: S01 checks pass and GitNexus scope is low risk.

## Files Likely Touched

- api/embed/onnx_test.go
- api/embed/hf_tokenizer_native_test.go
- api/main_test.go
- api/handlers/embeddings_integration_test.go
- tools/evaluate_legal_retrieval.py
- benchmark.py
- docs/onnx-artifacts/user-bge-m3-dense-fp32.json
- docs/onnx-artifacts/PROVISIONING.md
- docs/onnx-artifacts/README.md
- benchmark-results/fd-onnx-target-runtime-validation-contract-m037-s01.txt
