---
id: T01
parent: S06
milestone: M008-6hnowu
key_files: []
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T16:51:35.103Z
blocker_discovered: false
---

# T01: Verified ONNX CPU provider landscape: default CPU EP is baseline; oneDNN/OpenVINO are official but Intel-oriented; ZenDNN needs caution due older AMD-specific evidence.

**Verified ONNX CPU provider landscape: default CPU EP is baseline; oneDNN/OpenVINO are official but Intel-oriented; ZenDNN needs caution due older AMD-specific evidence.**

## What Happened

Researched ONNX Runtime CPU execution provider options for AMD EPYC/VPS relevance. Official ONNX Runtime docs list default CPU EP, Intel oneDNN/DNNL EP, Intel OpenVINO EP, XNNPACK, and other providers; oneDNN and OpenVINO have official documentation and build/runtime guidance. OpenVINO has current packages such as `onnxruntime-openvino` and explicit CPU/GPU/NPU provider options, but it is Intel-oriented and must be benchmarked on AMD rather than assumed beneficial. oneDNN is also Intel-oriented but open-source and CPU-targeted; it requires an ONNX Runtime build with DNNL provider and registration. ZenDNN evidence exists from AMD documentation, including tuning topics like env vars, thread affinity, NUMA/numactl, THP, batch size, and memory allocators, but the found AMD guide is old and tied to ONNX Runtime v1.12.1, so it is not yet a low-friction current option. Practical first candidates for a future benchmark matrix are: current TEI/Candle baseline, ONNX Runtime default CPU EP, ONNX Runtime with oneDNN if buildable, OpenVINO EP only as an experiment if it runs correctly on target hardware, and ZenDNN only if a current supported AMD package/build path is found.

## Verification

Read ONNX Runtime execution provider docs, oneDNN/OpenVINO docs, current OpenVINO package info, and AMD ONNX Runtime-ZenDNN guide evidence.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Search query: ONNX Runtime ZenDNN execution provider AMD EPYC oneDNN DNNL OpenVINO CPU execution provider build options 2026` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Read: https://onnxruntime.ai/docs/execution-providers/` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Read: https://onnxruntime.ai/docs/execution-providers/oneDNN-ExecutionProvider.html` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Read: https://onnxruntime.ai/docs/execution-providers/OpenVINO-ExecutionProvider.html` | -1 | unknown (coerced from string) | 0ms |
| 5 | `Read: https://pypi.org/project/onnxruntime-openvino/` | -1 | unknown (coerced from string) | 0ms |
| 6 | `Read: AMD ONNX Runtime-ZenDNN User Guide 4.0 PDF snippets` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

ZenDNN evidence found so far is AMD's ONNX Runtime-ZenDNN User Guide from 2023 referencing ONNX Runtime v1.12.1 and AMD-specific setup. It does not yet look like a standard current ONNX Runtime prebuilt execution provider comparable to CPU EP, oneDNN, or OpenVINO. Needs further verification before recommending.

## Files Created/Modified

None.
