## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2025-02-12 - Zero-copy Float32/Byte Array Serialization
**Learning:** In Go, converting `[]float32` to `[]byte` (and vice-versa) using manual allocations, loops, and `binary.LittleEndian` conversions causes unnecessary CPU and memory overhead during serialization in hot paths (like embedding encoding).
**Action:** Use `unsafe.Slice` with `unsafe.Pointer` to perform direct, zero-copy casting between byte slices and float32 slices to drastically reduce per-request latency.
