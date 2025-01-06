package common

// import (
// 	"os"
// 	"testing"
// )

// func TestMust(t *testing.T) {
// 	Must(nil)
// }

// func TestMust2(t *testing.T) {
// 	in := Must2(os.Open("/dev/null"))
// 	if in == nil {
// 		t.Fatalf("in is nil")
// 	}
// }

// func TestMust2WithMessage(t *testing.T) {
// 	in := WithMessage2("failed!").Must2(os.Open("/dev/null"))
// 	if in == nil {
// 		t.Fatalf("in is nil")
// 	}
// }
