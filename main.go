package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TODO: What other option do I have here other than using a global?
var fileMatches []string

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
		log.Println("Removing test file:", file)
		if err := os.Remove(file); err != nil {
			log.Fatal(fmt.Sprintf("Failed to remove %s: %s", file, err))
		}
	}

}

// WalkFunc is the type of the function called for each file or directory
// visited by Walk. The path argument contains the argument to Walk as a
// prefix; that is, if Walk is called with "dir", which is a directory
// containing the file "a", the walk function will be called with argument
// "dir/a". The info argument is the os.FileInfo for the named path.

// If there was a problem walking to the file or directory named by path, the
// incoming error will describe the problem and the function can decide how to
// handle that error (and Walk will not descend into that directory). In the
// case of an error, the info argument will be nil. If an error is returned,
// processing stops. The sole exception is when the function returns the
// special value SkipDir. If the function returns SkipDir when invoked on a
// directory, Walk skips the directory's contents entirely. If the function
// returns SkipDir when invoked on a non-directory file, Walk skips the
// remaining files in the containing directory.
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
		// Q: Why `return nil` instead of negating (`if ! info.IsDir(){}`)
		// and processing inside the brackets?
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
			fileMatches = append(fileMatches, path)

			return nil
		}

		log.Println("Skipping file:", path)

	}

	return nil
}
