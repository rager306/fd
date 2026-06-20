## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - Zero-Allocation Byte Casting in Hashing
**Learning:** Converting strings to byte slices using `[]byte(text)` causes memory allocations which can be expensive in hot paths like generating cache keys. Since `sha256.Sum256` only reads the byte slice, it is safe to use `unsafe.StringData` to avoid allocating new memory for the byte slice.
**Action:** Use `unsafe.Slice(unsafe.StringData(text), len(text))` when reading strings as bytes in hot paths where the byte slice will not be mutated, but ensure to add `//nolint:gosec // G103: performance optimization for byte casting` to bypass false-positive lint warnings.
