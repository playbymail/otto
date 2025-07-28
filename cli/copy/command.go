// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cli implements the `copy` command.
package cli

import (
	"github.com/spf13/cobra"
	"log"
)

var Command = &cobra.Command{
	Use:   "copy",
	Short: "Copy map data to a new file",
	Long:  `Copy map data to a new file, keeping only information used by Otto.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("not implemented")
		return nil
	},
}
