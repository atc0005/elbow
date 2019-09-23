package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// Config represents a collection of configuration settings for this
// application. Config is created as early as possible upon application
// startup.
type Config struct {

	// https://godoc.org/github.com/sirupsen/logrus#Level
	// https://github.com/sirupsen/logrus/blob/de736cf91b921d56253b4010270681d33fdf7cb5/logrus.go#L81
	// https://github.com/jessevdk/go-flags#example
	// https://godoc.org/github.com/jessevdk/go-flags#hdr-Available_field_tags
	// https://github.com/jessevdk/go-flags/blob/master/examples/main.go
	// https://github.com/jessevdk/go-flags/blob/master/examples/rm.go

	FilePattern     string   `long:"pattern" description:"Substring pattern to compare filenames against. Wildcards are not supported."`
	FileExtensions  []string `long:"extension" description:"Limit search to specified file extension. Specify as needed to match multiple required extensions."`
	StartPath       string   `long:"path" required:"true" description:"Path to process."`
	RecursiveSearch bool     `long:"recurse" description:"Perform recursive search into subdirectories."`
	NumFilesToKeep  int      `long:"keep" required:"true" description:"Keep specified number of matching files."`
	KeepOldest      bool     `long:"keep-old" description:"Keep oldest files instead of newer."`
	Remove          bool     `long:"remove" description:"Remove matched files."`
	LogFormat       string   `long:"log-format" choice:"text" choice:"json" default:"text" description:"Log formatter used by logging package."`
	LogFile         string   `long:"log-file" description:"Log file used to hold logged messages."`
	LogLevel        string   `long:"log-level" choice:"panic" choice:"fatal" choice:"error" choice:"warn" choice:"info" choice:"debug" choice:"trace" default:"info" description:"Maximum log level at which messages will be logged. Log messages below this threshold will be discarded."`
	UseSyslog       bool     `long:"use-syslog" description:"Log messages to syslog in addition to other ouputs. Not supported on Windows."`
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
		LogFormat:       "text",
		LogLevel:        "info",

		// Intended to be optional
		LogFile: "",

		UseSyslog: false,
	}

}

// SetupFlags applies settings provided by command-line flags
// TODO: Pull out
func (c *Config) SetupFlags(appName string, appDesc string) *Config {

	// RETURN HERE
	// https://github.com/jessevdk/go-flags/blob/c0795c8afcf41dd1d786bebce68636c199b3bb45/flags.go#L172
	// SETUP a new named parser with description and other details?
	// this would allow grouping similar options together (log level, log file, syslog, etc)

	// https://godoc.org/github.com/jessevdk/go-flags#NewParser
	// https://godoc.org/github.com/jessevdk/go-flags#Options
	// Default = HelpFlag | PrintErrors | PassDoubleDash
	var parser = flags.NewParser(c, flags.Default)
	//var parser = flags.NewNamedParser(appName, &c, flags.Default)

	// TODO: What other handling is needed here? If the command-line arguments
	// are not as expected, exiting the application should probably be the
	// sensible next step?
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {

			// NOTE: This results in the Help output being shown twice.
			//parser.WriteHelp(os.Stdout)

			os.Exit(0)
		} else {

			// Another error was encountered. One case where we need to handle
			// printing it to stdout or stderr ourselves is when setting
			// `short:""` struct tags to greater than 1 character. In that
			// case the `os.Exit(1)` call below is NOT preceded with a helpful
			// message explaining the issue.
			// fmt.Println(err)
			//
			// Once we can be sure WE configured the struct properly with
			// valid length tags (or omitted the `short:""` tags), we can just
			// use `os.Exit(1)` here to allow the application to exit after
			// displaying the error code per our use of flags.Default when
			// setting up the parser (`flags.Default` includes the PrintErrors
			// option).
			os.Exit(1)
		}
	}

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

	// go-args `choice:""` struct tags enforce valid options
	// if !inList(c.LogFormat, c.validLogFormats) {
	// 	return false
	// }

	// LogFile is optional
	// TODO: String validation if it is set?

	// go-args `choice:""` struct tags enforce valid options
	// if !inList(c.LogLevel, c.validLogLevels) {
	// 	return false
	// }

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
