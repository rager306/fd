## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-28 - Zero-Allocation Hex Encoding
**Learning:** In highly frequent hot paths like Request ID generation, using `hex.EncodeToString()` dynamically allocates memory on the heap and creates noticeable overhead when combined with string concatenation.
**Action:** Replace `hex.EncodeToString` and concatenation with direct encoding into a statically-sized stack-allocated byte array (e.g. `var buf [36]byte; hex.Encode(buf[:], bytes[:])`), then return it as a single string conversion `string(buf[:])` to minimize allocations and CPU overhead.
