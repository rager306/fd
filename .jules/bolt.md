## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-06-30 - Optimize Request ID Hex Encoding
**Learning:** In Go, `hex.EncodeToString` allocates memory on every call. Using `hex.Encode` into a fixed-size stack-allocated buffer containing dashes (e.g., `var buf [36]byte`), and manually formatting the UUID parts avoids intermediate slice allocations and provides explicit zero-allocation string creation. This significantly improves performance in highly-frequent hot paths like header middleware.
**Action:** Prefer `hex.Encode` into a carefully constructed stack buffer containing formatting characters (like dashes) and converting the entire buffer to a string `string(buf[:])` over `hex.EncodeToString` and incremental string building for performance-critical string formats.
