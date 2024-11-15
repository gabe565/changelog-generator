package main

import (
	"os"

	"gabe565.com/changelog-generator/cmd"
)

//nolint:gochecknoglobals
var (
	version = "beta"
	commit  = ""
)

func main() {
	if err := cmd.New(version, commit).Execute(); err != nil {
		os.Exit(1)
	}
}
