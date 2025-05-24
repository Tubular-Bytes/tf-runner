package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Tubular-Bytes/tf-runner/pkg/logexporter"
	"github.com/Tubular-Bytes/tf-runner/pkg/tofu"
	"github.com/go-git/go-git/v5"
)

func (r *RunCmd) Run(ctx *Context) error {
	logWriter := logexporter.NewLogWriter()

	store, err := logexporter.New(logexporter.ExporterConfig{
		Endpoint:        r.Endpoint,
		AccessKeyID:     r.AccessKey,
		SecretAccessKey: r.SecretKey,
		UseSSL:          true,
		Repository:      r.repoName(),
	})
	if err != nil {
		slog.Error("failed to create log exporter", "error", err)

		return err
	}

	defer func() {
		store.Flush(logWriter.Data(), true)
		cleanUp()
	}()

	if err := mustHaveTofu(); err != nil {
		slog.Error("tofu not found", "error", err)

		return err
	}

	if err := os.MkdirAll("workspace", os.ModePerm); err != nil {
		slog.Error("failed to create workspace", "error", err)

		return err
	}

	slog.Info("cloning repo", "repo", r.RepoUrl, "dir", r.workingDir())

	if _, err := git.PlainClone(r.workingDir(), false, &git.CloneOptions{
		URL: r.RepoUrl,
	}); err != nil {
		slog.Error("failed to clone repo", "error", err)

		return err
	}

	if err := r.writeOverride(); err != nil {
		slog.Error("failed to write override file", "error", err)

		return err
	}

	if err := r.pipeline(logWriter); err != nil {
		slog.Error("pipeline failed", "error", err)

		return err
	}

	return nil
}

func initLogger() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}

func cleanUp() {
	slog.Debug("cleaning up")

	if err := os.RemoveAll("workspace"); err != nil {
		slog.Warn("failed to remove workspace", "error", err)
	}
}

func mustHaveTofu() error {
	if _, err := exec.LookPath("tofu"); err != nil {
		slog.Error("tofu not found", "error", err)

		return fmt.Errorf("tofu not found: %w", err)
	}

	return nil
}

func (r *RunCmd) writeOverride() error {
	override, err := tofu.Render(r.StateURL, r.LockURL, r.UnlockURL)
	if err != nil {
		return err
	}

	overrideFile := filepath.Join(r.workingDir(), "main_override.tf")
	slog.Info("writing override file",
		"file", overrideFile,
		"state_url", r.StateURL,
		"lock_url", r.LockURL,
		"unlock_url", r.UnlockURL,
	)

	if err := os.WriteFile(overrideFile, override, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (r *RunCmd) repoName() string {
	repoParts := strings.Split(r.RepoUrl, "/")

	return strings.TrimSuffix(repoParts[len(repoParts)-1:][0], ".git")
}

func (r *RunCmd) workingDir() string {
	return filepath.Join("workspace", r.repoName())
}

func (r *RunCmd) pipeline(logWriter *logexporter.LogWriter) error {
	if err := tofu.Init(r.workingDir(), logWriter); err != nil {
		slog.Error("init failed", "error", err)

		return err
	}

	if err := tofu.Plan(r.workingDir(), logWriter); err != nil {
		slog.Error("plan failed", "error", err)

		return err
	}

	if r.Apply {
		if err := tofu.Apply(r.workingDir(), logWriter); err != nil {
			slog.Error("plan failed", "error", err)

			return err
		}
	} else {
		slog.Info("apply is disabled, skipping", "dir", r.workingDir())
	}

	return nil
}
