## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - Hash Allocation in Hot Paths
**Learning:** In Go, performing standard `hex.EncodeToString` on full arrays inside frequently called hashing utilities (like caching keys) allocates both the intermediate slice and the string backing array. Furthermore, typical string-to-byte casting allocates new memory.
**Action:** Use `unsafe.Slice(unsafe.StringData(str), len(str))` for zero-allocation read-only string casting before hashing. Optimize `hex.Encode` by writing only required bytes directly into a stack-allocated byte array (`var dst [N]byte`) to avoid heap allocations.
