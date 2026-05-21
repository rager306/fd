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

See `docs/onnx-artifacts/PROVISIONING.md` for the full external artifact provisioning/cache contract. See `docs/onnx-artifacts/OPERATIONS.md` for the operational diagnostics and rollout/rollback contract.

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


## Local ONNX Docker image proof

Use the dedicated opt-in packaging script after local artifacts are provisioned and verified:

```bash
IMAGE_TAG=fd-api:onnx1024-local tools/build_onnx_image.sh
```

The script creates a temporary Docker context under `.gsd/runtime/docker/`, copies only the API source plus verified local artifacts, and builds with `-tags "onnx hf_tokenizers"`. The default `api/Dockerfile` remains the TEI/default image and must not require these artifacts.

## CI boundary

The regular Go Quality workflow now runs only artifact-free ONNX packaging checks:

- `tools/verify_onnx_artifacts.py --allow-missing` validates manifest shape and safety metadata without requiring local binaries;
- a binary hygiene check fails if `.onnx`, `libtokenizers.a`, or `libonnxruntime.so` are tracked;
- default Go tests and lint remain TEI/default-path checks.

Full ONNX image CI is intentionally not claimed yet. It requires an external artifact store/cache/provisioning step for:

- `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`;
- `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a`;
- `tei-models/deepvk--USER-bge-m3/tokenizer.json` or an equivalent pinned tokenizer artifact;
- `libonnxruntime.so.1.26.0` or an equivalent pinned ONNX Runtime distribution.

Only after those artifacts are provisioned and verified should CI run `IMAGE_TAG=... tools/build_onnx_image.sh`, tagged ONNX tests, packaged legal quality, and packaged performance benchmarks.

A manual workflow skeleton exists at `.github/workflows/onnx-packaging.yml`. It is intentionally `workflow_dispatch` only and requires explicit artifact source inputs before it can provision artifacts, run tagged tests, and build the ONNX image. Do not use signed or secret-bearing URLs as plain workflow inputs; use masked secrets or a non-secret immutable artifact URL/cache key.

Manual workflow input policy is defined in `docs/onnx-artifacts/PROVISIONING.md`. In short: `onnx_source_url` and `native_tokenizer_source_url` are required; `tokenizer_json_source_url`, `onnx_runtime_source_url`, and `onnx_runtime_sha256` are optional. When `onnx_runtime_source_url` is supplied without `onnx_runtime_sha256`, provisioning uses `source_contract.onnx_runtime.library_sha256` from `user-bge-m3-dense-fp32.json`. The exact ONNX model binary source remains the blocking required input.

M035 defines a planned exact-binary hosting contract for that blocker: the current `.onnx` must be mirrored/uploaded as the exact `1432482908` byte binary with sha256 `28538a17a99302e144149732d73fb273cd7c7a0468dc59167caa5a2d5ff2a3d4`. The recommended object key and release filename are documented, but they are not real sources until the artifact is actually uploaded/mirrored and reverified. Do not use local developer paths, mutable `latest` URLs, or signed/plain-secret URLs as hosted proof.

## Future Docker/CI gate

A future packaging milestone should define how CI or Docker builds obtain these artifacts without committing binaries. That gate should verify:

1. artifact download/cache/provisioning path;
2. checksum validation before build or runtime use;
3. tagged `hf_tokenizers` tokenizer tests and `onnx hf_tokenizers` backend build/tests in the provisioned environment;
4. packaged legal quality rerun at `ONNX_MAX_SEQUENCE_LENGTH=1024`;
5. packaged performance benchmark against TEI baseline;
6. startup health and actionable failure messages for missing or mismatched artifacts.
