---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Check Go runtime prerequisites

Check local prerequisites for Go target-runtime smoke: ONNX artifact, native tokenizer, tokenizer JSON, ONNX Runtime shared library, Redis availability, and clean local ports. Do not start services yet except read-only probes.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`
- `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a`

## Expected Output

- `Task summary`

## Verification

Prerequisite check reports exact availability or blocker without leaking secrets.

## Observability Impact

Distinguishes real blockers from missing setup.
