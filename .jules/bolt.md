## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-01-28 - Zero-Copy String to Byte Slice Casting
**Learning:** In highly-frequent hot paths (like cache lookups), `[]byte(string)` creates unnecessary allocations due to copying the underlying string data. Replacing this with `unsafe.Slice(unsafe.StringData(text), len(text))` provides zero-copy string-to-byte casting, reducing allocations and improving performance.
**Action:** Use the `unsafe` approach for string-to-byte conversions when reading strings (like for SHA256 hashing) in hot paths. Remember to append the comment `//nolint:gosec // G103: performance optimization for byte casting` to satisfy the linter.
