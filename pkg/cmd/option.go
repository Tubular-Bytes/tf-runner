package cmd

import "io"

type Option func(*Command)

func WithDir(dir string) Option {
	return func(c *Command) {
		c.setDir(dir)
	}
}

func WithStdout(w io.Writer) Option {
	return func(c *Command) {
		c.setStdout(w)
	}
}

func WithStderr(w io.Writer) Option {
	return func(c *Command) {
		c.setStderr(w)
	}
}

func WithArgs(args ...string) Option {
	return func(c *Command) {
		if c.command != nil {
			c.SetArgs(args...)
		}
	}
}

func WithDebug(debug bool) Option {
	return func(c *Command) {
		if c.command != nil {
			c.SetDebug(debug)
		}
	}
}
