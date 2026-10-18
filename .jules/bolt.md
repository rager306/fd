## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2026-10-17 - Go HTTP Middleware Initialization Overhead
**Learning:** In Go, any operations placed directly inside an HTTP middleware closure (e.g., computing a SHA-256 hash of a globally configured API key, string trimming, constant value checking) are needlessly evaluated on every incoming request. Redundant string slicing using strings.TrimPrefix also causes overhead.
**Action:** Precompute static values, constant conditional logic, or expensive operations outside the returned closure during initialization. Use fast-path string slicing instead of `strings.TrimPrefix` for known prefixes.
