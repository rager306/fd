# S04: Benchmark synthesis and decision — UAT

**Milestone:** M014-vjfs9f
**Written:** 2026-05-20T04:36:29.968Z

# S04 UAT — Benchmark synthesis and decision

## Checks

- [x] TEI vs tagged ONNX comparison artifact exists.
- [x] Comparison states faster/slower by scenario.
- [x] Caveats include fixed-probe correctness, local host, native packaging, and production-default boundaries.
- [x] Decision D010 recorded.
- [x] Default Go tests passed: 78 tests in 4 packages.
- [x] Pinned GolangCI-Lint passed: 0 issues.
- [x] Tagged tests passed: 20 tests in 1 package.
- [x] Artifact hygiene passed: raw probe text leaks 0.
- [x] Tracked native/ONNX binaries: 0.
- [x] Default API health ok; tagged ONNX server stopped.
- [x] GitNexus scope check passed.

## Recommendation

Continue tagged ONNX as opt-in experimental work. Do not switch production/default from TEI yet.

## UAT Result

Pass.

