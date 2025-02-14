package fileutils

import (
	"os"
)

func FileExists(file string) bool {
	stat, err := os.Stat(file)
	return err == nil && ((stat.Mode() & os.ModeDir) == 0)
}

func DirExists(file string) bool {
	stat, err := os.Stat(file)
	return err == nil && ((stat.Mode() & os.ModeDir) != 0)
}
