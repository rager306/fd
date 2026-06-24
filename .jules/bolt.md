## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2025-02-20 - Hex encoding allocations in hot paths
**Learning:** Using `hex.EncodeToString` allocates multiple times because it allocates memory for the returned string and internal byte slice.
**Action:** In high-frequency paths, prefer using `hex.Encode` with a stack-allocated buffer (e.g., `var buf [64]byte`) and then casting to a string (e.g., `string(buf[:])`), which reduces allocations to 1.
