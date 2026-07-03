---
estimated_steps: 1
estimated_files: 5
skills_used: []
---

# T01: Собран инвентарь существующих тестов и verification scripts.

Собрать полный список текущих Go tests, root integration tests, verification Python scripts and CI test commands. Классифицировать их по уровню: unit, in-process integration, live runtime, static contract, perf or artifact verifier. Не менять код в этом task.

## Inputs

- `M050-rfqm1p roadmap`
- `current repository tree`

## Expected Output

- `benchmark-results/m050-s01-test-actuality.md`

## Verification

Artifact contains complete inventory counts and categories; no code changes required.

## Observability Impact

Показывает, какие проверки реально существуют и какие только лежат в repo.
