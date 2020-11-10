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
	"errors"
	"fmt"
	"os"

	"github.com/atc0005/elbow/config"
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

	// Checked at the end to determine if any non-fatal issues were
	// encountered during app run
	problemsEncountered := false

	// If this fails, the application should immediately exit.
	appConfig, err := config.NewConfig(version)
	if err != nil {
		// NOTE: We're not using `log` here as the user-specified
		// configuration could be too botched to use reliably.

		// TODO: Any point in setting this? Perhaps wire it into a function
		// that handles setting the flag and other useful spindown tasks?
		problemsEncountered = true

		// Provide user with error and valid usage details
		fmt.Printf("Failed to process configuration:\n%s\n\n", err)
		config.WriteDefaultHelpText(config.DefaultAppName)
		os.Exit(1)

	}

	log := appConfig.GetLogger()

	log.Debug("Config object created")

	// https://www.joeshaw.org/dont-defer-close-on-writable-files/
	if appConfig.GetLogFileHandle() != nil {
		log.Debug("Deferring closure of log file")
		defer func() {
			if err := appConfig.GetLogFileHandle().Close(); err != nil {
				// Ignore "file already closed" errors
				if !errors.Is(err, os.ErrClosed) {
					log.Errorf(
						"failed to close log file %q: %s",
						appConfig.GetLogFilePath(),
						err.Error(),
					)
				}
			}
		}()
	}

	log.WithFields(logrus.Fields{
		"paths":        appConfig.GetPaths(),
		"file_pattern": appConfig.GetFilePattern(),
		"extensions":   appConfig.GetFileExtensions(),
		"file_age":     appConfig.GetFileAge(),
	}).Info("Starting evaluation of paths list")

	// Used as a global counter/bucket for presentation/logging purposes
	var appResults paths.ProcessingResults

	var pass int
	var totalPaths int = len(appConfig.GetPaths())
	for _, path := range appConfig.GetPaths() {

		pass++

		log.WithFields(logrus.Fields{
			"total_paths":   totalPaths,
			"iteration":     pass,
			"ignore_errors": appConfig.GetIgnoreErrors(),
		}).Infof("Beginning processing of path %q (%d of %d)",
			path, pass, totalPaths)

		log.Debug("Confirm that requested path actually exists")
		pathExists, pathErr := paths.PathExists(path)
		if pathErr != nil {
			// checked at end of application run for summary report
			problemsEncountered = true

			log.WithFields(logrus.Fields{
				"ignore_errors": appConfig.GetIgnoreErrors(),
				"iteration":     pass,
			}).Error("error:", pathErr)

			if !appConfig.GetIgnoreErrors() {
				log.WithFields(logrus.Fields{
					"ignore_errors": appConfig.GetIgnoreErrors(),
				}).Warn("Error encountered and option to ignore errors not set. Exiting")
				return
			}
			log.Warn("Error encountered, but continuing as requested.")
		}

		if !pathExists {

			// checked at end of application run for summary report
			problemsEncountered = true

			log.WithFields(logrus.Fields{
				"ignore_errors": appConfig.GetIgnoreErrors(),
			}).Errorf("Requested path not found: %q", path)

			if appConfig.GetIgnoreErrors() {
				log.WithFields(logrus.Fields{
					"ignore_errors": appConfig.GetIgnoreErrors(),
				}).Warn("Error encountered, but continuing as requested.")
				continue
			} else {
				log.WithFields(logrus.Fields{
					"ignore_errors": appConfig.GetIgnoreErrors(),
				}).Warn("Error encountered and option to ignore errors not set. Exiting")

				return
			}
		}

		fileMatches, err := paths.ProcessPath(appConfig, path)
		if err != nil {

			// checked at end of application run for summary report
			problemsEncountered = true

			log.WithFields(logrus.Fields{
				"ignore_errors": appConfig.GetIgnoreErrors(),
				"iteration":     pass,
			}).Error("error:", err)

			if !appConfig.GetIgnoreErrors() {
				log.WithFields(logrus.Fields{
					"ignore_errors": appConfig.GetIgnoreErrors(),
				}).Warn("Error encountered and option to ignore errors not set. Exiting")
				return
			}
			log.Warn("Error encountered, but continuing as requested.")
		}

		log.Debugf("Length of matches slice: %d", len(fileMatches))

		log.Debugf("Early exit if no matching files were found.")
		if len(fileMatches) == 0 {

			log.WithFields(logrus.Fields{
				"path":         path,
				"file_pattern": appConfig.GetFilePattern(),
				"extensions":   appConfig.GetFileExtensions(),
				"file_age":     appConfig.GetFileAge(),
				"iteration":    pass,
			}).Info("No matches found")

			log.WithFields(logrus.Fields{
				"total_paths":   totalPaths,
				"iteration":     pass,
				"ignore_errors": appConfig.GetIgnoreErrors(),
			}).Infof("Ending processing of path %q (%d of %d)",
				path, pass, totalPaths)
			if pass < totalPaths {
				log.Debugf("Continuing to next available path")
			}
			continue

		}

		log.WithFields(logrus.Fields{
			"path":            path,
			"file_pattern":    appConfig.GetFilePattern(),
			"extensions":      appConfig.GetFileExtensions(),
			"file_age":        appConfig.GetFileAge(),
			"total_file_size": fileMatches.TotalFileSize(),
			"iteration":       pass,
		}).Infof("%d files eligible for removal (%s)",
			len(fileMatches),
			fileMatches.TotalFileSizeHR())

		appResults.EligibleRemove += len(fileMatches)
		appResults.EligibleFileSize += fileMatches.TotalFileSize()

		log.WithFields(logrus.Fields{
			"keep_oldest": appConfig.GetKeepOldest(),
			"iteration":   pass,
		}).Infof("%d files to keep as requested", appConfig.GetNumFilesToKeep())

		filesToPrune := fileMatches.FilesToPrune(appConfig)

		if len(filesToPrune) == 0 {
			log.Info("Nothing to prune")
			log.WithFields(logrus.Fields{
				"total_paths":   totalPaths,
				"iteration":     pass,
				"ignore_errors": appConfig.GetIgnoreErrors(),
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
			"iteration":       pass,
		}).Debug("Calling cleanPath")
		log.Infof("Ignoring file removal errors: %t", appConfig.GetIgnoreErrors())
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
				"iteration":      pass,
			}).Info(file.Path)
		}

		log.Infof("%d files failed to remove (%s)",
			len(removalResults.FailedRemovals),
			removalResults.FailedRemovals.TotalFileSizeHR())
		for _, file := range removalResults.FailedRemovals {
			log.WithFields(logrus.Fields{
				"failed_removal": true,
				"file_size":      file.SizeHR(),
				"iteration":      pass,
			}).Info(file.Path)
		}

		// this is the error checking for paths.CleanPath()
		if err != nil {

			// checked at end of application run for summary report
			problemsEncountered = true

			log.Warnf("Error encountered while processing %s: %s", path, err)

			if !appConfig.GetIgnoreErrors() {
				log.WithFields(logrus.Fields{
					"ignore_errors": appConfig.GetIgnoreErrors(),
					"iteration":     pass,
				}).Warn("Error encountered and option to ignore errors not set. Exiting")
				return
			}
			log.Warn("Error encountered, but continuing as requested.")
			continue
		}

		log.WithFields(logrus.Fields{
			"total_paths":   totalPaths,
			"iteration":     pass,
			"ignore_errors": appConfig.GetIgnoreErrors(),
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
		summaryLogger.Warnf("%s completed, but issues were encountered.", appConfig.GetAppName())
	} else {
		summaryLogger.Infof("%s successfully completed.", appConfig.GetAppName())
	}

}
