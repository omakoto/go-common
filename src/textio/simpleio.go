package textio

import (
	"io"
	"os"
)

func WriteStringToFile(file string, s string, mode os.FileMode) error {
	wr, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer wr.Close()
	_, err = wr.WriteString(s)
	if err != nil {
		return err
	}
	return nil
}

func ReadStringFromFile(file string) (string, error) {
	rd, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer rd.Close()
	b, err := io.ReadAll(rd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
