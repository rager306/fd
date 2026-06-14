# M045 S01 TEI Startup Recon

Captured: 2026-06-14T11:23:49Z

## Safety Boundary

This recon was read-only. No `docker restart`, `docker compose restart`, `docker compose up`, `docker compose down`, `docker run`, or container recreate was performed.

## Current Runtime State

| Container | Image | Status | Health | Command | Safe env subset |
|---|---|---|---|---|---|
| fd_tei | ghcr.io/huggingface/text-embeddings-inference:cpu-1.9 | running | healthy | `text-embeddings-router --model-id deepvk/USER-bge-m3` | `HF_HOME=/data, HUGGINGFACE_HUB_CACHE=/data, PORT=80` |
| fd_api | fd-api | running | healthy | `/api` | `LOG_LEVEL=info, MODEL_ID=deepvk/USER-bge-m3, OTEL_SERVICE_NAME=fd-api, PORT=8000` |
| fd_redis | redis/redis-stack:latest | running | healthy | `redis-server --maxmemory 2gb --maxmemory-policy allkeys-lru --save 300 1 --appendonly no --protected-mode no` | `` |

TEI image digest observed by `docker inspect`: `sha256:ad950d30878eceb72aaf32024d26fa2b1d04a75304fa0b4776b49aa1941fea07`.
TEI started at `2026-06-14T09:24:48.743364832Z`.

## fd and TEI Smoke Evidence

fd `/health` compact result:

```json
{
  "http_status": 200,
  "body": {
    "status": "ok",
    "time": "2026-06-14T11:23:48Z",
    "model_loaded": true,
    "warmup_done": true,
    "device": "cpu",
    "last_inference_at": "2026-06-14T11:18:32Z",
    "in_flight_requests": 0,
    "runtime": {
      "backend": "tei",
      "model": "deepvk/USER-bge-m3",
      "dimensions": 1024,
      "production_default": true,
      "cache_namespace": "v2"
    }
  }
}
```

fd `/ready` compact result:

```json
{
  "http_status": 200,
  "body": {
    "status": "ready",
    "time": "2026-06-14T11:23:48Z"
  }
}
```

fd `/v1/embeddings` cache-miss smoke:

```json
{
  "http_status": 200,
  "latency_ms": 428.51,
  "model": "deepvk/USER-bge-m3",
  "data_len": 1,
  "embedding_len": 1024,
  "usage": {
    "prompt_tokens": 14,
    "total_tokens": 14
  }
}
```

Direct TEI `/embeddings` smoke:

```json
{
  "http_status": 200,
  "latency_ms": 289.14,
  "object": "list",
  "data_len": 1,
  "embedding_len": 1024
}
```

## Startup Log Timeline

The current container log tail still contains startup events from `fd_tei` start. Key lines:

```
fd_tei  | 2026-06-14T09:24:48.969199Z  INFO text_embeddings_router: router/src/main.rs:216: Args { model_id: "dee***/****-**e-m3", revision: None, tokenization_workers: None, dtype: None, served_model_name: None, pooling: None, max_concurrent_requests: 512, max_batch_tokens: 16384, max_batch_requests: None, max_client_batch_size: 32, auto_truncate: true, default_prompt_name: None, default_prompt: None, dense_path: None, hf_api_token: None, hf_token: None, hostname: "a97365ebc54d", port: 80, uds_path: "/tmp/text-embeddings-inference-server", huggingface_hub_cache: Some("/data"), payload_limit: 2000000, api_key: None, json_output: false, disable_spans: false, otlp_endpoint: None, otlp_service_name: "text-embeddings-inference.server", prometheus_port: 9000, cors_allow_origin: None }
fd_tei  | 2026-06-14T09:24:49.151051Z  INFO text_embeddings_router: router/src/lib.rs:271: Starting model backend
fd_tei  | 2026-06-14T09:42:51.735476Z  WARN text_embeddings_backend: backends/src/lib.rs:710: Could not download `onnx/model.onnx`: request error: error sending request for url (https://huggingface.co/deepvk/USER-bge-m3/resolve/main/onnx/model.onnx)
fd_tei  | 2026-06-14T09:42:51.735512Z  INFO text_embeddings_backend: backends/src/lib.rs:711: Downloading `model.onnx`
fd_tei  | 2026-06-14T09:51:52.407388Z  WARN text_embeddings_backend: backends/src/lib.rs:715: Could not download `model.onnx`: request error: error sending request for url (https://huggingface.co/deepvk/USER-bge-m3/resolve/main/model.onnx)
fd_tei  | 2026-06-14T10:00:53.079400Z  WARN text_embeddings_backend: backends/src/lib.rs:724: Could not download `onnx/model.onnx_data`: request error: error sending request for url (https://huggingface.co/deepvk/USER-bge-m3/resolve/main/onnx/model.onnx_data)
fd_tei  | 2026-06-14T10:00:53.079427Z  INFO text_embeddings_backend: backends/src/lib.rs:725: Downloading `model.onnx_data`
fd_tei  | 2026-06-14T10:09:53.751411Z  WARN text_embeddings_backend: backends/src/lib.rs:729: Could not download `model.onnx_data`: request error: error sending request for url (https://huggingface.co/deepvk/USER-bge-m3/resolve/main/model.onnx_data)
fd_tei  | 2026-06-14T10:09:53.751455Z ERROR text_embeddings_backend: backends/src/lib.rs:388: Model ONNX files not found in the repository. You can easily create ONNX files using the following scripts: https://gist.github.com/tomaarsen/4b00b0e3be8884efa64cfab9230b161f, or use this Space: https://huggingface.co/spaces/sentence-transformers/backend-export
fd_tei  | 2026-06-14T10:09:53.752185Z ERROR text_embeddings_backend: backends/src/lib.rs:412: Could not start ORT backend: Could not start backend: File at `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae/onnx/model.onnx` does not exist
fd_tei  | 2026-06-14T10:09:53.752198Z  INFO text_embeddings_backend: backends/src/lib.rs:595: Downloading `model.safetensors`
fd_tei  | 2026-06-14T10:09:55.334030Z  INFO text_embeddings_router: router/src/lib.rs:289: Warming up model
fd_tei  | 2026-06-14T10:12:58.805445Z  INFO text_embeddings_router::http::server: router/src/http/server.rs:1881: Ready
```

Observed pattern:

1. TEI starts as `text-embeddings-router --model-id deepvk/USER-bge-m3` with `dtype: None` and Hugging Face cache under `/data`.
2. TEI downloads or reads safetensors first, then enters an ORT/ONNX probe sequence for `onnx/model.onnx`, `model.onnx`, `onnx/model.onnx_data`, and `model.onnx_data`.
3. Missing ONNX files produce warnings/errors and `Could not start ORT`.
4. TEI then falls back to Candle/safetensors, warms the model, and reaches ready state.

## Upstream Source Findings

Sources inspected:

- `https://raw.githubusercontent.com/huggingface/text-embeddings-inference/main/router/src/main.rs`
- `https://raw.githubusercontent.com/huggingface/text-embeddings-inference/main/backends/src/lib.rs`
- `https://huggingface.co/docs/huggingface_hub/en/package_reference/environment_variables`

Relevant findings:

- CLI args include `--revision`, `--tokenization-workers`, `--dtype`, `--served-model-name`, `--pooling`, and `--dense-path`.
- No documented CLI arg named like `--backend`, `--disable-ort`, `--disable-onnx`, or `--force-candle` was found in the inspected router args.
- Backend init source has an `#[cfg(feature = "ort")]` section that attempts ONNX download/start when an API repo is present, before Candle fallback.
- Hugging Face Hub supports `HF_HUB_OFFLINE=1`: when set, no HTTP calls are made and only cached files are used. This may turn slow missing-file network probes into fast cache misses if all required Candle files are already cached.

## Candidate Mitigations for S02

| Candidate | Expected effect | Risk | S02 recommendation |
|---|---|---|---|
| `HF_HUB_OFFLINE=1` with existing `/data` cache | Avoid remote HTTP calls for missing ONNX artifacts; TEI should use cached safetensors/tokenizer files or fail fast if cache is incomplete. | If cache lacks required Candle files, TEI may fail startup. Requires controlled restart proof. | Strongest candidate to test first, after confirming cached required files exist. |
| Use a local model path as `--model-id` instead of Hub repo ID | May avoid `api_repo` remote download path and reduce ONNX probing against Hugging Face Hub. | Requires a complete local model directory layout and compose change; could break model discovery if layout differs. | Secondary candidate if offline mode is insufficient. |
| Add ONNX files to the model cache | Would satisfy ORT probe instead of failing it. | Reintroduces ONNX operational artifact scope that M042 explicitly removed from fd active product path. | Reject for current fd scope. |
| Switch TEI image or build a Candle-only TEI image | Could remove ORT feature entirely if such an image/build is available. | Larger supply-chain/build change; no existing documented image variant found in S01. | Future fallback only if simple env/config fails. |
| `--dtype` changes | May affect model dtype/backend behavior but does not appear to control ORT probing. | Could change numerical behavior and was already removed from compose because it was not active/effective. | Do not use as startup mitigation without stronger evidence. |

## Recommendation for S02

Proceed with a cautious offline-cache mitigation design:

1. Inventory `/data` cached files needed for Candle startup without printing secrets.
2. If required model/tokenizer/config/safetensors files are present, prepare a compose candidate with `HF_HUB_OFFLINE=1` for TEI.
3. Before any destructive proof, define capture and rollback commands, expected maximum wait, and success criteria.
4. If offline-cache proof is not safe, document the external TEI limitation rather than reintroducing fd ONNX runtime scope.

## Verification

- Current containers are running and healthy.
- fd `/health`, fd `/ready`, fd `/v1/embeddings`, and direct TEI `/embeddings` all returned HTTP 200.
- fd runtime remains TEI-only: backend `tei`, model `deepvk/USER-bge-m3`, dimensions `1024`, cache namespace `v2`.
- No restart or recreate was performed in S01.
