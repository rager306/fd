# ONNX Operational Diagnostics and Rollout Contract

This runbook defines operational expectations for the opt-in ONNX backend. It does not approve production/default promotion. TEI remains the default runtime until a separate production rollout decision is made.

## Runtime mode boundaries

| Mode | Backend | Expected state |
|---|---|---|
| Default production | TEI | `EMBEDDING_BACKEND` unset or `tei`; default `api/Dockerfile`; no ONNX artifacts required. |
| Opt-in ONNX validation | ONNX | `EMBEDDING_BACKEND=onnx`, `ONNX_MAX_SEQUENCE_LENGTH=1024`, build tags `onnx hf_tokenizers`, verified artifacts available. |
| Rollback | TEI | Reset backend to TEI/default image and stop ONNX container/job. Redis cache namespace may remain but must not be reused for cross-backend claims. |

## Startup preflight contract

Before an ONNX runtime serves traffic it must validate:

1. `production_default` is false unless a future production decision explicitly changes it.
2. `ONNX_ARTIFACT_MANIFEST` exists and parses as JSON.
3. ONNX artifact exists at the manifest local path or provisioned image path.
4. ONNX artifact size and sha256 match the manifest.
5. Runtime metadata matches expectations:
   - output name: `dense_vecs`;
   - dimensions: `1024`;
   - normalized vectors: `true`;
   - validated max sequence length: `1024`.
6. `ONNX_MAX_SEQUENCE_LENGTH` is present and does not exceed the validated contract being claimed.
7. `ONNX_TOKENIZER_PATH` exists and corresponds to the model/tokenizer manifest checksum.
8. ONNX Runtime shared library exists and is loadable.
9. Native tokenizer build path is opt-in only and never required for default builds.
10. Redis cache namespace is explicit for any TEI-vs-ONNX comparison.

## Actionable failure messages

Failures must be hard failures before inference, not silent fallbacks.

| Failure | Required diagnostic fields |
|---|---|
| Missing manifest | env key, configured path, backend. |
| Manifest parse error | path, JSON error type, artifact label. |
| Missing ONNX artifact | artifact_id, expected path, provisioning doc path. |
| Checksum mismatch | artifact_id, expected sha256, actual sha256, path. |
| Size mismatch | artifact_id, expected size, actual size, path. |
| Sequence length mismatch | configured value, validated value, manifest path. |
| Missing tokenizer JSON | configured tokenizer path, model id. |
| Missing ONNX Runtime library | configured library path, backend. |
| Unsupported provider | configured provider, available providers when available. |
| Production-default misuse | manifest artifact_id, `production_default`, required decision reference. |

Never log raw embedding input text. Never log signed artifact URLs or secret values.

## Health and observability expectations

Minimum surfaces before rollout:

- startup log event: backend, runtime label, artifact_id, model id, model revision, max sequence length, dimensions, cache namespace, verification status;
- health endpoint remains `ok` only after successful preflight and model/runtime initialization;
- failed preflight exits non-zero for containers and records a clear final error;
- benchmark artifacts record sanitized effective config and artifact metadata;
- CI logs show provisioning and verification phases separately.

Recommended future health detail for opt-in ONNX:

```json
{
  "backend": "onnx",
  "artifact_id": "user-bge-m3-dense-fp32",
  "validated_max_sequence_length": 1024,
  "dimensions": 1024,
  "artifact_verified": true,
  "production_default": false
}
```

This should expose only metadata, never raw input text or secrets.

## Implemented diagnostics status

Implemented in M026:

- default `/health` response remains compatible (`status` + `time`) for TEI/default runtime;
- ONNX runtime can attach safe `/health.runtime` metadata after successful startup;
- ONNX health metadata includes backend, model, artifact_id, dimensions, configured max sequence length, validated max sequence length, artifact verification state, production_default flag, and cache namespace;
- ONNX health metadata intentionally excludes manifest path, tokenizer path, runtime library path, raw input text, and signed URLs;
- startup config rejects `ONNX_MAX_SEQUENCE_LENGTH` values above manifest `runtime.validated_max_sequence_length`;
- startup logs safe ONNX preflight metadata after manifest verification;
- Redis connection logs the effective cache namespace.

Not yet implemented:

- tokenizer JSON checksum preflight in Go startup;
- ONNX Runtime shared library sha256 preflight in Go startup;
- provider availability diagnostics;
- a richer health endpoint status object for failed preflight, because failed preflight exits before serving HTTP;
- staging rollout/rollback execution proof.

## Rollout stages

1. Local artifact verification: `tools/verify_onnx_artifacts.py` strict mode passes.
2. Local Docker packaging: `tools/build_onnx_image.sh` builds the image.
3. Packaged legal quality: M023-style legal evaluator passes with isolated cache namespace.
4. Packaged performance: M024-style benchmark passes with container-specific restart command.
5. Hosted artifact provisioning: immutable artifact source/cache exists and verifies checksums in CI.
6. Hosted full ONNX CI: manual workflow builds image, runs tagged tests, and optionally runs quality/performance jobs.
7. Staging opt-in: run ONNX behind explicit backend/env flag while TEI remains default rollback.
8. Production decision: only after evidence review, decide whether to route any production traffic to ONNX.

## Rollback contract

Rollback must be boring and immediate:

1. stop ONNX container/job;
2. set backend back to TEI/default or deploy the default image;
3. use a distinct Redis `EMBEDDING_CACHE_VERSION` for post-rollback checks if comparing outputs;
4. verify `/health` on the TEI/default endpoint;
5. rerun a smoke embedding request against TEI;
6. record rollback reason and last ONNX preflight/health error.

Do not delete verified artifacts during rollback unless the artifact itself is suspected corrupt; preserving them helps post-mortem diagnosis.

## Current open gaps

- Tokenizer JSON checksum preflight is not yet implemented in Go startup.
- ONNX Runtime shared library sha256/provider diagnostics are not yet implemented in Go startup.
- Hosted artifact provisioning/cache is not selected.
- Hosted full ONNX CI is not enabled as a required workflow.
- Security review for artifact path handling and signed URL handling remains future work.
- Production rollout and rollback have not been tested in staging.
