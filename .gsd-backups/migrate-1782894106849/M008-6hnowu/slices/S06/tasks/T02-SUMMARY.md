---
id: T02
parent: S06
milestone: M008-6hnowu
key_files:
  - benchmark-results/fd-environment-inxi-m008.txt
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:01:45.638Z
blocker_discovered: false
---

# T02: Defined ONNX threading/NUMA benchmark controls: start with ORT defaults, then fixed intra-op/thread-affinity/numactl variants only with topology evidence.

**Defined ONNX threading/NUMA benchmark controls: start with ORT defaults, then fixed intra-op/thread-affinity/numactl variants only with topology evidence.**

## What Happened

Researched ONNX Runtime thread/NUMA controls. Official ORT docs describe intra-op threads for parallelism inside operators, inter-op threads for parallelism across graph nodes when execution mode is ORT_PARALLEL, graph optimization level, spinning controls, spin duration/backoff, and explicit intra-op thread affinity. Defaults use physical cores and enable some affinity when intra_op_num_threads is left at 0; explicit thread counts remove default affinity unless affinity is configured. Since ORT 1.14, ORT thread pools can use physical cores across NUMA nodes; docs recommend testing settings because cross-NUMA cooperation can incur cache-miss overhead, and an example showed single-NUMA affinity improving performance versus split-NUMA by nearly 20%. For fd, benchmark controls should compare default ORT CPU behavior, fixed intra-op counts such as 1/2/4/8/12, sequential vs parallel execution only if the graph has enough branches, spinning on/off or bounded spin, and optional numactl/cpuset only after topology is recorded. Benchmark artifacts must record ORT session settings, environment variables, CPU affinity, NUMA topology, container CPU limits, and model/batch/sequence length.

## Verification

Read ONNX Runtime thread management docs, NUMA tuning section, AMD ZenDNN guide snippets, and current environment snapshot.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Fetched: https://onnxruntime.ai/docs/performance/tune-performance/threading.html` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Search query: ONNX Runtime CPU threading intra_op inter_op thread affinity NUMA performance tuning environment variables OMP_NUM_THREADS numactl` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Artifact: benchmark-results/fd-environment-inxi-m008.txt` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

The current environment snapshot reports KVM/QEMU 12-core AMD EPYC with no explicit NUMA topology from `inxi`; a future implementation spike should also record `lscpu`, `numactl -H` if available, container CPU limits, and Docker CPU affinity. NUMA optimization may be irrelevant if the VPS exposes a single NUMA node or virtualized topology hides NUMA.

## Files Created/Modified

- `benchmark-results/fd-environment-inxi-m008.txt`
