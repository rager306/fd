## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-07-01 - Avoid allocations when hashing strings
**Learning:** In Go, converting strings to byte slices and byte slices to hex strings during hashing can cause unnecessary memory allocations in hot paths like cache lookups.
**Action:** Use `unsafe.Slice(unsafe.StringData(text), len(text))` for zero-copy string-to-byte conversion (with a `!= ""` check to avoid panics, and the `//nolint:gosec // G103: performance optimization for byte casting` comment). For hex encoding, use `hex.Encode` into a stack-allocated fixed buffer (e.g. `var buf [64]byte`) and return a string slice from it instead of using `hex.EncodeToString`.
