---
id: S02
parent: M045-d0e5xq
milestone: M045-d0e5xq
provides:
  - S03 candidate compose config.
  - Cache completeness evidence for offline proof.
requires:
  - slice: S01
    provides: Read-only recon and candidate list.
affects:
  []
key_files:
  - docker-compose.yaml
  - docs/same-host-embedding-service-contract.md
  - documents/tei-startup-mitigation-m045.md
key_decisions:
  - Use `HF_HUB_OFFLINE=1` as the S03 candidate because required USER-bge-m3 safetensors/tokenizer files are cached and ONNX artifacts remain absent.
patterns_established:
  - Stage startup config separately from applying it to running containers.
  - Prove TEI startup mitigations with controlled restart only after cache completeness evidence.
observability_surfaces:
  - documents/tei-startup-mitigation-m045.md records cache inventory, selected candidate, rollback plan, and smoke evidence.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-14T11:38:52.674Z
blocker_discovered: false
---

# S02: Select safe TEI startup mitigation

**Selected and staged `HF_HUB_OFFLINE=1` as the TEI startup mitigation candidate without restarting the running service.**

## What Happened

S02 inventoried the TEI `/data` cache and found the USER-bge-m3 snapshot includes config, pooling config, modules, tokenizer_config, tokenizer.json, sentence-transformers config, and model.safetensors, with zero ONNX files. Based on S01 source findings and Hugging Face Hub offline mode docs, S02 selected `HF_HUB_OFFLINE=1` as the smallest safe candidate to avoid slow remote probes for missing ONNX artifacts. `docker-compose.yaml` now stages `HF_HUB_OFFLINE=1` and explicit `HUGGINGFACE_HUB_CACHE=/data`; docs clarify this is startup mitigation only, not fd fallback behavior. The running container remains unchanged until S03 controlled proof.

## Verification

S02 UAT PASS via gsd_uat_exec: artifact `a25a25c0-371a-4a11-9de5-f30f7d330344`, compose candidate `226151ba-ead1-4f52-9fb5-d78c157cf60e`, fd runtime smoke `92f7bd31-41c2-447b-8091-814dd295bba9`, and running-container unchanged check `03edc3bc-0a5b-42fb-8798-a02ee267b323`. No restart/recreate occurred.

## Requirements Advanced

- R028 — Advanced by selecting and staging the startup mitigation candidate; validation remains pending controlled proof.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

The mitigation is staged but not proven. S03 must perform a controlled restart/startup proof to validate whether offline mode removes or bounds the ORT/ONNX probe delay.

## Follow-ups

S03 should execute the controlled restart proof with capture and rollback. If offline mode fails due missing cache metadata, rollback by removing `HF_HUB_OFFLINE=1` and document the TEI external limitation.

## Files Created/Modified

- `docker-compose.yaml` — Added TEI `HF_HUB_OFFLINE=1` and explicit `HUGGINGFACE_HUB_CACHE=/data` candidate env.
- `docs/same-host-embedding-service-contract.md` — Documented TEI offline cache mode as startup mitigation only.
- `documents/tei-startup-mitigation-m045.md` — New mitigation selection and cache inventory artifact.
