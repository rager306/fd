## 2026-06-14 - Optimize Redis Cache Serialization and Hashing
**Learning:** In hot paths doing cache lookup, memory allocation overhead from text hashing (`[]byte(text)` casting) and byte-to-float loops (`unmarshalEmbedding`, `marshalEmbedding`) dominates CPU time.
**Action:** Use `unsafe.Slice` to alias memory for zero-copy hashing and embedding byte conversions. Added dynamic architecture checking `isLittleEndian` to ensure safe operation when directly copying memory for float representation.

## 2026-06-14 - Fix unsafe slice vulnerability
**Learning:** Using `unsafe.Slice` to map a byte array to float slice requires strict bounds checking, otherwise out-of-bounds dim leads to out-of-bounds heap reading/information disclosure vulnerabilities.
**Action:** Add bounds check `dim <= len(embedding)` when creating a fast-path slice with `unsafe.Slice`.
