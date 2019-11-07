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
	"runtime"
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
	FilePattern    string   `toml:"pattern" arg:"--pattern,env:ELBOW_FILE_PATTERN" help:"Substring pattern to compare filenames against. Wildcards are not supported."`
	FileExtensions []string `toml:"file_extensions" arg:"--extensions,env:ELBOW_EXTENSIONS" help:"Limit search to specified file extensions. Specify as space separated list to match multiple required extensions."`
	FileAge        int      `toml:"file_age" arg:"--age,env:ELBOW_FILE_AGE" help:"Limit search to files that are the specified number of days old or older."`
	NumFilesToKeep int      `toml:"files_to_keep" arg:"--keep,required,env:ELBOW_KEEP" help:"Keep specified number of matching files per provided path."`
	KeepOldest     bool     `toml:"keep_oldest" arg:"--keep-old,env:ELBOW_KEEP_OLD" help:"Keep oldest files instead of newer per provided path."`
	Remove         bool     `toml:"remove" arg:"--remove,env:ELBOW_REMOVE" help:"Remove matched files per provided path."`
	IgnoreErrors   bool     `toml:"ignore_errors" arg:"--ignore-errors,env:ELBOW_IGNORE_ERRORS" help:"Ignore errors encountered during file removal."`
}

// Search represents options specific to controlling how this application
// performs searches in the filesystem
type Search struct {
	Paths           []string `toml:"paths" arg:"--paths,required,env:ELBOW_PATHS" help:"List of comma or space-separated paths to process."`
	RecursiveSearch bool     `toml:"recursive_search" arg:"--recurse,env:ELBOW_RECURSE" help:"Perform recursive search into subdirectories per provided path."`
}

// Logging represents options specific to how this application handles
// logging.
type Logging struct {
	LogLevel      string `toml:"log_level" arg:"--log-level,env:ELBOW_LOG_LEVEL" help:"Maximum log level at which messages will be logged. Log messages below this threshold will be discarded."`
	LogFormat     string `toml:"log_format" arg:"--log-format,env:ELBOW_LOG_FORMAT" help:"Log formatter used by logging package."`
	LogFilePath   string `toml:"log_file_path" arg:"--log-file,env:ELBOW_LOG_FILE" help:"Optional log file used to hold logged messages. If set, log messages are not displayed on the console."`
	ConsoleOutput string `toml:"console_output" arg:"--console-output,env:ELBOW_CONSOLE_OUTPUT" help:"Specify how log messages are logged to the console."`
	UseSyslog     bool   `toml:"use_syslog" arg:"--use-syslog,env:ELBOW_USE_SYSLOG" help:"Log messages to syslog in addition to other outputs. Not supported on Windows."`
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
	ConfigFile string `toml:"config_file" arg:"--config-file,env:ELBOW_CONFIG_FILE" help:"Full path to optional TOML-formatted configuration file. See config.example.toml for a starter template."`
}

// NewConfig returns a newly configured object representing a collection of
// user-provided and default settings.
func NewConfig(appName, appDescription, appURL, appVersion string) *Config {

	// Note: The majority of the default settings are supplied via struct tags

	// The base configuration object that will be returned to the caller
	var config Config

	// Our baseline
	var defaultConfig Config

	// Settings provided via command-line flags and environment variables.
	// This object will always be set in some manner as either flags or env
	// vars will be needed to bootstrap the application. While we may support
	// using a configuration file to provide settings, it is not used by
	// default.
	var argsConfig Config

	// Settings provided via config file
	var fileConfig Config

	// Common metadata
	defaultConfig.AppName = appName
	defaultConfig.AppDescription = appDescription
	defaultConfig.AppURL = appURL
	defaultConfig.AppVersion = appVersion

	// Apply default settings that other configuration sources will be allowed
	// to override
	defaultConfig.FilePattern = ""
	defaultConfig.FileAge = 0
	defaultConfig.NumFilesToKeep = 0
	defaultConfig.KeepOldest = false
	defaultConfig.Remove = false
	defaultConfig.IgnoreErrors = false
	defaultConfig.RecursiveSearch = false
	defaultConfig.LogLevel = "info"
	defaultConfig.LogFormat = "text"
	defaultConfig.LogFilePath = ""
	defaultConfig.ConsoleOutput = "stdout"
	defaultConfig.UseSyslog = false
	defaultConfig.ConfigFile = ""

	// Apply baseline before loading custom settings
	config = defaultConfig
	fileConfig = defaultConfig
	argsConfig = defaultConfig

	// Initialize logger "handle" for later use
	config.Logger = logrus.New()

	// Bundle the returned `*.arg.Parser` for later use from `main()` so that
	// we can explicitly display usage or help details should the
	// user-provided settings fail validation.
	config.FlagParser = arg.MustParse(&argsConfig)

	// If user specified a config file, let's try to use it
	if argsConfig.ConfigFile != "" {
		// Check for a configuration file and load it if found.
		// FIXME: The method needs to be updated to reference the path provided
		// via environment variable or command-line flag
		if err := fileConfig.LoadConfigFile(argsConfig.ConfigFile); err != nil {
			fmt.Printf("Error loading config file: %s", err)
		}
	}

	// At this point `config` is our base config. We merge the other
	// configuration objects into it to create a unified configuration object
	// that we return to the caller.
	if err := MergeConfig(&config, fileConfig, defaultConfig); err != nil {
		fmt.Printf("Error merging config file settings with base config: %s", err)
	}

	if err := MergeConfig(&config, argsConfig, defaultConfig); err != nil {
		fmt.Printf("Error merging args config settings with base config: %s", err)
	}

	_, _, line, _ := runtime.Caller(0)
	fmt.Printf("Line %d\n", line)
	fmt.Println("The config object that we are returning:", config.String())

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

// MergeConfig receives source, destination and default Config objects and
// merges non-default field values from the source Config object to the
// destination config object, overwriting any field values already present.
// The goal is to respect the current documented configuration precedence for
// multiple configuration sources (e.g., config file and command-line flags).
func MergeConfig(destination *Config, source Config, defaultConfig Config) error {

	// FIXME: How can we get all field names programatically so we don't have to
	// manually reference each field?

	fmt.Println("MergeConfig called")
	fmt.Printf("Source struct: %+v\n", source)
	fmt.Printf("Dest struct: %+v\n", *destination)
	fmt.Printf("Default struct: %+v\n", defaultConfig)

	// Copy over select source struct field values if destination struct field
	// values are empty or some other invalid state. These fields are not
	// supported by `default` value logic.
	if len(source.Paths) > len(defaultConfig.Paths) {
		destination.Paths = source.Paths
	}

	if len(source.FileExtensions) > len(defaultConfig.FileExtensions) {
		destination.FileExtensions = source.FileExtensions
	}

	if source.FilePattern != defaultConfig.FilePattern {
		destination.FilePattern = source.FilePattern
	}

	// source and destination config structs already have usable default
	// values upon creation using our NewConfig() constructor; only copy if
	// source struct has a different value
	if source.FileAge > defaultConfig.FileAge {
		destination.FileAge = source.FileAge
	}

	if source.NumFilesToKeep > defaultConfig.NumFilesToKeep {
		destination.NumFilesToKeep = source.NumFilesToKeep
	}

	// TODO: any reason to check this? Perhaps just direct copy for boolean
	// variables?
	if source.KeepOldest != defaultConfig.KeepOldest {
		destination.KeepOldest = source.KeepOldest
	}

	if source.Remove != defaultConfig.Remove {
		destination.Remove = source.Remove
	}

	if source.IgnoreErrors != defaultConfig.IgnoreErrors {
		destination.IgnoreErrors = source.IgnoreErrors
	}

	if source.RecursiveSearch != defaultConfig.RecursiveSearch {
		destination.RecursiveSearch = source.RecursiveSearch
	}

	// only copy source field value if non-default
	if source.LogLevel != defaultConfig.LogLevel {
		destination.LogLevel = source.LogLevel
	}

	if source.LogFormat != defaultConfig.LogFormat {
		destination.LogFormat = source.LogFormat
	}

	if source.LogFilePath != defaultConfig.LogFilePath {
		destination.LogFilePath = source.LogFilePath
	}

	if source.ConsoleOutput != defaultConfig.ConsoleOutput {
		destination.ConsoleOutput = source.ConsoleOutput
	}

	if source.UseSyslog != defaultConfig.UseSyslog {
		destination.UseSyslog = source.UseSyslog
	}

	// FIXME: Placeholder
	// FIXME: What useful error code would we return from this function?
	return nil
}

// LoadConfigFile reads and unmarshals a configuration file in TOML format
func (c *Config) LoadConfigFile(filename string) error {

	// Read file to byte slice
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(configFile, c); err != nil {
		return err
	}

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
