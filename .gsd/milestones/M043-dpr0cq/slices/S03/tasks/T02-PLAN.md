---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: CI govulncheck step added after golangci-lint in go-quality workflow

.github/workflows/go-quality.yml: добавить step "Run govulncheck" между "Run GolangCI-Lint" и финальным job completion. Использовать `go run golang.org/x/vuln/cmd/govulncheck@latest` для reproducibility (не требует manual install). Поместить после lint step чтобы vuln scan запускается только если lint passed. Failure on any reported vuln. Verify CI workflow YAML valid.

## Inputs

- None specified.

## Expected Output

- `.github/workflows/go-quality.yml (govulncheck step added)`

## Verification

CI workflow YAML valid (parseable). Step добавлен в правильное место. В локальном dry-run (если возможно) показывает step.
