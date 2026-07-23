package embed
import (
	"testing"
	"unsafe"
)
func TestUnsafeEmpty(t *testing.T) {
	var s1 []float32
	var s2 = make([]float32, 0)
	p1 := unsafe.SliceData(s1)
	p2 := unsafe.SliceData(s2)
	_ = p1
	_ = p2
}
