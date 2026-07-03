# S01: Benchmark matrix and metadata harness

**Goal:** Design the TEI vs tagged ONNX benchmark matrix and extend tooling only enough to capture required metadata and runtime modes.
**Demo:** After this, the benchmark matrix and harness changes are defined before running expensive measurements.

## Must-Haves

- Benchmark matrix names cold/warm/batch/cache/startup/memory dimensions.
- Required metadata fields are defined.
- Harness can target arbitrary API URL and record tagged ONNX metadata.
- Raw probe text remains excluded.
- No runtime switch occurs.

## Proof Level

- This slice proves: Plan/harness proof plus dry-run checks.

## Integration Closure

Prevents ad hoc benchmark runs that cannot be compared.

## Verification

- Defines config snapshot fields for build tags, native artifact, ONNX artifact, ORT library, Redis namespace, and startup/memory.

## Tasks

- [x] **T01: Define benchmark matrix and metadata contract** `est:small`
  Inspect current benchmark.py config snapshot and M013 artifacts to decide the minimal metadata additions needed for tagged ONNX benchmark comparability.
  - Files: `benchmark.py`
  - Verify: Task summary lists scenarios and metadata fields.

- [x] **T02: Add tagged ONNX metadata snapshot fields** `est:medium`
  Extend benchmark.py snapshot to optionally include tagged ONNX/native metadata from env vars and manifest files, while preserving existing TEI behavior and raw text exclusion.
  - Files: `benchmark.py`
  - Verify: `uv run --python 3.13 --with requests --with redis python -m py_compile benchmark.py` and a lightweight snapshot function check pass.

- [x] **T03: Verify benchmark harness metadata hygiene** `est:small`
  Run a dry-run or lightweight artifact parser check to confirm benchmark output still has config snapshot and no raw probe text leakage after metadata changes.
  - Files: `benchmark.py`
  - Verify: py_compile, targeted snapshot check, raw-text leakage guard, and GitNexus detect_changes pass.

## Files Likely Touched

- benchmark.py
