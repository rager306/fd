## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-07-09 - Fast-path Float32 to Base64 String Encoding
**Learning:** Using `base64.StdEncoding.EncodeToString` allocates a new string. Additionally, manually converting `[]float32` to `[]byte` via `binary.LittleEndian.PutUint32` inside a loop is slow. On little-endian architectures, memory aliasing the float slice to a byte slice via `unsafe.Slice` and using `copy()` is significantly faster.
**Action:** Use `unsafe.Slice` and `copy()` to fast-path little-endian slice conversions. Use `base64.StdEncoding.Encode` into a pre-allocated byte buffer and convert it to a string without allocation using `unsafe.String(unsafe.SliceData(buf), len(buf))`.
