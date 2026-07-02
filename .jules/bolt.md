## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2026-07-02 - Avoid hex.EncodeToString in hot paths
**Learning:** `hex.EncodeToString` allocates a new string every time. In hot paths, it's faster to use `hex.Encode` into a stack-allocated buffer (e.g., `var buf [12]byte`) and then convert that buffer to a string (`string(buf[:])`). Furthermore, zero-copy string-to-byte conversion using `unsafe.Slice(unsafe.StringData(value), len(value))` avoids an additional allocation when feeding the string into a hash function like `sha256.Sum256`.
**Action:** Always prefer `hex.Encode` with a stack-allocated buffer and zero-copy string-to-byte conversions via `unsafe` (with bounds checks and linting exceptions) for hot-path hashing and encoding in Go.
