## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - Avoid Hex String Allocations
**Learning:** Using `hex.EncodeToString` in combination with string concatenation allocates multiple intermediate strings and performs unnecessary copying.
**Action:** To avoid allocations when hex-encoding strings in hot paths, prefer using `hex.Encode` into a fixed-size stack-allocated buffer (e.g., `var buf [64]byte`) followed by a standard string conversion (e.g., `string(buf[:])`) rather than `hex.EncodeToString`. Insert extra required characters (e.g., hyphens or quotes) directly into the byte array before conversion to eliminate further string concatenations.
