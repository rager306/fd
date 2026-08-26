## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - UUID String Allocation Overhead
**Learning:** In Go, string concatenation (`a + "-" + b`) combined with intermediate slice strings (`hex.EncodeToString`) creates unnecessary heap allocations in hot-path middleware. Fallback paths using `fmt.Sprintf` add further reflection overhead.
**Action:** When assembling fixed-length strings like UUIDs (especially in high-frequency paths like middleware), encode directly into a stack-allocated byte array (`var buf [36]byte`) using `hex.Encode` and convert to string once (`string(buf[:])`). This reduces allocations from 2 to 1 and halves generation time.
