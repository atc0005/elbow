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

// NewConfig returns a newly created Config object.
func NewConfig() Config {

	// Explicitly initialize with intended defaults
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
