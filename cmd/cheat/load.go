package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cheat/cheat/internal/config"
	"github.com/cheat/cheat/internal/sheet"
	"github.com/cheat/cheat/internal/sheets"
)

// loadSheets loads the cheatsheets on all cheatpaths, and filters them by tag
// if --tag was passed. It exits the program if the cheatsheets cannot be
// loaded.
func loadSheets(cmd *cobra.Command, conf config.Config) []map[string]sheet.Sheet {

	// load the cheatsheets
	cheatsheets, err := sheets.Load(conf.Cheatpaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list cheatsheets: %v\n", err)
		os.Exit(1)
	}

	// filter the cheatsheets by tag if --tag was passed
	if cmd.Flags().Changed("tag") {
		tagVal, _ := cmd.Flags().GetString("tag")
		cheatsheets = sheets.Filter(
			cheatsheets,
			strings.Split(tagVal, ","),
		)
	}

	return cheatsheets
}
