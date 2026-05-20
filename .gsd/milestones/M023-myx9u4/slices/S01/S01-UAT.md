# S01: Packaged ONNX legal quality gate — UAT

**Milestone:** M023-myx9u4
**Written:** 2026-05-20T10:58:52.537Z

# S01 UAT — Packaged ONNX legal quality gate

## Checks

- [x] TEI baseline `/health` passed.
- [x] Packaged ONNX Docker `/health` passed.
- [x] Packaged ONNX non-legal embedding smoke returned 1024 dimensions.
- [x] Evaluator ran with isolated ONNX cache namespace `m023-onnx-docker-legal`.
- [x] Legal gate verdict PASS.
- [x] Raw legal text leak check returned zero leaks.
- [x] ONNX container stopped and port 18000 clean.
- [x] CI binary hygiene false positive for `Dockerfile.onnx` fixed and actionlint passed.

## Key Metrics

- minimum cross-backend cosine: `0.99989883`
- top-1 agreement: `1.0`
- mean overlap@5: `0.997701`
- ONNX recall ratio: `1.0`

## Result

Pass. The packaged ONNX Docker image preserves the selected Russian/legal quality gate.

