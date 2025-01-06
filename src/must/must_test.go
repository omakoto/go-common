package must

import (
	"os"
	"testing"
)

func TestMust(t *testing.T) {
	Must(nil)
}

func TestMust2(t *testing.T) {
	in := Must2(os.Open("/dev/null"))
	if in == nil {
		t.Fatalf("in is nil")
	}
}

func TestWith2(t *testing.T) {
	file := "/dev/null"
	in := With2(os.Open(file)).Checkf("unable to open file '%s'", file)
	if in == nil {
		t.Fatalf("in is nil")
	}
}
