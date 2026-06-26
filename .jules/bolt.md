## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-06-26 - Hex Encoding Allocations Overhead
**Learning:** `hex.EncodeToString` creates extra heap allocations due to returning a new dynamically allocated string. In hot paths (like request IDs, cache keys, ETags), this overhead adds up.
**Action:** Replace `hex.EncodeToString` with `hex.Encode` into a stack-allocated buffer (e.g., `var buf [64]byte`) and use string conversion (e.g., `string(buf[:])`) to allow the Go compiler to optimize and reduce allocations.
