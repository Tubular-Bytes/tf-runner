package main

import (
	"log/slog"

	"github.com/Tubular-Bytes/tf-runner/pkg/version"
)

var CLI struct {
	Run     RunCmd     `cmd:"" help:"Run Terraform from a remote repository"`
	Version VersionCmd `cmd:"" help:"Show version"`
}

type RunCmd struct {
	RepoUrl      string `arg:"repo"       env:"TERRENCE_REPO"                     help:"Repository URL"`          //nolint:lll
	Apply        bool   `arg:"apply"      env:"TERRENCE_APPLY"                    help:"Apply the changes"`       //nolint:lll
	Endpoint     string `arg:"endpoint"   env:"TERRENCE_ENDPOINT"                 help:"Object store endpoint"`   //nolint:lll
	AccessKey    string `arg:"access-key" env:"TERRENCE_ACCESS_KEY"               help:"Object store access key"` //nolint:lll
	SecretKey    string `arg:"secret-key" env:"TERRENCE_SECRET_KEY"               help:"Object store secret key"`
	OverrideFile string `arg:"override-file" env:"TERRENCE_OVERRIDE_FILE" help:"Filename to write backend overrides into"`                        //nolint:lll
	StateURL     string `arg:"state"      default:"https://127.0.0.1:8080/state"  env:"TERRENCE_STATE_URL"       help:"URL for Terraform state"`  //nolint:lll
	LockURL      string `arg:"lock"       default:"https://127.0.0.1:8080/lock"   env:"TERRENCE_LOCK_URL"        help:"URL for Terraform lock"`   //nolint:lll
	UnlockURL    string `arg:"unlock"     default:"https://127.0.0.1:8080/unlock" env:"TERRENCE_UNLOCK_URL"      help:"URL for Terraform unlock"` //nolint:lll
}

type VersionCmd struct{}

func (v *VersionCmd) Run() error {
	initLogger()
	slog.Info("version information",
		"commit_hash", version.CommitHash,
		"version", version.Version,
		"build_date", version.BuildTime,
	)

	return nil
}
