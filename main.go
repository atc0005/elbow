package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/integrii/flaggy"
)

// TODO: What other option do I have here other than using globals?
// Use closures?
var config *Config

func main() {

	/*

		TODO: Collect these options (command-line / config file / env vars)

		1) [string] File pattern to match on
		2) [string] Starting path for processing
		3) [bool] Recursive search
		4) [int] Number of files to keep
		5) [bool] KeepYoungest (true, default)

	*/

	// create default configuration so that we can compare against it to
	// determine whether the user has provided flags
	//
	// TODO: Is this needed? We'll have to validate the flags either way?
	defaultConfig := NewConfig()

	// DEBUG
	fmt.Printf("Default configuration:\t%+v\n", defaultConfig)

	appName := "Elbow"
	appDesc := "Prune content matching specific patterns, either in a single directory or recursively through a directory tree."

	config = NewConfig().SetupFlags(appName, appDesc)
	fmt.Printf("Our configuration:\t%+v\n", config)

	// TODO: How can I reliably compare these?
	//  invalid operation: *defaultConfig != *config (struct containing []string cannot be compared)
	// if *defaultConfig != *config {
	// 	log.Println("User specified command-line options")
	// } else {
	// 	log.Println("User did not provide any command-line flags")
	// }

	// TODO: Is this even needed? Shouldn't I instead focus on whether the
	// values that are set (default or not) actually validate?
	if reflect.DeepEqual(*defaultConfig, *config) {
		// DEBUG
		log.Println("User did not provide command-line flags; current configuration matches default settings")

		// KEEP
		flaggy.ShowHelpAndExit("Required command-line options not provided.")
	} else {
		// DEBUG
		log.Println("User provided command-line flags, proceeding ...")
	}

	//fmt.Printf("%+v\n", *config)

	// Confirm that requested path actually exists
	if !pathExists(config.StartPath) {
		log.Fatalf("Error processing requested path: %q", config.StartPath)
	}

	// INFO
	log.Println("Processing path:", config.StartPath)

	matches, err := processPath(config)

	// TODO: How to handle errors from gathering removal candidates?
	if err != nil {
		log.Println("error:", err)
	}

	//os.Exit(0)

	// TODO: Does this in-place attempt at sorting work because slices are
	// reference types to begin with?
	matches.sortByModTimeAsc()

	// DEBUG
	log.Printf("Length of matches slice: %d\n", len(matches))

	// Early exit if no matching files were found.
	if len(matches) <= 0 {

		// INFO
		fmt.Printf("No matches found in path %q for files with substring pattern of %q and with extensions %v\n",
			config.StartPath, config.FilePattern, config.FileExtensions)

		// TODO: Not finding something is a valid outcome, so "normal" exit
		// code?
		os.Exit(0)
	}

	// TODO: Is it safe to allow use of the default 0 value here?
	//
	// Do we keep the oldest or the newest files (limited to
	// config.FilesToKeep) ?
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

	// Prune specified files, do NOT ignore errors
	filesRemoved, err := cleanPath(filesToPrune, false, config)

	// Show what we WERE able to successfully remove
	log.Printf("%d files successfully removed:\n", len(filesRemoved))
	log.Println("----------------------------")
	for _, file := range filesRemoved {
		log.Println("*", file.Name())
	}

	// Determine if we need to display error, exit with unsuccessful error code
	if err != nil {
		log.Fatalf("Errors encountered while processing %s: %s", config.StartPath, err)
	}

}
