package cmd_test

import (
	"bytes"
	"testing"

	"github.com/Tubular-Bytes/tf-runner/pkg/cmd"
	"github.com/stretchr/testify/require"
)

func TestSuite(t *testing.T) {
	t.Run("New", testNew)
	t.Run("NilCommand", testNilCommand)
	t.Run("Run", testRun)
	t.Run("RunNil", testRunNil)
	t.Run("RunError", testRunError)
}

func testNew(t *testing.T) {
	ob := bytes.NewBufferString("")
	eb := bytes.NewBufferString("")
	c := cmd.New(
		"echo",
		cmd.WithDir("./"),
		cmd.WithStdout(ob),
		cmd.WithStderr(eb),
		cmd.WithArgs("hello", "world"),
	)

	require.Equal(t, "/bin/echo hello world", c.String())
	require.Equal(t, "./", c.Dir())
	require.Equal(t, ob, c.Stdout().(*bytes.Buffer))
	require.Equal(t, eb, c.Stderr().(*bytes.Buffer))
}

func testNilCommand(t *testing.T) {
	c := &cmd.Command{}
	require.Equal(t, "", c.String())
	require.Nil(t, c.Stdout())
	require.Nil(t, c.Stderr())
	require.Equal(t, "", c.Dir())
}

func testRun(t *testing.T) {
	ob := bytes.NewBufferString("")
	eb := bytes.NewBufferString("")
	c := cmd.New(
		"echo",
		cmd.WithDir("./"),
		cmd.WithStdout(ob),
		cmd.WithStderr(eb),
		cmd.WithArgs("hello", "world"),
	)

	err := c.Run()
	require.NoError(t, err)
	require.Equal(t, "hello world\n", c.Stdout().(*bytes.Buffer).String())
	require.Equal(t, "", c.Stderr().(*bytes.Buffer).String())
}

func testRunNil(t *testing.T) {
	c := &cmd.Command{}
	err := c.Run()
	require.Error(t, err)
	require.Equal(t, cmd.ErrNilCommand, err)
}

func testRunError(t *testing.T) {
	ob := bytes.NewBufferString("")
	eb := bytes.NewBufferString("")
	c := cmd.New(
		"cp",
		cmd.WithDir("./"),
		cmd.WithStdout(ob),
		cmd.WithStderr(eb),
		cmd.WithArgs("hello", "world"),
	)

	err := c.Run()
	require.Error(t, err)
	require.Equal(t, "exit status 1", err.Error())
	require.Equal(t, "", c.Stdout().(*bytes.Buffer).String())
	require.Equal(t, "cp: hello: No such file or directory\n", c.Stderr().(*bytes.Buffer).String())
}
