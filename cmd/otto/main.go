// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements the `otto` command line application.
package main

import (
	"fmt"
	"github.com/playbymail/otto"
	cmdCopy "github.com/playbymail/otto/cmd/otto/copy"
	cmdInfo "github.com/playbymail/otto/cmd/otto/info"
	cmdVersion "github.com/playbymail/otto/cmd/otto/version"
	"github.com/playbymail/otto/config"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func main() {
	// versions hack
	for _, arg := range os.Args {
		if arg == "version" || arg == "-version" || arg == "--version" {
			fmt.Printf("%s\n", otto.Version().String())
			os.Exit(0)
		}
	}

	cfg := &config.Config_t{}

	cmdRoot := &cobra.Command{
		Use:   "otto",
		Short: "otto command line utility",
		Long:  `Otto is a tool for creating TribeNet maps.`,
	}

	cmdRoot.AddCommand(cmdCopy.Command)
	if err := cmdCopy.RegisterArgs(cfg); err != nil {
		log.Fatal(err)
	}
	cmdRoot.AddCommand(cmdInfo.Command)
	cmdRoot.AddCommand(cmdVersion.Command)

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
