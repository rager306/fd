## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-03-24 - Zero Allocation ETag Generation
**Learning:** `hex.EncodeToString` creates a new string and allocates memory. For generating ETags in hot paths (like `responseETag`), we can avoid multiple allocations and string concatenations (`"\"" + hex.EncodeToString(sum[:]) + "\""`) by directly writing to a stack-allocated byte array (`var dst [66]byte`) and using `hex.Encode(dst[1:65], sum[:])` followed by a single cast to `string(dst[:])`. This reduces allocations to 1 per operation and improves performance.
**Action:** Replace `hex.EncodeToString` concatenations in `responseETag` with direct stack-allocated byte array encoding and casting.
