---
id: T02
parent: S01
milestone: M014-vjfs9f
key_files:
  - benchmark.py
key_decisions:
  - `benchmark.py` snapshot version bumped to 2 because runtime metadata shape changed.
  - New metadata is optional and preserves TEI/default behavior when env vars are unset.
  - Artifact metadata records manifest SHA256 and artifact SHA256/path/existence; ONNX Runtime library metadata records path existence and SHA256.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:12:40.516Z
blocker_discovered: false
---

# T02: Added optional tagged ONNX/native/ORT metadata fields to the benchmark config snapshot without changing default benchmark behavior.

**Added optional tagged ONNX/native/ORT metadata fields to the benchmark config snapshot without changing default benchmark behavior.**

## What Happened

Extended `benchmark.py` effective configuration snapshot with optional runtime metadata for tagged ONNX benchmarks. New env vars include `BENCHMARK_RUNTIME_LABEL`, `BENCHMARK_BUILD_TAGS`, `BENCHMARK_ONNX_ARTIFACT_MANIFEST`, `BENCHMARK_NATIVE_TOKENIZER_MANIFEST`, and `BENCHMARK_ONNX_RUNTIME_LIBRARY`. The snapshot now includes parseable manifest metadata, artifact paths/checksums/existence, native tokenizer source metadata, ONNX Runtime library path/hash, and build tags. A lightweight snapshot check confirmed tagged ONNX metadata is present while raw input logging remains disabled.

## Verification

Python compile passed and a lightweight snapshot check with tagged ONNX env vars confirmed runtime metadata fields are populated.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile benchmark.py` | 0 | ✅ pass | 0ms |
| 2 | `BENCHMARK_RUNTIME_LABEL=tagged-onnx-hf BENCHMARK_BUILD_TAGS=hf_tokenizers BENCHMARK_ONNX_ARTIFACT_MANIFEST=... BENCHMARK_NATIVE_TOKENIZER_MANIFEST=... BENCHMARK_ONNX_RUNTIME_LIBRARY=... uv run --python 3.13 --with requests --with redis python snapshot check` | 0 | ✅ pass — snapshot_version=2; manifests parseable; ORT sha256 present | 0ms |

## Deviations

None.

## Known Issues

The snapshot records native/ONNX artifact metadata but does not enforce that actual artifact SHA256 matches manifest SHA256. Enforcement remains in artifact validators; benchmark snapshot is observational.

## Files Created/Modified

- `benchmark.py`
