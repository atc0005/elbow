package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/integrii/flaggy"
)

// TODO: What other option do I have here other than using globals?
// Use closures?
var matches FileMatches
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

	// TODO: Branch at this point based off of whether the recursive option
	// was chosen.
	if config.RecursiveSearch {
		// DEBUG
		log.Println("Recursive option is enabled")
		log.Printf("%v", config)
		os.Exit(0)
	}

	// If RecursiveSearch is not enabled, process just the provided StartPath
	// NOTE: The same cleanPath() function is used in either case, the
	// difference is in how the FileMatches slice is populated

	//os.Exit(0)

	//
	// TODO: Refactor filepath.Walk() call below; split into at least two
	// functions, one to do what is being done now (recursive work), another
	// to use `ioutil.ReadDir(path)` to gather matches from specific
	// directory.
	//

	// Walk walks the file tree rooted at root, calling crawlPath for each
	// file or directory in the tree, including root. All errors that arise
	// visiting files and directories are filtered by crawlPath. The files
	// are walked in lexical order, which makes the output deterministic but
	// means that for very large directories Walk can be inefficient. Walk
	// does not follow symbolic links.
	err := filepath.Walk(config.StartPath, crawlPath)
	if err != nil {
		log.Println("error:", err)
	}

	// TODO: Does this in-place attempt at sorting work because slices are
	// reference types to begin with?
	matches.sortByModTimeAsc()

	// DEBUG
	log.Printf("Length of matches slice: %d", len(matches))

	//pruneFilesStartPoint := 2
	if len(matches) <= 0 {

		// INFO
		fmt.Printf("No matches found in path %q for %v",
			config.StartPath, config.FilePattern)

		// TODO: Not finding something is a valid outcome, so "normal" exit
		// code?
		os.Exit(0)
	}

	// TODO: Is it safe to allow use of the default 0 value here?
	//
	// Do we keep the oldest or the newest files (limited to
	// config.FilesToKeep) ?
	var filesToPrune FileMatches
	if config.KeepOldest {
		filesToPrune = matches[config.FilesToKeep:]
	} else {
		filesToPrune = matches[:config.FilesToKeep]
	}

	// Prune specified files, do NOT ignore errors
	filesRemoved, err := cleanPath(filesToPrune, false)

	// Show what we WERE able to successfully remove
	fmt.Println("%d files successfully removed:", len(filesRemoved))
	for _, file := range filesRemoved {
		fmt.Println("*", file)
	}

	// Determine if we need to display error, exit with unsuccessful error code
	if err != nil {
		log.Fatalf("Errors encountered while processing %s: %s", config.StartPath, err)
	}

}
