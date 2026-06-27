package embed

import "unsafe"

var isLittleEndian bool

func init() {
	var i int32 = 0x01020304
	//nolint:gosec // G103: used to detect architecture endianness at startup
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	isLittleEndian = (b == 0x04)
}
