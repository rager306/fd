---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T03: Integrated LocalCache shutdown, ran gates, wrote evidence, and validated R031.

Close LocalCache in main shutdown/error paths. Run full Go suite, cache race tests, lint, govulncheck, and static proof that LocalCache no longer uses sync.Map or separate size counter. Write S05 evidence artifact and update remediation plan/R031.

## Inputs

- `api/cache/local.go`
- `api/main.go`

## Expected Output

- `benchmark-results/m046-s05-localcache-correctness.md`
- `documents/issue-3-audit-remediation-plan-m046.md`

## Verification

cd api && go test ./... && cd api && go test -race ./cache -run TestLocalCache && lint && govulncheck

## Observability Impact

Evidence artifact captures race/gate results.
