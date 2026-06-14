---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T02: Applied local snapshot TEI command and measured successful startup.

Update compose and override so TEI uses `--model-id /data/models--deepvk--USER-bge-m3/snapshots/<revision>` with the cached local snapshot path. Run `docker compose up -d tei` to apply the command, then poll Docker health until healthy or timeout. Capture logs from the new container start window.

## Inputs

- `docker-compose.yaml`
- `docker-compose.override.yaml`

## Expected Output

- `benchmark-results/m045-tei-local-path-startup-proof.md`

## Verification

Artifact records compose command exit, start timestamp, healthy timestamp or timeout, and TEI logs. Logs should show local path model_id and avoid Hub ONNX download attempts if the theory is correct.

## Observability Impact

Measures whether local path avoids Hub/api_repo/ONNX probe delay.
