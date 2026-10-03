package cache

import "strconv"

// localCacheKey returns a string of the form "key:ddim".
// It is optimized for the common case where dim is 1024, 768, 512, or 256.
func localCacheKey(key string, dim int) string {
	if dim == 1024 {
		return key + ":d1024"
	}
	if dim == 768 {
		// Fast paths for common dimensions to avoid strconv.Itoa heap allocations
		return key + ":d768"
	}
	if dim == 512 {
		return key + ":d512"
	}
	if dim == 256 {
		// Fast paths for common dimensions to avoid strconv.Itoa heap allocations
		return key + ":d256"
	}
	return key + ":d" + strconv.Itoa(dim)
}
