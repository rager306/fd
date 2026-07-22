## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - Safe Unsafe Memory Casting
**Learning:** When using `unsafe` in Go to cast `[]byte` to slices of types with stricter alignment (like `[]float32`), casting the byte slice's backing pointer directly can cause unaligned read panics (SIGBUS) on architectures enforcing strict memory alignment.
**Action:** Instead of casting `[]byte` up to `[]float32`, allocate the destination `[]float32` first (guaranteeing 4-byte alignment), cast it down to `[]byte`, and use `copy()` to move the data.
