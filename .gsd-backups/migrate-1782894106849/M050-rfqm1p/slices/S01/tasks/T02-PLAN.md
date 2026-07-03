---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Запущены текущие test commands и выявлены устаревшие ожидания root integration layer.

Запустить актуальные test commands: `cd api && go test ./...`; попробовать root integration path and verification scripts in safe dry or help mode where possible. Сравнить failures/skips с текущими контрактами auth/runtime/cache. Зафиксировать stale vs legitimate failures.

## Inputs

- `T01 inventory`

## Expected Output

- `benchmark-results/m050-s01-test-actuality.md`

## Verification

Fresh command outputs recorded with exit codes; stale expectations identified with file references.

## Observability Impact

Отделяет настоящие failures от тестов, которые больше не соответствуют продукту.
