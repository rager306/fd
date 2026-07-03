# S03: Controlled startup proof — UAT

**Milestone:** M045-d0e5xq
**Written:** 2026-06-14T12:50:13.274Z

# S03: Controlled startup proof — UAT

**Milestone:** M045-d0e5xq
**Written:** 2026-06-14

## UAT Type

- UAT mode: mixed
- Why this mode is sufficient: This backend runtime slice has no interactive product UI, but it must prove both live runtime behavior and the localhost HTTP surface. Runtime checks verify container command, TEI health, fd health, fd embedding smoke, and effective compose config; browser checks verify the localhost `/health` JSON surface exposed to operators.

## Preconditions

- `fd_tei` is running from the updated compose command.
- `fd_api` is running on `http://localhost:8000`.
- The cached USER-bge-m3 snapshot exists at `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae`.

## Smoke Test

Open `http://localhost:8000/health` and verify the visible JSON contains `"backend":"tei"`, `"model":"deepvk/USER-bge-m3"`, and `"dimensions":1024`.

## Test Cases

### 1. TEI container uses local snapshot path

1. Inspect the running `fd_tei` container command.
2. Confirm the command contains `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae`.
3. Confirm container health is `healthy`.
4. **Expected:** TEI is healthy and uses the local snapshot path instead of the Hub ID.

### 2. fd runtime identity remains stable

1. Request `http://localhost:8000/health`.
2. Confirm runtime backend is `tei`.
3. Confirm runtime model is `deepvk/USER-bge-m3`.
4. Confirm dimensions are `1024`.
5. **Expected:** fd reports the same public runtime contract while TEI uses a local model path internally.

### 3. Embedding smoke succeeds

1. Send a `/v1/embeddings` request through fd.
2. Confirm the returned vector length is `1024`.
3. **Expected:** fd embedding clients still receive valid 1024-dimensional embeddings.

### 4. Rejected offline candidate is absent

1. Render effective compose config for `tei`.
2. Confirm the local snapshot path is present.
3. Confirm `HF_HUB_OFFLINE` is absent.
4. **Expected:** effective runtime config uses the validated local path mitigation and does not keep the failed offline-env candidate.

## Edge Cases

### Missing ONNX artifact warning does not block startup

1. Inspect controlled proof logs.
2. Confirm TEI can log a local ORT/ONNX miss for missing `onnx/model.onnx`.
3. Confirm Candle backend warms up and TEI becomes healthy.
4. **Expected:** local ONNX miss remains a warning/diagnostic signal, not a blocking remote Hub probe.

## Failure Signals

- `fd_tei` health is not `healthy`.
- `fd_tei` command contains `deepvk/USER-bge-m3` as Hub ID instead of the local `/data/.../snapshots/...` path.
- `docker compose config tei` contains `HF_HUB_OFFLINE`.
- fd `/health` does not expose backend `tei`, model `deepvk/USER-bge-m3`, dimensions `1024`.
- fd `/v1/embeddings` does not return a 1024-dimensional vector.

## Requirements Proved By This UAT

- R028 — TEI startup is bounded by a controlled local snapshot path proof and no longer relies on Hub ID startup resolution.

## Not Proven By This UAT

- It does not prove a generic solution for arbitrary future snapshot revisions.
- It does not remove TEI's local ONNX probing log; it proves the probe no longer blocks startup through remote Hub resolution.

## Notes for Tester

Evidence:
- `benchmark-results/m045-tei-local-path-startup-proof.md`
- Browser timeline: `/root/fd/.artifacts/browser/2026-06-14T12-23-36-082Z-session/m045-s03-browser-health-timeline.json`
- gsd_uat_exec IDs: `335f897b-97a6-4a0d-bb21-44e2df9fc8cc`, `a55d4583-cd51-4b8f-943a-b56c1a045b1e`, `708de145-0232-4cc6-9346-75b837562a25`, `298cc0b1-45c5-48cd-9748-eca5e1854320`.
