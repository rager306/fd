---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Run final M015 verification gates

Run final verification gates after last changes: py_compile evaluator, dry-run hygiene, artifact hygiene, Go tests, pinned lint, tagged tests, runtime cleanup, GitNexus detect.

## Inputs

- `tools/evaluate_legal_retrieval.py`
- `api`

## Expected Output

- `Task summary with final evidence`

## Verification

All commands pass or expected quality-fail artifact is verified; no background process remains.

## Observability Impact

Prevents closing milestone without fresh verification evidence.
