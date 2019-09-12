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
func cleanPath(files FileMatches, ignoreErrors bool) (FileMatches, error) {

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

	for _, file := range files {

		log.Println("Removing test file:", file.Name())
		err := os.Remove(file.Name())
		if err != nil {
			log.Println(fmt.Errorf("Failed to remove %s: %s", file, err))

			// Record failed removal, proceed to the next file
			failedRemovals = append(failedRemovals, file)
			continue
		}

		// Record successful removal
		successfulRemovals = append(successfulRemovals, file)
	}

	// DEBUG
	for _, file := range failedRemovals {
		log.Println("Failed to remove:", file)
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

	// TODO: Branch at this point based off of whether the recursive option
	// was chosen.
	if config.RecursiveSearch {
		// DEBUG
		log.Println("Recursive option is enabled")
		log.Printf("%v", config)

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

				if hasValidExtension(path, config) {
					fileMatch := FileMatch{FileInfo: info, Path: path}
					matches = append(matches, fileMatch)
				}

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
			if hasValidExtension(file.Name(), config) {
				fileMatch := FileMatch{FileInfo: file, Path: config.StartPath}
				matches = append(matches, fileMatch)
			}
		}
	}

	return matches, err
}
