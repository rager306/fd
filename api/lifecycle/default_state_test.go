package lifecycle

import "testing"

func TestDefaultStateReturnsSingleton(t *testing.T) {
	first := DefaultState()
	second := DefaultState()
	if first == nil {
		t.Fatal("DefaultState returned nil")
	}
	if first != second {
		t.Fatal("DefaultState did not return the same singleton")
	}
}
