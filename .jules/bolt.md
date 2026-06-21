## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - Fast-path Float32 to Byte Encoding
**Learning:** Manual iteration over `[]float32` to perform bit conversion and assignment into `[]byte` via `binary.LittleEndian.PutUint32` is an unnecessary overhead on little-endian architectures. Using `unsafe.Pointer` combined with `copy()` reduces slice-to-byte encode times by >15% without sacrificing safety if we dynamically check endianness.
**Action:** Always prefer `unsafe.Slice` + `copy` for encoding primitive slice types (`[]float32`, `[]int`) into bytes on compatible architectures, ensuring to fallback for big-endian systems.
