# ONNX Artifact Provisioning Contract

This contract defines how the opt-in ONNX backend may obtain large/native artifacts for local Docker builds, hosted CI, and deployment packaging without committing binaries to git.

TEI remains the production/default runtime. ONNX remains explicit opt-in and experimental until artifact provisioning, hosted CI, operational diagnostics, rollout, and rollback gates all pass.

## Required artifacts

| Logical artifact | Manifest / contract source | Destination path | Required verification | Current source status |
|---|---|---|---|---|
| ONNX model | `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx` | size + sha256 from manifest | **Blocked for hosted CI**: exact exported binary has no immutable external source yet. |
| Native tokenizer static library | `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json` | `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` | archive sha256 + extracted size + extracted sha256 | **Pinned candidate selected**: `daulet/tokenizers` `v1.27.0` release asset; hosted proof still required. |
| Tokenizer JSON | `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` `model.source_files.tokenizer.json` | `tei-models/deepvk--USER-bge-m3/tokenizer.json` or CI cache equivalent | size + sha256 from manifest source_files | **Pinned candidate selected**: Hugging Face `deepvk/USER-bge-m3` revision `0cc6cfe48e260fb0474c753087a69369e88709ae`; hosted proof still required. |
| ONNX Runtime shared library | `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` `source_contract.onnx_runtime` | `.gsd/runtime/onnxruntime/libonnxruntime.so.1.26.0` or image `/opt/onnxruntime/libonnxruntime.so.1.26.0` | wheel sha256 + extracted library sha256 | **Pinned candidate selected**: PyPI `onnxruntime==1.26.0` CP313 manylinux x86_64 wheel; hosted proof still required. |

## Provisioning rules

1. Never commit `.onnx`, `libtokenizers.a`, `libonnxruntime.so*`, or model cache binaries.
2. Every provisioned artifact must be verified before build or runtime use.
3. Checksum mismatch is a hard failure. Do not run tests, quality gates, or benchmarks with an unverified artifact.
4. Missing source URLs are blockers, not defaults. A CI job must fail early with a message naming the missing logical artifact.
5. `--allow-missing` is valid only for metadata-shape checks in default CI. It is not evidence of ONNX runtime readiness.
6. Artifact URLs must not contain secrets. If signed URLs are used, pass them through masked CI secrets and never print them.
7. Native tokenizer artifacts must remain isolated behind `hf_tokenizers`; the ONNX runtime remains isolated behind `onnx`.
8. Default `api/Dockerfile` and default Go Quality CI must remain artifact-free and TEI/default-path oriented.

## Remote source safety policy

M029 hardens `tools/provision_onnx_artifacts.py` for remote artifact sources before hosted workflow output can be treated as rollout evidence.

Remote URL rules:

- local file sources remain supported for local/offline proof;
- remote sources must use `https`;
- redirects are rejected instead of followed;
- private, loopback, link-local, reserved, multicast, and unspecified resolved addresses are blocked by default;
- optional host allowlisting is available with repeated `--allowed-artifact-host` or `FD_ONNX_ALLOWED_ARTIFACT_HOSTS=host1,host2`;
- private artifact hosts can be enabled only for trusted local testing with `--allow-private-artifact-hosts` or `FD_ONNX_ALLOW_PRIVATE_ARTIFACT_HOSTS=true`;
- downloads are bounded by `--max-download-bytes` or `FD_ONNX_MAX_DOWNLOAD_BYTES`;
- URL display redacts query strings through the helper's summarized source display.

Archive rules:

- native tokenizer archive extraction selects `libtokenizers.a` by basename, writes to a fixed destination, and does not use `extractall`;
- the selected member must be a regular file;
- member size is checked before copy against manifest `size_bytes` when available, otherwise `--max-archive-member-bytes` / `FD_ONNX_MAX_ARCHIVE_MEMBER_BYTES`;
- extraction streams through the same bounded-copy helper.

Remaining security work:

- M028 LOW findings are remediated in M030 for default tool/startup behavior: manifest artifact paths are constrained to approved repo-relative roots, and default diagnostics prefer repo-relative or basename-safe path display over absolute host paths.
- Hosted ONNX packaging still needs immutable artifact sources and a real workflow run before it is rollout evidence.

## Artifact path policy

M030 constrains manifest and tooling artifact paths to approved repo-relative roots. Approved roots are:

- `.gsd/runtime/onnx` for ONNX model artifacts;
- `.gsd/runtime/tokenizers` for native tokenizer artifacts;
- `.gsd/runtime/onnxruntime` for runtime shared libraries provisioned by tooling;
- `tei-models` for tokenizer JSON/model cache inputs.

Path rules:

- manifest `artifact.local_path` values must be repo-relative;
- absolute paths and `..` traversal are rejected;
- paths outside approved roots are rejected;
- default diagnostics should use repo-relative paths or basename-safe placeholders such as `.../artifact.onnx` for absolute local inputs;
- build-script missing-artifact messages name the relevant env var instead of printing its full configured value.

## Source selection status

M031 selected pinned source candidates for every supporting artifact except the exported ONNX model binary.

| Artifact | Source status | Source | Verification |
|---|---|---|---|
| ONNX model `user-bge-m3-dense.onnx` | `blocked` | No immutable external source yet. Exact binary must be mirrored/uploaded or reproduced under a separate validated export gate. | Required size `1432482908`, sha256 `28538a17a99302e144149732d73fb273cd7c7a0468dc59167caa5a2d5ff2a3d4`. |
| Native tokenizer archive | `immutable_selected` candidate | `https://github.com/daulet/tokenizers/releases/download/v1.27.0/libtokenizers.linux-amd64.tar.gz` | Archive sha256 `72556cdca798dd4ea7cdaba308e5f0d68a8cb93b67c96edf485b7a0edd7b07f4`; extracted `libtokenizers.a` sha256 `e6862b31745bb7d07980fcee70e49cd3b4318097609180f5d2d3fb394f305d50`. |
| Tokenizer JSON | `immutable_selected` candidate | `https://huggingface.co/deepvk/USER-bge-m3/resolve/0cc6cfe48e260fb0474c753087a69369e88709ae/tokenizer.json` | Size `3327728`, sha256 `068d9f7ed9dd190a00a567e5f7750fdc591b93bc623072ac8050a303c25f5937`. |
| ONNX Runtime CP313 wheel | `immutable_selected` candidate | PyPI `onnxruntime-1.26.0-cp313-cp313-manylinux_2_27_x86_64.manylinux_2_28_x86_64.whl` | Wheel sha256 `c07af6fc6d5557835f2b6ee7a96d8b3235d0c57a8e230efdedaee106a8a3cbc6`; extracted `libonnxruntime.so.1.26.0` sha256 `50775d390eb55e7abd9f6d734da103a04f0e5342ef0a76b1c6ec795544439295`. |

`immutable_selected` here means a non-secret, pinned, checksum-matched source candidate has been selected. It does not mean hosted CI has run or that ONNX is production-ready.

## Recommended cache layout

Local and CI caches should use ignored paths:

```text
.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx
.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a
.gsd/runtime/docker/onnx1024-context/
```

Hosted CI can map these paths from an artifact cache, object storage download, or release asset download, then run:

```bash
python3 tools/provision_onnx_artifacts.py --dry-run \
  --onnx-manifest docs/onnx-artifacts/user-bge-m3-dense-fp32.json \
  --native-tokenizer-manifest docs/onnx-artifacts/hf-tokenizers-linux-amd64.json

python3 tools/verify_onnx_artifacts.py \
  --onnx-manifest docs/onnx-artifacts/user-bge-m3-dense-fp32.json \
  --native-tokenizer-manifest docs/onnx-artifacts/hf-tokenizers-linux-amd64.json

IMAGE_TAG=fd-api:onnx1024-ci tools/build_onnx_image.sh
```

For hosted proof, `.github/workflows/onnx-packaging.yml` provides a manual `workflow_dispatch` skeleton. It requires explicit artifact source inputs, then runs provisioning, strict verifier, tagged tokenizer tests, ONNX smoke tests, and the Docker image build. The workflow is intentionally not triggered by `push` or `pull_request` until immutable artifact sources exist.

Do not pass signed or secret-bearing URLs as plain workflow inputs. Use non-secret immutable URLs, release assets, object keys, or masked secrets wired into a future hardened workflow.

The dry-run command is safe when sources are missing; full provisioning requires explicit `--onnx-source` and `--native-tokenizer-source` arguments or equivalent environment variables.

## Source selection recommendation

Preferred order:

1. immutable object storage keys with sha256 in the object name or metadata;
2. immutable GitHub release assets with pinned tag, never `latest`;
3. CI cache warmed by a trusted workflow and keyed by sha256;
4. local developer files for local-only proof.

Current blocker:

- The ONNX model artifact has only a local ignored path. Hosted CI/full deployment cannot be truthful until this exact binary has an immutable external source, or until a separate reproducible-export workflow is created and revalidates quality/performance from the regenerated artifact.
- Pinned source candidates exist for native tokenizer, tokenizer JSON, and ONNX Runtime, but they still require a real hosted workflow run before they are rollout evidence.

## Local export contract verifier

M032 adds `tools/verify_onnx_export_contract.py` as a local verifier for the current ignored ONNX export artifact. It checks the tracked ONNX manifest, M010 source provenance, M010 export metadata, and the local `.onnx` file.

The verifier's claim scope is intentionally narrow:

```text
existing_artifact_contract_verification_not_regenerated_export
```

It verifies:

- `production_default=false`;
- local ONNX artifact size and sha256;
- M010 model revision and source file checksums;
- M010 export toolchain pins, including `transformers==4.51.3`;
- export metadata output checksum and CPU provider/1024-dimension metadata.

It does not regenerate the ONNX binary. It must not be treated as byte-for-byte reproducible export proof. The next gate must choose one of two paths:

1. exact-binary hosting: mirror/upload the existing `.onnx` binary to a non-secret immutable source and verify size/sha256; or
2. reproducible-export workflow: regenerate the ONNX binary from pinned source/toolchain and rerun legal quality, performance, packaging, and hosted proof gates.

## Failure and diagnostics contract

Provisioning failures must name:

- logical artifact label;
- manifest path;
- destination path;
- whether the failure was missing source, missing destination, size mismatch, or checksum mismatch.

Provisioning output must not print raw embedding input text or secrets.
