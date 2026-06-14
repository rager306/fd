package cache

import "strconv"

// localCacheKey returns a string of the form "key:ddim".
// It is optimized for the common case where dim is 1024 or 512.
func localCacheKey(key string, dim int) string {
	if dim == 1024 {
		return key + ":d1024"
	}
	if dim == 512 {
		return key + ":d512"
	}
	return key + ":d" + strconv.Itoa(dim)
}
