package main

import "testing"

func TestGetEnvIntReturnsDefaultWhenUnset(t *testing.T) {
	t.Setenv("FD_TEST_INT", "")

	got := getEnvInt("FD_TEST_INT", 50)
	if got != 50 {
		t.Fatalf("getEnvInt unset = %d, want 50", got)
	}
}

func TestGetEnvIntParsesPositiveInteger(t *testing.T) {
	t.Setenv("FD_TEST_INT", "75")

	got := getEnvInt("FD_TEST_INT", 50)
	if got != 75 {
		t.Fatalf("getEnvInt = %d, want 75", got)
	}
}

func TestGetEnvIntReturnsDefaultForInvalidValue(t *testing.T) {
	t.Setenv("FD_TEST_INT", "12x")

	got := getEnvInt("FD_TEST_INT", 50)
	if got != 50 {
		t.Fatalf("getEnvInt invalid = %d, want 50", got)
	}
}
