# M040-pbp9z1: M040-pbp9z1: Same-host embedding service readiness

**Vision:** Prepare `fd` as a same-host local HTTP embedding service for neighboring services with excellent Russian/legal-domain quality and optimal speed on the current host; ONNX is a candidate runtime, not the goal itself.

## Success Criteria

- Same-host local HTTP embedding service contract exists and is grounded in current code/runtime evidence.
- Packaged Docker restart/cache behavior is proven or truthfully blocked with exact evidence.
- Legal-domain quality remains no-regression for any runtime/candidate included in the recommendation.
- A bounded alternative model quick gate is completed without open-ended experimentation.
- Final artifact recommends TEI vs ONNX, or explicitly defers recommendation, using the evidence envelope.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, local services have a concrete contract for calling fd and interpreting runtime readiness.

- [ ] **S02: S02** `risk:high` `depends:[]`
  > After this: After this, packaged ONNX restart/cache behavior is measured instead of skipped.

- [ ] **S03: Bounded legal model quick gate** `risk:medium` `depends:[S01]`
  > After this: After this, alternative model scope is bounded by legal-domain evidence and cannot hijack the service-readiness milestone.

- [ ] **S04: Runtime recommendation and operating contract** `risk:medium` `depends:[S02,S03]`
  > After this: After this, the user has a TEI-vs-ONNX same-host runtime recommendation with evidence and operating contract.

## Boundary Map

## Boundary Map

### S01 → S02

Produces:
- Same-host local HTTP service contract: endpoints, env/runtime requirements, health metadata, timeout/retry guidance, no-silent-fallback rule.

Consumes:
- Existing M038/M039 runtime evidence and current API surfaces.

### S02 → S04

Produces:
- Packaged Docker restart/cache benchmark evidence with Redis L2 restart behavior and sanitized config.

Consumes:
- S01 service contract and packaged ONNX runtime requirements.

### S03 → S04

Produces:
- Bounded legal-domain candidate model quick-gate result: keep current, defer candidate, or reject candidate.

Consumes:
- S01 scope boundary and R001/R008 quality constraints.
