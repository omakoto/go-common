package fileinput

import (
	"bufio"
	"iter"
	"os"

	"github.com/omakoto/go-common/src/common"
)

type Options struct {
	InlineReplace bool
	Files         []string
}

type FileInfo struct {
	// filename of the current input file.
	Filename string

	// line number within the file. (0 based)
	Line int
}

type FileInputSeq func(yield func(text, file string, line int) bool)

func FileInput() iter.Seq2[string, FileInfo] {
	return FileInputOption(Options{})
}

func FileInputOption(options Options) iter.Seq2[string, FileInfo] {
	// inline := options.InlineReplace
	files := options.Files
	if len(files) == 0 {
		files = os.Args[1:]
		if len(files) == 0 {
			files = append(files, "/dev/stdin")
		}
	}

	return func(yield func(text string, info FileInfo) bool) {
		for _, file := range files {
			in, err := os.Open(file)
			common.Checkf(err, "Unable to open file '%s'", file)
			defer in.Close()

			sc := bufio.NewScanner(in)

			fi := FileInfo{file, 0}

			for sc.Scan() {
				if !yield(sc.Text(), fi) {
					return
				}
				fi.Line++
			}
		}
	}
}
