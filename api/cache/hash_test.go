package cache

import (
	"testing"
)

func TestShortHash(t *testing.T) {
	if shortHash("test") != "9f86d081884c" {
		t.Fail()
	}
}
