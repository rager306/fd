## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-30 - ETag Generation Memory Optimization
**Learning:** Generating ETags dynamically for every request via standard string concatenation (`'"' + hex.EncodeToString(sum) + '"'`) causes multiple hidden heap allocations on the hot path.
**Action:** Use a stack-allocated byte array (`var dst [66]byte`) and `hex.Encode` to assemble the hex string and literal quotes directly into a single continuous buffer, reducing allocations from 3 to 1 per request.
