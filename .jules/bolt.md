## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-11-20 - UUID Generation Allocations
**Learning:** In hot-path middleware (like request-id generation), using `hex.EncodeToString` combined with string slicing/concatenation creates unnecessary heap allocations per request.
**Action:** Assemble string components directly into a stack-allocated byte array (`var dst [36]byte`) and cast to string at the end (`string(dst[:])`) to halve allocations and improve generation speed.
