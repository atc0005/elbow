package main

import (
	"fmt"

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

// NewConfig returns a new Config pointer that can be chained with builder
// methods to set multiple configuration values inline without using pointers.
func NewConfig() *Config {

	// Explicitly initialize with intended defaults
	return &Config{
		StartPath:   "",
		FilePattern: "",
		// NOTE: This creates an empty slice (not nil since there is an
		// underlying array of zero length) FileExtensions:  []string{},
		//
		// Leave at default value of nil slice instead by not providing a
		// value here
		// FileExtensions:  []string,
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
	flaggy.String(&c.FilePattern, "fp", "pattern", "Substring pattern to compare filenames against. Wildcards are not supported.")
	flaggy.StringSlice(&c.FileExtensions, "e", "extension", "Limit search to specified file extension. Specify as needed to match multiple required extensions.")
	flaggy.Int(&c.FilesToKeep, "k", "keep", "Keep specified number of matching files")
	flaggy.Bool(&c.RecursiveSearch, "r", "recurse", "Perform recursive search into subdirectories")
	flaggy.Bool(&c.KeepOldest, "ko", "keep-old", "Keep oldest files instead of newer")
	flaggy.Bool(&c.Remove, "rm", "remove", "Remove matched files")

	// Parse the flags
	flaggy.Parse()

	// https://github.com/atc0005/elbow/issues/2#issuecomment-524032239
	//
	// For flags, you can easily just check the value after calling
	// flaggy.Parse(). If the value is set to something other than the
	// default, then the caller supplied it. If it was the default value (set
	// by you or the language), then it was not used.

	return c

}

// String() satisfies the Stringer{} interface. This is intended for non-JSON
// formatting if using the TextFormatter logrus formatter.
func (c *Config) String() string {
	return fmt.Sprintf("FilePattern=%q, FileExtensions=%q, StartPath=%q, RecursiveSearch=%t, FilesToKeep=%d, KeepOldest=%t, Remove=%t",

		c.FilePattern,
		c.FileExtensions,
		c.StartPath,
		c.RecursiveSearch,
		c.FilesToKeep,
		c.KeepOldest,
		c.Remove,
	)
}
