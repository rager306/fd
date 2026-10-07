## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2026-10-07 - Fast paths for string conversion
**Learning:** In high-throughput hot paths (like Redis cache key generation), `strconv.Itoa` causes heap allocations for numbers over 99. Simple string concatenation using constants enables single-allocation optimization by the Go compiler.
**Action:** Add fast paths with string literals for commonly used dimension sizes (e.g., 512 and 1024) to eliminate `strconv.Itoa` heap allocations.
## 2026-10-07 - Fast paths and Go string concatenation
**Learning:** Using intermediate variables (e.g. `base := a + b`) during string concatenation defeats Go's `concatstrings` compiler optimization and forces multiple allocations.
**Action:** When creating fast paths, write the full concatenation on a single line (e.g., `return a + b + "suffix"`) to let the Go compiler optimize it into a single allocation.
