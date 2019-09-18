package main

import (
	"fmt"
	"log"
	"os"

	"github.com/integrii/flaggy"
	"github.com/r3labs/diff"
)

func main() {

	// create default configuration so that we can compare against it to
	// determine whether the user has provided flags
	defaultConfig := NewConfig()
	//fmt.Printf("Default configuration:\t%+v\n", defaultConfig)

	appName := "Elbow"
	appDesc := "Prune content matching specific patterns, either in a single directory or recursively through a directory tree."

	config := NewConfig().SetupFlags(appName, appDesc)
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

// NOTE: I've found multiple examples that all return a pointer in order to
// support "chaining" where the new config feeds directly into the next
// method
// https://github.com/go-sql-driver/mysql/blob/877a9775f06853f611fb2d4e817d92479242d1cd/dsn.go#L67
// https://github.com/aws/aws-sdk-go/blob/10878ad0389c5b3069815112ce888b191c8cd325/aws/config.go#L251
// https://github.com/aws/aws-sdk-go/blob/master/aws/config.go
//
// https://github.com/aws/aws-sdk-go/blob/10878ad0389c5b3069815112ce888b191c8cd325/awstesting/integration/performance/s3GetObject/config.go#L25
// https://github.com/aws/aws-sdk-go/blob/10878ad0389c5b3069815112ce888b191c8cd325/awstesting/integration/performance/s3GetObject/main.go#L25
//
//
/*
func NewConfig() *Config {
	return &Config{}
}
// WithRegion sets a config Region value returning a Config pointer for
// chaining.
func (c *Config) WithRegion(region string) *Config {
	c.Region = &region
	return c
}
*/

// NewConfig returns a new Config pointer that can be chained with builder
// methods to set multiple configuration values inline without using pointers.
func NewConfig() *Config {

	// Explicitly initialize with intended defaults
	// Note: We compare against the default values in order to determine
	// whether the user has specified a particular flag
	return &Config{
		StartPath:       "",
		FilePattern:     "",
		FileExtensions:  []string{},
		FilesToKeep:     0,
		RecursiveSearch: false,
		KeepOldest:      false,
		Remove:          false,
	}

}

// SetupFlags applies settings provided by command-line flags
// TODO: Pull out
func (c *Config) SetupFlags(appName string, appDesc string) *Config {

	flaggy.SetName(appName)
	flaggy.SetDescription(appDesc)

	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	// Add flags
	flaggy.String(&c.StartPath, "p", "path", "Path to process")
	flaggy.String(&c.FilePattern, "fp", "pattern", "File pattern to match against")
	flaggy.StringSlice(&c.FileExtensions, "e", "extension", "Limit search to specified file extension")
	flaggy.Int(&c.FilesToKeep, "k", "keep", "Keep specified number of matching files")
	flaggy.Bool(&c.RecursiveSearch, "r", "recurse", "Perform recursive search into subdirectories")
	flaggy.Bool(&c.KeepOldest, "ko", "keep-old", "Keep oldest files instead of newer")
	flaggy.Bool(&c.Remove, "rm", "remove", "Remove matched files")

	// Parse the flags
	flaggy.Parse()

	return c

}
