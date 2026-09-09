## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - Redis Cache Key Generation Memory Allocations
**Learning:** Constructing complex strings like cache keys using standard concatenation (`+`) alongside function calls like `hex.EncodeToString` and `strconv.Itoa` causes multiple heap allocations per call (e.g. 5 allocs, 292B).
**Action:** For hot-paths, encode hex representations and other strings directly into a single pre-allocated stack/heap byte array, using `unsafe.String` and `unsafe.Slice` for zero-copy conversions. This reduces allocations to 1 (96B).
