## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-31 - Fast-path string formatting
**Learning:** Using strconv.Itoa in hot paths (like cache key generation) introduces heap allocations that can be avoided for known, frequent values (like 512, 1024). Using fast-path literals allows the compiler to optimize multi-part string concatenations into a single allocation.
**Action:** When creating string identifiers, hardcode fast-paths for highly frequent configurations to avoid strconv.Itoa or fmt.Sprintf.
