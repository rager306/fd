## 2026-07-28 - String Slicing Memory Leaks
**Learning:** Slicing a string (e.g., `hex.EncodeToString(h[:])[:12]`) keeps the entire original backing array alive in Go's memory. This was causing a memory leak for long-lived cache keys because a 64-byte array was retained in memory just to hold 12 characters.
**Action:** When truncating encoded output, encode only the required prefix into a precisely-sized stack-allocated buffer (e.g., `var dst [12]byte`), encode into it, and then cast it directly to a string. This avoids the memory leak and minimizes heap allocations.
