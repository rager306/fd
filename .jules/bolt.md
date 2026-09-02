## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-11-20 - Hash String Shortening Overhead
**Learning:** In Go, creating a shortened hex representation of a SHA256 hash by encoding the entire byte array to a string first (`hex.EncodeToString(hash[:])[:12]`) causes unnecessary heap allocations for the full 64-character string.
**Action:** Encode only the needed bytes into a stack-allocated byte array (`var dst [12]byte; hex.Encode(dst[:], h[:6])`) to minimize heap allocations and CPU overhead.

## 2023-11-20 - UUID Generation Overhead
**Learning:** In Go, when generating custom UUIDs using `hex.EncodeToString` and string concatenations (`+`), it causes multiple heap allocations.
**Action:** Encode directly into a single stack-allocated byte array of the exact required length (e.g. `var buf [36]byte`) and convert it to a string once using `string(buf[:])`.
