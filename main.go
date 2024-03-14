package main

import (
	"os"

	"github.com/gabe565/changelog-generator/cmd"
)

var (
	version = "beta"
	commit  = ""
)

func main() {
	if err := cmd.New(version, commit).Execute(); err != nil {
		os.Exit(1)
	}
}
