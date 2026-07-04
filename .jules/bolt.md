## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2026-07-04 - Optimize hex encoding and cache keys
**Learning:** In hot paths, `hex.EncodeToString` and subsequent string manipulations (like concatenating hyphens or quotes) cause unnecessary memory allocations. Allocating a fixed-size byte array on the stack (e.g. `var buf [64]byte`), using `hex.Encode` into it, and converting it to a string once `string(buf[:])` avoids these allocations.
**Action:** Prefer `hex.Encode` into a stack-allocated buffer followed by string conversion instead of `hex.EncodeToString` when the result needs further string concatenation or manipulation.
