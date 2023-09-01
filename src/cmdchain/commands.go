package cmdchain

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/textio"
	"io"
	"os"
	"os/exec"
	"strings"
)

func openForRead(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDONLY, 0)
}

func openForWrite(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0)
}

// MustOpenForRead opens a file for reading.
func MustOpenForRead(filename string) *os.File {
	ret, err := openForRead(filename)
	common.Check(err, fmt.Sprintf("Unable to open file %s for reading", filename))
	return ret
}

// MustOpenForWrite opens (or creates) a file for writing.
func MustOpenForWrite(filename string) *os.File {
	ret, err := openForWrite(filename)
	common.Check(err, fmt.Sprintf("Unable to open file %s for writing", filename))
	return ret
}

type StringConsumer func(s string)
type BytesConsumer func(data []byte)

type BytesReader struct {
	data []byte
}

func NewBytesReader() *BytesReader {
	return &BytesReader{}
}

func (b *BytesReader) Get() []byte {
	if b.data == nil {
		panic("CommandChain hasn't been waited yet.")
	}
	return b.data
}

type commandValidator func(cmd *exec.Cmd, waitError error) error

func extractStatusCode(waitError error) int {
	var e *exec.ExitError
	if errors.As(waitError, &e) {
		return e.ExitCode()
	}
	return -1
}

func standardValidator(cmd *exec.Cmd, waitError error) error {
	return waitError
}

func ensureNilAndSet[T any](a *T, value T, format string, args ...any) {
	if a == nil {
		panic(fmt.Sprintf(format, args...))
	}
	*a = value
}

func arrayGet[T any](array []T, index int) T {
	if index >= 0 {
		return array[index]
	}
	return array[len(array)+index]
}

func arraySet[T any](array []T, index int, value T) {
	if index >= 0 {
		array[index] = value
	} else {
		array[len(array)+index] = value
	}
}

// CommandChain is a chain of exec.Cmd.
type CommandChain struct {
	deferredError error
	defaultStdout io.Writer
	defaultStderr io.Writer

	nextStdin io.Reader

	prevErrToOut bool

	Commands     []*exec.Cmd
	validators   []commandValidator
	stderrReader []*BytesReader
}

// ChainWaiter is a handle that can be wait()'ed on.
type ChainWaiter struct {
	Chain *CommandChain
}

// ChainResult provides the overall result of the command chain.
// (For now it only provides a reference to the originating command CommandChain.)
type ChainResult struct {
	Chain *CommandChain
}

// New creates a new CommandChain.
func New() *CommandChain {
	return &CommandChain{}
}

// WithStdIn creates a new CommandChain, with a given io.Reader as strin.
func WithStdIn(reader io.Reader) *CommandChain {
	ret := New()
	ret.nextStdin = reader
	return ret
}

// WithStdInFile creates a new CommandChain, with a given file as stdin.
func WithStdInFile(filename string) *CommandChain {
	in, err := openForRead(filename)
	ret := WithStdIn(in)
	ret.setDeferredError(err)
	return ret
}

// WithStdInString creates a new CommandChain, with a given string as stdin.
func WithStdInString(text string) *CommandChain {
	return WithStdInBytes([]byte(text))
}

// WithStdInBytes creates a new CommandChain, with a given []byte as stdin.
func WithStdInBytes(data []byte) *CommandChain {
	return WithStdIn(bytes.NewReader(data))
}

func (c *CommandChain) getDefaultStdout() io.Writer {
	if c.defaultStdout != nil {
		return c.defaultStdout
	}
	return os.Stdout
}

func (c *CommandChain) getDefaultStderr() io.Writer {
	if c.defaultStderr != nil {
		return c.defaultStderr
	}
	return os.Stderr
}

func (c *CommandChain) getCommandDescription(index int) string {
	if index < 0 {
		index += len(c.Commands)
	}
	cmd := c.Commands[index]
	return fmt.Sprintf("\"%s\" at index %d in the chain", cmd.Path, index)
}

func (c *CommandChain) fixUpLastCommand() {
	cmd := c.lastCommand()
	if c.prevErrToOut {
		ensureNilAndSet(&cmd.Stderr, cmd.Stdout, "Stderr has already been set to command %s", c.getCommandDescription(-1))
	}
	if cmd.Stdout == nil {
		cmd.Stdout = c.getDefaultStdout()
	}
	if cmd.Stderr == nil {
		cmd.Stderr = c.getDefaultStderr()
	}
	if arrayGet(c.validators, -1) == nil {
		arraySet(c.validators, -1, standardValidator)
	}
}

func (c *CommandChain) setDeferredError(err error) *CommandChain {
	if c.deferredError == nil && err != nil {
		common.Warnf("Error detected: %v", err)
		c.deferredError = err
	}
	return c
}

func (c *CommandChain) lastCommand() *exec.Cmd {
	if len(c.Commands) == 0 {
		panic("No command is set yet.")
	}
	return c.Commands[len(c.Commands)-1]
}

// Command adds a new command to a CommandChain.
func (c *CommandChain) Command(name string, args ...string) *CommandChain {
	common.Debugf("Command: %s", name)

	cmd := exec.Command(name, args...)

	if len(c.Commands) > 0 {
		c.fixUpLastCommand()
	}

	c.Commands = append(c.Commands, cmd)
	c.validators = append(c.validators, nil)
	c.stderrReader = append(c.stderrReader, nil)

	if c.nextStdin != nil {
		cmd.Stdin = c.nextStdin
		c.nextStdin = nil
	} else {
		if len(c.Commands) == 1 {
			cmd.Stdin = os.Stdin
		} else {
			c.deferredError = fmt.Errorf("duplicate command \"%s\" detected without a pipe", name)
		}
	}

	return c
}

// CommandWithEnv adds a new command to a CommandChain with environmental variables.
func (c *CommandChain) CommandWithEnv(env map[string]string, name string, args ...string) *CommandChain {
	e := make([]string, len(env))

	i := 0
	for name, value := range env {
		e[i] = name + "=" + value
		i++
	}

	c.lastCommand().Env = e
	return c
}

func (c *CommandChain) setValidator(validator commandValidator) {
	ensureNilAndSet(&c.validators[len(c.validators)-1], validator, "command status code validator is already set to command %s", c.getCommandDescription(-1))
}

// AllowAnyStatus will allow the previous command to return any exit status code.
func (c *CommandChain) AllowAnyStatus(actualStatus *int) *CommandChain {
	c.setValidator(func(cmd *exec.Cmd, waitError error) error {
		status := extractStatusCode(waitError)
		if status < 0 {
			return waitError
		}
		if actualStatus != nil {
			*actualStatus = status
		}
		return nil
	})
	return c
}

// AllowStatus will allow the previous command to return any of specified exit status codes.
// At least one status code must be provided.
func (c *CommandChain) AllowStatus(actualStatus *int, allowed ...int) *CommandChain {
	if len(allowed) == 0 {
		panic("AllowStatus expects 1 or more allowed status codes.")
	}
	c.setValidator(func(cmd *exec.Cmd, waitError error) error {
		status := extractStatusCode(waitError)
		if status < 0 {
			return waitError
		}
		if actualStatus != nil {
			*actualStatus = status
		}
		for _, a := range allowed {
			if a == status {
				return nil
			}
		}
		return waitError
	})
	return c
}

// SetDefaultOut Set the default stdout to the subsequent commands.
func (c *CommandChain) SetDefaultOut(writer io.Writer) *CommandChain {
	c.defaultStdout = writer
	return c
}

// SetDefaultErr Set the default stderr to the subsequent commands.
func (c *CommandChain) SetDefaultErr(writer io.Writer) *CommandChain {
	c.defaultStderr = writer
	return c
}

// SetStdout sets a writer to the stdout of the last command.
func (c *CommandChain) SetStdout(writer io.Writer) *CommandChain {
	ensureNilAndSet(&c.lastCommand().Stdout, writer, "Stdout already set to command %s", c.getCommandDescription(-1))
	return c
}

// SetStderr sets a writer to the stderr of the last command.
func (c *CommandChain) SetStderr(writer io.Writer) *CommandChain {
	ensureNilAndSet(&c.lastCommand().Stderr, writer, "Stderr already set to command %s", c.getCommandDescription(-1))
	return c
}

// ErrToOut redirects stderr of the last command to stdout
func (c *CommandChain) ErrToOut() *CommandChain {
	if c.prevErrToOut {
		panic(fmt.Sprintf("ErrToOut has already been called on command %s", c.getCommandDescription(-1)))
	}
	c.prevErrToOut = true
	return c
}

// SetStdoutFile sets a file to the stdout of the last command.
func (c *CommandChain) SetStdoutFile(filename string) *CommandChain {
	f, err := openForWrite(filename)
	if err != nil {
		c.SetStdout(io.Writer(f))
	}
	c.setDeferredError(err)
	return c
}

// SetStderrFile sets a file to the stderr of the last command.
func (c *CommandChain) SetStderrFile(filename string) *CommandChain {
	f, err := openForWrite(filename)
	if err != nil {
		c.SetStderr(io.Writer(f))
	}
	c.setDeferredError(err)
	return c
}

// getStdoutPipe gets a pipe from stdout of the last command.
// This should only be used internally.
func (c *CommandChain) getStdoutPipe(reader **io.ReadCloser) *CommandChain {
	cmd := c.lastCommand()
	p, err := cmd.StdoutPipe()
	if err != nil {
		c.setDeferredError(fmt.Errorf("StdoutPipe() failed on %s: %w", c.getCommandDescription(-1), err))
		return c
	}
	*reader = &p
	return c
}

// SaveStderr saves stderr to a tempfile and sends it to con when all the commands are done.
func (c *CommandChain) SaveStderr(r *BytesReader) *CommandChain {
	temp, err := os.CreateTemp(os.TempDir(), "stderr*.dat")
	if err != nil {
		c.setDeferredError(fmt.Errorf("CreateTemp() failed on %s: %w", c.getCommandDescription(-1), err))
		return c
	}
	c.SetStderr(temp)
	arraySet(c.stderrReader, -1, r)
	return c
}

func (c *CommandChain) setNextStdin(rd io.ReadCloser) {
	ensureNilAndSet(&c.nextStdin, io.Reader(rd), "Pipe() has already been called on command %s", c.getCommandDescription(-1))
}

// Pipe prepares a pipe from stdout of the last command to the next command.
// It should be followed by Command()
func (c *CommandChain) Pipe() *CommandChain {
	var rd *io.ReadCloser
	c.getStdoutPipe(&rd)
	c.setNextStdin(*rd)
	return c
}

func (c *CommandChain) validateBeforeRun() error {
	if c.deferredError != nil {
		return c.deferredError
	}
	if len(c.Commands) == 0 {
		panic("Must have at least 1 command")
	}
	if c.nextStdin != nil {
		panic("Expecting next command to consume stdin")
	}
	return nil
}

// Run starts a CommandChain.
func (c *CommandChain) Run() (*ChainWaiter, error) {
	err := c.validateBeforeRun()
	if err != nil {
		return nil, err
	}
	c.fixUpLastCommand()

	for i, cmd := range c.Commands {
		err := cmd.Start()
		if err != nil {
			return nil, fmt.Errorf("unable to execute command \"%s\" (command #%d): %s", cmd.Path, i+1, err.Error())
		}
	}
	return &ChainWaiter{Chain: c}, nil
}

// MustRun starts a CommandChain.
func (c *CommandChain) MustRun() *ChainWaiter {
	cw, err := c.Run()
	common.Check(err, "Unable to execute command(s)")
	return cw
}

// MustRunAndWait starts a CommandChain and wait().
func (c *CommandChain) MustRunAndWait() *ChainResult {
	return c.MustRun().MustWait()
}

// MustRunAndGetReader starts a CommandChain and get stdout of the last command as an io.Reader.
func (c *CommandChain) MustRunAndGetReader() (io.Reader, *ChainWaiter) {
	err := c.validateBeforeRun()
	common.Check(err, "Unable to execute command(s)")

	var rd *io.ReadCloser
	c.getStdoutPipe(&rd)

	cw := c.MustRun()

	return *rd, cw
}

const defaultBufSize = 4096

// MustRunAndGetBufferedReader starts a CommandChain and get stdout of the last command as an bufio.Reader.
func (c *CommandChain) MustRunAndGetBufferedReader() (*bufio.Reader, *ChainWaiter) {
	return c.MustRunAndGetBufferedReaderBufSize(defaultBufSize)
}

// MustRunAndGetBufferedReaderBufSize starts a CommandChain and get stdout of the last command as an bufio.Reader.
func (c *CommandChain) MustRunAndGetBufferedReaderBufSize(bufSize int) (*bufio.Reader, *ChainWaiter) {
	rd, cw := c.MustRunAndGetReader()
	return bufio.NewReaderSize(rd, bufSize), cw
}

// MustRunAndGetBytes starts a CommandChain and return stdout of the last command as []byte,
// and it also calls MustWait().
func (c *CommandChain) MustRunAndGetBytes() []byte {
	rd, waiter := c.MustRunAndGetReader()

	data, err := io.ReadAll(rd)
	common.Check(err, "Error while reading from commands")

	waiter.MustWait()

	return data
}

// MustRunAndGetString starts a CommandChain and return stdout of the last command as string,
// and it also calls MustWait().
func (c *CommandChain) MustRunAndGetString() string {
	return string(c.MustRunAndGetBytes())
}

// MustRunAndGetStrings starts a CommandChain and return stdout of the last command as []string,
// and it also calls MustWait().
func (c *CommandChain) MustRunAndGetStrings() []string {
	return strings.Split(textio.StringChomp(c.MustRunAndGetString()), "\n")
}

// MustRunAndStreamStrings starts a CommandChain, read stdout of the last command line by line
// and feed them to con.
func (c *CommandChain) MustRunAndStreamStrings(con StringConsumer) {
	c.MustRunAndStreamStringsBufSize(con, defaultBufSize)
}

// MustRunAndStreamStringsBufSize starts a CommandChain, read stdout of the last command line by line
// and feed them to con, using bufSize for buffered reading.
func (c *CommandChain) MustRunAndStreamStringsBufSize(con StringConsumer, bufSize int) {
	rd, cw := c.MustRunAndGetBufferedReaderBufSize(bufSize)
	defer cw.MustWait()
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Sprintf("ReadString failed. %s", err))
		}
		con(line)
	}
}

// MustRunAndStreamBytes starts a CommandChain, read stdout of the last command line by line
// and feed them to con.
func (c *CommandChain) MustRunAndStreamBytes(con BytesConsumer) {
	c.MustRunAndStreamBytesBufSize(con, defaultBufSize)
}

// MustRunAndStreamBytesBufSize starts a CommandChain, read stdout of the last command line by line
// and feed them to con, using bufSize for buffered reading.
func (c *CommandChain) MustRunAndStreamBytesBufSize(con BytesConsumer, bufSize int) {
	rd, cw := c.MustRunAndGetBufferedReaderBufSize(bufSize)
	defer cw.MustWait()
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Sprintf("ReadString failed. %s", err))
		}
		con(line)
	}
}

// Wait wait() on all commands in a CommandChain.
func (cw *ChainWaiter) Wait() (*ChainResult, error) {

	var firstError error
	for i, cmd := range cw.Chain.Commands {
		err := cmd.Wait()
		err = cw.Chain.validators[i](cmd, err)

		if err != nil {
			if firstError == nil {
				firstError = fmt.Errorf("failed to wait on command %s: %w", cmd.Path, err)
			}
			continue
		}

		// See if there's any stderr consumers.
		ser := cw.Chain.stderrReader[i]
		if ser != nil {
			errf := cmd.Stderr.(*os.File)
			errf.Seek(0, 0)

			data, err := io.ReadAll(errf)
			if err != nil {
				if firstError == nil {
					firstError = fmt.Errorf("failed to read from tempfile %s: %w", errf.Name(), err)
					continue
				}
			}
			ser.data = data
		}
	}
	if firstError != nil {
		return nil, firstError
	}

	return &ChainResult{Chain: cw.Chain}, nil
}

// MustWait wait() on all commands in a CommandChain.
func (cw *ChainWaiter) MustWait() *ChainResult {
	cr, err := cw.Wait()
	common.Checke(err)
	return cr
}
