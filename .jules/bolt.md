## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2025-02-23 - UUID String Generation Overhead
**Learning:** In Go, generating UUID strings using `hex.EncodeToString` and multiple string concatenations causes unnecessary allocations.
**Action:** Use a pre-allocated stack byte array (e.g., `var buf [36]byte`), write hex components and hyphens directly into the buffer, and then perform a single string conversion (`string(buf[:])`) to reduce execution time and memory allocations in hot paths.
