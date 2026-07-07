## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2026-07-07 - Hex Encoding and String Concatenation Overhead
**Learning:** Using `hex.EncodeToString` followed by slicing and string concatenation to build a UUID causes unnecessary memory allocations and CPU overhead in hot paths.
**Action:** Use `hex.Encode` to write directly into a fixed-size stack-allocated byte array, insert any formatting characters (like hyphens) into the array, and then perform a single `string(buf[:])` conversion at the end.
