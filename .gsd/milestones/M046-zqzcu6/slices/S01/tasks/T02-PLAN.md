---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Validated key issue #3 P0/P1 defect signals with safe static probes.

For each P0/P1 cluster, use read-only code inspection and safe local tests or lightweight scripts to determine whether the issue exists in current code. Avoid destructive load tests; prefer unit-level or handler-level evidence.

## Inputs

- `api/main.go`
- `api/middleware`
- `api/handlers`
- `api/cache`

## Expected Output

- `benchmark-results/m046-s01-audit-validation.md`

## Verification

Safe probes or code evidence exist for batch guardrails, auth exposure posture, and LocalCache correctness. Commands do not mutate source files.

## Observability Impact

Captures objective validation evidence without triggering abusive load.
