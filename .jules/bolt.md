## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-07-06 - Hex Encoding Overhead
**Learning:** In Go hot paths, using `hex.EncodeToString` creates unnecessary allocations because it internally allocates a byte slice before converting it to a string.
**Action:** To avoid allocations when hex-encoding strings in hot paths, prefer using `hex.Encode` into a fixed-size stack-allocated buffer (e.g., `var buf [64]byte`) followed by a standard string conversion (e.g., `string(buf[:])`) rather than `hex.EncodeToString`. Avoid `unsafe.String` for this conversion.
