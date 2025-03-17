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

func NewerThan(file1 string, file2 string) (bool, error) {
	stat1, err := os.Stat(file1)
	if err != nil {
		return false, nil // If we can't stat file 1, return false
	}

	stat2, err := os.Stat(file2)
	if err != nil {
		return false, err
	}
	return stat1.ModTime().After(stat2.ModTime()), nil
}
