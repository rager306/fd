package endian

import "unsafe"

// IsLittleEndian is true if the system architecture is little-endian.
var IsLittleEndian bool

func init() {
	var i int32 = 0x01020304
	//nolint:gosec // G103: checking system endianness
	if *(*byte)(unsafe.Pointer(&i)) == 0x01 {
		IsLittleEndian = false
	} else {
		IsLittleEndian = true
	}
}
