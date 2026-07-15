## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - Zero-Allocation String to Hex Hash Conversion
**Learning:** In Go, string to byte slice conversion (`[]byte(str)`) and hex encoding (`hex.EncodeToString`) cause heap allocations. In hot paths like cache lookups, this creates measurable overhead and garbage collection pressure.
**Action:** Use `unsafe.Slice(unsafe.StringData(text), len(text))` for zero-copy string-to-byte conversion (checking for empty string to avoid panics). For hex encoding, use `hex.Encode` into a stack-allocated buffer (e.g., `var buf [64]byte`) and convert to string.
