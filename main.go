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

// version represents the version of this application. Set externally during
// build by our Makefile
var version = "x.y.z"

func main() {

	// Checked at the end to determine if any issues were encountered during
	// app run. There is likely a much better way to handle this
	problemsEncountered := false

	appName := "Elbow"
	appDescription := "prunes content matching specific patterns, either in a single directory or recursively through a directory tree."
	appURL := "https://github.com/atc0005/elbow"

	log.Debug("Constructing config object")

	// If this fails, the application will immediately exit.
	config := NewConfig(appName, appDescription, appURL, version)
	defaultConfig := NewConfig(appName, appDescription, appURL, version)

	log.Debug("Config object created")

	// Validate configuration
	if ok, err := config.Validate(); !ok {
		// NOTE: We're not using `log` here as the user-specified
		// configuration could be too botched to use reliably.

		// Provide user with error and valid usage details
		fmt.Printf("\nERROR: configuration validation failed\n%s\n\n", err)
		config.FlagParser.WriteUsage(os.Stdout)
		problemsEncountered = true
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

	log.WithFields(logrus.Fields{
		"paths":        config.Paths,
		"file_pattern": config.FilePattern,
		"extensions":   config.FileExtensions,
		"file_age":     config.FileAge,
	}).Info("Starting evaluation of paths list")

	var pass int
	var totalPaths int = len(config.Paths)
	for _, path := range config.Paths {

		pass++

		log.WithFields(logrus.Fields{
			"total_paths": totalPaths,
			"iteration":   pass,
		}).Infof("Beginning processing of path %q (%d of %d)",
			path, pass, totalPaths)

		log.Debug("Confirm that requested path actually exists")
		if !pathExists(path) {

			problemsEncountered = true

			log.WithFields(logrus.Fields{
				"ignore_errors": config.IgnoreErrors,
			}).Errorf("Requested path not found: %q", path)

			if config.IgnoreErrors {
				log.WithFields(logrus.Fields{
					"ignore_errors": config.IgnoreErrors,
				}).Warn("Error encountered, but continuing as requested.")
				continue
			} else {
				log.WithFields(logrus.Fields{
					"ignore_errors": config.IgnoreErrors,
				}).Warn("Error encountered and option to ignore errors not set. Exiting")

				os.Exit(1)
			}
		}

		matches, err := processPath(config, path)
		if err != nil {

			problemsEncountered = true

			log.WithFields(logrus.Fields{
				"ignore_errors": config.IgnoreErrors,
			}).Error("error:", err)

			if !config.IgnoreErrors {
				log.WithFields(logrus.Fields{
					"ignore_errors": config.IgnoreErrors,
				}).Warn("Error encountered and option to ignore errors not set. Exiting")
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
				"path":         path,
				"file_pattern": config.FilePattern,
				"extensions":   config.FileExtensions,
				"file_age":     config.FileAge,
			}).Info("No matches found")

			log.Debugf("Ending processing of path %d", pass)
			if pass < totalPaths {
				log.Debugf("Continuing to next available path")
			}
			continue

		}

		var filesToPrune FileMatches

		log.WithFields(logrus.Fields{
			"path":         path,
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
			log.Info("Nothing to prune")
			log.Debugf("Ending processing of path %d", pass)
			if pass < totalPaths {
				log.Debugf("Continuing to next available path")
			}
			continue
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

			log.Warnf("Error encountered while processing %s: %s", path, err)

			if !config.IgnoreErrors {
				log.WithFields(logrus.Fields{
					"ignore_errors": config.IgnoreErrors,
				}).Warn("Error encountered and option to ignore errors not set. Exiting")
			}
			log.Warn("Error encountered, but continuing as requested.")
			continue
		}

		log.WithFields(logrus.Fields{
			"total_paths": totalPaths,
			"iteration":   pass,
		}).Infof("Ending processing of %q (%d of %d)", path, pass, totalPaths)

	}

	if problemsEncountered {
		log.Warnf("%s completed, but issues were encountered.", config.AppName)
	} else {
		log.Infof("%s successfully completed.", config.AppName)
	}

}
