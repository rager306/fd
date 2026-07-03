# S03: Recommend optimization path — UAT

**Milestone:** M008-6hnowu
**Written:** 2026-05-19T17:14:56.897Z

# UAT: S03 Recommend optimization path

## Evidence

- S03 research artifact ranks options and proposes the next implementation milestone.

## Acceptance

- Next work is measurement/cache foundation first.
- Redis MGET/pipeline is conditional on measured round-trip pressure.
- ONNX FP32 dense-only is later and gated by quality/config evidence.
- INT8, provider tuning, Rust sidecar, C service, and model replacement are deferred or rejected for now.

