package loader

import (
	"fmt"
	"os"
	"testing"

	"github.com/yigithanbalci/gonfik/pkg/loader"
)

// TestSimpleConfig calls blah blah
func TestSimpleConfig(t *testing.T) {
	err := os.Setenv("", "")
	if err != nil {
		t.Fatalf("Error setting environment variable: ", err)
		return
	}
	want := "config"
	konfik := loader.GlobalConfig()

}
