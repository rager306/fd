## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-02-28 - Zero-Copy Float32/Bytes Conversion with unsafe.Slice
**Learning:** For converting `[]float32` to `[]byte` (like for base64 embeddings), copying values using binary put functions introduces unnecessary overhead. By using `unsafe.Slice((*byte)(unsafe.Pointer(&slice[0])), len(slice)*4)` and checking system endianness dynamically, we can skip allocating the slice when encoding base64.
**Action:** Use `unsafe.Slice` alongside a dynamic endianness check (`*(*byte)(unsafe.Pointer(&x)) == 0xFF` equivalence) for high-frequency primitive slice castings where endianness aligns, saving allocations and time.
