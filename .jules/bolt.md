## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2026-06-25 - Unsafe Slice Memory Casting
**Learning:** To accelerate converting `[]float32` to little-endian `[]byte` in hot paths, directly cast the slice memory using `src := unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(slice))), len(slice)*4)` after verifying the system's endianness, then `copy(dst, src)` to a new byte slice. This avoids slow element-by-element iteration (`binary.LittleEndian.PutUint32`).
**Action:** When converting large numeric slices to byte slices on hot paths, check endianness and use `unsafe.Slice` + `copy` for a ~30-50% speedup.
