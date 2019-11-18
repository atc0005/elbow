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
	AppName        *string `toml:"app_name" arg:"-"`
	AppDescription *string `toml:"app_description" arg:"-"`
	AppVersion     *string `toml:"app_version" arg:"-"`
	AppURL         *string `toml:"app_url" arg:"-"`
}

// FileHandling represents options specific to how this application
// handles files.
type FileHandling struct {
	FilePattern    *string   `toml:"pattern" arg:"--pattern,env:ELBOW_FILE_PATTERN" help:"Substring pattern to compare filenames against. Wildcards are not supported."`
	FileExtensions *[]string `toml:"file_extensions" arg:"--extensions,env:ELBOW_EXTENSIONS" help:"Limit search to specified file extensions. Specify as space separated list to match multiple required extensions."`
	FileAge        *int      `toml:"file_age" arg:"--age,env:ELBOW_FILE_AGE" help:"Limit search to files that are the specified number of days old or older."`
	NumFilesToKeep *int      `toml:"files_to_keep" arg:"--keep,env:ELBOW_KEEP" help:"Keep specified number of matching files per provided path."`
	KeepOldest     *bool     `toml:"keep_oldest" arg:"--keep-old,env:ELBOW_KEEP_OLD" help:"Keep oldest files instead of newer per provided path."`
	Remove         *bool     `toml:"remove" arg:"--remove,env:ELBOW_REMOVE" help:"Remove matched files per provided path."`
	IgnoreErrors   *bool     `toml:"ignore_errors" arg:"--ignore-errors,env:ELBOW_IGNORE_ERRORS" help:"Ignore errors encountered during file removal."`
}

// Search represents options specific to controlling how this application
// performs searches in the filesystem
type Search struct {
	Paths           *[]string `toml:"paths" arg:"--paths,env:ELBOW_PATHS" help:"List of comma or space-separated paths to process."`
	RecursiveSearch *bool     `toml:"recursive_search" arg:"--recurse,env:ELBOW_RECURSE" help:"Perform recursive search into subdirectories per provided path."`
}

// Logging represents options specific to how this application handles
// logging.
type Logging struct {
	LogLevel      *string `toml:"log_level" arg:"--log-level,env:ELBOW_LOG_LEVEL" help:"Maximum log level at which messages will be logged. Log messages below this threshold will be discarded."`
	LogFormat     *string `toml:"log_format" arg:"--log-format,env:ELBOW_LOG_FORMAT" help:"Log formatter used by logging package."`
	LogFilePath   *string `toml:"log_file_path" arg:"--log-file,env:ELBOW_LOG_FILE" help:"Optional log file used to hold logged messages. If set, log messages are not displayed on the console."`
	ConsoleOutput *string `toml:"console_output" arg:"--console-output,env:ELBOW_CONSOLE_OUTPUT" help:"Specify how log messages are logged to the console."`
	UseSyslog     *bool   `toml:"use_syslog" arg:"--use-syslog,env:ELBOW_USE_SYSLOG" help:"Log messages to syslog in addition to other outputs. Not supported on Windows."`
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
	ConfigFile *string `toml:"config_file" arg:"--config-file,env:ELBOW_CONFIG_FILE" help:"Full path to optional TOML-formatted configuration file. See config.example.toml for a starter template."`
}

// DefaultConfig returns a configuration object with baseline settings applied
// for further extension by the caller.
func DefaultConfig(appName, appDescription, appURL, appVersion string) Config {

	// Our baseline. The majority of the default settings were previously
	// supplied via struct tags

	var defaultConfig Config

	// Common metadata
	*defaultConfig.AppName = appName
	*defaultConfig.AppDescription = appDescription
	*defaultConfig.AppURL = appURL
	*defaultConfig.AppVersion = appVersion

	// Apply default settings that other configuration sources will be allowed
	// to (and for a few settings MUST) override
	*defaultConfig.FilePattern = ""
	*defaultConfig.FileAge = 0
	*defaultConfig.NumFilesToKeep = -1
	*defaultConfig.KeepOldest = false
	*defaultConfig.Remove = false
	*defaultConfig.IgnoreErrors = false
	*defaultConfig.RecursiveSearch = false
	*defaultConfig.LogLevel = "info"
	*defaultConfig.LogFormat = "text"
	*defaultConfig.LogFilePath = ""
	*defaultConfig.ConsoleOutput = "stdout"
	*defaultConfig.UseSyslog = false
	*defaultConfig.ConfigFile = ""

	return defaultConfig

}

// NewConfig returns a pointer to a newly configured object representing a
// collection of user-provided and default settings.
func NewConfig(appName, appDescription, appURL, appVersion string) *Config {

	// Baseline collection of settings before loading custom config sources
	defaultConfig := DefaultConfig(appName, appDescription, appURL, appVersion)

	// The base configuration object that will be returned to the caller
	baseConfig := DefaultConfig(appName, appDescription, appURL, appVersion)

	// Settings provided via config file. Intentionally using uninitialized
	// struct here so that we can check for nil pointers to indicate whether
	// a field has been populated with configuration values.
	fileConfig := Config{}

	// Settings provided via command-line flags and environment variables.
	// This object will always be set in some manner as either flags or env
	// vars will be needed to bootstrap the application. While we may support
	// using a configuration file to provide settings, it is not used by
	// default.
	argsConfig := Config{}

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
	if argsConfig.ConfigFile != nil {
		// Check for a configuration file and load it if found.
		if err := fileConfig.LoadConfigFile(*argsConfig.ConfigFile); err != nil {
			logBuffer.Add(logging.LogRecord{
				Level:   logrus.ErrorLevel,
				Message: fmt.Sprintf("Error loading config file: %s", err),
				Fields:  logrus.Fields{"config_file": argsConfig.ConfigFile},
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

// MergeConfig receives source, destination and default Config objects and
// merges select, non-nil field values from the source Config object to
// the destination config object, overwriting any field value already present.
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

	if source.Paths != nil {
		*destination.Paths = *source.Paths
	}

	if source.FileExtensions != nil {
		*destination.FileExtensions = *source.FileExtensions
	}

	if source.FilePattern != nil {
		*destination.FilePattern = *source.FilePattern
	}

	if source.FileAge != nil {
		*destination.FileAge = *source.FileAge
	}

	if source.NumFilesToKeep != nil {
		*destination.NumFilesToKeep = *source.NumFilesToKeep
	}

	if source.KeepOldest != nil {
		*destination.KeepOldest = *source.KeepOldest
	}

	if source.Remove != nil {
		*destination.Remove = *source.Remove
	}

	if source.IgnoreErrors != nil {
		*destination.IgnoreErrors = *source.IgnoreErrors
	}

	if source.RecursiveSearch != nil {
		*destination.RecursiveSearch = *source.RecursiveSearch
	}

	if source.LogLevel != nil {
		*destination.LogLevel = *source.LogLevel
	}

	if source.LogFormat != nil {
		*destination.LogFormat = *source.LogFormat
	}

	if source.LogFilePath != nil {
		*destination.LogFilePath = *source.LogFilePath
	}

	if source.ConsoleOutput != nil {
		*destination.ConsoleOutput = *source.ConsoleOutput
	}

	if source.UseSyslog != nil {
		*destination.UseSyslog = *source.UseSyslog
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

	return fmt.Sprintf("%s %s", *c.AppName, *c.AppDescription)
}

// Version provides a version string that appears at the top of the
// application Help output
func (c Config) Version() string {

	versionString := fmt.Sprintf("%s %s\n%s",
		strings.ToTitle(*c.AppName), *c.AppVersion, *c.AppURL)

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

	if len(*c.Paths) == 0 {
		return false, fmt.Errorf("one or more paths not provided")
	}

	// recursiveSearch is optional

	// NumFilesToKeep is optional, but if specified we should make sure it is
	// a non-negative number. AFAIK, this is not currently enforced any other
	// way.
	if *c.NumFilesToKeep < 0 {
		return false, fmt.Errorf("invalid value provided for files to keep")
	}

	// We only want to work with positive file modification times 0 is
	// acceptable as it is the default value and indicates that the user has
	// not chosen to use the flag (or has chosen improperly and it will be
	// ignored).
	if *c.FileAge < 0 {
		return false, fmt.Errorf("negative number for file age not supported")
	}

	// keepOldest is optional
	// Remove is optional
	// ignoreErrors is optional

	switch *c.LogFormat {
	case "text":
	case "json":
	default:
		return false, fmt.Errorf("invalid option %q provided for log format", *c.LogFormat)
	}

	// logFilePath is optional
	// TODO: String validation if it is set?

	// Do nothing for valid choices, return false if invalid value specified
	switch *c.ConsoleOutput {
	case "stdout":
	case "stderr":
	default:
		return false, fmt.Errorf("invalid option %q provided for console output destination", *c.ConsoleOutput)
	}

	switch *c.LogLevel {
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
		return false, fmt.Errorf("invalid option %q provided for log level", *c.LogLevel)
	}

	// useSyslog is optional

	// Optimist
	return true, nil

}

// String() satisfies the Stringer{} interface. This is intended for non-JSON
// formatting if using the TextFormatter logrus formatter.
func (c *Config) String() string {
	return fmt.Sprintf("AppName=%q, AppDescription=%q, AppVersion=%q, AppURL=%q, FilePattern=%q, FileExtensions=%q, Paths=%v, RecursiveSearch=%t, FileAge=%d, NumFilesToKeep=%d, KeepOldest=%t, Remove=%t, IgnoreErrors=%t, LogFormat=%q, LogFilePath=%q, ConfigFile=%q, ConsoleOutput=%q, LogLevel=%q, UseSyslog=%t, logger=%v, flagParser=%v,  logFileHandle=%v",

		*c.AppName,
		*c.AppDescription,
		*c.AppVersion,
		*c.AppURL,
		*c.FilePattern,
		*c.FileExtensions,
		*c.Paths,
		*c.RecursiveSearch,
		*c.FileAge,
		*c.NumFilesToKeep,
		*c.KeepOldest,
		*c.Remove,
		*c.IgnoreErrors,
		*c.LogFormat,
		*c.LogFilePath,
		*c.ConfigFile,
		*c.ConsoleOutput,
		*c.LogLevel,
		*c.UseSyslog,
		c.logger,
		c.flagParser,
		c.logFileHandle,
	)
}

// SetLoggerConfig applies chosen configuration settings that control logging
// output.
func (c *Config) SetLoggerConfig() {

	logging.SetLoggerFormatter(c.logger, *c.LogFormat)
	logging.SetLoggerConsoleOutput(c.logger, *c.ConsoleOutput)

	if fileHandle, err := logging.SetLoggerLogFile(c.logger, *c.LogFilePath); err == nil {
		c.logFileHandle = fileHandle
	} else {
		// Need to collect the error for display later
		logBuffer.Add(logging.LogRecord{
			Level:   logrus.ErrorLevel,
			Message: fmt.Sprintf("%s", err),
			Fields:  logrus.Fields{"log_file": c.LogFilePath},
		})
	}

	logging.SetLoggerLevel(c.logger, *c.LogLevel)

	// https://godoc.org/github.com/sirupsen/logrus#New
	// https://godoc.org/github.com/sirupsen/logrus#Logger

	// make sure that the user actually requested syslog logging as it is
	// currently supported on UNIX only.
	if *c.UseSyslog {
		logBuffer.Add(logging.LogRecord{
			Level:   logrus.InfoLevel,
			Message: "Syslog logging requested, attempting to enable it",
			Fields:  logrus.Fields{"use_syslog": c.UseSyslog},
		})

		if err := logging.EnableSyslogLogging(c.logger, &logBuffer, *c.LogLevel); err != nil {
			// TODO: Is this sufficient cause for failing? Perhaps if a local
			// log file is not also set consider it a failure?

			logBuffer.Add(logging.LogRecord{
				Level:   logrus.ErrorLevel,
				Message: fmt.Sprintf("Failed to enable syslog logging: %s", err),
				Fields:  logrus.Fields{"use_syslog": c.UseSyslog},
			})

			logBuffer.Add(logging.LogRecord{
				Level:   logrus.WarnLevel,
				Message: "Proceeding without syslog logging",
				Fields:  logrus.Fields{"use_syslog": c.UseSyslog},
			})
		}
	} else {
		logBuffer.Add(logging.LogRecord{
			Level:   logrus.DebugLevel,
			Message: "Syslog logging not requested, not enabling",
			Fields:  logrus.Fields{"use_syslog": c.UseSyslog},
		})

	}

}

// GetAppName returns the appName field if it's non-nil, zero value otherwise.
func (c *Config) GetAppName() string {
	if c == nil || c.AppName == nil {
		return ""
	}
	return *c.AppName
}

// GetAppDescription returns the appDescription field if it's non-nil, zero value otherwise.
func (c *Config) GetAppDescription() string {
	if c == nil || c.AppDescription == nil {
		return ""
	}
	return *c.AppDescription

}

// GetAppVersion returns the appVersion field if it's non-nil, zero value otherwise.
func (c *Config) GetAppVersion() string {
	if c == nil || c.AppVersion == nil {
		return ""
	}
	return *c.AppVersion
}

// GetAppURL returns the appURL field if it's non-nil, zero value otherwise.
func (c *Config) GetAppURL() string {
	if c == nil || c.AppURL == nil {
		return ""
	}
	return *c.AppURL
}

// GetFilePattern returns the filePattern field if it's non-nil, zero value otherwise.
func (c *Config) GetFilePattern() string {
	if c == nil || c.FilePattern == nil {
		return ""
	}
	return *c.FilePattern
}

// GetFileExtensions returns the fileExtensions field if it's non-nil, zero value
// otherwise.
// TODO: Double check this one; how should we safely handle returning an
// empty/zero value?
// As an example, the https://github.com/google/go-github package has a
// `Issue.GetAssignees()` method that returns nil if the `Issue.Assignees`
// field is nil. This seems to suggest that this is all we really can do here?
//
func (c *Config) GetFileExtensions() []string {
	if c == nil || c.FileExtensions == nil {
		// FIXME: Isn't the goal to avoid returning nil?
		return nil
	}
	return *c.FileExtensions
}

// GetFileAge returns the fileAge field if it's non-nil, zero value otherwise.
func (c *Config) GetFileAge() int {
	if c == nil || c.FileAge == nil {
		return 0
	}
	return *c.FileAge
}

// GetNumFilesToKeep returns the numFilesToKeep field if it's non-nil, zero value
// otherwise.
func (c *Config) GetNumFilesToKeep() int {
	if c == nil || c.NumFilesToKeep == nil {
		return 0
	}
	return *c.NumFilesToKeep
}

// GetKeepOldest returns the keepOldest field if it's non-nil, zero value
// otherwise.
func (c *Config) GetKeepOldest() bool {
	if c == nil || c.KeepOldest == nil {
		return false
	}
	return *c.KeepOldest
}

// GetRemove returns the remove field if it's non-nil, zero value otherwise.
func (c *Config) GetRemove() bool {
	if c == nil || c.Remove == nil {
		return false
	}
	return *c.Remove
}

// GetIgnoreErrors returns the ignoreErrors field if it's non-nil, zero value
// otherwise.
func (c *Config) GetIgnoreErrors() bool {
	if c == nil || c.IgnoreErrors == nil {
		return false
	}
	return *c.IgnoreErrors
}

// GetPaths returns the paths field if it's non-nil, zero value otherwise.
func (c *Config) GetPaths() []string {
	if c == nil || c.Paths == nil {
		return nil
	}
	return *c.Paths
}

// GetRecursiveSearch returns the recursiveSearch field if it's non-nil, zero
// value otherwise.
func (c *Config) GetRecursiveSearch() bool {
	if c == nil || c.RecursiveSearch == nil {
		return false
	}
	return *c.RecursiveSearch
}

// GetLogLevel returns the logLevel field if it's non-nil, zero value otherwise.
func (c *Config) GetLogLevel() string {
	if c == nil || c.LogLevel == nil {
		return ""
	}
	return *c.LogLevel
}

// GetLogFormat returns the logFormat field if it's non-nil, zero value otherwise.
func (c *Config) GetLogFormat() string {
	if c == nil || c.LogFormat == nil {
		return ""
	}
	return *c.LogFormat
}

// GetLogFilePath returns the logFilePath field if it's non-nil, zero value
// otherwise.
func (c *Config) GetLogFilePath() string {
	if c == nil || c.LogFilePath == nil {
		return ""
	}
	return *c.LogFilePath
}

// GetConsoleOutput returns the consoleOutput field if it's non-nil, zero value
// otherwise.
func (c *Config) GetConsoleOutput() string {
	if c == nil || c.ConsoleOutput == nil {
		return ""
	}
	return *c.ConsoleOutput
}

// GetUseSyslog returns the useSyslog field if it's non-nil, zero
// value otherwise.
func (c *Config) GetUseSyslog() bool {
	if c == nil || c.UseSyslog == nil {
		return false
	}
	return *c.UseSyslog
}

// GetConfigFile returns the configFile field if it's non-nil, zero value
// otherwise.
func (c *Config) GetConfigFile() string {
	if c == nil || c.ConfigFile == nil {
		return ""
	}
	return *c.ConfigFile
}

// GetLogger returns the logger field if it's non-nil, zero value otherwise.
func (c *Config) GetLogger() *logrus.Logger {
	if c == nil || c.logger == nil {
		return nil
	}
	return c.logger
}

// GetFlagParser returns the flagParser field if it's non-nil, zero value otherwise.
func (c *Config) GetFlagParser() *arg.Parser {
	if c == nil || c.flagParser == nil {
		return nil
	}
	return c.flagParser
}

// GetLogFileHandle returns the logFileHandle field if it's non-nil, zero value otherwise.
func (c *Config) GetLogFileHandle() *os.File {
	if c == nil || c.logFileHandle == nil {
		return nil
	}
	return c.logFileHandle
}
