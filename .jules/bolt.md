## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-11-20 - Hex Encoding Allocations in Hot Paths
**Learning:** `hex.EncodeToString` and `fmt.Sprintf` cause measurable allocation overhead when generating frequent values like request correlation IDs.
**Action:** Use `hex.Encode` with stack-allocated byte arrays `var buf [36]byte` and simple string conversions `string(buf[:])` instead.
