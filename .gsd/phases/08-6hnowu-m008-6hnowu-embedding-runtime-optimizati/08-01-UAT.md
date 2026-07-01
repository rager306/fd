# S01: Verify model preserving embedding runtime options — UAT

**Milestone:** M008-6hnowu
**Written:** 2026-05-19T17:05:14.297Z

# UAT: S01 Verify model preserving embedding runtime options

## Evidence

- T01 verified `go-bge-m3-embed` source and limitations.
- T02 verified MiniLM Go ONNX path as implementation reference only.
- T03 verified `yalue/onnxruntime_go` path and native artifact requirements.
- T04 defined Russian legal corpus benchmark gate.

## Acceptance

- Current model remains the primary optimization target.
- Model replacements require Russian legal retrieval metrics.
- ONNX path should start with dense output equivalence to current TEI/Candle.

