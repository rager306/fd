# M045 TEI Startup Mitigation Outcome

Updated: 2026-06-14T12:16:00Z

## Summary

The validated startup mitigation is to run HuggingFace TEI with the cached local USER-bge-m3 snapshot path as `--model-id`, not the Hub repo ID.

Validated command:

```text
text-embeddings-router --model-id /data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae --max-batch-tokens 8192
```

Why this works:

- Official TEI CLI docs say `--model-id` may be a local directory containing model files saved by Transformers or Sentence Transformers.
- Indexed TEI source shows local paths take the local-model branch and set `api_repo=None`; Hub IDs use `ApiBuilder`, `download_artifacts`, and later allow the ORT/ONNX probe path.
- The cached `/data` snapshot contains the required USER-bge-m3 safetensors/tokenizer/config files.
- The local-path proof reached TEI healthy and passed fd/direct TEI embedding smoke.

## Cache Inventory

`fd_tei` mounts Docker volume `fd_tei_data` at `/data`. Symlink-aware inventory found the current model snapshot at:

```text
/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae
```

Important files present:

| Path | Size bytes | Kind |
|---|---:|---|
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/config.json` | 697 | symlink |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/1_Pooling/config.json` | 297 | symlink |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/model.safetensors` | 1436151696 | symlink |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/modules.json` | 349 | symlink |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/tokenizer_config.json` | 1362 | symlink |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/tokenizer.json` | 3327728 | symlink |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/sentence_bert_config.json` | 54 | symlink |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/config_sentence_transformers.json` | 195 | symlink |

ONNX files found: `0`.

## Candidate Outcomes

| Candidate | Outcome | Reason |
|---|---|---|
| Local snapshot path as `--model-id` | Validated | Reached TEI healthy in the controlled proof and avoids the Hub `api_repo` path. |
| `HF_HUB_OFFLINE=1` with Hub ID | Rejected | TEI stayed unhealthy after the 15-minute proof timeout and still entered `Downloading onnx/model.onnx`. |
| Add ONNX artifacts | Rejected | Reintroduces ONNX operational artifact scope removed in M042. |
| Change `--dtype` | Rejected | Does not control backend selection and may alter numerical behavior. |
| Build/switch Candle-only image | Future fallback | Larger supply-chain/build change; unnecessary after local-path proof passed. |

## Proof Result

Local-path proof artifact: `benchmark-results/m045-tei-local-path-startup-proof.md`.

Result summary:

- Container command: `--model-id /data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae --max-batch-tokens 8192`.
- TEI container started: `2026-06-14T12:12:14Z`.
- TEI became healthy: `2026-06-14T12:15:15Z`.
- Approximate time to healthy from container start: about 3 minutes.
- fd `/health`: HTTP 200, backend `tei`, model `deepvk/USER-bge-m3`, dimensions `1024`.
- fd `/v1/embeddings`: HTTP 200, 1024-dimensional embedding.
- direct TEI `/embeddings`: HTTP 200, 1024-dimensional embedding.

TEI still logs an immediate local ORT failure because ONNX files are intentionally absent:

```text
Could not start ORT backend: File at `/data/.../onnx/model.onnx` does not exist
```

That failure is now local and fast. The prior long remote Hub probe/download path is avoided by using the local snapshot directory.

## Operational Guidance

- Keep the TEI Docker volume `tei_data` populated before using the local path command.
- If the model snapshot changes, update the compose command to the new snapshot directory and run a controlled proof.
- Do not add ONNX artifacts for this fd product path; ONNX remains future research only.
- If TEI fails to start with the local path, rollback to the previous Hub ID command only as a temporary recovery path and expect slower startup behavior.
