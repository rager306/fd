## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-02-14 - Optimize string hashing in hot paths
**Learning:** Using `[]byte(string)` and `hex.EncodeToString` creates unnecessary heap allocations in Go hot paths.
**Action:** Use `unsafe.Slice(unsafe.StringData(text), len(text))` for zero-copy string-to-byte casting (with `//nolint:gosec // G103: performance optimization for byte casting`) and stack-allocated arrays with `hex.Encode` to eliminate these allocations.
