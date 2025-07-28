// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cli implements the `version` command.
package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "version",
	Short: "Show application version",
	Long:  `Show application version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// apologies for this, but the version command is implemented in the
		// project's `~/main.go` file. this is here just to force Cobra to
		// show the command line options and help.
		return fmt.Errorf("not implemented")
	},
}
