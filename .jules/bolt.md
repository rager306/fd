## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-07-09 - Fast-path float32 slice to bytes casting
**Learning:** Using `base64.StdEncoding.EncodeToString` combined with manual little-endian byte iteration for float32 arrays creates excessive allocations in the encoding loop for embeddings. By tracking `isLittleEndian` dynamically during `init` and utilizing `unsafe.Slice` alongside `copy`, and then writing the base64 encoding directly into a pre-allocated byte slice and converting via `unsafe.String`, we can bypass memory allocations and dramatically speed up formatting strings.
**Action:** When handling hot-path embedding conversions in Go on little-endian architectures, use `unsafe.Pointer` to safely map between byte boundaries and pre-allocate string arrays.
