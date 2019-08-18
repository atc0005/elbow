package main

import (
	"fmt"
	"os"
	"path/filepath"
)

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
		fmt.Println("error:", err)
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
		if info.IsDir() {
			fmt.Println("Directory:", path)
		} else {
			fmt.Println("File:", path)
		}
	}

	return nil
}
