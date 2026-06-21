## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-05-18 - String to byte slice conversion overhead
**Learning:** In Go, casting a string to a byte slice `[]byte(string)` causes an allocation. This creates a bottleneck in performance-critical paths (e.g. hashing embedding texts to build a cache key).
**Action:** Use `unsafe.Slice(unsafe.StringData(text), len(text))` in Go to perform zero-copy byte casting. Include `//nolint:gosec // G103: performance optimization for byte casting` to avoid linter failure.
