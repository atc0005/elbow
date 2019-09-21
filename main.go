package main

import (
	"fmt"
	"os"

	// Use `log` if we are going to override the default `log`, otherwise
	// import without an "override" if we want to use the `logrus` name.
	// https://godoc.org/github.com/sirupsen/logrus
	// "github.com/jessevdk/go-flags"

	"github.com/sirupsen/logrus"
)

func main() {

	log := newLogger()

	// TODO: Can this info be set using go-flags? An interface for this?
	appName := "Elbow"
	appDesc := "Prune content matching specific patterns, either in a single directory or recursively through a directory tree."

	// If this fails, the application will immediately exit.
	config := NewConfig().SetupFlags(appName, appDesc)

	// Log as JSON instead of the default ASCII formatter.
	// TODO: Use command-line option here
	log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// TODO: Accept command-line parameter to determine this level
	// TODO: Setup mapping between command-line options and valid logrus levels
	// so that they can be referenced here
	log.SetLevel(logrus.DebugLevel)

	defaultConfig := NewConfig()
	log.WithFields(logrus.Fields{
		"defaultConfig": defaultConfig,
	}).Debug("Default configuration")

	log.WithFields(logrus.Fields{
		"config": config,
	}).Debug("Our configuration")

	// Validate configuration
	// TODO: How much of this work does go-flags handle for us?
	if ok := config.Validate(); !ok {
		fmt.Println("configuration validation failed")
		os.Exit(1)
	}

	log.Debug("Confirm that requested path actually exists")
	if !pathExists(config.StartPath) {
		fmt.Printf("Error processing requested path: %q", config.StartPath)
		os.Exit(1)
	}

	log.Info("Processing path:", config.StartPath)

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

	log.Debugf("Length of matches slice: %d", len(matches))

	log.Debugf("Early exit if no matching files were found.")
	if len(matches) <= 0 {

		noMatchesMessage := fmt.Sprintf("No matches found in path %q for files with substring pattern of %q and with extensions %v",
			config.StartPath, config.FilePattern, config.FileExtensions)

		// TODO: How to handle duplicate output to the console from both
		// commands?
		log.Info(noMatchesMessage)
		fmt.Println(noMatchesMessage)

		// TODO: Not finding something is a valid outcome, so "normal" exit
		// code?
		os.Exit(0)
	}

	var filesToPrune FileMatches

	// DEBUG
	log.Debugf("%d total items in matches", len(matches))
	log.Debugf("%d items to keep per config.NumFilesToKeep", config.NumFilesToKeep)

	if config.KeepOldest {
		// DEBUG
		log.Debug("Keeping older files")
		log.Debug("start at specified number to keep, go until end of slice")
		filesToPrune = matches[config.NumFilesToKeep:]
	} else {
		// DEBUG
		log.Debug("Keeping newer files")
		log.Debug("start at beginning, go until specified number to keep")
		filesToPrune = matches[:(len(matches) - config.NumFilesToKeep)]
	}

	log.Info("%d items to prune", len(filesToPrune))

	log.Info("Prune specified files, do NOT ignore errors")
	// TODO: Add support for ignoring errors (though I cannot immediately
	// think of a good reason to do so)
	removalResults, err := cleanPath(filesToPrune, false, config)

	// Show what we WERE able to successfully remove
	// TODO: Refactor this into a function to handle displaying results?
	log.Infof("%d files successfully removed", len(removalResults.SuccessfulRemovals))
	for _, file := range removalResults.SuccessfulRemovals {
		log.Info(file.Name())
	}

	log.Info("%d files failed to remove", len(removalResults.FailedRemovals))
	for _, file := range removalResults.FailedRemovals {
		log.Println("*", file.Name())
	}

	if err != nil {
		log.Errorf("Errors encountered while processing %s: %s", config.StartPath, err)
		os.Exit(1)
	}

	log.Infof("%s successfully completed.", appName)

}
