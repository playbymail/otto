// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cli implements the `info` command.
package cli

import (
	"github.com/spf13/cobra"
	"log"
)

var Command = &cobra.Command{
	Use:   "info",
	Short: "Show map information",
	Long:  `Info displays metadata from a map like  the Worldographer version, height, and width.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("not implemented")
		return nil
	},
}
