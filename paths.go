package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func crawlPath(path string, info os.FileInfo, err error) error {

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

		// process specific files
		// TODO: This should be specified by command-line
		ext := filepath.Ext(path)

		if inFileExtensionsPatterns(strings.ToLower(ext), config.FileExtensions) {
			log.Printf("Adding %s to fileMatches\n", path)
			// Created test files via:
			// touch {1..10}.test
			fileInfo, err := os.Stat(path)

			// unknown field 'os.FileInfo' in struct literal of type FileMatch (but does have FileInfo)
			//fileMatch := FileMatch{os.FileInfo: fileInfo, Path: path}

			// Positional initialization
			//fileMatch := FileMatch{fileInfo, path}

			// "If we need to refer to an embedded field directly, the type
			// name of the field, ignoring the package qualifier, serves as a
			// field name"
			// https://golang.org/doc/effective_go.html#embedding
			//
			// Explicit initialization
			fileMatch := FileMatch{FileInfo: fileInfo, Path: path}

			if err != nil {
				return fmt.Errorf("Unable to stat %s: %s", path, err)
			}
			matches = append(matches, fileMatch)

			return nil
		}

		log.Println("Skipping file:", path)

	}

	return nil
}
