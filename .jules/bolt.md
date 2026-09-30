## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-31 - Cache Key Allocation Hot Path
**Learning:** In Go, string concatenations in hot paths that include functions like `strconv.Itoa` trigger unnecessary heap allocations when integers exceed small bounds (>99). Hardcoded fast-paths for standard parameters (e.g., embedding dimensions 512, 1024) can eliminate these allocations, allowing the compiler to optimize the remaining concatenation into a single allocation.
**Action:** When working in high-frequency network/cache utilities, look for small string conversion functions and implement fast-paths for highly common parameters to eliminate memory overhead.
