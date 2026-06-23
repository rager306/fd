## 2024-05-15 - Prevent Length-Based Timing Attacks in ConstantTimeCompare
**Vulnerability:** Comparing an arbitrary-length user-provided token against a secret using `subtle.ConstantTimeCompare` without hashing.
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the two byte slices differ. This exposes the length of the secret API key to a timing attack, potentially allowing an attacker to brute-force the secret more efficiently or identify the secret's format.
**Prevention:** Always hash both inputs (e.g., with `crypto/sha256`) before passing them to `subtle.ConstantTimeCompare` to guarantee they are the same length, regardless of the user's input.
