# M001-h8xt3d: Review remediation

**Vision:** Turn the project review findings into verified fixes while preserving the compact embedding service architecture.

## Success Criteria

- No high-severity review findings remain unaddressed or undocumented.
- All modified behavior has targeted tests.
- `cd api && go test ./...` passes.
- Logical commits record the remediation work.

## Slices

- [x] **S01: S01** `risk:high` `depends:[]`
  > After this: Cache layer cannot cross-contaminate 512d and 1024d requests; short TEI vectors return errors rather than panics.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: Invalid batch dimensions or encoding formats return HTTP 400; production handler paths are tested directly.

- [ ] **S03: Local cache semantics** `risk:medium` `depends:[S01]`
  > After this: L1 cache enforces documented size/overwrite semantics with tests.

- [ ] **S04: Runtime configuration hardening** `risk:medium` `depends:[S02,S03]`
  > After this: Docker Compose healthchecks work with runtime image and config names match application expectations.

## Boundary Map

Not provided.
