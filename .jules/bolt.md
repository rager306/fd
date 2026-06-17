## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - Zero-Copy String to Byte Slice Conversion
**Learning:** `sha256.Sum256([]byte(string))` causes an unnecessary heap allocation because `[]byte` forces a copy of the string contents. Since `sha256.Sum256` only reads the byte slice, it is safe to use `unsafe.StringData` to avoid the allocation.
**Action:** Use `unsafe.Slice(unsafe.StringData(text), len(text))` when casting immutable strings to read-only byte slices in performance-sensitive paths. Add `//nolint:gosec // G103: performance optimization for byte casting` to bypass lint checks.
