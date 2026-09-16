1. Update `HashText` in `api/cache/redis.go` to use zero-allocation hex encoding using stack arrays, avoiding `hex.EncodeToString` heap allocations.
2. Update `key` in `api/cache/redis.go` to use fast paths for common `dim` arguments (1024, 512), avoiding `strconv.Itoa` heap allocations.
3. Update `newRequestID` in `api/middleware/headers.go` to avoid `hex.EncodeToString` and slice concatenation which creates multiple heap allocations.
4. Run Go tests to make sure there are no breakages.
5. Add a journal entry to `.jules/bolt.md` documenting this CRITICAL learning on how to optimize string serialization in Go.
6. Commit the changes and open PR.
