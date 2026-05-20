# ONNX artifact provisioning contract

This directory tracks metadata for the opt-in ONNX backend. It does **not** track the ONNX model binary or the native Hugging Face tokenizer static library.

## Current status

- TEI remains the production/default runtime.
- ONNX remains opt-in experimental.
- `user-bge-m3-dense-fp32.json` is the tracked runtime contract for the local FP32 dense-only `deepvk/USER-bge-m3` ONNX artifact.
- The artifact was exported with dynamic sequence axes and has been validated locally at `ONNX_MAX_SEQUENCE_LENGTH=1024`.
- Packaging, CI provisioning, and operational rollout are still future gates.

## Required local artifacts

| Artifact | Manifest | Local ignored path | Purpose |
|---|---|---|---|
| ONNX model | `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx` | Dense FP32 embedding runtime |
| HF tokenizer static library | `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json` | `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` | Tagged Go tokenizer parity path |

Both local artifact paths are intentionally outside tracked source. Do not commit `.onnx` files or `libtokenizers.a`.

## Verification

Run this before tagged ONNX tests, local packaging experiments, or benchmark evidence collection:

```bash
python3 tools/verify_onnx_artifacts.py \
  --onnx-manifest docs/onnx-artifacts/user-bge-m3-dense-fp32.json \
  --native-tokenizer-manifest docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
```

Expected local result when artifacts are provisioned:

- both artifacts are present;
- manifest checksums match local files;
- manifest sizes match local files;
- `production_default` is false;
- `artifact.git_tracked` is false;
- no artifact path is tracked by git.

For default CI jobs that intentionally do not provision ONNX/native artifacts, use `--allow-missing` only as a contract-shape check. It must not be used as evidence for tagged ONNX runtime readiness.

## Runtime environment for validated ONNX 1024 path

The currently validated local runtime path uses:

```bash
EMBEDDING_BACKEND=onnx
ONNX_ARTIFACT_MANIFEST=docs/onnx-artifacts/user-bge-m3-dense-fp32.json
ONNX_TOKENIZER_PATH=tei-models/deepvk--USER-bge-m3/tokenizer.json
ONNX_MAX_SEQUENCE_LENGTH=1024
CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64'
go run -tags 'onnx hf_tokenizers' .
```

The exact ONNX Runtime shared library path is environment-specific and should be supplied by the local/packaged environment. Benchmark artifacts should record it via `BENCHMARK_ONNX_RUNTIME_LIBRARY`.

## Evidence already collected

- Legal quality PASS at 1024: `benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt`
- Local performance viability at 1024: `benchmark-results/fd-benchmark-m019-onnx1024.txt`
- Runtime contract decision: `.gsd/DECISIONS.md` decision D018

## Future Docker/CI gate

A future packaging milestone should define how CI or Docker builds obtain these artifacts without committing binaries. That gate should verify:

1. artifact download/cache/provisioning path;
2. checksum validation before build or runtime use;
3. tagged `hf_tokenizers` tokenizer tests and `onnx hf_tokenizers` backend build/tests in the provisioned environment;
4. packaged legal quality rerun at `ONNX_MAX_SEQUENCE_LENGTH=1024`;
5. packaged performance benchmark against TEI baseline;
6. startup health and actionable failure messages for missing or mismatched artifacts.
