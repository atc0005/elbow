package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// cleanPath receives a slice of FileMatch objects and removes each file. Any
// errors encountered while removing files may optionally be ignored (default
// is to return immediately upon first error). The total number of files
// successfully removed is returned along with an error code (nil if no errors
// were encountered).
func cleanPath(files FileMatches, ignoreErrors bool, config *Config) (FileMatches, error) {

	// DEBUG
	for _, file := range files {

		//fmt.Println("Details of file ...")
		//fmt.Printf("%T / %+v\n", file, file)
		//fmt.Println(file.ModTime().Format("2006-01-02 15:04:05"))

		// DEBUG
		log.Printf("Full path: %s, ShortPath: %s, Size: %d, Modified: %v\n",
			file.Path,
			file.Name(),
			file.Size(),
			file.ModTime().Format("2006-01-02 15:04:05"))
	}

	// NOTE: I considered cloning the existing slice and just removing
	// elements matching failed removals, but wasn't 100% sure of costs
	// involved. This might still be the best way to go, so noting the
	// idea here for future review.
	// https://yourbasic.org/golang/delete-element-slice/
	// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
	// https://github.com/golang/go/wiki/SliceTricks
	var successfulRemovals FileMatches
	var failedRemovals FileMatches

	if !config.Remove {

		// INFO
		log.Println("File removal not enabled.")

		// DEBUG
		log.Println("listing what WOULD be removed")
		log.Println("----------------------------")
		for _, file := range files {
			log.Println("*", file.Name())
		}

		// Nothing to show for this, but since neither of the FileMatches
		// slices have been touched we can return them
		return successfulRemovals, nil
	}

	for _, file := range files {

		filename := file.Name()

		// INFO
		log.Println("Removing file:", filename)

		//
		err := os.Remove(filename)

		if err != nil {
			log.Println(fmt.Errorf("Failed to remove %s: %s", filename, err))

			// Record failed removal, proceed to the next file
			failedRemovals = append(failedRemovals, file)
			continue
		}

		// Record successful removal
		successfulRemovals = append(successfulRemovals, file)
	}

	// DEBUG
	for _, file := range failedRemovals {
		log.Println("Failed to remove:", file.Name())
	}

	return successfulRemovals, nil

}

func pathExists(path string) bool {

	// Make sure path isn't empty
	if strings.TrimSpace(path) == "" {
		return false
	}

	// https://gist.github.com/mattes/d13e273314c3b3ade33f
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false

}

func processPath(config *Config) (FileMatches, error) {

	var matches FileMatches
	var err error

	if config.RecursiveSearch {
		// DEBUG
		log.Println("Recursive option is enabled")
		//log.Printf("%v", config)

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

				if !hasValidExtension(path, config) {
					return nil
				}

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

		// DEBUG
		log.Println("Recursive option is NOT enabled")
		log.Printf("%v", config)

		files, err := ioutil.ReadDir(config.StartPath)

		// TODO: Do we really want to exit early at this point if there are
		// failures evaluating some of the files?
		// Is it possible to partially evaluate some of the files?
		// TODO: Wrap error?
		if err != nil {
			log.Fatal("Error from ioutil.ReadDir():", err)
		}

		// Use []os.FileInfo returned from ioutil.ReadDir() to build slice of
		// FileMatch objects
		for _, file := range files {

			filename := file.Name()

			// Apply validity checks against filename. If validity fails,
			// go to the next file in the list.

			if !hasValidExtension(filename, config) {
				continue
			}

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
