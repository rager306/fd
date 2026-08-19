## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - Redis Cache Key Generation Hot Path
**Learning:** Constructing complex strings with standard string concatenation (`+`) and intermediate hex encoding (`hex.EncodeToString`) in hot paths like Redis cache key generation creates multiple intermediate allocations.
**Action:** Use a pre-sized stack-allocated byte array and write the string segments directly into it (e.g. using `hex.Encode` instead of `hex.EncodeToString`) and combine with zero-allocation `unsafe` string-to-byte conversion for the sha256 hash. This reduces allocations from 5 to 1.
