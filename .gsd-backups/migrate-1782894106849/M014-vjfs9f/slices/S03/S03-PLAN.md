# S03: S03

**Goal:** Run tagged ONNX+HF-tokenizer benchmark using validated native/ONNX artifacts and isolated Redis namespace.
**Demo:** After this, tagged ONNX performance evidence exists under comparable benchmark conditions.

## Must-Haves

- Tagged ONNX server starts with `hf_tokenizers`.
- Redis namespace is isolated.
- Correctness/cosine gate is referenced or rerun before interpreting speed.
- Benchmark artifact includes native/ONNX/ORT metadata.
- Server cleanup verified.

## Proof Level

- This slice proves: Tagged server run plus benchmark artifact and cleanup proof.

## Integration Closure

Produces performance evidence for the benchmark-ready tagged ONNX path.

## Verification

- Records tagged startup, memory/RSS where possible, native/ONNX/ORT checksums, Redis namespace, and cache behavior.

## Tasks

- [x] **T01: Preflight confirmed tagged ONNX artifacts, checksums, and tagged tests are ready for the benchmark.** `est:small`
  Validate local ONNX/native/ORT artifact availability, confirm tagged test still passes, and record intended env vars including isolated cache namespace.
  - Verify: Artifact checks, tagged Go test, and namespace/env summary pass.

- [x] **T02: Started and verified the tagged ONNX benchmark server on port 18000.** `est:medium`
  Start tagged ONNX API on port 18000 with `hf_tokenizers`, isolated Redis namespace, and runtime env. Capture startup duration and memory/RSS where practical.
  - Verify: Health endpoint returns ok on port 18000 and process metadata is captured.

- [x] **T03: Run tagged ONNX benchmark with restart-aware harness** `est:medium`
  Run `benchmark.py` against tagged ONNX API with snapshot v3 tagged metadata, configurable API restart command, and write `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`.
  - Files: `benchmark.py`, `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`
  - Verify: Benchmark command exits 0 and artifact includes tagged runtime metadata plus configured restart behavior.

- [x] **T04: Verify tagged ONNX artifact and cleanup** `est:small`
  Verify tagged ONNX artifact for required sections, snapshot v3 metadata, raw text hygiene, correctness gate reference, and cleanup the tagged server.
  - Files: `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`
  - Verify: Parser/leak checks pass, health/process cleanup confirmed, GitNexus detect_changes pass.

## Files Likely Touched

- benchmark.py
- benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt
