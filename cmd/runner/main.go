package main

import (
	"github.com/alecthomas/kong"
)

func main() {
	ctx := kong.Parse(&CLI,
		kong.Description("Terrence Runner"),
		kong.UsageOnError(),
		kong.Name("terrence"),
		kong.Vars{"version": "0.1.0"},
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
