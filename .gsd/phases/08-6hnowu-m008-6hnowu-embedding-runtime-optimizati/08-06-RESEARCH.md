# S06 Research: ONNX CPU Acceleration and Quantization

## Scope

Verify ONNX Runtime CPU optimization paths for current-model inference:

- Default ONNX Runtime CPU EP
- oneDNN/DNNL EP
- OpenVINO EP
- ZenDNN / AMD-specific ORT path
- NUMA and threading controls
- INT8 quantization
- BGE-M3 dense/sparse/ColBERT output implications

This research does **not** recommend changing the default runtime yet. Any change requires A/B benchmark evidence against the current TEI/Candle baseline and, when output quality may change, Russian legal corpus quality evidence.

## Current environment baseline

Environment snapshot saved at:

- `benchmark-results/fd-environment-inxi-m008.txt`

Key facts:

- Ubuntu 24.04.4 LTS
- Kernel `6.8.0-117-generic`
- KVM/QEMU virtual machine
- 12-core AMD EPYC with IBPB
- 48 GiB RAM
- No swap detected
- Docker/virtio networking present

Future ONNX benchmarks should refresh this with:

- `inxi -F -c0`
- `lscpu`
- `numactl -H` if available
- Docker image IDs
- container CPU/memory limits
- git commit
- sanitized effective env/config snapshot

## Provider classification

### 1. Current TEI/Candle baseline

Status: current default.

Role: control baseline for all comparisons.

Keep unless a measured alternative beats it on latency/throughput/resource usage while preserving dense embedding compatibility and Russian/legal quality.

### 2. ONNX Runtime default CPU EP

Status: first ONNX candidate.

Why:

- standard ONNX Runtime path;
- lowest provider complexity;
- easiest to reproduce;
- suitable first baseline for BGE-M3 ONNX dense output.

Benchmark first before oneDNN/OpenVINO/ZenDNN.

### 3. oneDNN/DNNL EP

Status: benchmark candidate if build/package path is practical.

Evidence:

- Official ONNX Runtime oneDNN EP docs exist.
- It is Intel-originated/open-source CPU optimization.
- Requires ORT build/package with DNNL provider and registration.

Risk:

- It is Intel-oriented; benefit on AMD EPYC must be measured.
- Build/deploy complexity may outweigh benefit.

### 4. OpenVINO EP

Status: experimental candidate.

Evidence:

- Official ONNX Runtime OpenVINO EP docs exist.
- Current `onnxruntime-openvino` packages exist.
- OpenVINO EP exposes options like `device_type`, `num_of_threads`, `num_streams`, `cache_dir`, and precision hints.

Risk:

- Intel-oriented CPU/GPU/NPU stack.
- On AMD EPYC, do not assume benefit.

### 5. ZenDNN

Status: caution / verify current support before benchmarking.

Evidence:

- AMD ONNX Runtime-ZenDNN User Guide exists.
- Guide is old: January 2023 and references ONNX Runtime v1.12.1.
- It includes tuning topics: env vars, thread affinity, NUMA/numactl, THP, batch size, memory allocators.

Risk:

- Does not currently look like a standard modern prebuilt ONNX Runtime EP comparable to CPU/oneDNN/OpenVINO.
- May require AMD-specific or custom ORT build.

Recommendation: do not treat `ONNX Runtime + ZenDNN` as a simple switch until a current supported build path is verified.

## Threading and NUMA controls

Relevant ONNX Runtime controls:

- `intra_op_num_threads`: parallelism inside operators.
- `inter_op_num_threads`: parallelism across graph nodes when `ORT_PARALLEL` execution mode is used.
- `execution_mode`: default `ORT_SEQUENTIAL`; `ORT_PARALLEL` can help branching graphs but can hurt linear graphs.
- `graph_optimization_level`: default all optimizations.
- `session.intra_op.allow_spinning` and `session.inter_op.allow_spinning`.
- `session.intra_op.spin_duration_us` / `session.inter_op.spin_duration_us`.
- `session.intra_op.spin_backoff_max` / `session.inter_op.spin_backoff_max`.
- `session.intra_op_thread_affinities` for explicit CPU affinity.

Important ORT behavior:

- `intra_op_num_threads=0` uses physical CPU cores and enables some default affinity.
- Explicit thread counts may disable default affinity unless configured explicitly.
- Since ORT 1.14, ORT thread pools can use cores across NUMA nodes.
- ORT docs show that pinning threads to one NUMA node can materially improve performance in some cases, but this must be benchmarked.

Future benchmark matrix:

- default ORT settings;
- fixed intra-op: 1, 2, 4, 8, 12;
- optional inter-op only if model graph benefits;
- spinning enabled vs disabled/bounded;
- numactl/cpuset only after topology is recorded;
- batch sizes relevant to chunking workload.

## INT8 quantization

Status: later experiment, not first implementation.

Evidence:

- ONNX Runtime supports 8-bit quantization via QOperator and QDQ formats.
- ORT supports dynamic, static, and QAT quantization APIs.
- ORT docs recommend preprocessing/shape inference and transformer-specific optimization paths.
- `gpahal/bge-m3-onnx-int8` exists and claims BGE-M3 ONNX Runtime INT8 with dense/sparse/ColBERT outputs plus O2 optimization.

Risks:

- Quantization is not lossless.
- Static quantization requires representative calibration data.
- For this project, representative data means Russian legal corpus chunks/queries, not generic English samples.
- INT8 can change dense embedding similarities and retrieval ranking.

Required gate before adoption:

- FP32 ONNX baseline works and is compared to current TEI/Candle.
- INT8 output similarity measured against FP32/current output.
- Retrieval metrics measured on Russian legal corpus.
- Config snapshot records quantization method, calibration corpus, provider, ORT version, artifact hash, tokenizer hash.

## BGE-M3 output implications

Verified claims:

- BGE-M3 supports dense retrieval, sparse lexical retrieval, and multi-vector/ColBERT-style retrieval.
- ONNX variants such as `aapot/bge-m3-onnx` and `philipchung/bge-m3-onnx` claim dense, sparse, and ColBERT outputs in one pass.

fd compatibility:

- Current fd API/cache is dense-only.
- `api/embed/types.go` represents OpenAI-compatible dense embeddings.
- `api/cache/redis.go` stores one dense float32 vector per text/dimension key.

Recommendation:

- First ONNX spike should consume **dense output only** for compatibility.
- Sparse/ColBERT should be ignored or stored in a separate experimental namespace until a hybrid retrieval milestone exists.
- Do not expose sparse/ColBERT through `/v1/embeddings` by default.

## Recommended benchmark order

1. Current TEI/Candle baseline with config/environment snapshot.
2. BGE-M3 ONNX FP32 dense-only using default CPU EP.
3. ORT CPU EP threading matrix.
4. Optional NUMA/cpuset/numactl matrix if topology is visible.
5. oneDNN EP if practical to build/run.
6. OpenVINO EP only as an experiment.
7. ZenDNN only after current supported build path is verified.
8. INT8 only after FP32 ONNX baseline and Russian legal quality gate.

## Required config snapshot fields

- `ONNX_PROVIDER`
- `ONNX_RUNTIME_VERSION`
- `ONNX_MODEL_PATH`
- `ONNX_MODEL_SHA256`
- `TOKENIZER_PATH`
- `TOKENIZER_SHA256`
- `ONNX_GRAPH_OPT_LEVEL`
- `ONNX_QUANTIZATION`
- `ONNX_QUANTIZATION_METHOD`
- `ONNX_CALIBRATION_CORPUS`
- `ORT_INTRA_OP_THREADS`
- `ORT_INTER_OP_THREADS`
- `ORT_EXECUTION_MODE`
- `ORT_ALLOW_SPINNING`
- `ORT_SPIN_DURATION_US`
- `ORT_SPIN_BACKOFF_MAX`
- `ORT_THREAD_AFFINITY`
- `NUMACTL_ARGS`
- `OMP_NUM_THREADS`
- `OMP_PROC_BIND`
- `OMP_PLACES`
- `BATCH_SIZE`
- `MAX_SEQUENCE_LENGTH`
- `BENCHMARK_CORPUS_VERSION`
- `GIT_COMMIT`
- Docker image IDs and container CPU/memory limits

## Stop criteria

Stop or reject a provider/quantization variant if:

- dense output cannot be matched to current semantics;
- tokenizer/model artifact provenance is unclear;
- Russian legal retrieval quality regresses beyond agreed threshold;
- operational build complexity is high and latency/throughput gain is marginal;
- memory use exceeds VPS budget;
- benchmark results cannot be reproduced with recorded config.
