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

// Package paths provides various functions and types related to processing
// paths in the filesystem, often for the purpose of removing older/unwanted
// files.
package paths

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/atc0005/elbow/config"
	"github.com/atc0005/elbow/matches"
	"github.com/sirupsen/logrus"
)

// ProcessingResults is used to collect execution results for use in logging
// and output summary presentation to the user
type ProcessingResults struct {

	// Number of files eligible for removal. This is before files are excluded
	// by request per user request.
	EligibleRemove int

	// Number of files successfully removed.
	SuccessRemoved int

	// Number of files failed to remove.
	FailedRemoved int

	// Size of all files eligible for removal.
	EligibleFileSize int64

	// Size of all files successfully removed.
	SuccessTotalFileSize int64

	// Size of all files failed to remove.
	FailedTotalFileSize int64

	// Size of all files successfully and unsuccessfully removed. This is
	// essentially the size of eligible files to be removed minus any files
	// that are excluded by user request.
	TotalProcessedFileSize int64
}

// PathPruningResults represents the number of files that were successfully
// removed and those that were not. This is used in various calculations and
// to provide a brief summary of results to the user at program completion.
type PathPruningResults struct {
	SuccessfulRemovals matches.FileMatches
	FailedRemovals     matches.FileMatches
}

// CleanPath receives a slice of FileMatch objects and removes each file. Any
// errors encountered while removing files may optionally be ignored via
// command-line flag(default is to return immediately upon first error). The
// total number of files successfully removed is returned along with an error
// code (nil if no errors were encountered).
func CleanPath(files matches.FileMatches, config *config.Config) (PathPruningResults, error) {

	log := config.GetLogger()

	for _, file := range files {
		log.WithFields(logrus.Fields{
			"fullpath":        strings.TrimSpace(file.Path),
			"shortpath":       file.Name(),
			"size":            file.Size(),
			"modified":        file.ModTime().Format("2006-01-02 15:04:05"),
			"removal_enabled": config.GetRemove(),
		}).Debug("Matching file")
	}

	var removalResults PathPruningResults

	if !config.GetRemove() {

		log.Info("File removal not enabled, not removing files")

		// Nothing to show for this yet, but since the initial state reflects
		// that we can return it as-is
		return removalResults, nil
	}

	for _, file := range files {

		log.WithFields(logrus.Fields{
			"removal_enabled": config.GetRemove(),

			// fully-qualified path to the file
			"file": file.Path,
		}).Debug("Removing file")

		// We need to reference the full path here, not the short name since
		// the current working directory may not be the same directory
		// where the file is located
		err := os.Remove(file.Path)
		if err != nil {
			log.WithFields(logrus.Fields{

				// Include full details for troubleshooting purposes
				"file": file,
			}).Errorf("Error encountered while removing file: %s", err)

			// Record failed removal, proceed to the next file
			removalResults.FailedRemovals = append(removalResults.FailedRemovals, file)

			// Confirm that we should ignore errors (likely enabled)
			if !config.GetIgnoreErrors() {
				remainingFiles := len(files) - len(removalResults.FailedRemovals) - len(removalResults.SuccessfulRemovals)
				log.Debugf("Abandoning removal of %d remaining files", remainingFiles)
				break
			}

			log.Debug("Ignoring error as requested")
			continue
		}

		// Record successful removal
		removalResults.SuccessfulRemovals = append(removalResults.SuccessfulRemovals, file)
	}

	return removalResults, nil

}

// PathExists confirms that the specified path exists
// FIXME: Update this to break reliance on config.Config; use bare args
func PathExists(path string, config *config.Config) bool {

	log := config.GetLogger()

	// Make sure path isn't empty
	if strings.TrimSpace(path) == "" {
		log.Debugf("path is empty string")
		return false
	}

	// https://gist.github.com/mattes/d13e273314c3b3ade33f
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.WithFields(logrus.Fields{
			"path": path,
		}).Debug("path found")
		return true
	}

	return false

}

// ProcessPath accepts a configuration object and a path to process and
// returns a slice of FileMatch objects
func ProcessPath(config *config.Config, path string) (matches.FileMatches, error) {

	log := config.GetLogger()

	var fileMatches matches.FileMatches
	var err error

	log.WithFields(logrus.Fields{
		"recursive_search": config.GetRecursiveSearch(),
	}).Debugf("Recursive search: %t", config.GetRecursiveSearch())

	if config.GetRecursiveSearch() {

		// Walk walks the file tree rooted at root, calling the anonymous function
		// for each file or directory in the tree, including root. All errors that
		// arise visiting files and directories are filtered by the anonymous
		// function. The files are walked in lexical order, which makes the output
		// deterministic but means that for very large directories Walk can be
		// inefficient. Walk does not follow symbolic links.
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

			// If an error is received, return it. If we return a non-nil error, this
			// will stop the filepath.Walk() function from continuing to walk the
			// path, and your main function will immediately move to the next line.
			if err != nil {
				return err
			}

			// make sure we're not working with the root directory itself
			if path != "." {

				// ignore directories
				if info.IsDir() {
					return nil
				}

				// ignore non-matching extension (only applies if user chose
				// one or more extensions to match against)
				if !matches.HasMatchingExtension(path, config) {
					return nil
				}

				// ignore non-matching filename pattern (only applies if user
				// specified a filename pattern)
				if !matches.HasMatchingFilenamePattern(path, config) {
					return nil
				}

				// ignore non-matching modification age
				if !matches.HasMatchingAge(info, config) {
					return nil
				}

				// If we made it to this point, then we must assume that the file
				// has met all criteria to be removed by this application.
				fileMatch := matches.FileMatch{FileInfo: info, Path: path}
				fileMatches = append(fileMatches, fileMatch)

			}

			return err
		})

	} else {

		// If RecursiveSearch is not enabled, process just the provided StartPath
		// NOTE: The same cleanPath() function is used in either case, the
		// difference is in how the FileMatches slice is populated

		// err is already declared earlier at a higher scope, so do not
		// redeclare here
		var files []os.FileInfo
		files, err = ioutil.ReadDir(path)

		if err != nil {
			// TODO: Do we really want to exit early at this point if there are
			// failures evaluating some of the files?
			// Is it possible to partially evaluate some of the files?
			// TODO: Wrap error?
			log.Errorf("Error from ioutil.ReadDir(): %s", err)
		}

		// Use []os.FileInfo returned from ioutil.ReadDir() to build slice of
		// FileMatch objects
		for _, file := range files {

			// ignore directories
			if file.IsDir() {
				continue
			}

			filename := file.Name()

			// Apply validity checks against filename. If validity fails,
			// go to the next file in the list.

			// ignore invalid extensions (only applies if user chose one
			// or more extensions to match against)
			if !matches.HasMatchingExtension(filename, config) {
				continue
			}

			// ignore invalid filename patterns (only applies if user
			// specified a filename pattern)
			if !matches.HasMatchingFilenamePattern(filename, config) {
				continue
			}

			// ignore non-matching modification age
			if !matches.HasMatchingAge(file, config) {
				continue
			}

			// If we made it to this point, then we must assume that the file
			// has met all criteria to be removed by this application.
			fileMatch := matches.FileMatch{FileInfo: file, Path: filename}
			fileMatches = append(fileMatches, fileMatch)
		}
	}

	return fileMatches, err
}
