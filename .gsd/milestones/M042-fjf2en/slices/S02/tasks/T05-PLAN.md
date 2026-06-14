---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T05: Final TEI-only gates and milestone scope decision

Run mandatory M043 gates and final TEI-only checks: `go test ./...`, golangci-lint v2.12.2, govulncheck, and a small runtime/config smoke if Docker service is healthy. Record whether R021 async chunking is deferred or implemented separately. Validate R027. Write final S02 evidence artifacts.

## Inputs

- `api/`
- `documents/onnx-deactivation-inventory-m042.md`

## Expected Output

- `benchmark-results/m042-s02-go-test.txt`
- `benchmark-results/m042-s02-lint.txt`
- `benchmark-results/m042-s02-govulncheck.txt`
- `benchmark-results/m042-s02-tei-only-check.txt`

## Verification

All mandatory gates pass; R027 validated; R021 either validated with evidence or deferred with explicit rationale.

## Observability Impact

Final evidence proves TEI-only active posture.
