---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T01: Confirmed issue #7 runtime contract debt exists before S02 fixes.

Run static checks proving issue #7 #26/#29/#30 are present before fixes: ONNX-only RuntimeHealth fields, duplicate interfaces, default singleton.

## Inputs

- `documents/issue-7-current-m048.md`

## Expected Output

- `api/handlers/health.go`
- `api/handlers/embeddings.go`
- `api/lifecycle/warmup.go`
- `api/lifecycle/state.go`

## Verification

Static gsd_exec check should PASS for pre-fix presence.

## Observability Impact

Documents exact pre-fix contract surfaces.
