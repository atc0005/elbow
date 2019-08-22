package main

import (
	"fmt"
	"log"
	"os"

	"github.com/integrii/flaggy"
	"github.com/r3labs/diff"
)

var config Config

func main() {

	// create default configuration so that we can compare against it to
	// determine whether the user has provided flags
	defaultConfig := Config{}
	//fmt.Printf("Default configuration:\t%+v\n", defaultConfig)

	config = NewConfig()
	//fmt.Printf("Our configuration:\t%+v\n", config)

	changelog, err := diff.Diff(defaultConfig, config)
	if err != nil {
		log.Fatal(err)
	}

	if len(changelog) > 0 {
		log.Println("User specified command-line options")
		fmt.Printf("Changes to default settings: %+v\n", changelog)
		//fmt.Println("Changes to default settings:", changelog)
	} else {
		log.Println("User did not provide any command-line flags")
	}

	// TODO: Print error message and exit since (evidently) the target
	// starting path does not exist.
	//
	// https://gist.github.com/mattes/d13e273314c3b3ade33f
	//
	// if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
	// 	// path/to/whatever does not exist
	// }

	// if _, err := os.Stat("/path/to/whatever"); !os.IsNotExist(err) {
	// 	// path/to/whatever exists
	// }

	log.Println("Processing path:", config.StartPath)

	os.Exit(0)

}

// Config represents a collection of configuration settings for this
// application. Config is created as early as possible upon application
// startup.
type Config struct {
	FilePattern     string   `diff:"filepattern"`
	FileExtensions  []string `diff:"filextensions"`
	StartPath       string   `diff:"startpath"`
	RecursiveSearch bool     `diff:"recursivesearch"`
	FilesToKeep     int      `diff:"filestokeep"`
	KeepOldest      bool     `diff:"keepoldest"`
	Remove          bool     `diff:"remove"`
}

// NewConfig returns a newly created Config object.
func NewConfig() Config {

	// Explicitly initialize with intended defaults
	// Note: We compare against the default values in order to determine
	// whether the user has specified a particular flag
	config := Config{
		StartPath:       "",
		FilePattern:     "",
		FileExtensions:  []string{},
		FilesToKeep:     0,
		RecursiveSearch: false,
		KeepOldest:      false,
		Remove:          false,
	}

	flaggy.SetName("Elbow")
	flaggy.SetDescription("Prune content matching specific patterns, either in a single directory or recursively through a directory tree.")

	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	// Add flags
	flaggy.String(&config.StartPath, "p", "path", "Path to process")
	flaggy.String(&config.FilePattern, "fp", "pattern", "File pattern to match against")
	flaggy.StringSlice(&config.FileExtensions, "e", "extension", "Limit search to specified file extension")
	flaggy.Int(&config.FilesToKeep, "k", "keep", "Keep specified number of matching files")
	flaggy.Bool(&config.RecursiveSearch, "r", "recurse", "Perform recursive search into subdirectories")
	flaggy.Bool(&config.KeepOldest, "ko", "keep-old", "Keep oldest files instead of newer")
	flaggy.Bool(&config.Remove, "rm", "remove", "Remove matched files")

	// Parse the flag
	flaggy.Parse()

	//config.Set()

	return config
}
