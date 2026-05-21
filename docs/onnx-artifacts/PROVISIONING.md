# ONNX Artifact Provisioning Contract

This contract defines how the opt-in ONNX backend may obtain large/native artifacts for local Docker builds, hosted CI, and deployment packaging without committing binaries to git.

TEI remains the production/default runtime. ONNX remains explicit opt-in and experimental until artifact provisioning, hosted CI, operational diagnostics, rollout, and rollback gates all pass.

## Required artifacts

| Logical artifact | Manifest / contract source | Destination path | Required verification | Current source status |
|---|---|---|---|---|
| ONNX model | `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx` | size + sha256 from manifest | **Blocked for hosted CI**: no external canonical URL is recorded yet. |
| Native tokenizer static library | `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json` | `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` | size + sha256 from manifest | Upstream URL exists but uses `latest`; production should pin an immutable release asset or mirror. |
| Tokenizer JSON | `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` `model.source_files.tokenizer.json` | `tei-models/deepvk--USER-bge-m3/tokenizer.json` or CI cache equivalent | size + sha256 from manifest source_files | Needs explicit source URL/cache entry before hosted CI. |
| ONNX Runtime shared library | benchmark/runtime config | packaged image path such as `/opt/onnxruntime/libonnxruntime.so.1.26.0` | pinned version + sha256 | Needs pinned distribution source or base image contract before hosted CI. |

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

- LOW findings from M028 remain open: default log/path output sanitization and manifest artifact path root policy.
- Hosted ONNX packaging still needs immutable artifact sources and a real workflow run before it is rollout evidence.

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

- The ONNX model artifact has only a local ignored path. Hosted CI/full deployment cannot be truthful until this binary has an immutable external source.
- The native tokenizer manifest currently references an upstream `latest` URL; this is acceptable only when checksum verification remains mandatory, but production should pin or mirror it.

## Failure and diagnostics contract

Provisioning failures must name:

- logical artifact label;
- manifest path;
- destination path;
- whether the failure was missing source, missing destination, size mismatch, or checksum mismatch.

Provisioning output must not print raw embedding input text or secrets.
