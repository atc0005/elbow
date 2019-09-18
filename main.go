package main

import (
	"fmt"
	"log"
	"os"

	"github.com/integrii/flaggy"
)

func main() {

	// DEBUG
	// TODO: Enable this once leveled logging has been implemented.
	//defaultConfig := NewConfig()
	//fmt.Printf("Default configuration:\t%+v\n", defaultConfig)

	appName := "Elbow"
	appDesc := "Prune content matching specific patterns, either in a single directory or recursively through a directory tree."

	config := NewConfig().SetupFlags(appName, appDesc)

	// DEBUG
	// TODO: Enable this once leveled logging has been implemented.
	//fmt.Printf("Our configuration:\t%+v\n", config)

	// DEBUG
	log.Println("Confirm that requested path actually exists")
	if !pathExists(config.StartPath) {
		flaggy.ShowHelpAndExit(fmt.Sprintf("Error processing requested path: %q", config.StartPath))
	}

	// INFO
	log.Println("Processing path:", config.StartPath)

	matches, err := processPath(config)

	// TODO
	// How to handle errors from gathering removal candidates?
	// Add optional flag to allow ignoring errors, fail immediately otherwise?
	if err != nil {
		log.Println("error:", err)
	}

	// NOTE: If this sort order changes, make sure to update the later logic
	// which retains the top or bottom X items (specific flag to preserve X
	// number of files while pruning the others)
	matches.sortByModTimeAsc()

	// DEBUG
	log.Printf("Length of matches slice: %d\n", len(matches))

	// DEBUG
	log.Println("Early exit if no matching files were found.")
	if len(matches) <= 0 {

		// INFO
		fmt.Printf("No matches found in path %q for files with substring pattern of %q and with extensions %v\n",
			config.StartPath, config.FilePattern, config.FileExtensions)

		// TODO: Not finding something is a valid outcome, so "normal" exit
		// code?
		os.Exit(0)
	}

	var filesToPrune FileMatches

	// DEBUG
	log.Printf("%d total items in matches", len(matches))
	log.Printf("%d items to keep per config.FilesToKeep", config.FilesToKeep)

	if config.KeepOldest {
		// DEBUG
		log.Println("Keeping older files")
		log.Println("start at specified number to keep, go until end of slice")
		filesToPrune = matches[config.FilesToKeep:]
	} else {
		// DEBUG
		log.Println("Keeping newer files")
		log.Println("start at beginning, go until specified number to keep")
		filesToPrune = matches[:(len(matches) - config.FilesToKeep)]
	}

	// DEBUG, INFO?
	log.Printf("%d items to prune", len(filesToPrune))

	log.Println("Prune specified files, do NOT ignore errors")
	// TODO: Add support for ignoring errors (though I cannot immediately
	// think of a good reason to do so)
	removalResults, err := cleanPath(filesToPrune, false, config)

	// Show what we WERE able to successfully remove
	// TODO: Refactor this into a function to handle displaying results?
	log.Printf("%d files successfully removed\n", len(removalResults.SuccessfulRemovals))
	log.Println("----------------------------")
	for _, file := range removalResults.SuccessfulRemovals {
		log.Println("*", file.Name())
	}

	log.Printf("%d files failed to remove\n", len(removalResults.FailedRemovals))
	log.Println("----------------------------")
	for _, file := range removalResults.FailedRemovals {
		log.Println("*", file.Name())
	}

	// Determine if we need to display error, exit with unsuccessful error code
	if err != nil {
		log.Fatalf("Errors encountered while processing %s: %s", config.StartPath, err)
	}

	log.Printf("%s successfully completed.", appName)

}
