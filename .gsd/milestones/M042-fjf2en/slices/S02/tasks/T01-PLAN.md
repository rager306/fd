---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Mapped active ONNX source/config/docs/tooling surfaces and defined the TEI-only removal boundary.

Map active ONNX references across source, tests, Dockerfile, compose, docs, tools, CI, and requirements. Distinguish historical research artifacts from active runtime/build surfaces. Produce `documents/onnx-deactivation-inventory-m042.md` with a remove/keep table and risk notes. Do not edit code in this task.

## Inputs

- `documents/te-perf-root-cause-m042.md`
- `.gsd/REQUIREMENTS.md`
- `api/go.mod`
- `api/main.go`
- `api/embed/`

## Expected Output

- `documents/onnx-deactivation-inventory-m042.md`

## Verification

Inventory artifact exists and lists active source/config/docs surfaces plus explicit keep/remove decisions.

## Observability Impact

Clarifies which runtime self-description fields/docs must change.
