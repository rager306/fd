## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-07-11 - Fast-path float32 slice to byte slice conversion
**Learning:** Element-by-element float32 to byte conversion via `math.Float32bits` and `binary.LittleEndian.PutUint32` is a bottleneck. Using `unsafe.Slice` and `copy()` provides a ~30% faster alternative on little-endian architectures, but memory aliasing must be avoided by allocating a new slice instead of returning the casted slice directly.
**Action:** When converting large slices across types on compatible architectures, allocate a new slice and use `unsafe.Slice` combined with `copy()` for fast performance, always checking endianness dynamically and avoiding memory aliasing.
