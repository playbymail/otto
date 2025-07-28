// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements the `otto` command line application.
package main

import (
	"fmt"
	"github.com/maloquacious/semver"
	cmdCopy "github.com/playbymail/otto/cli/copy"
	cmdInfo "github.com/playbymail/otto/cli/info"
	cmdVersion "github.com/playbymail/otto/cli/version"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	version = semver.Version{Minor: 1}
)

func main() {
	// versions hack
	for _, arg := range os.Args {
		if arg == "version" || arg == "-version" || arg == "--version" {
			fmt.Printf("%s\n", version.String())
			os.Exit(0)
		}
	}

	cmdRoot := &cobra.Command{
		Use:   "otto",
		Short: "otto command line utility",
		Long:  `Otto is a tool for creating TribeNet maps.`,
	}

	cmdRoot.AddCommand(cmdCopy.Command)
	cmdRoot.AddCommand(cmdInfo.Command)
	cmdRoot.AddCommand(cmdVersion.Command)

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
