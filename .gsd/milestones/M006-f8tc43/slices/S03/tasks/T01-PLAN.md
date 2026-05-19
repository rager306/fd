---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Add GolangCI-Lint config

Add root GolangCI-Lint config with Staticcheck and common analyzers suitable for the existing Go module under api/.

## Inputs

- `api/go.mod`

## Expected Output

- `.golangci.yml`

## Verification

Config file exists and references staticcheck.

## Observability Impact

Defines repeatable lint gate.
