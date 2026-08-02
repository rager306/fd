## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2025-02-12 - strconv vs fmt.Sprintf allocations
**Learning:** While `strconv.Itoa` paired with string concatenation is generally faster than `fmt.Sprintf` (e.g. ~100ns vs ~160ns), it can surprisingly result in more distinct memory allocations (2 allocs vs 1 alloc) in Go depending on the size of the final concatenated string due to how Go escapes and allocates concatenated slices on the heap.
**Action:** Do not blindly replace `fmt.Sprintf` with `strconv.Itoa + "..."` in cold paths (like error handlers) just for micro-optimizations, as the allocation trade-off may not be strictly beneficial. For true hot-path zero-allocation string building, use stack-allocated `[]byte` buffers and `strconv.AppendInt`.

## 2025-02-12 - Endian-aware Unsafe Casting for []float32
**Learning:** Element-by-element iteration using `binary.LittleEndian.PutUint32` or `Uint32` to convert `[]float32` <-> `[]byte` is slow in hot paths (like Redis cache backfilling/parsing). Using `unsafe.Slice` and `copy` is nearly 2x faster, but memory layout depends on system architecture.
**Action:** When optimizing binary conversions, implement an architecture-safe fast path: if the system is little-endian (e.g. `isLittleEndian`), directly cast the slice memory using `unsafe.SliceData` and `copy()`. Fall back to element-by-element on big-endian systems to prevent mutability bugs or corrupted data.
