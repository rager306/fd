## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-08-20 - ETag String Assembly Optimization
**Learning:** In Go, combining hex-encoded strings with literals using `+` operator forces multiple heap allocations (one for the hex string, one for the concatenated result).
**Action:** Assemble the components directly into a single stack-allocated byte array before casting to a string to limit allocations to just one for the final string.
