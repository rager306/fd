# M050-rfqm1p: Current test truth and stronger gates

**Vision:** Сделать тестовую систему fd честной относительно текущей версии сервиса: сначала выявить и исправить устаревшие существующие тесты, затем оформить реальный Docker Compose e2e suite и начальный mutation baseline для критичных backend-пакетов.

## Success Criteria

- Все существующие тесты классифицированы как актуальные, исправленные, выведенные из регулярного запуска или явно отложенные с причиной.
- `go test ./...` в `api/` проходит после ревизии тестов и не содержит устаревших ожиданий старого API/auth/runtime.
- Реальный Docker Compose e2e suite запускается одной командой и проверяет основные HTTP/runtime/cache/metrics потоки текущего сервиса.
- Mutation-testing baseline существует для выбранных критичных пакетов или документированно ограничен техническим blocker с воспроизводимым evidence.
- CI или документация ясно объясняет, какие тестовые уровни регулярные, какие требуют Docker/секретов, и как запускать каждый уровень.

## Slices

- [x] **S01: Existing test actuality audit** `risk:high` `depends:[]`
  > After this: После этого есть проверенный инвентарь всех текущих тестов и исправленные stale tests; `api` suite зелёный на текущем контракте.

- [x] **S02: Docker e2e suite for current service** `risk:high` `depends:[S01]`
  > After this: После этого одна команда запускает auth-aware black-box проверку реального Compose runtime.

- [x] **S03: Mutation baseline for critical packages** `risk:medium` `depends:[S01]`
  > After this: После этого есть измеренный mutation baseline или честный documented blocker для выбранного Go mutation runner.

- [x] **S04: Test gates documentation and closure** `risk:medium` `depends:[S02,S03]`
  > After this: После этого будущий агент понимает, какие test levels запускать, когда и почему.

## Boundary Map

| Boundary | In scope | Out of scope |
|---|---|---|
| Existing tests | Audit and fix stale expectations | Large feature rewrites unrelated to tests |
| E2E | Current Docker Compose fd API plus Redis plus TEI | External hosted environments |
| Mutation | Bounded critical package baseline | Immediate repo-wide mandatory mutation gate if too slow |
| CI/docs | Clear commands and optional/manual heavy gates | Pushing remote changes without explicit permission |
