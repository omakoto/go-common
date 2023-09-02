package cmdchain

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/textio"
	"github.com/omakoto/go-common/src/utils"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
)

// TODO Clean up temp files.

const (
	StateBuilding = iota
	StateRunning
	StateWaiting
	StateSucceeded // Successfully finished waiting on all commands
	StateFailed    // Failed in or after Run().
)

func openForRead(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDONLY, 0)
}

func openForWrite(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0)
}

// MustOpenForRead opens a file for reading.
func MustOpenForRead(filename string) *os.File {
	ret, err := openForRead(filename)
	common.CheckPanicf(err, "Unable to open file %s for reading", filename)
	return ret
}

// MustOpenForWrite opens (or creates) a file for writing.
func MustOpenForWrite(filename string) *os.File {
	ret, err := openForWrite(filename)
	common.Checkf(err, "Unable to open file %s for writing", filename)
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

func standardValidator(_ *exec.Cmd, waitError error) error {
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
	state int32

	deferredError error
	defaultStdout io.Writer
	defaultStderr io.Writer

	nextStdin io.Reader

	prevErrToOut bool

	Commands     []*exec.Cmd
	validators   []commandValidator
	stderrReader []*BytesReader

	tempFiles []*os.File

	cleanupMu sync.Mutex
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

func (c *CommandChain) ensureBuilding() {
	if atomic.LoadInt32(&(c.state)) != StateBuilding {
		panic("Invalid operation on CommandChain. It's already running.")
	}
}

func (c *CommandChain) moveToRunning() {
	if !atomic.CompareAndSwapInt32(&(c.state), StateBuilding, StateRunning) {
		panic("Invalid operation Run() on CommandChain. It's already running.")
	}
}

func (c *CommandChain) moveToWaiting() {
	if !atomic.CompareAndSwapInt32(&(c.state), StateRunning, StateWaiting) {
		panic("Invalid operation Wait() on CommandChain. It's not running.")
	}
}

func (c *CommandChain) moveToSucceeded() {
	if !atomic.CompareAndSwapInt32(&(c.state), StateWaiting, StateSucceeded) {
		panic("Invalid operation Wait() on CommandChain. It's not running.")
	}
	c.cleanUp()
}

func (c *CommandChain) moveToFailed() {
	atomic.StoreInt32(&(c.state), StateFailed)
	c.cleanUp()
}

func (c *CommandChain) cleanUp() {
	c.cleanupMu.Lock()
	defer c.cleanupMu.Unlock()

	for _, f := range c.tempFiles {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}
}

func (c *CommandChain) ensureHasCommand() {
	if len(c.Commands) == 0 {
		panic("No command is set yet.")
	}
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
	c.ensureHasCommand()
	return c.Commands[len(c.Commands)-1]
}

// Command adds a new command to a CommandChain.
func (c *CommandChain) Command(name string, args ...string) *CommandChain {
	c.ensureBuilding()
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
	c.ensureBuilding()

	c.Command(name, args...)

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
	c.ensureBuilding()
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
	c.ensureBuilding()
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
	c.ensureBuilding()
	c.defaultStdout = writer
	return c
}

// SetDefaultErr Set the default stderr to the subsequent commands.
func (c *CommandChain) SetDefaultErr(writer io.Writer) *CommandChain {
	c.ensureBuilding()
	c.defaultStderr = writer
	return c
}

// SetStdout sets a writer to the stdout of the last command.
func (c *CommandChain) SetStdout(writer io.Writer) *CommandChain {
	c.ensureBuilding()
	ensureNilAndSet(&c.lastCommand().Stdout, writer, "Stdout already set to command %s", c.getCommandDescription(-1))
	return c
}

// SetStderr sets a writer to the stderr of the last command.
func (c *CommandChain) SetStderr(writer io.Writer) *CommandChain {
	c.ensureBuilding()
	ensureNilAndSet(&c.lastCommand().Stderr, writer, "Stderr already set to command %s", c.getCommandDescription(-1))
	return c
}

// ErrToOut redirects stderr of the last command to stdout
func (c *CommandChain) ErrToOut() *CommandChain {
	c.ensureBuilding()
	if c.prevErrToOut {
		panic(fmt.Sprintf("ErrToOut has already been called on command %s", c.getCommandDescription(-1)))
	}
	c.prevErrToOut = true
	return c
}

// SetStdoutFile sets a file to the stdout of the last command.
func (c *CommandChain) SetStdoutFile(filename string) *CommandChain {
	c.ensureBuilding()
	f, err := openForWrite(filename)
	if err == nil {
		c.SetStdout(io.Writer(f))
	}
	c.setDeferredError(err)
	return c
}

// SetStderrFile sets a file to the stderr of the last command.
func (c *CommandChain) SetStderrFile(filename string) *CommandChain {
	c.ensureBuilding()
	f, err := openForWrite(filename)
	if err == nil {
		c.SetStderr(io.Writer(f))
	}
	c.setDeferredError(err)
	return c
}

// getStdoutPipe gets a pipe from stdout of the last command.
// This should only be used internally.
func (c *CommandChain) getStdoutPipe(reader **io.ReadCloser) *CommandChain {
	c.ensureBuilding()
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
	c.ensureBuilding()
	temp, err := os.CreateTemp(os.TempDir(), "stderr*.dat")
	if err != nil {
		c.setDeferredError(fmt.Errorf("CreateTemp() failed on %s: %w", c.getCommandDescription(-1), err))
		return c
	}
	c.tempFiles = append(c.tempFiles, temp)

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
	c.ensureBuilding()

	var rd *io.ReadCloser
	c.getStdoutPipe(&rd)
	c.setNextStdin(*rd)
	return c
}

func (c *CommandChain) validateBeforeRun() error {
	if c.deferredError != nil {
		return c.deferredError
	}
	c.ensureHasCommand()
	if c.nextStdin != nil {
		panic("Expecting next command to consume stdin")
	}
	return nil
}

// Run starts a CommandChain.
func (c *CommandChain) Run() (*ChainWaiter, error) {
	c.moveToRunning()

	err := c.validateBeforeRun()
	if err != nil {
		c.moveToFailed()
		return nil, err
	}
	c.fixUpLastCommand()

	for i, cmd := range c.Commands {
		err := cmd.Start()
		if err != nil {
			c.moveToFailed()
			return nil, fmt.Errorf("unable to execute command \"%s\" (command #%d): %s", cmd.Path, i+1, err.Error())
		}
	}
	return &ChainWaiter{Chain: c}, nil
}

// MustRun starts a CommandChain.
func (c *CommandChain) MustRun() *ChainWaiter {
	cw, err := c.Run()
	common.CheckPanic(err, "Unable to execute command(s)")
	return cw
}

// MustRunAndWait starts a CommandChain and wait().
func (c *CommandChain) MustRunAndWait() *ChainResult {
	return c.MustRun().MustWait()
}

// MustRunAndGetReader starts a CommandChain and get stdout of the last command as an io.Reader.
func (c *CommandChain) MustRunAndGetReader() (io.Reader, *ChainWaiter) {
	err := c.validateBeforeRun()
	common.CheckPanic(err, "Unable to execute command(s)")

	var rd *io.ReadCloser
	c.getStdoutPipe(&rd)

	cw := c.MustRun()

	return *rd, cw
}

const defaultBufSize = 4096

// MustRunAndGetBufferedReader starts a CommandChain and get stdout of the last command as an bufio.Reader.
func (c *CommandChain) MustRunAndGetBufferedReader() (*bufio.Reader, *ChainWaiter) {
	return c.mustRunAndGetBufferedReaderBufSize(defaultBufSize)
}

// mustRunAndGetBufferedReaderBufSize starts a CommandChain and get stdout of the last command as an bufio.Reader.
func (c *CommandChain) mustRunAndGetBufferedReaderBufSize(bufSize int) (*bufio.Reader, *ChainWaiter) {
	rd, cw := c.MustRunAndGetReader()
	return bufio.NewReaderSize(rd, bufSize), cw
}

// MustRunAndGetBytes starts a CommandChain and return stdout of the last command as []byte,
// and it also calls MustWait().
func (c *CommandChain) MustRunAndGetBytes() []byte {
	rd, waiter := c.MustRunAndGetReader()

	data, err := io.ReadAll(rd)
	common.CheckPanic(err, "Error while reading from commands")

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

func (c *CommandChain) MustRunAndGetStringsIter() func() *string {
	return utils.Iter(c.MustRunAndGetStrings())
}

// MustRunAndStreamStrings starts a CommandChain, read stdout of the last command line by line
// and feed them to con.
func (c *CommandChain) MustRunAndStreamStrings(con StringConsumer) {
	c.mustRunAndStreamStringsBufSize(con, defaultBufSize)
}

// mustRunAndStreamStringsBufSize starts a CommandChain, read stdout of the last command line by line
// and feed them to con, using bufSize for buffered reading.
func (c *CommandChain) mustRunAndStreamStringsBufSize(con StringConsumer, bufSize int) {
	rd, cw := c.mustRunAndGetBufferedReaderBufSize(bufSize)
	defer cw.MustWait()
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Sprintf("ReadString failed. %s", err))
		}
		line = textio.StringChomp(line)
		con(line)
	}
}

// MustRunAndStreamBytes starts a CommandChain, read stdout of the last command line by line
// and feed them to con.
func (c *CommandChain) MustRunAndStreamBytes(con BytesConsumer) {
	c.mustRunAndStreamBytesBufSize(con, defaultBufSize)
}

// mustRunAndStreamBytesBufSize starts a CommandChain, read stdout of the last command line by line
// and feed them to con, using bufSize for buffered reading.
func (c *CommandChain) mustRunAndStreamBytesBufSize(con BytesConsumer, bufSize int) {
	rd, cw := c.mustRunAndGetBufferedReaderBufSize(bufSize)
	defer cw.MustWait()
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Sprintf("ReadString failed. %s", err))
		}
		line = textio.Chomp(line)
		con(line)
	}
}

func (c *CommandChain) mustRunAndStreamBytesIterInner() (it func() (*[]byte, bool), cl func()) {
	rd, cw := c.mustRunAndGetBufferedReaderBufSize(defaultBufSize)

	it = func() (e *[]byte, ok bool) {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			return nil, false
		}
		if err != nil {
			panic(fmt.Sprintf("ReadString failed. %s", err))
		}
		line = textio.Chomp(line)
		return &line, true
	}
	cl = func() {
		cw.MustWait()
	}
	return
}

func (c *CommandChain) MustRunAndStreamBytesIter() *utils.Iterator[[]byte] {
	it, cl := c.mustRunAndStreamBytesIterInner()

	return utils.NewIterable(it, cl)
}

func (c *CommandChain) MustRunAndStreamStringsIter() *utils.Iterator[string] {
	it, cl := c.mustRunAndStreamBytesIterInner()

	return utils.NewIterable(
		func() (*string, bool) {
			data, ok := it()
			if ok {
				s := string(*data)
				return &s, ok
			} else {
				return nil, false
			}
		},
		cl,
	)
}

// Wait wait() on all commands in a CommandChain.
func (cw *ChainWaiter) Wait() (*ChainResult, error) {
	cw.Chain.moveToWaiting()

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
		cw.Chain.moveToFailed()
		return nil, firstError
	}
	cw.Chain.moveToSucceeded()

	return &ChainResult{Chain: cw.Chain}, nil
}

// MustWait wait() on all commands in a CommandChain.
func (cw *ChainWaiter) MustWait() *ChainResult {
	cr, err := cw.Wait()
	common.CheckPanice(err)
	return cr
}
