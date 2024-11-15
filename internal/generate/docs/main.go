package main

import (
	"fmt"
	"log"
	"os"

	"gabe565.com/changelog-generator/cmd"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra/doc"
)

func main() {
	var err error
	output := "./docs"

	err = os.RemoveAll(output)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to remove existing dia: %w", err))
	}

	err = os.MkdirAll(output, 0o755)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to mkdir: %w", err))
	}

	rootCmd := cmd.New(cobrax.WithVersion("beta"))

	err = doc.GenMarkdownTree(rootCmd, output)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to generate markdown: %w", err))
	}
}
