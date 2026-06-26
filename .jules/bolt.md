## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2025-02-23 - Zero-allocation Hash Generation in Go Hot Paths
**Learning:** In Go, string to byte slice conversions in hot paths (like cache hash generation) cause unnecessary heap allocations. Using `hex.EncodeToString` similarly allocates new strings for each conversion, introducing measurable overhead when generating short hashes frequently.
**Action:** Replace `[]byte(string)` cast with zero-allocation `unsafe.Slice(unsafe.StringData(str), len(str))` and avoid `hex.EncodeToString` by using `hex.Encode` into a stack-allocated buffer (e.g. `var buf [64]byte`), converting only the necessary slice to a string to dramatically reduce GC pressure.
