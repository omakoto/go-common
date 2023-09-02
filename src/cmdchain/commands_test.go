package cmdchain

import (
	"github.com/omakoto/go-common/src/common"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func mustMakeTempFile(content string) string {
	wr, err := os.CreateTemp("/tmp", "temp*.txt")
	common.Check(err, "openForWrite failed")

	_, err = wr.Write([]byte(content))
	common.Check(err, "Write failed")

	err = wr.Close()
	common.Check(err, "Close failed")

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

		_, err := cw.Wait()
		assert.Equal(t, 3, extractStatusCode(err))
	}

	{
		var status int
		out := New().Command("bash", "-c", "echo ok; exit 3").AllowStatus(&status, 3).MustRunAndGetString()
		assert.Equal(t, "ok\n", out)
		assert.Equal(t, 3, status)
	}

	{
		var status int
		_, cw := New().Command("bash", "-c", "echo ok; exit 2").AllowStatus(&status, 3, 4).MustRunAndGetReader()
		_, err := cw.Wait()
		assert.Equal(t, 2, extractStatusCode(err))
		assert.Equal(t, 2, status)
	}

	{
		var status int
		out := New().Command("bash", "-c", "echo ok; exit 33").AllowAnyStatus(&status).MustRunAndGetString()
		assert.Equal(t, "ok\n", out)
		assert.Equal(t, 33, status)
	}

	{
		out := New().Command("bash", "-c", "echo ok; exit 1").AllowAnyStatus(nil).MustRunAndGetString()
		assert.Equal(t, "ok\n", out)
	}

	{
		out := New().Command("bash", "-c", "echo ok; exit 1").AllowStatus(nil, 1).MustRunAndGetString()
		assert.Equal(t, "ok\n", out)
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
		out := New().Command("bash", "-c", "echo out; echo err 1>&2").ErrToOut().Pipe().Command("wc", "-l").MustRunAndGetString()
		assert.Equal(t, "2\n", out)
	}

	{
		// Same as above but ErrToOut() is at a different position.
		out := New().Command("bash", "-c", "echo out; echo err 1>&2").Pipe().ErrToOut().Command("wc", "-l").MustRunAndGetString()
		assert.Equal(t, "2\n", out)
	}

	{
		out := New().Command("bash", "-c", "echo out; echo err 1>&2").ErrToOut().MustRunAndGetStrings()
		assert.New(t).ElementsMatch(out, []string{"out", "err"})
	}

	{
		i := New().Command("bash", "-c", "echo out; echo err 1>&2").ErrToOut().MustRunAndGetStringsIter()
		assert.Equal(t, "out", *i())
		assert.Equal(t, "err", *i())
		assert.Equal(t, (*string)(nil), i())
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
		WithStdInFile(temp1).Command("cat", "-An").SetStdoutFile(temp2).SetStdoutFile(temp2).MustRunAndWait()
		assert.Equal(t, "     1\tabc$\n     2\tdef$\n", mustReadAllFileAsString(temp2))
	}

	{
		out := WithStdInBytes([]byte("abc\ndef\n")).Command("cat", "-An").MustRunAndGetString()
		assert.Equal(t, "     1\tabc$\n     2\tdef$\n", out)
	}

	{
		_, err := WithStdInString("abc").Command("cat").Command("cat", "-An").Run()
		assert.ErrorContains(t, err, "duplicate command")
	}

	{
		assert.PanicsWithValue(t, "No command is set yet.", func() {
			WithStdInString("abc").Run()
		}, "Expected panic")
	}

	{
		assert.PanicsWithValue(t, "Expecting next command to consume stdin", func() {
			WithStdInString("abc").Command("cat").Pipe().Run()
		}, "Expected panic")
	}

	{
		assert.PanicsWithValue(t, "No command is set yet.", func() {
			New().Pipe().Run()
		}, "Expected panic")
	}

	{
		assert.PanicsWithValue(t, "No command is set yet.", func() {
			WithStdInString("abc").Pipe().Run()
		}, "Expected panic")
	}

	{
		erd := NewBytesReader()
		rd, cw := New().Command("bash", "-c", "echo out1; echo err1 1>&2; exit 0").SaveStderr(erd).MustRunAndGetReader()
		assert.Equal(t, "out1\n", mustReadAllAsString(rd))

		cw.MustWait()

		assert.Equal(t, "err1\n", string(erd.Get()))
	}

	{
		erd1 := NewBytesReader()
		erd2 := NewBytesReader()
		rd, cw := New().Command("bash", "-c", "echo out1; echo err1 1>&2; exit 0").SaveStderr(erd1).Pipe().Command("bash", "-c", "cat ; echo out2; echo err2 1>&2; exit 0").SaveStderr(erd2).MustRunAndGetReader()

		assert.Equal(t, "out1\nout2\n", mustReadAllAsString(rd))

		cw.MustWait()

		assert.Equal(t, "err1\n", string(erd1.Get()))
		assert.Equal(t, "err2\n", string(erd2.Get()))
	}

	{
		var rd *io.ReadCloser
		cw, err := New().Command("bash", "-c", "echo out").getStdoutPipe(&rd).Run()
		assert.NoErrorf(t, err, "Run")

		assert.Equal(t, "out\n", mustReadAllAsString(*rd))

		cw.MustWait()
	}

	{
		env := map[string]string{
			"ARG1": "arg1",
			"ARG2": "arg2",
		}
		out := New().CommandWithEnv(env, "bash", "-c", "echo $ARG1 $ARG2").MustRunAndGetString()
		assert.Equal(t, "arg1 arg2\n", out)
	}

	// This doesn't work. Not sure why.
	//{
	//	var start, end time.Time
	//
	//	con := func(s string) {
	//		switch s {
	//		case "out1":
	//			start = time.Now()
	//		case "out2":
	//			end = time.Now()
	//		}
	//	}
	//	New().Command("bash", "-c", "stdbuf -oL /usr/bin/echo out1; sleep 0.5; stdbuf -oL /usr/bin/echo out2").MustRunAndStreamStringsBufSize(con, 1)
	//
	//	assert.LessOrEqualf(t, time.Duration(500_000_000), end.Sub(start), "Didn't received the strings 500ms apart. (start=%s, end=%s)", start, end)
	//}
}
