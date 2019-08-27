package main

import (
	"fmt"
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
	for _, file := range filesToPrune {

		//fmt.Println("Details of file ...")
		//fmt.Printf("%T / %+v\n", file, file)
		//fmt.Println(file.ModTime().Format("2006-01-02 15:04:05"))
		log.Printf("Full path: %s, ShortPath: %s, Size: %d, Modified: %v\n",
			file.Path,
			file.Name(),
			file.Size(),
			file.ModTime().Format("2006-01-02 15:04:05"))
	}

	for _, file := range files {

		// TODO: Accumulate successful removals, return that with an error code

		//log.Println("Removing test file:", file.Name())
		//if err := os.Remove(file.Name()); err != nil {
		//log.Fatal(fmt.Errorf("Failed to remove %s: %s", file, err))
		//}
	}


	// TODO: Flesh this out
	return 0, nil


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
