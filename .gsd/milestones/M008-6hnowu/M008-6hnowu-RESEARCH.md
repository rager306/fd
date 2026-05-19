# M008 Environment Snapshot

## Purpose

Baseline environment evidence for M008 optimization research and future benchmark comparability.

## Artifact

- Environment snapshot file: `benchmark-results/fd-environment-inxi-m008.txt`
- Command recorded: `inxi -F -c0`
- Git commit recorded in artifact: `c1d2453b47ab973169bca88305b25eddb3b719ab`
- Branch: `master`

## Key environment facts

- OS: Ubuntu 24.04.4 LTS (Noble Numbat)
- Kernel: Linux 6.8.0-117-generic x86_64
- Virtualization: KVM/QEMU, `pc-i440fx-9.0`
- CPU: 12-core AMD EPYC (with IBPB), average reported speed 3195 MHz
- Memory: 48 GiB total, about 47 GiB available at capture time
- Storage: 250 GiB QEMU disk, ext4 root filesystem
- Swap: none detected
- Network: Virtio network plus Docker bridges; active Docker bridge links reported at 10 Gbps
- Headless environment: no display server; llvmpipe OpenGL renderer
- `inxi`: 3.3.34

## Benchmark implications

- Future ONNX/Redis/language benchmark artifacts must include this environment snapshot or a refreshed equivalent.
- No swap means memory-pressure tests must record RSS/container memory carefully; OOM behavior may differ from swap-enabled hosts.
- KVM/QEMU and virtio networking should be considered when interpreting CPU affinity, NUMA, network, and Redis latency results.
- AMD EPYC CPU makes AMD/CPU-provider claims relevant, but ZenDNN/oneDNN/OpenVINO must still be benchmarked rather than assumed.
