## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - Redis Cache Key Generation Optimization
**Learning:** `hex.EncodeToString` causes unnecessary heap allocations. Using a stack-allocated buffer (e.g. `var buf [64]byte`) and explicitly passing strings saves an allocation and time. Also, dynamically allocating strings for dimensions (`strconv.Itoa(dim)`) on every cache hit adds significant overhead in hot paths.
**Action:** Replace `hex.EncodeToString` with stack allocation and `string()` conversion. Apply fast-paths to common integer concatenations to reduce allocations.
