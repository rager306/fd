---
id: T01
parent: S01
milestone: M010-84qfzu
key_files:
  - docker-compose.yaml
  - docker-compose.override.yaml
  - .gitignore
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T18:24:49.189Z
blocker_discovered: false
---

# T01: Local TEI cache has the current safetensors/tokenizer artifacts and revision, but no ONNX; future ONNX output must stay in ignored runtime/artifact storage.

**Local TEI cache has the current safetensors/tokenizer artifacts and revision, but no ONNX; future ONNX output must stay in ignored runtime/artifact storage.**

## What Happened

Inspected local model artifact storage. TEI stores HuggingFace cache under the Docker volume `fd_tei_data`, mounted at `/data` in `fd_tei`. The local `deepvk/USER-bge-m3` revision is `0cc6cfe48e260fb0474c753087a69369e88709ae`. Available artifacts include `model.safetensors` (~1.44GB), `tokenizer.json`, `config.json`, `modules.json`, `sentence_bert_config.json`, and pooling config. No ONNX file is present in the local snapshot. `.gitignore` already excludes `tei-models/`, but does not yet explicitly exclude a general ONNX artifact directory; future S03 should use an untracked local path such as `.gsd/runtime/onnx/` or add a dedicated ignored `onnx-artifacts/` before downloading/exporting large files. Redis and TEI named volumes are `fd_redis_data` and `fd_tei_data`; the stack is currently healthy.

## Verification

Read Compose/gitignore and inspected Docker volume/container artifact paths and hashes without staging large files.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `TEI volume: fd_tei_data mounted at /data` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Model revision: 0cc6cfe48e260fb0474c753087a69369e88709ae` | -1 | unknown (coerced from string) | 0ms |
| 3 | `model.safetensors size: 1436151696 bytes, sha256 e6aa9c8e51a60ff383186a2f28f658305ba4ad23d2fa24296607885458ef2733` | -1 | unknown (coerced from string) | 0ms |
| 4 | `tokenizer.json sha256 068d9f7ed9dd190a00a567e5f7750fdc591b93bc623072ac8050a303c25f5937` | -1 | unknown (coerced from string) | 0ms |
| 5 | `docker compose ps reported api/redis/tei healthy` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

TEI's local cache contains PyTorch/safetensors and tokenizer artifacts, not ONNX. The model safetensors artifact is large (~1.44GB), so future export output must stay outside git-tracked paths.

## Files Created/Modified

- `docker-compose.yaml`
- `docker-compose.override.yaml`
- `.gitignore`
