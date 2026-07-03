# S01: Existing test actuality audit

**Goal:** Проверить все существующие тесты и verification scripts на актуальность последней версии сервиса, исправить stale expectations и зафиксировать текущий честный baseline.
**Demo:** После этого есть проверенный инвентарь всех текущих тестов и исправленные stale tests; `api` suite зелёный на текущем контракте.

## Must-Haves

- Инвентарь покрывает `api/**/*_test.go`, `tests/integration`, `tools/verify*.py`, relevant CI workflow commands.
- Каждый элемент классифицирован как current, fixed, stale-deferred или obsolete with reason.
- Исправлены тесты, которые должны оставаться частью регулярной проверки текущего сервиса.
- `cd api && go test ./...` проходит после исправлений.
- Evidence сохранён в `benchmark-results/m050-s01-test-actuality.md`.

## Proof Level

- This slice proves: static inventory plus fresh command execution

## Integration Closure

S02 может начинаться только после подтверждения, что существующий baseline не содержит молча устаревших тестов.

## Verification

- Фиксирует drift между тестами и текущим runtime, чтобы будущие e2e/mutation gates строились на актуальных контрактах.

## Tasks

- [x] **T01: Собран инвентарь существующих тестов и verification scripts.** `est:45m`
  Собрать полный список текущих Go tests, root integration tests, verification Python scripts and CI test commands. Классифицировать их по уровню: unit, in-process integration, live runtime, static contract, perf or artifact verifier. Не менять код в этом task.
  - Files: `api/**/*_test.go`, `tests/integration/api_test.go`, `tools/verify*.py`, `.github/workflows/go-quality.yml`, `README.md`
  - Verify: Artifact contains complete inventory counts and categories; no code changes required.

- [x] **T02: Запущены текущие test commands и выявлены устаревшие ожидания root integration layer.** `est:60m`
  Запустить актуальные test commands: `cd api && go test ./...`; попробовать root integration path and verification scripts in safe dry or help mode where possible. Сравнить failures/skips с текущими контрактами auth/runtime/cache. Зафиксировать stale vs legitimate failures.
  - Files: `benchmark-results/m050-s01-test-actuality.md`
  - Verify: Fresh command outputs recorded with exit codes; stale expectations identified with file references.

- [x] **T03: Исправлен stale root integration test layer под текущий auth/runtime контракт.** `est:90m`
  Исправить тесты, которые должны оставаться актуальными, особенно root integration tests и auth/runtime expectations. Если тест не должен запускаться регулярно, вывести его в documented manual/legacy status вместо молчаливого broken path. Не добавлять новый e2e suite в этом task сверх минимальных исправлений актуальности.
  - Files: `tests/integration/api_test.go`, `api/**/*_test.go`, `tools/verify*.py`, `README.md`
  - Verify: `cd api && go test ./...` passes; any touched external integration test has a runnable documented command or explicit deferred reason.

- [x] **T04: S01 baseline artifact and requirement evidence recorded.** `est:30m`
  Сохранить финальный artifact по actuality audit, включая команды, результаты, классификацию, исправления и deferred items. Обновить requirement R043 evidence and complete the slice if verification passes.
  - Files: `benchmark-results/m050-s01-test-actuality.md`, `.gsd/REQUIREMENTS.md`
  - Verify: Final artifact exists; `cd api && go test ./...` fresh pass; R043 has validation notes.

## Files Likely Touched

- api/**/*_test.go
- tests/integration/api_test.go
- tools/verify*.py
- .github/workflows/go-quality.yml
- README.md
- benchmark-results/m050-s01-test-actuality.md
- .gsd/REQUIREMENTS.md
