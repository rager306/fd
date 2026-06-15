# S04: Test gates documentation and closure

**Goal:** Оформить команды и policy для unit, e2e и mutation test levels, затем выполнить финальную milestone validation.
**Demo:** После этого будущий агент понимает, какие test levels запускать, когда и почему.

## Must-Haves

- README or docs lists current test levels and exact commands.
- Heavy Docker e2e and mutation gates are clearly marked local/manual with prerequisites.
- Final verification commands pass.
- Milestone validation records R043-R045 evidence.

## Proof Level

- This slice proves: documentation plus command verification

## Integration Closure

S04 makes S01-S03 reusable by future agents and avoids hidden tribal knowledge.

## Verification

- Test failure modes and prerequisites become visible in durable docs.

## Tasks

- [x] **T01: README updated with current test levels, e2e modes, and mutation policy.** `est:45m`
  Update README with concise testing section: regular api suite, lint/govulncheck, root integration no-key mode, authenticated Docker e2e mode, bounded mutation baseline. Do not print secrets.
  - Files: `README.md`
  - Verify: README contains commands and secret-handling notes for each test level.

- [x] **T02: Final verification and closure artifact recorded for M050 test gates.** `est:45m`
  Run final regular and lightweight integration commands, record closure artifact if needed, and prepare milestone validation.
  - Files: `benchmark-results/m050-s04-test-gates-closure.md`
  - Verify: `cd api && go test ./...` and `cd tests/integration && go test -v .` pass; docs command snippets match verified commands.

## Files Likely Touched

- README.md
- benchmark-results/m050-s04-test-gates-closure.md
