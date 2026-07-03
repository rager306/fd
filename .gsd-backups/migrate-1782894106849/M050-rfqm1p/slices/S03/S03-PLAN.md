# S03: Mutation baseline for critical packages

**Goal:** Ввести начальный bounded mutation-testing baseline для критичных Go backend-пакетов или зафиксировать воспроизводимый blocker с альтернативным планом.
**Demo:** После этого есть измеренный mutation baseline или честный documented blocker для выбранного Go mutation runner.

## Must-Haves

- Mutation runner выбран после smoke/probe и причина выбора записана.
- Запуск ограничен критичными пакетами или файлами, чтобы быть воспроизводимым.
- Artifact содержит command, duration, pass/fail/mutation score or blocker, and survivors/follow-ups where available.
- `cd api && go test ./...` остаётся зелёным после любых test adjustments.
- R045 обновлён evidence.

## Proof Level

- This slice proves: mutation execution or reproducible blocker evidence

## Integration Closure

Mutation baseline дополняет S01/S02; если runner непригоден, milestone не притворяется, что mutation coverage есть.

## Verification

- Показывает силу assertions и surviving behavioral changes, которые statement coverage не видит.

## Tasks

- [x] **T01: Выбран рабочий Go mutation runner и подтверждён smoke-запуск.** `est:45m`
  Проверить доступные Go mutation runners без изменения production code. Начать с bounded help/smoke commands and choose one runner or document blocker.
  - Files: `benchmark-results/m050-s03-mutation-baseline.md`
  - Verify: Artifact records runner candidates, chosen command or blocker with exact output.

- [x] **T02: Bounded mutation baseline прошёл на критичных cache, handlers и lifecycle файлах.** `est:90m`
  Запустить выбранный mutation runner на одном или нескольких критичных пакетах with timeout. Prefer cache/handlers/lifecycle bounded scope. Capture score/survivors or exact failure.
  - Files: `api/cache`, `api/handlers`, `api/lifecycle`, `benchmark-results/m050-s03-mutation-baseline.md`
  - Verify: Mutation command exits successfully with baseline or fails with documented tool blocker; no production code changes required.

- [x] **T03: Mutation baseline policy recorded and R045 validated.** `est:30m`
  Сохранить final baseline, update R045, and decide whether mutation is informational or future CI candidate. Не включать hard CI gate, если baseline слишком медленный или noisy.
  - Files: `benchmark-results/m050-s03-mutation-baseline.md`, `.gsd/REQUIREMENTS.md`
  - Verify: R045 validation/notes updated; `cd api && go test ./...` passes.

## Files Likely Touched

- benchmark-results/m050-s03-mutation-baseline.md
- api/cache
- api/handlers
- api/lifecycle
- .gsd/REQUIREMENTS.md
