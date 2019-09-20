package main

import (
	"fmt"

	"github.com/integrii/flaggy"
)

// Config represents a collection of configuration settings for this
// application. Config is created as early as possible upon application
// startup.
type Config struct {

	// https://godoc.org/github.com/sirupsen/logrus#Level
	// https://github.com/sirupsen/logrus/blob/de736cf91b921d56253b4010270681d33fdf7cb5/logrus.go#L81
	// https://github.com/jessevdk/go-flags#example
	// https://godoc.org/github.com/jessevdk/go-flags#hdr-Available_field_tags

	FilePattern     string   `short:"" long:"" description:""`
	FileExtensions  []string `short:"" long:"" description:""`
	StartPath       string   `short:"" long:"" description:""`
	RecursiveSearch bool     `short:"" long:"" description:""`
	NumFilesToKeep  int      `short:"" long:"" description:""`
	KeepOldest      bool     `short:"" long:"" description:""`
	Remove          bool     `short:"" long:"" description:""`
	LogFormat       string   `short:"lf" long:"log-format" choice:"text" choice:"json" description:""`
	LogFile         string   `short:"log" long:"log-file" description:""`
	LogLevel        string   `short:"ll" long:"log-level" choice:"panic" choice:"fatal" choice:"error" choice:"warn" choice:"info" choice:"debug" choice:"trace" description:""`
	UseSyslog       bool     `short:"" long:"" description:""`
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
		NumFilesToKeep:  0,
		RecursiveSearch: false,
		KeepOldest:      false,
		Remove:          false,

		// All of these will require "x in y" type validation
		LogFormat:       "text",
		validLogFormats: []string{"text", "json"},

		LogLevel: "info",

		// Intended to be optional
		LogFile: "",

		UseSyslog: false,
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
	flaggy.Int(&c.NumFilesToKeep, "k", "keep", "Keep specified number of matching files")
	flaggy.Bool(&c.RecursiveSearch, "r", "recurse", "Perform recursive search into subdirectories")
	flaggy.Bool(&c.KeepOldest, "ko", "keep-old", "Keep oldest files instead of newer")
	flaggy.Bool(&c.Remove, "rm", "remove", "Remove matched files")

	// TODO: Is there any way to avoid listing the valid options for this flag?
	flaggy.String(&c.LogFormat, "lf", "log-format", "Log formatter used by logging package. text and json are the two currently supported formatters.")

	flaggy.String(&c.LogFile, "log", "log-file", "Log file used to hold logged messages.")

	// TODO: Is the word "above" or "below" in regards to the other log
	// messages which will be discarded?
	flaggy.String(&c.LogLevel, "ll", "log-level", "Maximum log level at which messages will be logged. Log messages below this threshold will be discarded. The default level is info.")

	flaggy.Bool(&c.UseSyslog, "sl", "use-syslog", "Log messages to syslog in addition to other ouputs. Not supported on Windows.")

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

// Validate verifies all struct fields have been provided accceptable
func (c *Config) Validate() bool {

	// FilePattern is optional

	// FileExtensions is optional
	// Discovered files are checked against FileExtensions later

	if len(c.StartPath) == 0 {
		return false
	}

	// RecursiveSearch is optional

	// NumFilesToKeep is optional, but if specified we should make sure it is
	// a non-negative number. AFAIK, this is not currently enforced any other
	// way.
	if c.NumFilesToKeep < 0 {
		return false
	}

	// KeepOldest is optional

	// Remove is optional

	if !inList(c.LogFormat, c.validLogFormats) {
		return false
	}

	// LogFile is optional
	// TODO: String validation if it is set?

	if !inList(c.LogLevel, c.validLogLevels) {
		return false
	}

	// UseSyslog is optional

	// Optimist
	return true

}

// String() satisfies the Stringer{} interface. This is intended for non-JSON
// formatting if using the TextFormatter logrus formatter.
func (c *Config) String() string {
	return fmt.Sprintf("FilePattern=%q, FileExtensions=%q, StartPath=%q, RecursiveSearch=%t, NumFilesToKeep=%d, KeepOldest=%t, Remove=%t, LogFormat=%q, LogFile=%q, UseSyslog=%t",

		c.FilePattern,
		c.FileExtensions,
		c.StartPath,
		c.RecursiveSearch,
		c.NumFilesToKeep,
		c.KeepOldest,
		c.Remove,
		c.LogFormat,
		c.LogFile,
		c.UseSyslog,
	)
}
