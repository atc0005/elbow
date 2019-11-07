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

	"github.com/atc0005/elbow/config"
	"github.com/atc0005/elbow/logging"
	"github.com/atc0005/elbow/matches"
	"github.com/atc0005/elbow/paths"
	"github.com/atc0005/elbow/units"

	// Use `log` if we are going to override the default `log`, otherwise
	// import without an "override" if we want to use the `logrus` name.
	// https://godoc.org/github.com/sirupsen/logrus
	"github.com/sirupsen/logrus"
)

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

	//log.Debug("Constructing config object")

	// If this fails, the application will immediately exit.
	appConfig := config.NewConfig(appName, appDescription, appURL, version)
	defaultConfig := config.NewConfig(appName, appDescription, appURL, version)

	// Validate configuration
	if ok, err := appConfig.Validate(); !ok {
		// NOTE: We're not using `log` here as the user-specified
		// configuration could be too botched to use reliably.

		// Provide user with error and valid usage details
		fmt.Printf("\nERROR: configuration validation failed\n%s\n\n", err)
		appConfig.FlagParser.WriteUsage(os.Stdout)
		problemsEncountered = true
		os.Exit(1)

		// failMessage := fmt.Sprint("configuration validation failed: ", err)
		// appConfig.FlagParser.Fail(failMessage)

	}

	logging.SetLoggerConfig(appConfig)

	log := appConfig.Logger

	log.Debug("Config object created")

	// Apply our custom logging settings

	fmt.Println("JSON unmarshal error here")
	fmt.Printf("defaultConfig struct: %+v\n", defaultConfig)
	log.WithFields(logrus.Fields{
		"defaultConfig": defaultConfig,
	}).Debug("Default configuration")

	fmt.Println("JSON unmarshal error here")
	fmt.Printf("appConfig struct: %+v\n", appConfig)
	log.WithFields(logrus.Fields{
		"config": appConfig,
	}).Debug("Our configuration")

	// FIXME: Remove this after finishing troubleshooting work
	fmt.Println("Where does this show?")

	// https://www.joeshaw.org/dont-defer-close-on-writable-files/
	if appConfig.LogFileHandle != nil {
		log.Debug("Deferring closure of log file")
		defer appConfig.LogFileHandle.Close()
	}

	log.WithFields(logrus.Fields{
		"paths":        appConfig.Paths,
		"file_pattern": appConfig.FilePattern,
		"extensions":   appConfig.FileExtensions,
		"file_age":     appConfig.FileAge,
	}).Info("Starting evaluation of paths list")

	// Used as a global counter/bucket for presentation/logging purposes
	var appResults paths.ProcessingResults

	var pass int
	var totalPaths int = len(appConfig.Paths)
	for _, path := range appConfig.Paths {

		pass++

		log.WithFields(logrus.Fields{
			"total_paths":   totalPaths,
			"iteration":     pass,
			"ignore_errors": appConfig.IgnoreErrors,
		}).Infof("Beginning processing of path %q (%d of %d)",
			path, pass, totalPaths)

		log.Debug("Confirm that requested path actually exists")
		if !paths.PathExists(path, appConfig) {

			problemsEncountered = true

			log.WithFields(logrus.Fields{
				"ignore_errors": appConfig.IgnoreErrors,
			}).Errorf("Requested path not found: %q", path)

			if appConfig.IgnoreErrors {
				log.WithFields(logrus.Fields{
					"ignore_errors": appConfig.IgnoreErrors,
				}).Warn("Error encountered, but continuing as requested.")
				continue
			} else {
				log.WithFields(logrus.Fields{
					"ignore_errors": appConfig.IgnoreErrors,
				}).Warn("Error encountered and option to ignore errors not set. Exiting")

				os.Exit(1)
			}
		}

		fileMatches, err := paths.ProcessPath(appConfig, path)
		if err != nil {

			problemsEncountered = true

			log.WithFields(logrus.Fields{
				"ignore_errors": appConfig.IgnoreErrors,
			}).Error("error:", err)

			if !appConfig.IgnoreErrors {
				log.WithFields(logrus.Fields{
					"ignore_errors": appConfig.IgnoreErrors,
				}).Warn("Error encountered and option to ignore errors not set. Exiting")
			}
			log.Warn("Error encountered, but continuing as requested.")
		}

		log.Debugf("Length of matches slice: %d", len(fileMatches))

		log.Debugf("Early exit if no matching files were found.")
		if len(fileMatches) <= 0 {

			log.WithFields(logrus.Fields{
				"path":         path,
				"file_pattern": appConfig.FilePattern,
				"extensions":   appConfig.FileExtensions,
				"file_age":     appConfig.FileAge,
			}).Info("No matches found")

			log.WithFields(logrus.Fields{
				"total_paths":   totalPaths,
				"iteration":     pass,
				"ignore_errors": appConfig.IgnoreErrors,
			}).Infof("Ending processing of path %q (%d of %d)",
				path, pass, totalPaths)
			if pass < totalPaths {
				log.Debugf("Continuing to next available path")
			}
			continue

		}

		log.WithFields(logrus.Fields{
			"path":            path,
			"file_pattern":    appConfig.FilePattern,
			"extensions":      appConfig.FileExtensions,
			"file_age":        appConfig.FileAge,
			"total_file_size": fileMatches.TotalFileSize(),
		}).Infof("%d files eligible for removal (%s)",
			len(fileMatches),
			fileMatches.TotalFileSizeHR())

		appResults.EligibleRemove += len(fileMatches)
		appResults.EligibleFileSize += fileMatches.TotalFileSize()

		log.WithFields(logrus.Fields{
			"keep_oldest": appConfig.KeepOldest,
		}).Infof("%d files to keep as requested", appConfig.NumFilesToKeep)

		// NOTE: If this sort order changes, make sure to update the logic for
		// `appConfig.KeepOldest` which retains the top or bottom X items
		// (specific flag to preserve X number of files while pruning the
		// others)
		fileMatches.SortByModTimeAsc()

		var pruneStartRange int
		var pruneEndRange int
		var filesToPrune matches.FileMatches

		switch {
		case appConfig.NumFilesToKeep > len(fileMatches):
			log.Debug("Specified number to keep is larger than total matches; will process all matches")
			pruneStartRange = 0
			pruneEndRange = len(fileMatches)

		case appConfig.KeepOldest:
			log.Debug("Keeping older files by skipping files towards the beginning of the list")
			log.Debug("Select matches from list at specified number to keep")
			pruneStartRange = appConfig.NumFilesToKeep
			pruneEndRange = len(fileMatches)

		case !appConfig.KeepOldest:
			log.Debug("Keeping newer files by skipping files towards the end of the list")
			log.Debug("Select matches from list starting with first item and ending with total length minus specified number to keep")
			pruneStartRange = 0
			pruneEndRange = (len(fileMatches) - appConfig.NumFilesToKeep)
		}

		log.WithFields(logrus.Fields{
			"start_range": pruneStartRange,
			"end_range":   pruneEndRange,
			"num_to_keep": appConfig.NumFilesToKeep,
		}).Debug("Building list of files to prune")
		filesToPrune = fileMatches[pruneStartRange:pruneEndRange]

		if len(filesToPrune) == 0 {
			log.Info("Nothing to prune")
			log.WithFields(logrus.Fields{
				"total_paths":   totalPaths,
				"iteration":     pass,
				"ignore_errors": appConfig.IgnoreErrors,
			}).Infof("Ending processing of path %q (%d of %d)",
				path, pass, totalPaths)
			if pass < totalPaths {
				log.Debugf("Continuing to next available path")
			}
			continue
		}

		log.WithFields(logrus.Fields{
			"files_to_prune":  len(filesToPrune),
			"total_file_size": filesToPrune.TotalFileSizeHR(),
		}).Debug("Calling cleanPath")
		log.Infof("Ignoring file removal errors: %t", appConfig.IgnoreErrors)
		removalResults, err := paths.CleanPath(filesToPrune, appConfig)

		appResults.SuccessRemoved += len(removalResults.SuccessfulRemovals)
		appResults.SuccessTotalFileSize += removalResults.SuccessfulRemovals.TotalFileSize()
		appResults.FailedRemoved += len(removalResults.FailedRemovals)
		appResults.FailedTotalFileSize += removalResults.FailedRemovals.TotalFileSize()

		// Show what we WERE able to successfully remove
		// TODO: Refactor this into a function to handle displaying results?
		log.Infof("%d files successfully removed (%s)",
			len(removalResults.SuccessfulRemovals),
			removalResults.SuccessfulRemovals.TotalFileSizeHR())
		for _, file := range removalResults.SuccessfulRemovals {
			log.WithFields(logrus.Fields{
				"failed_removal": false,
				"file_size":      file.SizeHR(),
			}).Info(file.Path)
		}

		log.Infof("%d files failed to remove (%s)",
			len(removalResults.FailedRemovals),
			removalResults.FailedRemovals.TotalFileSizeHR())
		for _, file := range removalResults.FailedRemovals {
			log.WithFields(logrus.Fields{
				"failed_removal": true,
				"file_size":      file.SizeHR(),
			}).Info(file.Path)
		}

		if err != nil {

			log.Warnf("Error encountered while processing %s: %s", path, err)

			if !appConfig.IgnoreErrors {
				log.WithFields(logrus.Fields{
					"ignore_errors": appConfig.IgnoreErrors,
				}).Warn("Error encountered and option to ignore errors not set. Exiting")
			}
			log.Warn("Error encountered, but continuing as requested.")
			continue
		}

		log.WithFields(logrus.Fields{
			"total_paths":   totalPaths,
			"iteration":     pass,
			"ignore_errors": appConfig.IgnoreErrors,
		}).Infof("Ending processing of path %q (%d of %d)",
			path, pass, totalPaths)

	}

	// Configure fields for execution summary results
	summaryLogger := log.WithFields(logrus.Fields{
		"success_removed": appResults.SuccessRemoved,
		"success_size":    units.ByteCountIEC(appResults.SuccessTotalFileSize),
		"failed_removed":  appResults.FailedRemoved,
		"failed_size":     units.ByteCountIEC(appResults.FailedTotalFileSize),
		"eligible_remove": appResults.EligibleRemove,
		"eligible_size":   units.ByteCountIEC(appResults.EligibleFileSize),

		// Not sure this "adds" anything to the summary and could be confusing
		// "total_processed": appResults.FailedRemoved + appResults.SuccessRemoved,
		// "total_size":      appResults.FailedTotalFileSize + appResults.SuccessTotalFileSize,

	})

	if problemsEncountered {
		summaryLogger.Warnf("%s completed, but issues were encountered.", appConfig.AppName)
	} else {
		summaryLogger.Infof("%s successfully completed.", appConfig.AppName)
	}

}
