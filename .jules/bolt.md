## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-05-18 - ETag String Assembly Overhead
**Learning:** Constructing ETags with string concatenation and `hex.EncodeToString` (`'"' + hex + '"'`) causes multiple heap allocations and overhead.
**Action:** Assemble the quoted hex string directly into a single stack-allocated byte array (`var dst [66]byte; hex.Encode(...)`) and perform a single cast to string to minimize heap allocations.
