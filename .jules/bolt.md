## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - Unsafe Pointer Alignment Risks
**Learning:** Directly casting a `[]byte` backing array (which is 1-byte aligned) to a stricter type like `[]float32` (which requires 4-byte alignment) using `unsafe.Pointer` violates Go's safety rules and can trigger hardware-level `SIGBUS` panics on certain architectures.
**Action:** When converting `[]byte` to `[]float32` using `unsafe`, you must instead allocate the target `[]float32` first (which guarantees proper alignment), cast that target slice down to a `[]byte`, and use `copy` to safely transfer the data from the source.
