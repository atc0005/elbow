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

	"github.com/atc0005/elbow/logging"

	"github.com/alexflint/go-arg"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

var logBuffer logging.LogBuffer

// AppMetadata represents data about this application that may be used in Help
// output, error messages and potentially log messages (e.g., AppVersion)
type AppMetadata struct {
	appName        *string `toml:"app_name" arg:"-"`
	appDescription *string `toml:"app_description" arg:"-"`
	appVersion     *string `toml:"app_version" arg:"-"`
	appURL         *string `toml:"app_url" arg:"-"`
}

// FileHandling represents options specific to how this application
// handles files.
type FileHandling struct {
	filePattern    *string   `toml:"pattern" arg:"--pattern,env:ELBOW_FILE_PATTERN" help:"Substring pattern to compare filenames against. Wildcards are not supported."`
	fileExtensions *[]string `toml:"file_extensions" arg:"--extensions,env:ELBOW_EXTENSIONS" help:"Limit search to specified file extensions. Specify as space separated list to match multiple required extensions."`
	fileAge        *int      `toml:"file_age" arg:"--age,env:ELBOW_FILE_AGE" help:"Limit search to files that are the specified number of days old or older."`
	numFilesToKeep *int      `toml:"files_to_keep" arg:"--keep,env:ELBOW_KEEP" help:"Keep specified number of matching files per provided path."`
	keepOldest     *bool     `toml:"keep_oldest" arg:"--keep-old,env:ELBOW_KEEP_OLD" help:"Keep oldest files instead of newer per provided path."`
	remove         *bool     `toml:"remove" arg:"--remove,env:ELBOW_REMOVE" help:"Remove matched files per provided path."`
	ignoreErrors   *bool     `toml:"ignore_errors" arg:"--ignore-errors,env:ELBOW_IGNORE_ERRORS" help:"Ignore errors encountered during file removal."`
}

// Search represents options specific to controlling how this application
// performs searches in the filesystem
type Search struct {
	paths           *[]string `toml:"paths" arg:"--paths,env:ELBOW_PATHS" help:"List of comma or space-separated paths to process."`
	recursiveSearch *bool     `toml:"recursive_search" arg:"--recurse,env:ELBOW_RECURSE" help:"Perform recursive search into subdirectories per provided path."`
}

// Logging represents options specific to how this application handles
// logging.
type Logging struct {
	logLevel      *string `toml:"log_level" arg:"--log-level,env:ELBOW_LOG_LEVEL" help:"Maximum log level at which messages will be logged. Log messages below this threshold will be discarded."`
	logFormat     *string `toml:"log_format" arg:"--log-format,env:ELBOW_LOG_FORMAT" help:"Log formatter used by logging package."`
	logFilePath   *string `toml:"log_file_path" arg:"--log-file,env:ELBOW_LOG_FILE" help:"Optional log file used to hold logged messages. If set, log messages are not displayed on the console."`
	consoleOutput *string `toml:"console_output" arg:"--console-output,env:ELBOW_CONSOLE_OUTPUT" help:"Specify how log messages are logged to the console."`
	useSyslog     *bool   `toml:"use_syslog" arg:"--use-syslog,env:ELBOW_USE_SYSLOG" help:"Log messages to syslog in addition to other outputs. Not supported on Windows."`
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
	// TODO: Should these be exposed or kept private?
	logFileHandle *os.File       `toml:"-" arg:"-"`
	logger        *logrus.Logger `toml:"-" arg:"-"`
	flagParser    *arg.Parser    `toml:"-" arg:"-"`

	// Path to (optional) configuration file
	configFile *string `toml:"config_file" arg:"--config-file,env:ELBOW_CONFIG_FILE" help:"Full path to optional TOML-formatted configuration file. See config.example.toml for a starter template."`
}

// DefaultConfig returns a configuration object with baseline settings applied
// for further extension by the caller.
func DefaultConfig(appName, appDescription, appURL, appVersion string) Config {

	// Our baseline. The majority of the default settings were previously
	// supplied via struct tags

	var defaultConfig Config

	// Common metadata
	*defaultConfig.appName = appName
	*defaultConfig.appDescription = appDescription
	*defaultConfig.appURL = appURL
	*defaultConfig.appVersion = appVersion

	// Apply default settings that other configuration sources will be allowed
	// to (and for a few settings MUST) override
	*defaultConfig.filePattern = ""
	*defaultConfig.fileAge = 0
	*defaultConfig.numFilesToKeep = -1
	*defaultConfig.keepOldest = false
	*defaultConfig.remove = false
	*defaultConfig.ignoreErrors = false
	*defaultConfig.recursiveSearch = false
	*defaultConfig.logLevel = "info"
	*defaultConfig.logFormat = "text"
	*defaultConfig.logFilePath = ""
	*defaultConfig.consoleOutput = "stdout"
	*defaultConfig.useSyslog = false
	*defaultConfig.configFile = ""

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
	baseConfig.logger = logrus.New()

	// Bundle the returned `*.arg.Parser` for later use from `main()` so that
	// we can explicitly display usage or help details should the
	// user-provided settings fail validation.
	baseConfig.flagParser = arg.MustParse(&argsConfig)

	fmt.Println("Dumping FlagParser for review")
	fmt.Printf("%+v", baseConfig.flagParser)

	/*************************************************************************
		At this point `baseConfig` is our baseline config object containing
		default settings and various handles to other resources. We do not
		apply those same resource handles to other config structs. We merge
		the other configuration objects into the baseConfig object to create
		a unified configuration object that we return to the caller.
	*************************************************************************/

	// If user specified a config file, let's try to use it
	if argsConfig.configFile != nil {
		// Check for a configuration file and load it if found.
		if err := fileConfig.LoadConfigFile(*argsConfig.configFile); err != nil {
			logBuffer.Add(logging.LogRecord{
				Level:   logrus.ErrorLevel,
				Message: fmt.Sprintf("Error loading config file: %s", err),
				Fields:  logrus.Fields{"config_file": argsConfig.configFile},
			})
		}

		logBuffer.Add(logging.LogRecord{
			Level:   logrus.DebugLevel,
			Message: "Processing fileConfig object with MergeConfig func",
		})

		if err := MergeConfig(&baseConfig, fileConfig, defaultConfig); err != nil {
			_, _, line, _ := runtime.Caller(0)
			logBuffer.Add(logging.LogRecord{
				Level:   logrus.ErrorLevel,
				Message: fmt.Sprintf("Error merging config file settings with base config: %s", err),
				Fields:  logrus.Fields{"line": line},
			})
		}

		if ok, err := baseConfig.Validate(); !ok {
			_, _, line, _ := runtime.Caller(0)
			logBuffer.Add(logging.LogRecord{
				Level:   logrus.ErrorLevel,
				Message: fmt.Sprintf("Error validating config after merging %s: %s", "fileConfig", err),
				Fields: logrus.Fields{
					"line":          line,
					"config_object": fmt.Sprintf("%+v", baseConfig),
				},
			})
		}

	}

	logBuffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: "Processing argsConfig object with MergeConfig func",
	})

	if err := MergeConfig(&baseConfig, argsConfig, defaultConfig); err != nil {
		_, _, line, _ := runtime.Caller(0)
		logBuffer.Add(logging.LogRecord{
			Level:   logrus.ErrorLevel,
			Message: fmt.Sprintf("Error merging args config settings with base config: %s", err),
			Fields:  logrus.Fields{"line": line},
		})
	}

	if ok, err := baseConfig.Validate(); !ok {
		_, _, line, _ := runtime.Caller(0)
		logBuffer.Add(logging.LogRecord{
			Level:   logrus.ErrorLevel,
			Message: fmt.Sprintf("Error validating config after merging %s: %s", "argsConfig", err),
			Fields:  logrus.Fields{"line": line},
		})
	}

	// Apply logging configuration
	baseConfig.SetLoggerConfig()

	logBuffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("The config object that we are returning: %+v", baseConfig),
	})

	logBuffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: "Empty queued up log messages from log buffer using user-specified logging settings",
	})
	logBuffer.Flush(baseConfig.logger)

	return &baseConfig

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

// FIXME: The description for this function is now broken
//
// MergeConfig receives source, destination and default Config objects and
// merges select, non-default field values from the source Config object to
// the destination config object, overwriting any field value already present.
//
// `source` and `destination` config structs already have usable default values
// upon creation using our `NewConfig()` constructor; only copy if source struct
// has a different value
//
// TODO: While this makes sense NOW, what is the best way to handle this if
// the default value becomes non-zero?
//
// The goal is to respect the current documented configuration precedence for
// multiple configuration sources (e.g., config file and command-line flags).
func MergeConfig(destination *Config, source Config, defaultConfig Config) error {

	// FIXME: How can we get all field names programatically so we don't have to
	// manually reference each field?

	logBuffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: "MergeConfig called",
	})
	logBuffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Source struct: %+v", source),
	})
	logBuffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Destination struct: %+v", *destination),
	})
	logBuffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Default struct: %+v", defaultConfig),
	})

	if source.paths != nil {
		*destination.paths = *source.paths
	}

	if source.fileExtensions != nil {
		*destination.fileExtensions = *source.fileExtensions
	}

	if source.filePattern != nil {
		*destination.filePattern = *source.filePattern
	}

	if source.fileAge != nil {
		*destination.fileAge = *source.fileAge
	}

	if source.numFilesToKeep != nil {
		*destination.numFilesToKeep = *source.numFilesToKeep
	}

	if source.keepOldest != nil {
		*destination.keepOldest = *source.keepOldest
	}

	if source.remove != nil {
		*destination.remove = *source.remove
	}

	if source.ignoreErrors != nil {
		*destination.ignoreErrors = *source.ignoreErrors
	}

	if source.recursiveSearch != nil {
		*destination.recursiveSearch = *source.recursiveSearch
	}

	if source.logLevel != nil {
		*destination.logLevel = *source.logLevel
	}

	if source.logFormat != nil {
		*destination.logFormat = *source.logFormat
	}

	if source.logFilePath != nil {
		*destination.logFilePath = *source.logFilePath
	}

	if source.consoleOutput != nil {
		*destination.consoleOutput = *source.consoleOutput
	}

	if source.useSyslog != nil {
		*destination.useSyslog = *source.useSyslog
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

	return fmt.Sprintf("%s %s", *c.appName, *c.appDescription)
}

// Version provides a version string that appears at the top of the
// application Help output
func (c Config) Version() string {

	versionString := fmt.Sprintf("%s %s\n%s",
		strings.ToTitle(*c.appName), *c.appVersion, *c.appURL)

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

	if len(*c.paths) == 0 {
		return false, fmt.Errorf("one or more paths not provided")
	}

	// recursiveSearch is optional

	// NumFilesToKeep is optional, but if specified we should make sure it is
	// a non-negative number. AFAIK, this is not currently enforced any other
	// way.
	if *c.numFilesToKeep < 0 {
		return false, fmt.Errorf("invalid value provided for files to keep")
	}

	// We only want to work with positive file modification times 0 is
	// acceptable as it is the default value and indicates that the user has
	// not chosen to use the flag (or has chosen improperly and it will be
	// ignored).
	if *c.fileAge < 0 {
		return false, fmt.Errorf("negative number for file age not supported")
	}

	// keepOldest is optional
	// Remove is optional
	// ignoreErrors is optional

	switch *c.logFormat {
	case "text":
	case "json":
	default:
		return false, fmt.Errorf("invalid option %q provided for log format", *c.logFormat)
	}

	// logFilePath is optional
	// TODO: String validation if it is set?

	// Do nothing for valid choices, return false if invalid value specified
	switch *c.consoleOutput {
	case "stdout":
	case "stderr":
	default:
		return false, fmt.Errorf("invalid option %q provided for console output destination", *c.consoleOutput)
	}

	switch *c.logLevel {
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
		return false, fmt.Errorf("invalid option %q provided for log level", *c.logLevel)
	}

	// useSyslog is optional

	// Optimist
	return true, nil

}

// String() satisfies the Stringer{} interface. This is intended for non-JSON
// formatting if using the TextFormatter logrus formatter.
func (c *Config) String() string {
	return fmt.Sprintf("AppName=%q, AppDescription=%q, AppVersion=%q, AppURL=%q, FilePattern=%q, FileExtensions=%q, Paths=%v, recursiveSearch=%t, fileAge=%d, NumFilesToKeep=%d, keepOldest=%t, Remove=%t, ignoreErrors=%t, logFormat=%q, logFilePath=%q, LogFileHandle=%v, consoleOutput=%q, logLevel=%q, useSyslog=%t, Logger=%v, FlagParser=%v), ConfigFile=%q, EndOfStringMethod",

		*c.appName,
		*c.appDescription,
		*c.appVersion,
		*c.appURL,
		*c.filePattern,
		*c.fileExtensions,
		*c.paths,
		*c.recursiveSearch,
		*c.fileAge,
		*c.numFilesToKeep,
		*c.keepOldest,
		*c.remove,
		*c.ignoreErrors,
		*c.logFormat,
		*c.logFilePath,
		*c.logFileHandle,
		*c.consoleOutput,
		*c.logLevel,
		*c.useSyslog,
		c.logger,
		*c.flagParser,
		*c.configFile,
	)
}

// SetLoggerConfig applies chosen configuration settings that control logging
// output.
func (c *Config) SetLoggerConfig() {

	logging.SetLoggerFormatter(c.logger, *c.logFormat)
	logging.SetLoggerConsoleOutput(c.logger, *c.consoleOutput)

	if fileHandle, err := logging.SetLoggerLogFile(c.logger, *c.logFilePath); err == nil {
		c.logFileHandle = fileHandle
	} else {
		// Need to collect the error for display later
		logBuffer.Add(logging.LogRecord{
			Level:   logrus.ErrorLevel,
			Message: fmt.Sprintf("%s", err),
			Fields:  logrus.Fields{"log_file": c.logFilePath},
		})
	}

	logging.SetLoggerLevel(c.logger, *c.logLevel)

	// https://godoc.org/github.com/sirupsen/logrus#New
	// https://godoc.org/github.com/sirupsen/logrus#Logger

	// make sure that the user actually requested syslog logging as it is
	// currently supported on UNIX only.
	if *c.useSyslog {
		logBuffer.Add(logging.LogRecord{
			Level:   logrus.InfoLevel,
			Message: "Syslog logging requested, attempting to enable it",
			Fields:  logrus.Fields{"use_syslog": c.useSyslog},
		})

		if err := logging.EnableSyslogLogging(c.logger, &logBuffer, *c.logLevel); err != nil {
			// TODO: Is this sufficient cause for failing? Perhaps if a local
			// log file is not also set consider it a failure?

			logBuffer.Add(logging.LogRecord{
				Level:   logrus.ErrorLevel,
				Message: fmt.Sprintf("Failed to enable syslog logging: %s", err),
				Fields:  logrus.Fields{"use_syslog": c.useSyslog},
			})

			logBuffer.Add(logging.LogRecord{
				Level:   logrus.WarnLevel,
				Message: "Proceeding without syslog logging",
				Fields:  logrus.Fields{"use_syslog": c.useSyslog},
			})
		}
	} else {
		logBuffer.Add(logging.LogRecord{
			Level:   logrus.DebugLevel,
			Message: "Syslog logging not requested, not enabling",
			Fields:  logrus.Fields{"use_syslog": c.useSyslog},
		})

	}

}

// AppName returns the appName field if it's non-nil, zero value otherwise.
func (c *Config) AppName() string {
	if c == nil || c.appName == nil {
		return ""
	}
	return *c.appName
}

// AppDescription returns the appDescription field if it's non-nil, zero value otherwise.
func (c *Config) AppDescription() string {
	if c == nil || c.appDescription == nil {
		return ""
	}
	return *c.appDescription

}

// AppVersion returns the appVersion field if it's non-nil, zero value otherwise.
func (c *Config) AppVersion() string {
	if c == nil || c.appVersion == nil {
		return ""
	}
	return *c.appVersion
}

// AppURL returns the appURL field if it's non-nil, zero value otherwise.
func (c *Config) AppURL() string {
	if c == nil || c.appURL == nil {
		return ""
	}
	return *c.appURL
}

// FilePattern returns the filePattern field if it's non-nil, zero value otherwise.
func (c *Config) FilePattern() string {
	if c == nil || c.filePattern == nil {
		return ""
	}
	return *c.filePattern
}

// FileExtensions returns the fileExtensions field if it's non-nil, zero value
// otherwise.
// TODO: Double check this one; how should we safely handle returning an
// empty/zero value?
// As an example, the https://github.com/google/go-github package has a
// `Issue.GetAssignees()` method that returns nil if the `Issue.Assignees`
// field is nil. This seems to suggest that this is all we really can do here?
//
func (c *Config) FileExtensions() []string {
	if c == nil || c.fileExtensions == nil {
		// FIXME: Isn't the goal to avoid returning nil?
		return nil
	}
	return *c.fileExtensions
}

// FileAge returns the fileAge field if it's non-nil, zero value otherwise.
func (c *Config) FileAge() int {
	if c == nil || c.fileAge == nil {
		return 0
	}
	return *c.fileAge
}

// NumFilesToKeep returns the numFilesToKeep field if it's non-nil, zero value
// otherwise.
func (c *Config) NumFilesToKeep() int {
	if c == nil || c.numFilesToKeep == nil {
		return 0
	}
	return *c.numFilesToKeep
}

// KeepOldest returns the keepOldest field if it's non-nil, zero value
// otherwise.
func (c *Config) KeepOldest() bool {
	if c == nil || c.keepOldest == nil {
		return false
	}
	return *c.keepOldest
}

// Remove returns the remove field if it's non-nil, zero value otherwise.
func (c *Config) Remove() bool {
	if c == nil || c.remove == nil {
		return false
	}
	return *c.remove
}

// IgnoreErrors returns the ignoreErrors field if it's non-nil, zero value
// otherwise.
func (c *Config) IgnoreErrors() bool {
	if c == nil || c.ignoreErrors == nil {
		return false
	}
	return *c.ignoreErrors
}

// Paths returns the paths field if it's non-nil, zero value otherwise.
func (c *Config) Paths() []string {
	if c == nil || c.paths == nil {
		return nil
	}
	return *c.paths
}

// RecursiveSearch returns the recursiveSearch field if it's non-nil, zero
// value otherwise.
func (c *Config) RecursiveSearch() bool {
	if c == nil || c.recursiveSearch == nil {
		return false
	}
	return *c.recursiveSearch
}

// LogLevel returns the logLevel field if it's non-nil, zero value otherwise.
func (c *Config) LogLevel() string {
	if c == nil || c.logLevel == nil {
		return ""
	}
	return *c.logLevel
}

// LogFormat returns the logFormat field if it's non-nil, zero value otherwise.
func (c *Config) LogFormat() string {
	if c == nil || c.logFormat == nil {
		return ""
	}
	return *c.logFormat
}

// LogFilePath returns the logFilePath field if it's non-nil, zero value
// otherwise.
func (c *Config) LogFilePath() string {
	if c == nil || c.logFilePath == nil {
		return ""
	}
	return *c.logFilePath
}

// ConsoleOutput returns the consoleOutput field if it's non-nil, zero value
// otherwise.
func (c *Config) ConsoleOutput() string {
	if c == nil || c.consoleOutput == nil {
		return ""
	}
	return *c.consoleOutput
}

// UseSyslog returns the useSyslog field if it's non-nil, zero
// value otherwise.
func (c *Config) UseSyslog() bool {
	if c == nil || c.useSyslog == nil {
		return false
	}
	return *c.useSyslog
}

// ConfigFile returns the configFile field if it's non-nil, zero value
// otherwise.
func (c *Config) ConfigFile() string {
	if c == nil || c.configFile == nil {
		return ""
	}
	return *c.configFile
}
