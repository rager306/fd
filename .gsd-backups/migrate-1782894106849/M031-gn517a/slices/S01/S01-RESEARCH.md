---
milestone: M031-gn517a
slice: S01
type: source-contract-research
status: draft
---

# M031 S01 — ONNX Artifact Source Contract Research

## Scope

Define truthful source status for every artifact required by the opt-in ONNX backend. This is not a hosted workflow run, not a binary upload, and not production/default promotion.

## Current artifact inventory

| Artifact | Required by | Local / destination path | Size | SHA256 | Current source status |
|---|---|---|---:|---|---|
| ONNX model `user-bge-m3-dense.onnx` | ONNX backend model artifact | `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx` | 1432482908 | `28538a17a99302e144149732d73fb273cd7c7a0468dc59167caa5a2d5ff2a3d4` | BLOCKED: local export only; no immutable external source recorded. |
| Native tokenizer `libtokenizers.a` | `hf_tokenizers` tagged build | `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` | 49776794 | `e6862b31745bb7d07980fcee70e49cd3b4318097609180f5d2d3fb394f305d50` | CANDIDATE/BLOCKED: manifest has upstream `latest` URL; needs pinned release asset or mirror. |
| Tokenizer JSON | ONNX native tokenizer runtime | `tei-models/deepvk--USER-bge-m3/tokenizer.json` | 3327728 | `068d9f7ed9dd190a00a567e5f7750fdc591b93bc623072ac8050a303c25f5937` | CANDIDATE: source can be the pinned Hugging Face model revision, but hosted packaging needs explicit URL/cache input. |
| ONNX Runtime shared library | tagged ONNX runtime/image | `.gsd/runtime/onnxruntime/libonnxruntime.so.1.26.0` or image `/opt/onnxruntime/libonnxruntime.so.1.26.0` | local observed size not tracked in manifest | `50775d390eb55e7abd9f6d734da103a04f0e5342ef0a76b1c6ec795544439295` | CANDIDATE/BLOCKED: version and sha are known locally, but no tracked runtime source manifest/source URL exists. |

## Existing source/provenance evidence

- ONNX model was exported locally from `deepvk/USER-bge-m3` at model revision `0cc6cfe48e260fb0474c753087a69369e88709ae` using `tools/export_user_bge_m3_dense_onnx.py`.
- Export packages included `python=3.13.12`, `torch=2.12.0`, `transformers=4.51.3`, `onnx=1.21.0`, `onnxruntime=1.26.0`, `safetensors=0.7.0`.
- Unpinned/latest `transformers 5.8.1` failed for this export path and must not be treated as equivalent.
- Tokenizer JSON and model source files already have local sha256 evidence in `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` and `.gsd/runtime/onnx/m010-s03/source-provenance.json`.
- Native tokenizer manifest currently records `https://github.com/daulet/tokenizers/releases/latest/download/libtokenizers.linux-amd64.tar.gz`, but `latest` is mutable and not acceptable as immutable rollout input.

## Source status definitions

- `immutable_selected`: source URL/key is non-secret, pinned, policy-compliant, and matches known size/sha in tracked contract.
- `candidate`: source strategy is plausible but still needs verification/pinning or a real hosted proof.
- `blocked`: no acceptable immutable source exists yet or required evidence is missing.

## Candidate source assessment

### ONNX model `user-bge-m3-dense.onnx`

Status: `blocked`.

The exact ONNX binary was exported locally from the pinned HF model revision and is ignored under `.gsd/runtime/onnx/m010-s03/`. No immutable external URL, release asset, object-storage key, or CI cache key exists yet for the exported binary itself.

Required next action before hosted proof: upload or mirror the exact ONNX binary to a non-secret immutable source and verify that the downloaded artifact has:

- size: `1432482908`
- sha256: `28538a17a99302e144149732d73fb273cd7c7a0468dc59167caa5a2d5ff2a3d4`

Do not treat the upstream HF model source files as a substitute for this exported ONNX artifact. Regenerating the ONNX model is a separate reproducibility path and must pin the export toolchain (`transformers==4.51.3`, `torch==2.12.0`, `onnx==1.21.0`, `onnxruntime==1.26.0`) and re-run legal/performance gates.

### Native tokenizer `libtokenizers.a`

Status: `immutable_selected` for the source candidate; rollout still needs hosted workflow proof.

Public GitHub release metadata shows a pinned asset:

- source: `https://github.com/daulet/tokenizers/releases/download/v1.27.0/libtokenizers.linux-amd64.tar.gz`
- archive size: `14331593`
- archive sha256: `72556cdca798dd4ea7cdaba308e5f0d68a8cb93b67c96edf485b7a0edd7b07f4`
- extracted member: `libtokenizers.a`
- extracted size: `49776794`
- extracted sha256: `e6862b31745bb7d07980fcee70e49cd3b4318097609180f5d2d3fb394f305d50`

Verification probe: `.gsd/exec/78224446-6bc5-4bcb-84f1-633ade016a0c.stdout`.

This replaces the mutable `latest` URL as the recommended source. Hosted workflow evidence is still required before rollout claims.

### Tokenizer JSON

Status: `immutable_selected` for the source candidate; rollout still needs hosted workflow proof.

The pinned Hugging Face revision URL is:

- source: `https://huggingface.co/deepvk/USER-bge-m3/resolve/0cc6cfe48e260fb0474c753087a69369e88709ae/tokenizer.json`
- size: `3327728`
- sha256: `068d9f7ed9dd190a00a567e5f7750fdc591b93bc623072ac8050a303c25f5937`

Verification probe: `.gsd/exec/24907b45-1ff7-436e-a6fa-b99bc5be5f06.stdout`.

This source is acceptable only with checksum verification. It does not authorize model replacement or ONNX production/default use.

### ONNX Runtime shared library

Status: `immutable_selected` for the source candidate; rollout still needs hosted workflow proof and manifest/docs wiring.

PyPI metadata for `onnxruntime==1.26.0` exposes a CPython 3.13 manylinux x86_64 wheel matching the local runtime library used in prior gates:

- source: `https://files.pythonhosted.org/packages/3d/26/4d09ddc755a84fc8d5e192991626b0e0680e8f6c5d58f4f1d05c42bc48cf/onnxruntime-1.26.0-cp313-cp313-manylinux_2_27_x86_64.manylinux_2_28_x86_64.whl`
- wheel size: `18185632`
- wheel sha256: `c07af6fc6d5557835f2b6ee7a96d8b3235d0c57a8e230efdedaee106a8a3cbc6`
- extracted member: `onnxruntime/capi/libonnxruntime.so.1.26.0`
- extracted library size: `23031768`
- extracted library sha256: `50775d390eb55e7abd9f6d734da103a04f0e5342ef0a76b1c6ec795544439295`

Verification probe: `.gsd/exec/24907b45-1ff7-436e-a6fa-b99bc5be5f06.stdout`.

This is Python-version/platform-specific. Future hosted workflow should either use this exact CP313 wheel for extracting the library or define a separate Linux x64 ONNX Runtime release source and verify the extracted library checksum.

## S01 conclusion

- The native tokenizer, tokenizer JSON, and ONNX Runtime now have pinned, checksum-matched immutable source candidates.
- The exported ONNX model binary remains blocked because there is no immutable external source for the exact local `.onnx` file.
- Hosted ONNX packaging cannot be truthful until the ONNX model binary is uploaded/mirrored to a non-secret immutable source or a separately validated reproducible-export path is built.
- TEI remains production/default. ONNX remains opt-in experimental.
