# M045-d0e5xq: Stabilize TEI startup

**Vision:** Make the TEI-only deployment operationally safe on restart by identifying and validating a bounded startup configuration for HuggingFace TEI, without reintroducing fd ONNX runtime scope.

## Success Criteria

- Effective current TEI image, command, env, mounted model/cache layout, and startup log timeline are documented without disrupting the running service.
- Candidate mitigation options are evaluated against TEI documentation or empirical local evidence, with clear accept/reject rationale.
- At least one controlled startup proof shows measured time to TEI health and fd readiness, or a documented blocker explains why destructive restart proof is unsafe.
- Operator docs and compose config are updated only if a validated safer TEI startup configuration exists.
- Mandatory Go gates remain green if repository code/config changes.

## Slices

- [x] **S01: Non destructive TEI startup recon** `risk:medium` `depends:[]`
  > After this: A recon artifact shows current TEI image command env model cache layout and startup log timeline, plus candidate knobs to test.

- [ ] **S02: Select safe TEI startup mitigation** `risk:high` `depends:[S01]`
  > After this: A selected mitigation is encoded as compose/docs changes or a documented no change decision with rationale.

- [ ] **S03: Controlled startup proof** `risk:high` `depends:[S02]`
  > After this: Startup proof artifact records time to TEI health, fd readiness, first embedding, and remaining ORT/ONNX warnings under the selected configuration.

## Boundary Map

| Boundary | In scope | Out of scope |
|---|---|---|
| fd Go API | Health/readiness smoke, docs/config consistency | Reintroducing ONNX backend |
| HuggingFace TEI | Image command/env/model layout, documented flags, startup logs | Modifying upstream TEI source |
| Docker Compose | Safe config changes and controlled proof | Unplanned destructive restarts |
| Redis/cache | Preserve namespace behavior during smoke checks | Cache benchmarking unless needed for startup proof |
