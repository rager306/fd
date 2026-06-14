---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Inventoried TEI `/data` cache and confirmed USER-bge-m3 Candle files are present.

Read-only inspect `/data` inside the running TEI container and host mounted cache paths if discoverable. Record required Candle/tokenizer/config/safetensors files, ONNX absence, file sizes, and whether `HF_HUB_OFFLINE=1` is likely safe to attempt. Do not print secrets or file contents.

## Inputs

- `documents/tei-startup-recon-m045.md`
- `docker inspect fd_tei`

## Expected Output

- `documents/tei-startup-mitigation-m045.md`

## Verification

Artifact lists cache files needed for offline Candle startup and states whether the cache appears complete.

## Observability Impact

Reduces risk of offline mode failing due incomplete cache.
