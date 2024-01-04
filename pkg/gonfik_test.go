package pkg

import (
	"os"
	"testing"
)

// TestSimpleConfig calls blah blah
func TestSimpleConfig(t *testing.T) {
	err := os.Setenv("GONFIK_DIR", "pkg")
	if err != nil {
		t.Fatalf("Error setting environment variable: %v", err)
		return
	}
	err = os.Setenv("GONFIK_PROD_FILE", "gonfik_test.json")
	if err != nil {
		t.Fatalf("Error setting environment variable: %v", err)
		return
	}
	want := "config"
	konfik, err := GlobalConfig()
	if err != nil {
		t.Fatalf("Error getting global config: %v", err)
		return
	}
	msg, isOk := konfik.Config("test.foo.bar")
	if msg != want || isOk != true {
		t.Fatalf(`konfik.Config("test.foo.bar") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}
