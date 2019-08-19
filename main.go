package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/peterbourgon/ff"
)

// TODO: What other option do I have here other than using globals?
var matches FileMatches

var config = NewConfig()

func main() {

	/*

		TODO: Collect these options (command-line / config file / env vars)

		1) [string] File pattern to match on
		2) [string] Starting path for processing
		3) [bool] Recursive search
		4) [int] Number of files to keep
		5) [bool] KeepYoungest (true, default)

	*/

	fs := flag.NewFlagSet("elbow", flag.ExitOnError)
	var (
		flagStartPath     = fs.String("path", "", "path to process")
		flagFilePattern   = fs.String("pattern", "", "file pattern to match against")
		flagFileExtension = fs.String("ext", "", "file extension to match against to limit search")
		flagRemove        = fs.Bool("remove", false, "remove matched files")
		flagRecurse       = fs.Bool("recurse", false, "recurse into subdirectories")
		flagKeep          = fs.Int("keep", 0, "keep specified number of matching files")
		flagKeepOldest    = fs.Bool("keep-oldest", false, "keep oldest files instead of newer")
		_                 = fs.String("config", "", "config file (optional)")
	)

	ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
		ff.WithEnvVarPrefix("ELBOW"),
	)

	fmt.Printf("startPath: %q\n", *flagStartPath)
	fmt.Printf("filePattern: %q\n", *flagFilePattern)
	fmt.Printf("fileExtension: %q\n", *flagFileExtension)
	fmt.Printf("remove: %v\n", *flagRemove)
	fmt.Printf("recurse: %v\n", *flagRecurse)
	fmt.Printf("keep: %v\n", *flagKeep)
	fmt.Printf("keepOldest: %v\n", *flagKeepOldest)

	log.Println("Processing path:", config.StartPath)

	os.Exit(0)

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
