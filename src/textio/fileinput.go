package textio

import (
	"io"
	"os"
)

type ReadFunc func(line []byte, lineNo int, filename string) error

func ReadReader(filename string, r io.Reader, f ReadFunc) error {
	n := 0
	lr := NewLineReader(r, false)
	for {
		n++
		line, err := lr.ReadLine()
		if line != nil {
			ferr := f(line, n, filename)
			if ferr != nil {
				return ferr
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func ReadFile(filename string, f ReadFunc) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	return ReadReader(filename, r, f)
}

func ReadFiles(files []string, f ReadFunc) error {
	if len(files) == 0 {
		return ReadReader("-", os.Stdin, f)
	} else {
		for _, file := range files {
			err := ReadFile(file, f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
