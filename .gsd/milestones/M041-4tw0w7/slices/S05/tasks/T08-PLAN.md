---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T08: Final 45-test acceptance suite и MILESTONE-UAT

tools/verify_fd_v2_contract.py: автоматизировать ВСЕ 45 test cases. Скрипт запускает каждый test case против running fd v2, проверяет HTTP status, body shape, headers, и пишет results в JSON. Final artifact: benchmark-results/fd-v2-validation-m041.md со всеми 45 test results, p95 perf numbers, и pass/fail summary. Если хоть 1 test fail — exit 1.

## Inputs

- None specified.

## Expected Output

- `tools/verify_fd_v2_contract.py`
- `benchmark-results/fd-v2-validation-m041.md`

## Verification

tools/verify_fd_v2_contract.py exit 0 на running fd v2, все 45 test cases pass. Artifact содержит pass/fail breakdown.
