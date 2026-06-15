---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T03: Исправлен stale root integration test layer под текущий auth/runtime контракт.

Исправить тесты, которые должны оставаться актуальными, особенно root integration tests и auth/runtime expectations. Если тест не должен запускаться регулярно, вывести его в documented manual/legacy status вместо молчаливого broken path. Не добавлять новый e2e suite в этом task сверх минимальных исправлений актуальности.

## Inputs

- `T02 stale list`

## Expected Output

- `tests/integration/api_test.go`
- `benchmark-results/m050-s01-test-actuality.md`

## Verification

`cd api && go test ./...` passes; any touched external integration test has a runnable documented command or explicit deferred reason.

## Observability Impact

Убирает ложную уверенность от устаревших тестов.
