// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cli implements the `copy` command.
package cli

import (
	"errors"
	"fmt"
	"github.com/maloquacious/wxx/xmlio"
	"github.com/playbymail/otto/config"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "copy",
	Short: "Copy map data to a new file",
	Long:  `Copy map data to a new file, keeping only information used by Otto.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// read the input
		from, err := cmd.Flags().GetString("from")
		if err != nil {
			return fmt.Errorf("could not read --from: %w", err)
		}
		w, err := xmlio.ReadFile(from)
		if err != nil {
			return errors.Join(fmt.Errorf("copy: xmlio.Read"), err)
		}
		fmt.Printf("\t%8d tiles high\n", w.Tiles.TilesHigh)
		fmt.Printf("\t%8d tiles wide\n", w.Tiles.TilesWide)
		fmt.Printf("\t%8d terrain tiles defined\n", len(w.TerrainMap.List))
		return fmt.Errorf("not implemented")
	},
}

func RegisterArgs(cfg *config.Config_t) error {
	Command.Flags().String("from", "", "name of map file to copy from")
	if err := Command.MarkFlagRequired("from"); err != nil {
		return errors.Join(fmt.Errorf("copy"), err)
	}
	Command.Flags().String("to", "", "name of map file to create")
	if err := Command.MarkFlagRequired("to"); err != nil {
		return errors.Join(fmt.Errorf("copy"), err)
	}
	return nil
}
