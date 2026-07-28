## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-07-28 - Optimizing hex string sub-slice memory leaks
**Learning:** Taking a sub-slice of a string (e.g., `str[:12]`) keeps the entire original backing string alive in memory indefinitely. For operations like short hash generation using `hex.EncodeToString`, this leaks a 64-byte string when only 12 bytes are needed.
**Action:** Always assemble shortened encoded values into stack-allocated fixed-size buffers (`var dst [64]byte`), perform the encoding, and manually cast the needed subset `string(dst[:12])` to guarantee a completely decoupled, smaller string allocation.
