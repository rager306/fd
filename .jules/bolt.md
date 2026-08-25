## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - Zero-allocation string assembly for ETag
**Learning:** String concatenation combined with hex encoding (e.g. `"` + hex.EncodeToString(sum[:]) + `"`) creates multiple heap allocations per call (3 allocs in this case) on every cached response path.
**Action:** Assemble string components directly in a stack-allocated array (like `[66]byte`) and cast to string once, reducing allocations to just 1 and saving CPU cycles in hot loops.
