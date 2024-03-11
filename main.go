package main

import (
	"os"

	"github.com/gabe565/changelog-generator/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
