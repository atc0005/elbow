package main

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"

	"github.com/integrii/flaggy"
)

// TODO: What other option do I have here other than using globals?
var matches FileMatches
var config *Config

func main() {

	/*

		TODO: Collect these options (command-line / config file / env vars)

		1) [string] File pattern to match on
		2) [string] Starting path for processing
		3) [bool] Recursive search
		4) [int] Number of files to keep
		5) [bool] KeepYoungest (true, default)

	*/

	// create default configuration so that we can compare against it to
	// determine whether the user has provided flags
	defaultConfig := NewConfig()
	fmt.Printf("Default configuration:\t%+v\n", defaultConfig)

	appName := "Elbow"
	appDesc := "Prune content matching specific patterns, either in a single directory or recursively through a directory tree."

	config = NewConfig().SetupFlags(appName, appDesc)
	fmt.Printf("Our configuration:\t%+v\n", config)

	// TODO: How can I reliably compare these?
	//  invalid operation: *defaultConfig != *config (struct containing []string cannot be compared)
	// if *defaultConfig != *config {
	// 	log.Println("User specified command-line options")
	// } else {
	// 	log.Println("User did not provide any command-line flags")
	// }

	if reflect.DeepEqual(*defaultConfig, *config) {
		log.Println("User did not provide command-line flags; current configuration matches default settings")
		flaggy.ShowHelpAndExit("Required command-line options not provided.")
	} else {
		log.Println("User provided command-line flags, proceeding ...")
	}

	//fmt.Printf("%+v\n", *config)

	// Confirm that requested path actually exists
	if !pathExists(config.StartPath) {
		log.Fatalf("Error processing requested path: %q", config.StartPath)
	}

	log.Println("Processing path:", config.StartPath)

	//os.Exit(0)

	//
	// TODO: Refactor filepath.Walk() call below; split into at least two
	// functions, one to do what is being done now (recursive work), another
	// to use `ioutil.ReadDir(path)` to gather matches from specific
	// directory.
	//

	// Walk walks the file tree rooted at root, calling crawlPath for each
	// file or directory in the tree, including root. All errors that arise
	// visiting files and directories are filtered by crawlPath. The files
	// are walked in lexical order, which makes the output deterministic but
	// means that for very large directories Walk can be inefficient. Walk
	// does not follow symbolic links.
	err := filepath.Walk(config.StartPath, crawlPath)
	if err != nil {
		log.Println("error:", err)
	}

	// TODO: Does this in-place attempt at sorting work because slices are
	// reference types to begin with?
	matches.sortByModTimeAsc()

	pruneFilesStartPoint := 2
	for _, file := range matches[pruneFilesStartPoint:] {
		//log.Println("Removing test file:", file.Name())
		//if err := os.Remove(file.Name()); err != nil {
		//log.Fatal(fmt.Errorf("Failed to remove %s: %s", file, err))
		//}

		//fmt.Println("Details of file ...")
		//fmt.Printf("%T / %+v\n", file, file)
		//fmt.Println(file.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Printf("Full path: %s, ShortPath: %s, Size: %d, Modified: %v\n",
			file.Path,
			file.Name(),
			file.Size(),
			file.ModTime().Format("2006-01-02 15:04:05"))

	}

}
