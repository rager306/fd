---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Выбран рабочий Go mutation runner и подтверждён smoke-запуск.

Проверить доступные Go mutation runners без изменения production code. Начать с bounded help/smoke commands and choose one runner or document blocker.

## Inputs

- `api/go.mod`
- `current Go version`
- `S01 coverage findings`

## Expected Output

- `benchmark-results/m050-s03-mutation-baseline.md`

## Verification

Artifact records runner candidates, chosen command or blocker with exact output.

## Observability Impact

Не даёт заявить mutation coverage без работающего runner.
