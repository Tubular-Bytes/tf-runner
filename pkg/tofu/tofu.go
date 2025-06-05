package tofu

import (
	"io"
	"log/slog"

	"github.com/Tubular-Bytes/tf-runner/pkg/cmd"
)

const (
	planOutput = "plan.out"
)

func tofu(workingDir string, stdout io.Writer, stderr io.Writer) *cmd.Command {
	return cmd.New("tofu",
		cmd.WithDir(workingDir),
		cmd.WithStdout(stdout),
		cmd.WithStderr(stderr),
		cmd.WithDebug(Debug()),
	)
}

func Init(workingDir string, output io.Writer) error {
	slog.Info("initializing workspace", "dir", workingDir)
	init := tofu(workingDir, output, output)
	init.SetArgs("init", "-no-color", "-input=false")

	return init.Run()
}

func Plan(workingDir string, output io.Writer) error {
	slog.Info("creating plan", "dir", workingDir, "output", planOutput)
	plan := tofu(workingDir, output, output)
	plan.SetArgs("plan", "-no-color", "-input=false", "-out="+planOutput)

	return plan.Run()
}

func Apply(workingDir string, output io.Writer) error {
	slog.Info("applying plan", "dir", workingDir)
	apply := tofu(workingDir, output, output)
	apply.SetArgs("apply", "-no-color", "-input=false", planOutput)

	return apply.Run()
}
