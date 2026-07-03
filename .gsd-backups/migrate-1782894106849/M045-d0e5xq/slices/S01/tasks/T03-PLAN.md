---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Identified safe TEI startup mitigation candidates from source/docs without destructive commands.

Use current TEI image help/docs and public documentation lookup to identify flags/env/model layout options that may force Candle, skip ORT, preload safetensors, set revision, or otherwise avoid ONNX probing. Prefer non-destructive documentation and help inspection. Do not run blocked destructive docker commands.

## Inputs

- `docker-compose.yaml`
- `TEI image metadata`
- `HuggingFace TEI docs`

## Expected Output

- `documents/tei-startup-recon-m045.md`

## Verification

Artifact includes candidate options with source, expected effect, risk, and whether it requires restart proof.

## Observability Impact

Makes mitigation decision auditable.
