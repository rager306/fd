---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Bounded mutation baseline прошёл на критичных cache, handlers и lifecycle файлах.

Запустить выбранный mutation runner на одном или нескольких критичных пакетах with timeout. Prefer cache/handlers/lifecycle bounded scope. Capture score/survivors or exact failure.

## Inputs

- `T01 chosen runner`

## Expected Output

- `benchmark-results/m050-s03-mutation-baseline.md`

## Verification

Mutation command exits successfully with baseline or fails with documented tool blocker; no production code changes required.

## Observability Impact

Выявляет слабые assertions или tool incompatibilities.
