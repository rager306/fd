## 2025-01-01 - Fast-path float32/byte conversion
**Learning:** Found a performance bottleneck in `api/embed/codec.go` where `Float32SliceToBytes` and `BytesToFloat32Slice` iteratively convert float32 to byte and vice-versa, which causes performance hit in base64 encoding/decoding.
**Action:** Use `unsafe.Slice` and `unsafe.Pointer` to perform fast-path slice conversions and utilize `copy` to prevent memory aliasing issues. Also add an `init()` block to verify the host architecture's endianness dynamically.
