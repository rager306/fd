## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - Zero-Copy String to Byte Slice Casting in Hashing
**Learning:** In Go, calling `[]byte(string)` to compute hashes in highly-frequent hot paths causes measurable overhead due to string-to-byte-slice memory allocation and copying.
**Action:** Replace `[]byte(string)` with zero-copy byte casting using `unsafe.Slice(unsafe.StringData(text), len(text))`. When using `unsafe` in the Go codebase for performance optimizations, append the comment `//nolint:gosec // G103: performance optimization for byte casting` to prevent linting failures.
