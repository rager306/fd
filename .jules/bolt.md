## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-11-06 - Hex Encoding Allocations
**Learning:** `hex.EncodeToString` causes unnecessary allocations. Casting string to bytes also causes allocations.
**Action:** Use `unsafe.Slice(unsafe.StringData(text), len(text))` to cast strings to bytes without allocations, and use `hex.Encode` into a fixed-size stack-allocated buffer `var buf [64]byte` followed by `string(buf[:])` instead of `hex.EncodeToString`.
