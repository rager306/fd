## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-07-24 - Optimizing Float32 <-> Bytes conversions
**Learning:** Found an opportunity to speed up embedding format processing in `api/embed/codec.go`. Iterating through float32 slices to encode/decode to `[]byte` via `binary.LittleEndian` involves high per-element calculation overhead on heavy execution paths. Replacing it with memory mapping via `unsafe.Slice` + `copy()` leverages native processing limits and fast-paths the operation on little-endian architectures, increasing conversion throughput and improving general memory access profiles.
**Action:** Always test system endianness properly using `unsafe.Pointer` on `uint16(0x00FF)` for mapping values into hardware byte representations, then map direct slice contents when feasible safely with `unsafe.Slice`.
