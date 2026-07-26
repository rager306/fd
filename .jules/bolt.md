## 2024-07-26 - [Go Memory Leak with String Sub-slicing]
**Learning:** Taking a sub-slice of a string (e.g., `hex.EncodeToString(h)[:12]`) keeps the entire original backing array alive in memory, preventing GC. In high-throughput cache key paths, this creates substantial memory pressure.
**Action:** Always encode hashes into a stack-allocated buffer (`var dst [64]byte`) and cast only the needed portion to a string (`string(dst[:12])`) to avoid pinning the larger array. Use `unsafe.Slice(unsafe.StringData(str), len(str))` for zero-allocation string-to-byte reads.
