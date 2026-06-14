# M046-zqzcu6: Audit remediation waves

**Vision:** Turn GitHub issue #3 from a broad audit report into verified, risk-ordered fixes. Each wave first proves whether a reported defect exists in current code, identifies the design decision or assumption that caused it, then applies the smallest safe remediation with tests and runtime proof.

## Success Criteria

- Issue #3 P0/P1 findings are verified against current code before fixes are applied.
- Confirmed P0 batch endpoint abuse paths are fixed with tests proving rejection before backend work.
- Default exposure/auth behavior matches the same-host contract and avoids accidental public unauthenticated inference.
- LocalCache correctness findings are either fixed with race evidence or falsified with proof.
- Mandatory Go gates and rebuilt-container smoke pass after remediation.
- Residual P2/P3 findings are explicitly triaged rather than silently ignored.

## Slices

- [x] **S01: Audit validation map** `risk:high` `depends:[]`
  > After this: A durable assessment maps issue #3 P0 and P1 findings to current-code evidence, false positives, root decisions, and fix waves.

- [x] **S02: Batch endpoint guardrails** `risk:medium` `depends:[S01]`
  > After this: Oversized or malformed requests to `/v1/batch` and `/embeddings/batch` are rejected before backend work, while valid batch smoke still passes.

- [x] **S03: Batch backend work shaping** `risk:high` `depends:[S02]`
  > After this: A mixed cache-miss batch triggers bounded TEI calls per chunk instead of one TEI call per input.

- [x] **S04: Exposure posture policy** `risk:high` `depends:[S01]`
  > After this: Default local compose remains usable for same-host development, while accidental non-loopback unauthenticated inference and sensitive diagnostics are blocked or explicitly opted in.

- [x] **S05: LocalCache correctness** `risk:medium` `depends:[S01]`
  > After this: Race-enabled cache tests prove bounded size/accounting and clean shutdown under concurrent access.

- [ ] **S06: Audit closure and residual plan** `risk:medium` `depends:[S02,S03,S04,S05]`
  > After this: Issue #3 has a validated closure matrix showing which P0/P1 findings are fixed, which lower-priority items remain, and what future work is needed.

## Boundary Map

Not provided.
