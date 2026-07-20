## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-05-24 - Zero-Allocation Hex Encoding
**Learning:** Using `hex.EncodeToString` combined with string concatenation or slicing (e.g., adding quotes or hyphens) causes multiple allocations and significant performance overhead in hot paths like cache key generation and request ID creation.
**Action:** Use `hex.Encode` into a stack-allocated byte array (e.g., `var buf [64]byte`), inserting extra characters directly into the array, and perform a single string conversion (`string(buf[:])`) to eliminate allocations and speed up encoding.
