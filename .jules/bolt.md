## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-08-22 - Fast-Path Cache Key Generation
**Learning:** In Go, string concatenation combined with `hex.EncodeToString` inside hot paths like cache key generation creates multiple unnecessary heap allocations per request.
**Action:** Use a stack-allocated byte array (`var buf [256]byte`), `hex.Encode` directly into the buffer, and copy constant strings/fast-path dimensions (like 1024 and 512) to reduce allocations to exactly 1 (for the final string cast).
