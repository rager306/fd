# Decisions Register

<!-- Append-only. Never edit or remove existing rows.
     To reverse a decision, add a new row that supersedes it.
     Read this file at the start of any planning or research phase. -->

| # | When | Scope | Decision | Choice | Rationale | Revisable? | Made By |
|---|------|-------|----------|--------|-----------|------------|---------|
| D001 | M005-pyn4x7 S02 | runtime-performance | How to handle TEI ONNX-missing fallback warnings for deepvk/USER-bge-m3 | Keep the current measured TEI CPU/Candle fallback runtime as acceptable for now; treat ONNX export/runtime changes as a future measured optimization that must be validated by A/B benchmarks before becoming default. | M003/M004 live validation showed the stack is correct and performant enough with the current runtime, while TEI logs indicate ONNX artifacts are missing. Exporting ONNX adds artifact and deployment complexity, so it should be justified by measured cold-latency/memory gains rather than treated as an immediate correctness fix. | Yes | agent |
