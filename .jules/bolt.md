## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-08-13 - Performance & String Assembly (Go)
**Learning:** To eliminate multiple heap allocations when constructing strings that combine hex-encoded data with byte literals (e.g., building a quoted ETag `"` + hex + `"`), assemble the components directly into a single stack-allocated byte array (e.g., `var dst [66]byte; dst[0] = '"'; hex.Encode(dst[1:65], sum); dst[65] = '"'`) before performing a single cast to `string(dst[:])`.
**Action:** Use stack-allocated byte arrays for fixed-size string concatenations, especially when hex-encoding hashes, to avoid unnecessary heap allocations in hot paths like HTTP middleware.
