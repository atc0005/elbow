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

// Package global intended to help "collect" log messages during configuration
// initialization in order to properly handle after configuration object is
// finalized.
var initLogger *logrus.Logger = logrus.New()

// AppMetadata represents data about this application that may be used in Help
// output, error messages and potentially log messages (e.g., AppVersion)
type AppMetadata struct {
	AppName        string `toml:"-" arg:"-"`
	AppDescription string `toml:"-" arg:"-"`
	AppVersion     string `toml:"-" arg:"-"`
	AppURL         string `toml:"-" arg:"-"`
}

// FileHandling represents options specific to how this application
// handles files.
type FileHandling struct {
	FilePattern    string   `toml:"pattern" arg:"--pattern,env:ELBOW_FILE_PATTERN" help:"Substring pattern to compare filenames against. Wildcards are not supported."`
	FileExtensions []string `toml:"file_extensions" arg:"--extensions,env:ELBOW_EXTENSIONS" help:"Limit search to specified file extensions. Specify as space separated list to match multiple required extensions."`
	FileAge        int      `toml:"file_age" arg:"--age,env:ELBOW_FILE_AGE" help:"Limit search to files that are the specified number of days old or older."`
	NumFilesToKeep int      `toml:"files_to_keep" arg:"--keep,env:ELBOW_KEEP" help:"Keep specified number of matching files per provided path."`
	KeepOldest     bool     `toml:"keep_oldest" arg:"--keep-old,env:ELBOW_KEEP_OLD" help:"Keep oldest files instead of newer per provided path."`
	Remove         bool     `toml:"remove" arg:"--remove,env:ELBOW_REMOVE" help:"Remove matched files per provided path."`
	IgnoreErrors   bool     `toml:"ignore_errors" arg:"--ignore-errors,env:ELBOW_IGNORE_ERRORS" help:"Ignore errors encountered during file removal."`
}

// Search represents options specific to controlling how this application
// performs searches in the filesystem
type Search struct {
	Paths           []string `toml:"paths" arg:"--paths,env:ELBOW_PATHS" help:"List of comma or space-separated paths to process."`
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
	LogFileHandle *os.File       `toml:"-" arg:"-"`
	Logger        *logrus.Logger `toml:"-" arg:"-"`
	FlagParser    *arg.Parser    `toml:"-" arg:"-"`

	// Path to (optional) configuration file
	ConfigFile string `toml:"config_file" arg:"--config-file,env:ELBOW_CONFIG_FILE" help:"Full path to optional TOML-formatted configuration file. See config.example.toml for a starter template."`
}

// DefaultConfig returns a configuration object with baseline settings applied
// for further extension by the caller.
func DefaultConfig(appName, appDescription, appURL, appVersion string) Config {

	// Our baseline. The majority of the default settings were previously
	// supplied via struct tags

	var defaultConfig Config

	// Common metadata
	defaultConfig.AppName = appName
	defaultConfig.AppDescription = appDescription
	defaultConfig.AppURL = appURL
	defaultConfig.AppVersion = appVersion

	// Apply default settings that other configuration sources will be allowed
	// to (and for a few settings MUST) override
	defaultConfig.FilePattern = ""
	defaultConfig.FileAge = 0
	defaultConfig.NumFilesToKeep = -1
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

	return defaultConfig

}

// NewConfig returns a pointer to a newly configured object representing a
// collection of user-provided and default settings.
func NewConfig(appName, appDescription, appURL, appVersion string) *Config {

	// Baseline collection of settings before loading custom config sources
	defaultConfig := DefaultConfig(appName, appDescription, appURL, appVersion)

	// The base configuration object that will be returned to the caller
	baseConfig := DefaultConfig(appName, appDescription, appURL, appVersion)

	// Settings provided via config file
	fileConfig := DefaultConfig(appName, appDescription, appURL, appVersion)

	// Settings provided via command-line flags and environment variables.
	// This object will always be set in some manner as either flags or env
	// vars will be needed to bootstrap the application. While we may support
	// using a configuration file to provide settings, it is not used by
	// default.
	argsConfig := DefaultConfig(appName, appDescription, appURL, appVersion)

	// Initialize logger "handle" for later use
	baseConfig.Logger = logrus.New()

	// Bundle the returned `*.arg.Parser` for later use from `main()` so that
	// we can explicitly display usage or help details should the
	// user-provided settings fail validation.
	baseConfig.FlagParser = arg.MustParse(&argsConfig)

	/*************************************************************************
		At this point `config` is our base config object containing default
		settings and various handles to other resources. We do not apply those
		same resource handles to other config structs. We merge the other
		configuration objects into it the base config object to create a
		unified configuration object that we return to the caller.
	*************************************************************************/

	// If user specified a config file, let's try to use it
	if argsConfig.ConfigFile != "" {
		// Check for a configuration file and load it if found.
		if err := fileConfig.LoadConfigFile(argsConfig.ConfigFile); err != nil {
			fmt.Printf("Error loading config file: %s\n", err)
		}
	}

	// Some number of commits back the following error messages were returned
	// when attempting to use logrus package to log the contents of the
	// struct. Remove this comment and the following error messages once this
	// code proves stable.
	//
	// Failed to obtain reader, failed to marshal fields to JSON, json: unsupported type: func([]string)
	// Failed to obtain reader, failed to marshal fields to JSON, json: unsupported type: func(*runtime.Frame) (string, string)
	fmt.Printf("\n\nProcessing fileConfig object with MergeConfig func\n")
	if err := MergeConfig(&baseConfig, fileConfig, defaultConfig); err != nil {
		_, _, line, _ := runtime.Caller(0)
		fmt.Printf("(line %d) Error merging config file settings with base config: %s\n", line, err)
	}

	if ok, err := baseConfig.Validate(); !ok {
		_, _, line, _ := runtime.Caller(0)
		fmt.Printf("(line %d) Error validating config after merging %s: %s\n",
			line, "fileConfig", err)
	}

	fmt.Printf("\n\nProcessing argsConfig object with MergeConfig func\n")
	if err := MergeConfig(&baseConfig, argsConfig, defaultConfig); err != nil {
		_, _, line, _ := runtime.Caller(0)
		fmt.Printf("(line %d) Error merging args config settings with base config: %s\n", line, err)
	}

	if ok, err := baseConfig.Validate(); !ok {
		_, _, line, _ := runtime.Caller(0)
		fmt.Printf("(line %d) Error validating config after merging %s: %s\n",
			line, "argsConfig", err)
	}

	fmt.Println("The config object that we are returning:", baseConfig)

	return &config

}

// GetStructTag returns the requested struct tag value, if set, and an error
// value indicating whether any problems were encountered.
func GetStructTag(c Config, fieldname string, tagName string) (string, bool) {

	t := reflect.TypeOf(c)

	var field reflect.StructField
	var ok bool
	var tagValue string

	if field, ok = t.FieldByName(fieldname); !ok {
		return "", false
	}

	if tagValue, ok = field.Tag.Lookup(tagName); !ok {
		return "", false
	}

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
		return false, fmt.Errorf("invalid value provided for files to keep")
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
	return fmt.Sprintf("AppName=%q, AppDescription=%q, AppVersion=%q, AppURL=%q, FilePattern=%q, FileExtensions=%q, Paths=%v, RecursiveSearch=%t, FileAge=%d, NumFilesToKeep=%d, KeepOldest=%t, Remove=%t, IgnoreErrors=%t, LogFormat=%q, LogFilePath=%q, LogFileHandle=%v, ConsoleOutput=%q, LogLevel=%q, UseSyslog=%t, Logger=%v, FlagParser=%v), ConfigFile=%q, EndOfStringMethod",

		c.AppName,
		c.AppDescription,
		c.AppVersion,
		c.AppURL,
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
		c.Logger,
		c.FlagParser,
		c.ConfigFile,
	)
}
