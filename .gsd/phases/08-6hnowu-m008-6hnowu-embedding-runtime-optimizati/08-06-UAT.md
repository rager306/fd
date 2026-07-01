# S06: Research ONNX CPU acceleration and quantization — UAT

**Milestone:** M008-6hnowu
**Written:** 2026-05-19T17:03:39.343Z

# UAT: S06 ONNX CPU Acceleration and Quantization

## Evidence

- `.gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md` exists and includes provider classification, threading/NUMA controls, INT8 gate, BGE-M3 multi-output implications, benchmark order, config fields, and stop criteria.
- `benchmark-results/fd-environment-inxi-m008.txt` captures the current environment.

## Acceptance

- Claims about ZenDNN are not accepted as simple switch-on behavior.
- Dense-only FP32 ONNX default CPU EP is first benchmark candidate.
- INT8 requires Russian legal quality benchmark.
- Sparse/ColBERT outputs are future hybrid retrieval scope, not current API behavior.

