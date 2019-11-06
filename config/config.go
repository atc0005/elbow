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

// Package config provides types and functions to collect, validate and apply
// user-provided settings.
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

// AppMetadata represents data about this application that may be used in Help
// output, error messages and potentially log messages (e.g., AppVersion)
type AppMetadata struct {
	AppName        string `arg:"-"`
	AppDescription string `arg:"-"`
	AppVersion     string `arg:"-"`
	AppURL         string `arg:"-"`
}

// FileHandling represents options specific to how this application
// handles files.
type FileHandling struct {
	FilePattern    string   `toml:"pattern" arg:"--pattern,env:ELBOW_FILE_PATTERN" default:"" help:"Substring pattern to compare filenames against. Wildcards are not supported."`
	FileExtensions []string `toml:"file_extensions" arg:"--extensions,env:ELBOW_EXTENSIONS" help:"Limit search to specified file extensions. Specify as space separated list to match multiple required extensions."`
	FileAge        int      `toml:"file_age" arg:"--age,env:ELBOW_FILE_AGE" default:"0" help:"Limit search to files that are the specified number of days old or older."`
	NumFilesToKeep int      `toml:"files_to_keep" arg:"--keep,required,env:ELBOW_KEEP" help:"Keep specified number of matching files per provided path."`
	KeepOldest     bool     `toml:"keep_oldest" arg:"--keep-old,env:ELBOW_KEEP_OLD" default:"false" help:"Keep oldest files instead of newer per provided path."`
	Remove         bool     `toml:"remove" arg:"--remove,env:ELBOW_REMOVE" default:"false" help:"Remove matched files per provided path."`
	IgnoreErrors   bool     `toml:"ignore_errors" arg:"--ignore-errors,env:ELBOW_IGNORE_ERRORS" default:"false" help:"Ignore errors encountered during file removal."`
}

// Search represents options specific to controlling how this application
// performs searches in the filesystem
type Search struct {
	Paths           []string `toml:"paths" arg:"--paths,required,env:ELBOW_PATHS" help:"List of comma or space-separated paths to process."`
	RecursiveSearch bool     `toml:"recursive_search" arg:"--recurse,env:ELBOW_RECURSE" default:"false" help:"Perform recursive search into subdirectories per provided path."`
}

// Logging represents options specific to how this application handles
// logging.
type Logging struct {
	LogLevel      string `toml:"log_level" arg:"--log-level,env:ELBOW_LOG_LEVEL" default:"info" help:"Maximum log level at which messages will be logged. Log messages below this threshold will be discarded."`
	LogFormat     string `toml:"log_format" arg:"--log-format,env:ELBOW_LOG_FORMAT" default:"text" help:"Log formatter used by logging package."`
	LogFilePath   string `toml:"log_file_path" arg:"--log-file,env:ELBOW_LOG_FILE" default:"" help:"Optional log file used to hold logged messages. If set, log messages are not displayed on the console."`
	ConsoleOutput string `toml:"console_output" arg:"--console-output,env:ELBOW_CONSOLE_OUTPUT" default:"stdout" help:"Specify how log messages are logged to the console."`
	UseSyslog     bool   `toml:"use_syslog" arg:"--use-syslog,env:ELBOW_USE_SYSLOG" default:"false" help:"Log messages to syslog in addition to other outputs. Not supported on Windows."`
}

// Config represents a collection of configuration settings for this
// application. Config is created as early as possible upon application
// startup.
type Config struct {

	// Embed other structs in an effort to better group related settings
	AppMetadata
	FileHandling
	Logging
	Search

	// Embedded to allow for easier carrying of "handles" between functions
	// TODO: Confirm that this is both needed and that it doesn't violate
	// best practices.
	LogFileHandle *os.File       `arg:"-"`
	Logger        *logrus.Logger `arg:"-"`
	FlagParser    *arg.Parser    `arg:"-"`

	// Path to (optional) configuration file
	ConfigFile string `arg:"-"`
}

// NewConfig returns a newly configured object representing a collection of
// user-provided and default settings.
func NewConfig(appName, appDescription, appURL, appVersion string) *Config {

	// Note: The majority of the default settings are supplied via struct tags

	// The base configuration object that will be returned to the caller
	var config Config

	// Settings provided via command-line flags and environment variables.
	// This object will always be set in some manner as either flags or env
	// vars will be needed to bootstrap the application. While we may support
	// using a configuration file to provide settings, it is not used by
	// default.
	var argsConfig Config

	// Settings provided via config file
	var fileConfig Config

	// Initialize logger "handle" for later use
	config.Logger = logrus.New()

	config.AppName = appName
	config.AppDescription = appDescription
	config.AppURL = appURL
	config.AppVersion = appVersion

	// Bundle the returned `*.arg.Parser` for later use from `main()` so that
	// we can explicitly display usage or help details should the
	// user-provided settings fail validation.
	config.FlagParser = arg.MustParse(&argsConfig)

	// Check for a configuration file and load it if found.
	// FIXME: The method needs to be updated to reference the path provided
	// via environment variable or command-line flag
	if err := fileConfig.LoadConfigFile("config.toml"); err != nil {
		fmt.Printf("Error loading config file: %s", err)
	}

	// At this point `config` is our base config. We merge the other
	// configuration objects into it to create a unified configuration object
	// that we return to the caller.
	if err := MergeConfig(&config, fileConfig); err != nil {
		fmt.Printf("Error merging config file settings with base config: %s", err)
	}

	if err := MergeConfig(&config, argsConfig); err != nil {
		fmt.Printf("Error merging args config settings with base config: %s", err)
	}

	fmt.Printf("The config object that we are returning:\n%+v", config)

	return &config

}

// GetStructTag returns the requested struct tag value, if set, and an error
// value indicating whether any problems were encountered.
func GetStructTag(c Config, fieldname string, tagName string) (string, bool) {

	t := reflect.TypeOf(c)

	// TODO: Rip out all of the print statements after confirming this works
	// as expected.

	var field reflect.StructField
	var ok bool
	var tagValue string

	// this struct field does not have a `default` tag
	//fmt.Printf("\nProcessing %s struct field ...\n", fieldname)
	//if field, ok = t.FieldByName("fieldname"); !ok {
	// FIXME: Are the quotes needed?
	if field, ok = t.FieldByName(fieldname); !ok {
		//return "", fmt.Errorf("%q field not found", fieldname)
		return "", false
	}

	//fmt.Println(field.Tag)
	if tagValue, ok = field.Tag.Lookup(tagName); !ok {
		// return "", fmt.Errorf("%q tag not found", tag)
		return "", false
	}

	//fmt.Printf("%q, %t\n", tagValue, ok)
	return tagValue, true

}

// MergeConfig receives a source and destination Config object and merges
// non-default field values from the source Config object to the destination
// config object, overwriting any current non-default field values. The goal
// is to respect the current documented configuration precedence for multiple
// configuration sources (e.g., config file and command-line flags).
func MergeConfig(destination *Config, source Config) error {

	var tagValue string
	var ok bool

	// FIXME: How can we get all field names programatically so we don't have to
	// manually reference each field?

	fmt.Println("MergeConfig called")
	fmt.Printf("Source struct: %+v\n", source)
	fmt.Printf("Dest struct: %+v\n", *destination)

	// Copy over select source struct field values if destination struct field
	// values are empty or some other invalid state. These fields are not
	// supported by `default` value logic.
	if len(destination.Paths) <= 0 {
		destination.Paths = source.Paths
	}

	if len(destination.FileExtensions) <= 0 {
		destination.FileExtensions = source.FileExtensions
	}

	if tagValue, ok = GetStructTag(*destination, "FilePattern", "default"); ok {
		// If we were able to get the value for the requested "destination"
		// struct tag, go ahead and convert the value for the source struct
		// field to a string for comparison purposes; this converted string is not
		// used for value assignment.
		//
		// Then, check to see if the destination field value is the configured
		// default. If it is, then take whatever is in the source struct field
		// and overwrite the destination struct field of the same name.
		if string(destination.FilePattern) == tagValue {
			destination.FilePattern = source.FilePattern
		}

	}

	if tagValue, ok = GetStructTag(*destination, "FileAge", "default"); ok {
		if string(destination.FileAge) == tagValue {
			destination.FileAge = source.FileAge
		}
	}

	if tagValue, ok = GetStructTag(*destination, "NumFilesToKeep", "default"); ok {
		if string(destination.NumFilesToKeep) == tagValue {
			destination.NumFilesToKeep = source.NumFilesToKeep
		}
	}

	if tagValue, ok = GetStructTag(*destination, "KeepOldest", "default"); ok {
		if defaultValue, err := strconv.ParseBool(tagValue); err == nil {
			if destination.KeepOldest == defaultValue {
				destination.KeepOldest = source.KeepOldest
			}
		}
	}

	if tagValue, ok = GetStructTag(*destination, "Remove", "default"); ok {
		if defaultValue, err := strconv.ParseBool(tagValue); err == nil {
			if destination.Remove == defaultValue {
				destination.Remove = source.Remove
			}
		}
	}

	if tagValue, ok = GetStructTag(*destination, "IgnoreErrors", "default"); ok {
		if defaultValue, err := strconv.ParseBool(tagValue); err == nil {
			if destination.IgnoreErrors == defaultValue {
				destination.IgnoreErrors = source.IgnoreErrors
			}
		}
	}

	if tagValue, ok = GetStructTag(*destination, "RecursiveSearch", "default"); ok {
		if defaultValue, err := strconv.ParseBool(tagValue); err == nil {
			if destination.RecursiveSearch == defaultValue {
				destination.RecursiveSearch = source.RecursiveSearch
			}
		}
	}

	if tagValue, ok = GetStructTag(*destination, "LogLevel", "default"); ok {
		if string(destination.LogLevel) == tagValue {
			destination.LogLevel = source.LogLevel
		}
	}

	if tagValue, ok = GetStructTag(*destination, "LogFormat", "default"); ok {

		fmt.Printf("destination.LogFormat: %v\n", destination.LogFormat)
		fmt.Printf("source.LogFormat: %v\n", source.LogFormat)
		fmt.Printf("tagValue: %v\n", tagValue)
		if string(destination.LogFormat) == tagValue {
			// FIXME: We need to take into consideration that the user may
			// have explicitly opted into using the same value as the
			// `default` struct tag value.
			destination.LogFormat = source.LogFormat
		}
	}

	if tagValue, ok = GetStructTag(*destination, "LogFilePath", "default"); ok {
		if string(destination.LogFilePath) == tagValue {
			destination.LogFilePath = source.LogFilePath
		}
	}

	if tagValue, ok = GetStructTag(*destination, "ConsoleOutput", "default"); ok {
		if string(destination.ConsoleOutput) == tagValue {
			destination.ConsoleOutput = source.ConsoleOutput
		}
	}

	if tagValue, ok = GetStructTag(*destination, "UseSyslog", "default"); ok {
		if defaultValue, err := strconv.ParseBool(tagValue); err == nil {
			if destination.UseSyslog == defaultValue {
				destination.UseSyslog = source.UseSyslog
			}
		}
	}

	// FIXME: Placeholder
	// FIXME: What useful error code would we return from this function?
	return nil
}

// LoadConfigFile is a stub method intended to help prototype the process of
// supporting multiple valid ways to load configuration settings
func (c *Config) LoadConfigFile(filename string) error {

	// Read file to byte slice
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal parses the TOML-encoded data and stores the result in the
	// value pointed to by v. Behavior is similar to the Go json encoder,
	// except that there is no concept of an Unmarshaler interface or
	// UnmarshalTOML function for sub-structs, and currently only definite
	// types can be unmarshaled to (i.e. no `interface{}`).
	//
	// The following struct annotations are supported:
	//
	// toml:"Field" Overrides the field's name to map to.
	// default:"foo" Provides a default value.
	// For default values, only fields of the following types are supported:
	//
	// * string
	// * bool
	// * int
	// * int64
	// * float64

	if err := toml.Unmarshal(configFile, c); err != nil {
		return err
	}

	// Is this supported?
	// fmt.Printf("%v\n", c)
	// fmt.Println("LogFormat:", c.LogFormat)

	return nil
}

// Description provides an overview as part of the application Help output
func (c Config) Description() string {

	return fmt.Sprintf("%s %s", c.AppName, c.AppDescription)
}

// Version provides a version string that appears at the top of the
// application Help output
func (c Config) Version() string {

	versionString := fmt.Sprintf("%s %s\n%s",
		strings.ToTitle(c.AppName), c.AppVersion, c.AppURL)

	//divider := strings.Repeat("-", len(versionString))

	// versionBlock := fmt.Sprintf("\n%s\n%s\n%s\n",
	// 	divider, versionString, divider)

	//return versionBlock

	return "\n" + versionString + "\n"
}

// Validate verifies all struct fields have been provided acceptable
func (c Config) Validate() (bool, error) {

	// FilePattern is optional
	// FileExtensions is optional
	//   Discovered files are checked against FileExtensions later

	if len(c.Paths) == 0 {
		return false, fmt.Errorf("one or more paths not provided")
	}

	// RecursiveSearch is optional

	// NumFilesToKeep is optional, but if specified we should make sure it is
	// a non-negative number. AFAIK, this is not currently enforced any other
	// way.
	if c.NumFilesToKeep < 0 {
		return false, fmt.Errorf("negative number not supported for files to keep")
	}

	// We only want to work with positive file modification times 0 is
	// acceptable as it is the default value and indicates that the user has
	// not chosen to use the flag (or has chosen improperly and it will be
	// ignored).
	if c.FileAge < 0 {
		return false, fmt.Errorf("negative number for file age not supported")
	}

	// KeepOldest is optional
	// Remove is optional
	// IgnoreErrors is optional

	// go-flags `choice:""` struct tags enforce valid options
	// if !inList(c.LogFormat, c.validLogFormats) {
	// 	return false
	// }

	switch c.LogFormat {
	case "text":
	case "json":
	default:
		return false, fmt.Errorf("invalid option %q provided for log format", c.LogFormat)
	}

	// LogFilePath is optional
	// TODO: String validation if it is set?

	// Do nothing for valid choices, return false if invalid value specified
	switch c.ConsoleOutput {
	case "stdout":
	case "stderr":
	default:
		return false, fmt.Errorf("invalid option %q provided for console output destination", c.ConsoleOutput)
	}

	// go-flags `choice:""` struct tags enforce valid options
	// if !inList(c.LogLevel, c.validLogLevels) {
	// 	return false
	// }

	switch c.LogLevel {
	case "emergency":
	case "alert":
	case "critical":
	case "panic":
	case "fatal":
	case "error":
	case "warn":
	case "info":
	case "notice":
	case "debug":
	case "trace":
	default:
		return false, fmt.Errorf("invalid option %q provided for log level", c.LogLevel)
	}

	// UseSyslog is optional

	// Optimist
	return true, nil

}

// String() satisfies the Stringer{} interface. This is intended for non-JSON
// formatting if using the TextFormatter logrus formatter.
func (c *Config) String() string {
	return fmt.Sprintf("AppName=%q, AppDescription=%q, AppVersion=%q, FilePattern=%q, FileExtensions=%q, Paths=%v, RecursiveSearch=%t, FileAge=%d, NumFilesToKeep=%d, KeepOldest=%t, Remove=%t, IgnoreErrors=%t, LogFormat=%q, LogFilePath=%q, LogFileHandle=%v, ConsoleOutput=%q, LogLevel=%q, UseSyslog=%t",

		// TODO: Finish syncing this against the config struct fields
		c.AppName,
		c.AppDescription,
		c.AppVersion,
		c.FilePattern,
		c.FileExtensions,
		c.Paths,
		c.RecursiveSearch,
		c.FileAge,
		c.NumFilesToKeep,
		c.KeepOldest,
		c.Remove,
		c.IgnoreErrors,
		c.LogFormat,
		c.LogFilePath,
		c.LogFileHandle,
		c.ConsoleOutput,
		c.LogLevel,
		c.UseSyslog,
	)
}
