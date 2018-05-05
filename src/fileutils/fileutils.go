package fileutils

import "os"

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func DirExists(file string) bool {
	stat, err := os.Stat(file)
	return err == nil && ((stat.Mode() & os.ModeDir) != 0)
}

