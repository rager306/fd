package envutil

import "testing"

func TestBoolOrDefaultFallbackWhenUnset(t *testing.T) {
	// Use a key we are confident is not set in any test runner.
	const key = "FD_TEST_BOOL_UNSET_KEY"
	t.Setenv(key, "")
	if got := BoolOrDefault(key, true); got != true {
		t.Fatalf("unset fallback=true got %v", got)
	}
	if got := BoolOrDefault(key, false); got != false {
		t.Fatalf("unset fallback=false got %v", got)
	}
}

func TestBoolOrDefaultTruthyTokens(t *testing.T) {
	cases := []string{"1", "true", "TRUE", "yes", "Yes", "on", "y", "t"}
	for _, value := range cases {
		t.Run(value, func(t *testing.T) {
			t.Setenv("FD_TEST_BOOL_OK", value)
			if got := BoolOrDefault("FD_TEST_BOOL_OK", false); got != true {
				t.Fatalf("value=%q want true got %v", value, got)
			}
		})
	}
}

func TestBoolOrDefaultFalsyTokens(t *testing.T) {
	cases := []string{"0", "false", "FALSE", "no", "off", "n", "f"}
	for _, value := range cases {
		t.Run(value, func(t *testing.T) {
			t.Setenv("FD_TEST_BOOL_OFF", value)
			if got := BoolOrDefault("FD_TEST_BOOL_OFF", true); got != false {
				t.Fatalf("value=%q want false got %v", value, got)
			}
		})
	}
}

func TestBoolOrDefaultUnknownPreservesFallback(t *testing.T) {
	t.Setenv("FD_TEST_BOOL_BOGUS", "perhaps")
	if got := BoolOrDefault("FD_TEST_BOOL_BOGUS", true); got != true {
		t.Fatalf("bogus+fallback=true want true got %v", got)
	}
	if got := BoolOrDefault("FD_TEST_BOOL_BOGUS", false); got != false {
		t.Fatalf("bogus+fallback=false want false got %v", got)
	}
}

func TestBoolOrDefaultWhitespaceAndCase(t *testing.T) {
	t.Setenv("FD_TEST_BOOL_WS", "  Yes  ")
	if got := BoolOrDefault("FD_TEST_BOOL_WS", false); got != true {
		t.Fatalf("padded yes want true got %v", got)
	}
}
