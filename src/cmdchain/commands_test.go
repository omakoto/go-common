package cmdchain

import (
	"github.com/omakoto/go-common/src/common"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

//	func TestSampleExecCode(t *testing.T) {
//		var err error
//		cat := exec.Command("cat", "-An")
//		//cat.Stdout = os.Stdout
//		cat.Stderr = os.Stderr
//		cat.Stdin = MustOpenForRead("/etc/fstab")
//
//		grep := exec.Command("grep", "#xx")
//		grep.Stdin, err = cat.StdoutPipe()
//		grep.Stdout = os.Stdout
//		grep.Stderr = os.Stderr
//
//		common.Checke(err)
//
//		common.Checke(cat.Start())
//		common.Checke(grep.Start())
//
//		CheckWaitError(cat.Wait())
//		CheckWaitError(grep.Wait())
//	}

func mustMakeTempFile(content string) string {
	wr, err := os.CreateTemp("/tmp", "temp*.txt")
	common.Check(err, "openForWrite failed")
	defer wr.Close()

	wr.Write([]byte(content))

	return wr.Name()
}

func mustReadAllAsString(rd io.Reader) string {
	out, err := io.ReadAll(rd)
	common.Check(err, "ReadAll failed")
	return string(out)
}

func mustReadAllFileAsString(filename string) string {
	rd, err := os.Open(filename)
	common.Check(err, "Open failed")
	out, err := io.ReadAll(rd)
	common.Check(err, "ReadAll failed")
	return string(out)
}

func TestBasic(t *testing.T) {
	New().Command("echo", "ok").MustRunAndWait()

	{
		rd, cw := New().Command("echo", "ok").MustRunAndGetReader()

		assert.Equal(t, "ok\n", mustReadAllAsString(rd))

		cw.MustWait()
	}

	{
		rd, cw := New().Command("bash", "-c", "echo ok; exit 3").MustRunAndGetReader()
		//defer assert.NoErrorf(t, rd.Close(), "Close failed") // not needed.

		assert.Equal(t, "ok\n", mustReadAllAsString(rd))

		_, err := cw.Wait() // TODO: check status code
		assert.Equal(t, 3, extractStatusCode(err))
	}

	{
		out := New().Command("bash", "-c", "echo ok").MustRunAndGetString()
		assert.Equal(t, "ok\n", out)
	}

	{
		out := New().Command("bash", "-c", "echo out; echo err 1>&2").MustRunAndGetString()
		assert.Equal(t, "out\n", out)
	}

	{
		out := New().Command("bash", "-c", "echo out; echo err 1>&2").ErrToOut().MustRunAndGetString()
		assert.Equal(t, "out\nerr\n", out)
	}

	{
		out := New().Command("bash", "-c", "echo out; echo err 1>&2").ErrToOut().PipeOutTo().Command("wc", "-l").MustRunAndGetString()
		assert.Equal(t, "2\n", out)
	}

	{
		out := New().Command("bash", "-c", "echo out; echo err 1>&2").ErrToOut().MustRunAndGetStrings()
		assert.New(t).ElementsMatch(out, []string{"out", "err"})
	}

	{
		temp := mustMakeTempFile("abc\ndef\n")
		out := WithStdInFile(temp).Command("cat", "-An").MustRunAndGetString()
		assert.Equal(t, "     1\tabc$\n     2\tdef$\n", out)
	}

	{
		out := WithStdInString("abc\ndef\n").Command("cat", "-An").MustRunAndGetString()
		assert.Equal(t, "     1\tabc$\n     2\tdef$\n", out)
	}

	{
		temp1 := mustMakeTempFile("abc\ndef\n")
		temp2 := mustMakeTempFile("")
		WithStdInFile(temp1).Command("cat", "-An").SetStdOutFile(temp2).MustRunAndWait()
		assert.Equal(t, "     1\tabc$\n     2\tdef$\n", mustReadAllFileAsString(temp2))
	}
}
