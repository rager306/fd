## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-09-26 - UUID String Concatenation Optimization
**Learning:** Replacing `hex.EncodeToString` and multi-part string concatenation (`+`) with a manual stack-allocated byte array for UUID generation surprisingly DIDN'T reduce memory allocations in Go. Both approaches resulted in 1 alloc and 48 bytes/op, as the Go compiler efficiently optimizes single-statement string concatenations.
**Action:** Do not sacrifice readability for complex manual byte-array string building when dealing with simple, single-statement string concatenations, as the compiler is likely already optimizing it to a single allocation.
