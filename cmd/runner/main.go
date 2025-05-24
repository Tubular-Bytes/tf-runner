package main

import (
	"github.com/alecthomas/kong"
)

/*
	EnvRepoUrl          = "TERRENCE_REPO"
	EnvApply            = "TERRENCE_APPLY"
	EnvStateURL         = "TERRENCE_STATE_URL"
	EnvLockURL          = "TERRENCE_LOCK_URL"
	EnvUnlockURL        = "TERRENCE_UNLOCK_URL"
	EnvStorageEndpoint  = "TERRENCE_ENDPOINT"
	EnvStorageAccessKey = "TERRENCE_ACCESS_KEY"
	EnvStorageSecretKey = "TERRENCE_SECRET_KEY"
*/

func main() {
	initLogger()

	ctx := kong.Parse(&CLI,
		kong.Description("Terrence Runner"),
		kong.UsageOnError(),
		kong.Name("terrence"),
		kong.Vars{"version": "0.1.0"},
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
