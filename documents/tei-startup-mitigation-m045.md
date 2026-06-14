# M045 S02 TEI Startup Mitigation Selection

Captured: 2026-06-14T11:35:29Z

## Safety Boundary

This slice remains non-destructive so far. No TEI restart, compose up/down, recreate, or image run was performed while producing this inventory.

## Cache Mount

`fd_tei` mounts Docker volume `fd_tei_data` at `/data`. The symlink-aware inventory found 33 files or symlinks under `/data`; total blob scan size from the initial file-only pass was 1440196876 bytes.

## USER-bge-m3 Cached Files

The current model snapshot contains the files needed for Candle/safetensors startup:

| Path | Size bytes | Kind |
|---|---:|---|
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/config.json` | 697 | L |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/1_Pooling/config.json` | 297 | L |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/model.safetensors` | 1436151696 | L |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/tokenizer_config.json` | 1362 | L |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/modules.json` | 349 | L |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/tokenizer.json` | 3327728 | L |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/sentence_bert_config.json` | 54 | L |
| `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/config_sentence_transformers.json` | 195 | L |

Required-file counts from symlink-aware scan:

```json
{
  "config.json": 4,
  "tokenizer.json": 2,
  "tokenizer_config.json": 1,
  "special_tokens_map.json": 0,
  "safetensors": 1,
  "pooling_config": 2,
  "modules.json": 1
}
```

ONNX files found: `0`.

## Candidate Decision

Selected candidate for S03 controlled proof: add `HF_HUB_OFFLINE=1` to the TEI service environment.

Rationale:

- M045 S01 found TEI currently starts from Hub repo ID and enters a slow ORT/ONNX missing-file probe before Candle fallback.
- Hugging Face Hub docs state `HF_HUB_OFFLINE=1` prevents HTTP calls and uses only cached files.
- The mounted `/data` cache contains the required USER-bge-m3 safetensors, tokenizer, pooling, module, and config files.
- ONNX files are absent, which is expected and should remain true under fd's TEI-only posture.

Rejected options:

| Option | Reason |
|---|---|
| Add ONNX artifacts | Reintroduces ONNX operational artifact scope that M042 explicitly removed. |
| Change `--dtype` | Does not appear to control ORT probing and may alter numerical behavior. |
| Local model path immediately | More invasive than offline mode; keep as fallback if offline proof still probes remotely or fails. |
| Build or switch to a Candle-only TEI image | Larger supply-chain/build change; no existing documented image variant was proven in S01. |

## S03 Proof Plan

Before performing any restart proof:

1. Record current compose config and container start times.
2. Apply the compose candidate only through tracked config, not ad-hoc container mutation.
3. Capture TEI logs from process start until ready or timeout.
4. Measure time to TEI health, fd readiness, and first fd embedding.
5. Roll back by removing `HF_HUB_OFFLINE=1` if TEI fails because required files are not found.

Success criteria:

- TEI reaches healthy state substantially faster than the prior tens-of-minutes ORT download/probe path, or ONNX/ORT remote download warnings disappear/fail fast.
- fd `/health` still reports backend `tei`, model `deepvk/USER-bge-m3`, dimensions `1024`.
- fd `/v1/embeddings` returns a 1024-dimensional embedding.

## Current Recommendation

Proceed to prepare `docker-compose.yaml` with `HF_HUB_OFFLINE=1` for the TEI service. Do not restart in S02. S03 will be the controlled proof.

## Compose Candidate Prepared

`docker-compose.yaml` now includes the S03 candidate environment for the TEI service:

```text
## compose tei environment candidate
    environment:
      HF_HOME: /data
      HF_HUB_OFFLINE: "1"
      HUGGINGFACE_HUB_CACHE: /data
    healthcheck:
\n## current running TEI env subset
HF_HOME=/data
HUGGINGFACE_HUB_CACHE=/data
\n## current container started
/fd_tei running health=healthy started=2026-06-14T09:24:48.743364832Z
/fd_api running health=healthy started=2026-06-14T08:30:46.620802937Z
/fd_redis running health=healthy started=2026-05-19T18:08:09.5269274Z

```

Important: the running `fd_tei` container does **not** yet have `HF_HUB_OFFLINE=1`; this is expected because S02 did not restart or recreate TEI. S03 is the controlled proof slice.

## S02 Non Destructive Smoke

```json
{
  "fd_health": {
    "http_status": 200,
    "runtime": {
      "backend": "tei",
      "model": "deepvk/USER-bge-m3",
      "dimensions": 1024,
      "production_default": true,
      "cache_namespace": "v2"
    },
    "status": "ok"
  },
  "fd_ready": {
    "http_status": 200,
    "body": {
      "status": "ready",
      "time": "2026-06-14T11:37:00Z"
    }
  },
  "fd_embedding": {
    "http_status": 200,
    "latency_ms": 492.09,
    "embedding_len": 1024,
    "model": "deepvk/USER-bge-m3"
  },
  "tei_embedding": {
    "http_status": 200,
    "latency_ms": 396.17,
    "embedding_len": 1024
  }
}
```

S02 verification result: current runtime remains healthy and TEI-only while the candidate compose file is staged for a future controlled restart proof.

