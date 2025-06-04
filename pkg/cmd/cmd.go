package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

var (
	ErrNilCommand = fmt.Errorf("command is nil")
)

type Command struct {
	command *exec.Cmd
	dir     string
	stdout  io.Writer
	stderr  io.Writer
}

func New(cmd string, opt ...Option) *Command {
	c := &Command{
		command: exec.Command(cmd),
		stdout:  bytes.NewBufferString(""),
		stderr:  bytes.NewBufferString(""),
	}

	for _, o := range opt {
		o(c)
	}

	return c
}

func (c *Command) Dir() string {
	if c.command == nil {
		return ""
	}

	return c.command.Dir
}

func (c *Command) Stdout() io.Writer {
	if c.command == nil {
		return nil
	}

	return c.stdout
}

func (c *Command) Stderr() io.Writer {
	if c.command == nil {
		return nil
	}

	return c.stderr
}

func (c *Command) String() string {
	if c.command == nil {
		return ""
	}

	return c.command.String()
}

func (c *Command) SetDebug(debug bool) {
	if c.command == nil {
		return
	}

	if debug {
		c.command.Env = append(c.command.Env, "TF_LOG=DEBUG")
	}
}

func (c *Command) SetArgs(args ...string) {
	if c.command != nil {
		c.command.Args = append(c.command.Args, args...)
	}
}

func (c *Command) setStdout(out io.Writer) {
	c.stdout = out
	if c.command != nil {
		c.command.Stdout = out
	}
}

func (c *Command) setStderr(err io.Writer) {
	c.stderr = err
	if c.command != nil {
		c.command.Stderr = err
	}
}

func (c *Command) setDir(dir string) {
	c.dir = dir
	if c.command != nil {
		c.command.Dir = dir
	}
}

func (c *Command) Run() error {
	if c.command == nil {
		return ErrNilCommand
	}

	return c.command.Run()
}
