## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-08-26 - ETag String Construction Overhead
**Learning:** In Go, concatenating small strings with dynamically generated hex strings (e.g., `"` + hex.EncodeToString(hash) + `"`) triggers multiple heap allocations per request in middleware hot paths.
**Action:** Use a stack-allocated byte array (e.g., `var dst [66]byte`) and `hex.Encode` directly into it, followed by a single cast to string, to significantly reduce GC pressure and allocations per request.
