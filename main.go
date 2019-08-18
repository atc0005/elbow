package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TODO: What other option do I have here other than using a global?
var fileMatches []os.FileInfo

func main() {

	// TODO: Need to reference the value passed from command-line
	root, _ := filepath.Abs(".")

	fmt.Println("Processing path", root)

	// Walk walks the file tree rooted at root, calling processPath for each
	// file or directory in the tree, including root. All errors that arise
	// visiting files and directories are filtered by processPath. The files
	// are walked in lexical order, which makes the output deterministic but
	// means that for very large directories Walk can be inefficient. Walk
	// does not follow symbolic links.
	err := filepath.Walk(root, processPath)
	if err != nil {
		log.Println("error:", err)
	}

	/*

		// https://stackoverflow.com/questions/46746862/list-files-in-a-directory-sorted-by-creation-time

		files, err := ioutil.ReadDir(path)
		// TODO: handle the error!
		sort.Slice(files, func(i,j int) bool{
			return files[i].ModTime().Before(files[j].ModTime())
		})

	*/

	// We need to first collect the entire list of files and then prune
	// all BUT the youngest three. Here with the prototype we can just
	// focus on files 4-10, but our real target is based on the file's
	// date/time values.

	// TODO: Sort prior to this point Skip the first three files from the
	// list. By this point they should have already been sorted based on age.
	// TODO: How do you make the starting point dynamic based on command-line
	// option chosen?
	pruneFilesStartPoint := 2
	for _, file := range fileMatches[pruneFilesStartPoint:] {
		log.Println("Removing test file:", file.Name())
		if err := os.Remove(file.Name()); err != nil {
			log.Fatal(fmt.Errorf("Failed to remove %s: %s", file, err))
		}
	}

}

func processPath(path string, info os.FileInfo, err error) error {

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

		if strings.ToLower(ext) == ".test" {
			log.Printf("Adding %s to fileMatches\n", path)
			// Created test files via:
			// touch {1..10}.test
			fileInfo, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("Unable to stat %s: %s", path, err)
			}
			fileMatches = append(fileMatches, fileInfo)

			return nil
		}

		log.Println("Skipping file:", path)

	}

	return nil
}
