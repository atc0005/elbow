package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	// 	"github.com/sirupsen/logrus"
)

// PathPruningResults represents the number of files that were successfully
// removed and those that were not. This is used in various calculations and
// to provide a brief summary of results to the user at program completion.
type PathPruningResults struct {
	SuccessfulRemovals FileMatches
	FailedRemovals     FileMatches
}

// cleanPath receives a slice of FileMatch objects and removes each file. Any
// errors encountered while removing files may optionally be ignored (default
// is to return immediately upon first error). The total number of files
// successfully removed is returned along with an error code (nil if no errors
// were encountered).
func cleanPath(files FileMatches, ignoreErrors bool, config *Config) (PathPruningResults, error) {

	// DEBUG
	for _, file := range files {

		//fmt.Println("Details of file ...")
		//fmt.Printf("%T / %+v\n", file, file)
		//fmt.Println(file.ModTime().Format("2006-01-02 15:04:05"))

		// DEBUG
		log.Debugf("Full path: %s, ShortPath: %s, Size: %d, Modified: %v\n",
			file.Path,
			file.Name(),
			file.Size(),
			file.ModTime().Format("2006-01-02 15:04:05"))
	}

	var removalResults PathPruningResults

	if !config.Remove {

		// INFO
		log.Info("File removal not enabled.")

		// DEBUG
		log.Debug("listing what WOULD be removed")
		for _, file := range files {
			log.Debug("*", file.Name())
		}

		// Nothing to show for this yet, but since the initial state reflects
		// that we can return it as-is
		return removalResults, nil
	}

	for _, file := range files {

		filename := file.Name()

		log.Info("Removing file:", filename)
		err := os.Remove(filename)

		if err != nil {
			log.Errorf("Failed to remove %s: %s", filename, err)

			// Record failed removal, proceed to the next file
			removalResults.FailedRemovals = append(removalResults.FailedRemovals, file)

			// Confirm that we should ignore errors (likely enabled)
			if !ignoreErrors {
				break
			}
			continue
		}

		// Record successful removal
		removalResults.SuccessfulRemovals = append(removalResults.SuccessfulRemovals, file)
	}

	// TODO: Is this needed? We display this in main() ...
	for _, file := range removalResults.FailedRemovals {
		log.Debug("Failed to remove:", file.Name())
	}

	return removalResults, nil

}

func pathExists(path string) bool {

	// Make sure path isn't empty
	if strings.TrimSpace(path) == "" {
		log.Debug("path is empty string")
		return false
	}

	// https://gist.github.com/mattes/d13e273314c3b3ade33f
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Debugf("path found: %s", err)
		return true
	}

	return false

}

func processPath(config *Config) (FileMatches, error) {

	var matches FileMatches
	var err error

	if config.RecursiveSearch {
		log.Debug("Recursive option is enabled")

		// Walk walks the file tree rooted at root, calling the anonymous function
		// for each file or directory in the tree, including root. All errors that
		// arise visiting files and directories are filtered by the anonymous
		// function. The files are walked in lexical order, which makes the output
		// deterministic but means that for very large directories Walk can be
		// inefficient. Walk does not follow symbolic links.
		err = filepath.Walk(config.StartPath, func(path string, info os.FileInfo, err error) error {

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

				// ignore invalid extensions (only applies if user chose one
				// or more extensions to match against)
				if !hasValidExtension(path, config) {
					return nil
				}

				// ignore invalid filename patterns (only applies if user
				// specified a filename pattern)
				if !hasValidFilenamePattern(path, config) {
					return nil
				}

				// If we made it to this point, then we must assume that the file
				// has met all criteria to be removed by this application.
				fileMatch := FileMatch{FileInfo: info, Path: path}
				matches = append(matches, fileMatch)

			}

			return err
		})

	} else {

		// If RecursiveSearch is not enabled, process just the provided StartPath
		// NOTE: The same cleanPath() function is used in either case, the
		// difference is in how the FileMatches slice is populated

		log.Debug("Recursive option is NOT enabled")
		log.Debugf("%v", config)

		// err is already declared earlier at a higher scope, so do not
		// redeclare here
		var files []os.FileInfo
		files, err = ioutil.ReadDir(config.StartPath)

		// TODO: Do we really want to exit early at this point if there are
		// failures evaluating some of the files?
		// Is it possible to partially evaluate some of the files?
		// TODO: Wrap error?
		// if err != nil {
		// 	log.Fatal("Error from ioutil.ReadDir():", err)
		// }

		// Use []os.FileInfo returned from ioutil.ReadDir() to build slice of
		// FileMatch objects
		for _, file := range files {

			filename := file.Name()

			// Apply validity checks against filename. If validity fails,
			// go to the next file in the list.

			// ignore invalid extensions (only applies if user chose one
			// or more extensions to match against)
			if !hasValidExtension(filename, config) {
				continue
			}

			// ignore invalid filename patterns (only applies if user
			// specified a filename pattern)
			if !hasValidFilenamePattern(filename, config) {
				continue
			}

			// If we made it to this point, then we must assume that the file
			// has met all criteria to be removed by this application.
			fileMatch := FileMatch{FileInfo: file, Path: filename}
			matches = append(matches, fileMatch)
		}
	}

	return matches, err
}
