## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - ETag Generation Overhead
**Learning:** `hex.EncodeToString` and string concatenation `"..."` cause multiple allocations on every cacheable request, adding unnecessary GC pressure in a high-frequency hot path (ETag calculation in middleware).
**Action:** Assemble the components directly into a single stack-allocated byte array (e.g., `var dst [66]byte; dst[0] = '"'; hex.Encode(dst[1:65], sum); dst[65] = '"'`) before performing a single cast to `string(dst[:])` to avoid multiple heap allocations.
