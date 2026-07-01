# S01: Benchmark matrix and metadata harness — UAT

**Milestone:** M014-vjfs9f
**Written:** 2026-05-20T04:13:56.252Z

# S01 UAT — Benchmark matrix and metadata harness

## Checks

- [x] Benchmark matrix defined: cold/warm/batch/cache/startup/memory/concurrency.
- [x] Snapshot version bumped to 2.
- [x] Runtime label/build tags can be recorded.
- [x] ONNX/native artifact manifests can be recorded.
- [x] ONNX Runtime library path/hash can be recorded.
- [x] Existing TEI behavior is preserved when env vars are unset.
- [x] Snapshot dry-run passes.
- [x] Raw probe text leak check passes.

## UAT Result

Pass. S02/S03 can now run comparable benchmark artifacts.

