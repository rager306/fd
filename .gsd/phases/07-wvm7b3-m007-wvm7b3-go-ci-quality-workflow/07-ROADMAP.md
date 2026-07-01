# M007-wvm7b3: M007-wvm7b3: Go CI quality workflow

**Vision:** Automate the Go quality gate added in M006 so tests and static analysis run consistently on future pushes and pull requests.

## Success Criteria

- Local CI workflow file added for Go tests/lint.
- Workflow mirrors M006 local quality commands.
- README documents CI/local parity.
- No remote push or GitHub state change is performed.
- Local verification gates pass and work is committed.

## Slices

- [x] **S01: S01** `risk:low` `depends:[]`
  > After this: After this, workflow syntax/action choices are grounded in current docs and repo constraints.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, `.github/workflows/go-quality.yml` runs tests and GolangCI-Lint locally-equivalent commands.

- [x] **S03: S03** `risk:low` `depends:[]`
  > After this: After this, CI docs are updated and the milestone is locally committed.

## Boundary Map

## Boundary Map

Not provided.
