## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2026-06-14 - Zero-allocation embedding codec
**Learning:** Manual loop-based `[]float32` to `[]byte` serialization causes significant CPU overhead and large heap allocations in the hot path. On little-endian architectures (x86/ARM), these slices can be directly memory-mapped.
**Action:** Use `unsafe.Slice` to directly map memory when encoding/decoding primitive slices to eliminate O(N) copying and heap allocations.
