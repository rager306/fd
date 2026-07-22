## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-07-22 - Fast Float32 to Byte Conversions
**Learning:** Serializing `[]float32` to `[]byte` via element-by-element iteration (`binary.LittleEndian.PutUint32`) is a bottleneck in the encoding path (for base64/redis). Using `unsafe.Slice` to cast the backing array and `copy` is significantly faster on little-endian systems. Refactoring system-wide endianness detection into a shared `internal/endian` package avoids logic duplication.
**Action:** Always check the system endianness and apply this `unsafe` casting optimization on little-endian systems for fast float serialization.
