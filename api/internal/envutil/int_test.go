package envutil

import (
	"strconv"
	"strings"
	"testing"
)

func TestIntParsesNonNegativeValues(t *testing.T) {
	t.Setenv("FD_TEST_ENV_INT", "42")
	if got := Int("FD_TEST_ENV_INT", -1); got != 42 {
		t.Fatalf("Int = %d, want 42", got)
	}

	t.Setenv("FD_TEST_ENV_INT", "0")
	if got := Int("FD_TEST_ENV_INT", -1); got != 0 {
		t.Fatalf("Int zero = %d, want 0", got)
	}
}

func TestIntFallsBackForInvalidValues(t *testing.T) {
	for _, value := range []string{"", "-1", "12ms", strings.Repeat("9", 100)} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("FD_TEST_ENV_INT", value)
			if got := Int("FD_TEST_ENV_INT", 7); got != 7 {
				t.Fatalf("Int(%q) = %d, want fallback", value, got)
			}
		})
	}
}

func TestIntTrimsWhitespace(t *testing.T) {
	t.Setenv("FD_TEST_ENV_INT", " 12 ")
	if got := Int("FD_TEST_ENV_INT", 7); got != 12 {
		t.Fatalf("Int trimmed = %d, want 12", got)
	}
}

func TestPositiveIntFallsBackForZeroAndInvalidValues(t *testing.T) {
	for _, value := range []string{"", "0", "-1", "bad"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("FD_TEST_POSITIVE_INT", value)
			if got := PositiveInt("FD_TEST_POSITIVE_INT", 9); got != 9 {
				t.Fatalf("PositiveInt(%q) = %d, want fallback", value, got)
			}
		})
	}
}

func TestPositiveIntParsesPositiveValues(t *testing.T) {
	t.Setenv("FD_TEST_POSITIVE_INT", strconv.Itoa(11))
	if got := PositiveInt("FD_TEST_POSITIVE_INT", 9); got != 11 {
		t.Fatalf("PositiveInt = %d, want 11", got)
	}
}
