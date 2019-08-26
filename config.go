package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/integrii/flaggy"
)

// Config represents a collection of configuration settings for this
// application. Config is created as early as possible upon application
// startup.
type Config struct {
	FilePattern     string
	FileExtensions  []string
	StartPath       string
	RecursiveSearch bool
	FilesToKeep     int
	KeepOldest      bool
	Remove          bool
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


// Set is a helper method used to configure various values for the application
// Config object.
// TODO: Handle setting based on provided values from user.
func (c *Config) Set() {

	// non-name c.StartPath on left side of :=
	//c.StartPath, err := GetStartPath()

	var err error
	c.StartPath, err = GetStartPath()
	if err != nil {
		log.Fatal(err)
	}

	c.FileExtensions, err = GetFileExtensionsPattern()
	if err != nil {
		log.Fatal(err)
	}
}

// GetStartPath is used to retrieve the starting point/path for processing.
func GetStartPath() (string, error) {

	// TODO: Replace this hard-coded path with a value from command-line
	path, ok := os.LookupEnv("TEMP")
	if !ok {
		return "", fmt.Errorf("Unable to retrieve TEMP environment variable")
	}

	startPath := filepath.FromSlash(path)

	return startPath, nil
}

// GetFileExtensionsPattern is used to match files to be pruned. This setting
// is complimentary to FilePattern and acts as a filter or constraint to limit
// the file matches.
func GetFileExtensionsPattern() ([]string, error) {
	// TODO: Hard-coded test values for now
	return []string{".tmp", ".test"}, nil
}

// Helper function to emulate Python's `if "x" in list:` functionality
func inFileExtensionsPatterns(ext string, exts []string) bool {
	for _, pattern := range exts {
		if ext == pattern {
			return true
		}
	}
	return false
}
