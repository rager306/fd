---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Mutation baseline policy recorded and R045 validated.

Сохранить final baseline, update R045, and decide whether mutation is informational or future CI candidate. Не включать hard CI gate, если baseline слишком медленный или noisy.

## Inputs

- `T02 output`

## Expected Output

- `benchmark-results/m050-s03-mutation-baseline.md`

## Verification

R045 validation/notes updated; `cd api && go test ./...` passes.

## Observability Impact

Оставляет будущим агентам точку отсчёта mutation quality.
