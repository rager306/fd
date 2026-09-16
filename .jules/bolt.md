## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - Hash Generation Overhead
**Learning:** In Go, dynamically encoding a hash into a string using `hex.EncodeToString(h[:])` and then immediately slicing it (e.g., `[:12]`) keeps the original full 64-byte backing array alive in memory. Additionally, converting a string to a byte slice using `[]byte(value)` causes an unnecessary heap allocation.
**Action:** Use `unsafe.Slice(unsafe.StringData(str), len(str))` for zero-allocation read-only string-to-byte conversion (guarded by length checks). Encode necessary hash segments directly into a properly sized stack-allocated byte array (e.g., `var dst [12]byte; hex.Encode(dst[:], h[:6])`) before casting to a string.
