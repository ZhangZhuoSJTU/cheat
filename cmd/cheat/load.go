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
		cheatsheets = sheets.Filter(cheatsheets, parseTags(tagVal))
	}

	return cheatsheets
}

// parseTags splits a comma-separated --tag value into individual tags. Each
// tag is trimmed of surrounding whitespace, and an optional leading '#' is
// stripped so that "#networking" and "networking" are treated identically.
func parseTags(tagVal string) []string {
	parts := strings.Split(tagVal, ",")
	tags := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = stripTagPrefix(p)
		tags = append(tags, p)
	}
	return tags
}

// stripTagPrefix removes a leading '#' from a tag. Tags prefixed with "##"
// (a convention for sub-category tags) have both hashes stripped.
func stripTagPrefix(tag string) string {
	if len(tag) > 0 && tag[0] == '#' {
		tag = tag[1:]
		if tag[0] == '#' {
			tag = tag[1:]
		}
	}
	return tag
}

// searchPattern builds a regex pattern for the search and list commands. When
// useRegex is false the phrase is wrapped in "(?i)" for case-insensitive
// matching; when true the raw phrase is returned unaltered so the caller can
// supply their own flags.
func searchPattern(phrase string, useRegex bool) string {
	if useRegex {
		return phrase
	}
	return "(?i)" + phrase
