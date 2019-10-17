// Copyright 2019 Adam Chalkley
//
// https://github.com/atc0005/elbow
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	// Use `log` if we are going to override the default `log`, otherwise
	// import without an "override" if we want to use the `logrus` name.
	// https://godoc.org/github.com/sirupsen/logrus
	"github.com/sirupsen/logrus"
)

// setup a shared logger object for use between various `main` package-level
// functions
var log = logrus.New()

// VERSION represents the version of this application. Set externally during
// build by our Makefile
var version = "x.y.z"

func main() {

	log.Debug("Constructing config object")

	// If this fails, the application will immediately exit.
	config := NewConfig()
	defaultConfig := NewConfig()

	log.Debug("Config object created")

	// Validate configuration
	if ok, err := config.Validate(); !ok {
		// NOTE: We're not using `log` here as the user-specified
		// configuration could be too botched to use reliably.

		// Provide user with error and valid usage details
		fmt.Printf("\nERROR: configuration validation failed\n%s\n\n", err)
		config.FlagParser.WriteUsage(os.Stdout)
		os.Exit(1)

		// failMessage := fmt.Sprint("configuration validation failed: ", err)
		// config.FlagParser.Fail(failMessage)

	}

	// Apply our custom logging settings on top of the existing global `log`
	// object (which uses default settings)
	setLoggerConfig(config, log)

	log.WithFields(logrus.Fields{
		"defaultConfig": defaultConfig,
	}).Debug("Default configuration")

	log.WithFields(logrus.Fields{
		"config": config,
	}).Debug("Our configuration")

	// https://www.joeshaw.org/dont-defer-close-on-writable-files/
	if config.LogFileHandle != nil {
		log.Debug("Deferring closure of log file")
		defer config.LogFileHandle.Close()
	}

	log.Debug("Confirm that requested path actually exists")
	if !pathExists(config.StartPath) {
		fmt.Printf("Error processing requested path: %q", config.StartPath)
		os.Exit(1)
	}

	log.WithFields(logrus.Fields{
		"path":         config.StartPath,
		"file_pattern": config.FilePattern,
		"extensions":   config.FileExtensions,
		"file_age":     config.FileAge,
	}).Info("Starting evaluation of path")

	matches, err := processPath(config)

	// TODO
	// How to handle errors from gathering removal candidates?
	// Add optional flag to allow ignoring errors, fail immediately otherwise?
	if err != nil {
		log.WithFields(logrus.Fields{
			"ignore_errors": config.IgnoreErrors,
		}).Error("error:", err)

		if !config.IgnoreErrors {
			log.Error("Error encountered, exiting")
			os.Exit(1)
		}
		log.Warn("Error encountered, but continuing as requested.")
	}

	// NOTE: If this sort order changes, make sure to update the later logic
	// which retains the top or bottom X items (specific flag to preserve X
	// number of files while pruning the others)
	matches.sortByModTimeAsc()

	log.Debugf("Length of matches slice: %d", len(matches))

	log.Debugf("Early exit if no matching files were found.")
	if len(matches) <= 0 {

		log.WithFields(logrus.Fields{
			"path":         config.StartPath,
			"file_pattern": config.FilePattern,
			"extensions":   config.FileExtensions,
			"file_age":     config.FileAge,
		}).Info("No matches found")

		// TODO: Not finding something is a valid outcome, so "normal" exit
		// code?
		os.Exit(0)
	}

	var filesToPrune FileMatches

	log.WithFields(logrus.Fields{
		"path":         config.StartPath,
		"file_pattern": config.FilePattern,
		"extensions":   config.FileExtensions,
		"file_age":     config.FileAge,
	}).Infof("%d files eligible for removal", len(matches))

	log.WithFields(logrus.Fields{
		"keep_oldest": config.KeepOldest,
	}).Infof("%d files to keep as requested", config.NumFilesToKeep)

	if config.KeepOldest {
		// TODO: Is debug output still useful?
		log.Debug("Keeping older files")
		log.Debug("start at specified number to keep, go until end of slice")
		filesToPrune = matches[config.NumFilesToKeep:]
	} else {
		log.Debug("Keeping newer files")
		log.Debug("start at beginning, go until specified number to keep")
		filesToPrune = matches[:(len(matches) - config.NumFilesToKeep)]
	}

	if len(filesToPrune) == 0 {
		log.Info("Nothing to prune, exiting")
		os.Exit(0)
	}

	log.WithFields(logrus.Fields{
		"files_to_prune": len(filesToPrune),
	}).Debug("Calling cleanPath")
	log.Infof("Ignoring file removal errors: %t", config.IgnoreErrors)
	removalResults, err := cleanPath(filesToPrune, config)

	// Show what we WERE able to successfully remove
	// TODO: Refactor this into a function to handle displaying results?
	log.Infof("%d files successfully removed", len(removalResults.SuccessfulRemovals))
	for _, file := range removalResults.SuccessfulRemovals {
		log.WithFields(logrus.Fields{
			"failed_removal": false,
		}).Info(file.Name())
	}

	log.Infof("%d files failed to remove", len(removalResults.FailedRemovals))
	for _, file := range removalResults.FailedRemovals {
		log.WithFields(logrus.Fields{
			"failed_removal": true,
		}).Info(file.Name())
	}

	if err != nil {
		log.Errorf("Errors encountered while processing %s: %s", config.StartPath, err)
		os.Exit(1)
	}

	log.Infof("%s successfully completed.", config.AppName)

}
