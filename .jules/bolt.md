## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-28 - String Assembly vs Concatenation in Hot Paths
**Learning:** In Go, string concatenation (`+`) combined with function calls like `hex.EncodeToString` in hot paths (e.g., generating ETags per response) causes multiple unnecessary heap allocations.
**Action:** To optimize string assembly in these paths, pre-calculate the required size, allocate a fixed-size byte array on the stack (e.g., `var dst [66]byte`), assemble the string bytes directly into the array, and perform a single cast to string (`string(dst[:])`) to reduce heap allocations to 1.
