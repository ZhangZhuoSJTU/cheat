package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cheat/cheat/internal/config"
	"github.com/cheat/cheat/internal/sheet"
	"github.com/cheat/cheat/internal/sheets"
)

// cmdRemove removes (deletes) a cheatsheet.
func cmdRemove(cmd *cobra.Command, _ []string, conf config.Config) {

	cheatsheet, _ := cmd.Flags().GetString("rm")

	// validate the cheatsheet name
	if err := sheet.Validate(cheatsheet); err != nil {
		fmt.Fprintf(os.Stderr, "invalid cheatsheet name: %v\n", err)
		os.Exit(1)
	}

	// load the cheatsheets, filtered by tag if --tag was passed
	cheatsheets := loadSheets(cmd, conf)

	// consolidate the cheatsheets found on all paths into a single map of
	// `title` => `sheet` (ie, allow more local cheatsheets to override less
	// local cheatsheets)
	consolidated := sheets.Consolidate(cheatsheets)

	// fail early if the requested cheatsheet does not exist
	sheet, ok := consolidated[cheatsheet]
	if !ok {
		fmt.Fprintf(os.Stderr, "No cheatsheet found for '%s'.\n", cheatsheet)
		os.Exit(2)
	}

	// fail early if the sheet is read-only
	if sheet.ReadOnly {
		fmt.Fprintf(os.Stderr, "cheatsheet '%s' is read-only.\n", cheatsheet)
		os.Exit(1)
	}

	// otherwise, attempt to delete the sheet
	if err := os.Remove(sheet.Path); err != nil {
		fmt.Fprintf(os.Stderr, "failed to delete sheet: %s, %v\n", sheet.Title, err)
		os.Exit(1)
	}
}
