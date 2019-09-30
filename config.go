// Copyright 2019 Adam Chalkley
//
// https://github.com/atc0005/elbow
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	FileAge         int      `long:"age" description:"Limit search to files that are the specified number of days old or older."`
	NumFilesToKeep  int      `long:"keep" required:"true" description:"Keep specified number of matching files."`
	KeepOldest      bool     `long:"keep-old" description:"Keep oldest files instead of newer."`
	Remove          bool     `long:"remove" description:"Remove matched files."`
	IgnoreErrors    bool     `long:"ignore-errors" description:"Ignore errors encountered during file removal."`
	LogFormat       string   `long:"log-format" choice:"text" choice:"json" default:"text" description:"Log formatter used by logging package."`
	LogFilePath     string   `long:"log-file" description:"Optional log file used to hold logged messages. If set, log messages are not displayed on the console."`
	LogFileHandle   *os.File `no-flag:"true"`
	ConsoleOutput   string   `long:"console-output" choice:"stdout" choice:"stderr" default:"stdout" description:"Specify how log messages are logged to the console."`
	LogLevel        string   `long:"log-level" choice:"emergency" choice:"alert" choice:"critical" choice:"panic" choice:"fatal" choice:"error" choice:"warn" choice:"info" choice:"notice" choice:"debug" choice:"trace" default:"info" description:"Maximum log level at which messages will be logged. Log messages below this threshold will be discarded."`
	UseSyslog       bool     `long:"use-syslog" description:"Log messages to syslog in addition to other outputs. Not supported on Windows."`
}

// NewConfig returns a new Config pointer that can be chained with builder
// methods to set multiple configuration values inline without using pointers.
func NewConfig() *Config {

	// Explicitly initialize with intended defaults
	// TODO: If we stay with go-flags (which applies defaults), is this
	// set of defaults still needed?
	return &Config{
		StartPath:   "",
		FilePattern: "",
		// NOTE: This creates an empty slice (not nil since there is an
		// underlying array of zero length) FileExtensions:  []string{},
		//
		// Leave at default value of nil slice instead by not providing a
		// value here
		// FileExtensions:  []string,
		FileAge:         0,
		NumFilesToKeep:  0,
		RecursiveSearch: false,
		KeepOldest:      false,
		Remove:          false,
		IgnoreErrors:    false,
		LogFormat:       "text",
		LogLevel:        "info",
		LogFilePath:     "",
		LogFileHandle:   nil,
		ConsoleOutput:   "stdout",
		UseSyslog:       false,
	}

}

// SetupFlags applies settings provided by command-line flags
// FIXME: go-flags doesn't use appName or appDesc. Keep?
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

// Validate verifies all struct fields have been provided acceptable
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

	// We only want to work with positive file modification times 0 is
	// acceptable as it is the default value and indicates that the user has
	// not chosen to use the flag (or has chosen improperly and it will be
	// ignored).
	if c.FileAge < 0 {
		return false
	}

	// KeepOldest is optional

	// Remove is optional

	// go-args `choice:""` struct tags enforce valid options
	// if !inList(c.LogFormat, c.validLogFormats) {
	// 	return false
	// }

	// LogFilePath is optional
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
	return fmt.Sprintf("FilePattern=%q, FileExtensions=%q, StartPath=%q, RecursiveSearch=%t, NumFilesToKeep=%d, KeepOldest=%t, Remove=%t, LogFormat=%q, LogFilePath=%q, UseSyslog=%t",

		c.FilePattern,
		c.FileExtensions,
		c.StartPath,
		c.RecursiveSearch,
		c.NumFilesToKeep,
		c.KeepOldest,
		c.Remove,
		c.LogFormat,
		c.LogFilePath,
		c.UseSyslog,
	)
}
