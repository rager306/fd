---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: CoalescingEmbedder + 44-FZ baseline

CoalescingEmbedder + unit tests + 44-FZ baseline + wiring в main.go — всё реализовано в одном таске.

## Inputs

- None specified.

## Expected Output

- `api/embed/coalescedembedder.go`
- `api/embed/coalescedembedder_test.go`
- `api/embed/coalesce_baseline_test.go`
- `api/main.go`

## Verification

go test ./api/... -race
