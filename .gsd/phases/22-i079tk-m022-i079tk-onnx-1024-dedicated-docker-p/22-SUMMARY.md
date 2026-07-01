---
id: M022-i079tk
title: "ONNX 1024 dedicated Docker packaging"
status: complete
completed_at: 2026-05-20T10:49:11.720Z
key_decisions:
  - D020: default CI runs artifact-free ONNX contract checks; full ONNX image CI is deferred until external artifact provisioning/cache exists.
key_files:
  - Dockerfile.onnx
  - tools/build_onnx_image.sh
  - .github/workflows/go-quality.yml
  - docs/onnx-artifacts/README.md
  - .gsd/DECISIONS.md
lessons_learned:
  - A generated staging Docker context under `.gsd/runtime/` cleanly bridges ignored local artifacts into a local image build without widening the default Docker context.
  - Hosted CI can still enforce safety metadata and binary hygiene even when full binary artifacts are unavailable.
  - The ONNX image build context is large (1.569GB), so CI should use artifact caching/provisioning before enabling full image builds.
---

# M022-i079tk: ONNX 1024 dedicated Docker packaging

**M022 proved a dedicated opt-in ONNX Docker image locally and added truthful CI-safe artifact contract checks without changing TEI defaults.**

## What Happened

M022 advanced ONNX 1024 from an artifact contract to a dedicated runnable Docker packaging proof. S01 implemented `Dockerfile.onnx` and `tools/build_onnx_image.sh`, which verifies manifests, stages only required source and artifacts under `.gsd/runtime/docker/`, builds with `onnx hf_tokenizers`, and packages the ONNX model, tokenizer JSON, manifests, and ONNX Runtime shared library into a dedicated image. The image built and served `/health` plus a 1024-dimensional embedding smoke response. S02 added hosted-CI-safe artifact contract checks to the Go Quality workflow: manifest metadata validation in `--allow-missing` mode and a binary hygiene gate. D020 and README docs explicitly defer full ONNX image CI until an external artifact provisioning/cache mechanism exists. Default TEI Docker behavior remains untouched and verified.

## Success Criteria Results

- Dedicated image path: PASS.
- Artifact verifier in build path: PASS.
- Default Docker guardrail: PASS.
- ONNX runtime smoke: PASS.
- CI boundary: PASS.
- Binary hygiene: PASS.

## Definition of Done Results

- Dedicated opt-in ONNX Docker packaging path exists: met.
- Packaged ONNX build verifies artifacts before build: met.
- Default Docker build remains passing: met.
- ONNX packaging proof builds/runs locally: met.
- No binary artifacts tracked: met.
- TEI remains default and ONNX remains opt-in experimental: met.
- CI boundary is truthful and artifact-free: met.

## Requirement Outcomes

- ONNX local Docker packaging: validated.
- Default TEI Docker safety: validated.
- CI metadata/binary hygiene contract: validated.
- Full hosted ONNX image CI: deferred with clear artifact provisioning blocker.
- Production ONNX promotion: still blocked pending packaged legal/performance gates and artifact store/rollout decisions.

## Deviations

None. Full ONNX hosted CI remains intentionally deferred because external artifact provisioning/cache does not yet exist.

## Follow-ups

Next milestone should either select an artifact provisioning mechanism (release assets/object storage/cache with checksums) or run packaged legal quality and performance gates against the M022 ONNX image locally. Hosted CI full image build should wait for artifact provisioning.
