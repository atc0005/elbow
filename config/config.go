// Copyright 2020 Adam Chalkley
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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/atc0005/elbow/logging"

	"github.com/alexflint/go-arg"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

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
	FilePattern    *string  `toml:"pattern" arg:"--pattern,env:ELBOW_FILE_PATTERN" help:"Substring pattern to compare filenames against. Wildcards are not supported."`
	FileExtensions []string `toml:"file_extensions" arg:"--extensions,env:ELBOW_EXTENSIONS" help:"Limit search to specified file extensions. Specify as space separated list to match multiple required extensions."`
	FileAge        *int     `toml:"file_age" arg:"--age,env:ELBOW_FILE_AGE" help:"Limit search to files that are the specified number of days old or older."`
	NumFilesToKeep *int     `toml:"files_to_keep" arg:"--keep,env:ELBOW_KEEP" help:"Keep specified number of matching files per provided path."`
	KeepOldest     *bool    `toml:"keep_oldest" arg:"--keep-old,env:ELBOW_KEEP_OLD" help:"Keep oldest files instead of newer per provided path."`
	Remove         *bool    `toml:"remove" arg:"--remove,env:ELBOW_REMOVE" help:"Remove matched files per provided path."`
	IgnoreErrors   *bool    `toml:"ignore_errors" arg:"--ignore-errors,env:ELBOW_IGNORE_ERRORS" help:"Ignore errors encountered during file removal."`
}

// Search represents options specific to controlling how this application
// performs searches in the filesystem
type Search struct {
	Paths           []string `toml:"paths" arg:"--paths,env:ELBOW_PATHS" help:"List of comma or space-separated paths to process."`
	RecursiveSearch *bool    `toml:"recursive_search" arg:"--recurse,env:ELBOW_RECURSE" help:"Perform recursive search into subdirectories per provided path."`
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

// NewDefaultConfig returns a newly constructed config object composed of
// default configuration settings.
func NewDefaultConfig(appVersion string) Config {

	// TODO: Is there a better way than creating a "throwaway" config object
	// just to make use of its methods for retrieving default values?
	c := Config{}
	defaultAppName := c.GetAppName()
	defaultAppDescription := c.GetAppDescription()
	defaultAppURL := c.GetAppURL()
	defaultFilePattern := c.GetFilePattern()
	defaultFileAge := c.GetFileAge()
	defaultNumFilesToKeep := c.GetNumFilesToKeep()
	defaultKeepOldest := c.GetKeepOldest()
	defaultRemove := c.GetRemove()
	defaultIgnoreErrors := c.GetIgnoreErrors()
	defaultRecursiveSearch := c.GetRecursiveSearch()
	defaultLogLevel := c.GetLogLevel()
	defaultLogFormat := c.GetLogFormat()
	defaultLogFilePath := c.GetLogFilePath()
	defaultConsoleOutput := c.GetConsoleOutput()
	defaultUseSyslog := c.GetUseSyslog()
	defaultConfigFile := c.GetConfigFile()

	defaultConfig := Config{
		AppMetadata: AppMetadata{
			AppName:        defaultAppName,
			AppDescription: defaultAppDescription,
			AppURL:         defaultAppURL,
			AppVersion:     appVersion,
		},
		FileHandling: FileHandling{
			FilePattern: &defaultFilePattern,
			//FileExtensions: &fileExtensions,
			FileAge:        &defaultFileAge,
			NumFilesToKeep: &defaultNumFilesToKeep,
			KeepOldest:     &defaultKeepOldest,
			Remove:         &defaultRemove,
			IgnoreErrors:   &defaultIgnoreErrors,
		},
		Logging: Logging{
			LogLevel:      &defaultLogLevel,
			LogFormat:     &defaultLogFormat,
			LogFilePath:   &defaultLogFilePath,
			ConsoleOutput: &defaultConsoleOutput,
			UseSyslog:     &defaultUseSyslog,
		},
		Search: Search{
			//Paths: ,
			RecursiveSearch: &defaultRecursiveSearch,
		},
		ConfigFile: &defaultConfigFile,
	}

	return defaultConfig
}

// NewConfig returns a pointer to a newly configured object representing a
// collection of user-provided and default settings.
func NewConfig(appVersion string) (*Config, error) {

	// fmt.Printf("os.Args quoted: %q\n", os.Args)
	// fmt.Printf("os.Args bare: %v\n", os.Args)
	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: "os.Args array contents",
		Fields: logrus.Fields{
			"line":    logging.GetLineNumber(),
			"os_args": os.Args,
		},
	})

	// Apply default settings that other configuration sources will be allowed
	// to (and for a few settings MUST) override
	baseConfig := NewDefaultConfig(appVersion)

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Current baseConfig after NewDefaultConfig() call: %+v\n", baseConfig),
		Fields:  logrus.Fields{"line": logging.GetLineNumber()},
	})

	// Settings provided via config file. Intentionally using uninitialized
	// struct here so that we can check for nil pointers to indicate whether
	// a field has been populated with configuration values.
	fileConfig := Config{}

	// Settings provided via command-line flags and environment variables.
	// This object will always be set in some manner as either flags or env
	// vars will be needed to bootstrap the application. While we may support
	// using a configuration file to provide settings, it is not used by
	// default. Directly add app version string here since this is the config
	// object referenced if the user requests either `--version` or `-h`
	// flags; the combined config object is not created in time to serve either
	// of those needs.
	argsConfig := Config{
		AppMetadata: AppMetadata{
			AppVersion: appVersion,
		},
	}

	// Initialize logger "handle" for later use
	baseConfig.logger = logrus.New()

	// Bundle the returned `*.arg.Parser` for later use from `main()` so that
	// we can explicitly display usage or help details should the
	// user-provided settings fail validation.
	baseConfig.flagParser = arg.MustParse(&argsConfig)

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Current argsConfig after MustParse() call: %+v\n", argsConfig),
		Fields:  logrus.Fields{"line": logging.GetLineNumber()},
	})

	/*************************************************************************
		At this point `baseConfig` is our baseline config object containing
		default settings and various handles to other resources. We do not
		apply those same resource handles to other config structs. We merge
		the other configuration objects into the baseConfig object to create
		a unified configuration object that we return to the caller.
	*************************************************************************/

	// If user specified a config file, let's try to use it
	// TODO: Fail if not found, or continue using defaults in its place?
	if argsConfig.ConfigFile != nil {
		// Check for a configuration file and load it if found.

		// path not found
		if _, err := os.Stat(*argsConfig.ConfigFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("requested config file not found: %v", err)
		}

		fh, err := os.Open(*argsConfig.ConfigFile)
		if err != nil {
			return nil, fmt.Errorf("unable to open config file: %v", err)
		}
		defer func() {
			if err := fh.Close(); err != nil {
				// Ignore "file already closed" errors
				if !errors.Is(err, os.ErrClosed) {
					logging.Buffer.Add(logging.LogRecord{
						Level: logrus.ErrorLevel,
						Message: fmt.Sprintf(
							"failed to close file %q: %s",
							*argsConfig.ConfigFile,
							err.Error(),
						),
						Fields: logrus.Fields{"line": logging.GetLineNumber()},
					})
				}
			}
		}()

		if err := fileConfig.LoadConfigFile(fh); err != nil {
			logging.Buffer.Add(logging.LogRecord{
				Level:   logrus.ErrorLevel,
				Message: fmt.Sprintf("Error loading config file: %s", err),
				Fields:  logrus.Fields{"config_file": argsConfig.ConfigFile},
			})

			// Application failure codepath. Dump collected log messages and
			// return control to the caller.
			if err := logging.Buffer.Flush(baseConfig.GetLogger()); err != nil {
				// if we're unable to flush the buffer, then something serious
				// has occurred and we should emit the error directly to the
				// console
				fmt.Printf("Failed to flush the log buffer: %v", err.Error())
			}

			// TODO: Wrap errors and return so they can be unpacked in main()
			return nil, fmt.Errorf("error loading configuration file: %s", err)
		}

		logging.Buffer.Add(logging.LogRecord{
			Level:   logrus.DebugLevel,
			Message: fmt.Sprintf("Current fileConfig after LoadConfigFile() call: %+v\n", fileConfig),
			Fields:  logrus.Fields{"line": logging.GetLineNumber()},
		})

		logging.Buffer.Add(logging.LogRecord{
			Level:   logrus.DebugLevel,
			Message: "Processing fileConfig object with MergeConfig func",
			Fields:  logrus.Fields{"line": logging.GetLineNumber()},
		})

		if err := MergeConfig(&baseConfig, fileConfig); err != nil {
			logging.Buffer.Add(logging.LogRecord{
				Level:   logrus.ErrorLevel,
				Message: fmt.Sprintf("Error merging config file settings with base config: %s", err),
				Fields: logrus.Fields{
					"line":        logging.GetLineNumber(),
					"base_config": fmt.Sprintf("%+v", baseConfig),
					"file_config": fmt.Sprintf("%+v", fileConfig),
				},
			})
		}

		// Don't fail the new configuration due to fileConfig not providing
		// all required values; we are *feathering* values, not replacing all
		// existing values in the config struct with ones from the next
		// configuration source.
		if err := baseConfig.Validate(); err != nil {
			logging.Buffer.Add(logging.LogRecord{
				Level:   logrus.DebugLevel,
				Message: fmt.Sprintf("Error validating config after merging %s: %s", "fileConfig", err),
				Fields: logrus.Fields{
					"line":        logging.GetLineNumber(),
					"base_config": fmt.Sprintf("%+v", baseConfig),
					"file_config": fmt.Sprintf("%+v", fileConfig),
				},
			})

			logging.Buffer.Add(logging.LogRecord{
				Level:   logrus.DebugLevel,
				Message: "Proceeding with evaluation of argsConfig",
				Fields: logrus.Fields{
					"line": logging.GetLineNumber(),
				},
			})

		}
	}

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: "Processing argsConfig object with MergeConfig func",
	})

	if err := MergeConfig(&baseConfig, argsConfig); err != nil {
		logging.Buffer.Add(logging.LogRecord{
			Level:   logrus.ErrorLevel,
			Message: fmt.Sprintf("Error merging args config settings with base config: %s", err),
			Fields: logrus.Fields{
				"line":        logging.GetLineNumber(),
				"base_config": fmt.Sprintf("%+v", baseConfig),
				"args_config": fmt.Sprintf("%+v", argsConfig),
			},
		})
	}

	if err := baseConfig.Validate(); err != nil {

		// ###################################################################
		// This code should only be reached if we were unable to properly
		// apply the configuration. At this point we cannot trust that our
		// settings are valid. We should ensure default settings are applied
		// to our logger instance, flush all held messages and then exit
		// immediately.
		// ###################################################################

		logging.Buffer.Add(logging.LogRecord{
			Level:   logrus.DebugLevel,
			Message: fmt.Sprintf("Error validating config after merging %s: %s", "argsConfig", err),
			Fields: logrus.Fields{
				"line":        logging.GetLineNumber(),
				"base_config": fmt.Sprintf("%+v", baseConfig),
				"args_config": fmt.Sprintf("%+v", argsConfig),
			},
		})

		// Application failure codepath. Dump collected log messages and
		// return control to the caller.
		if err := logging.Buffer.Flush(baseConfig.GetLogger()); err != nil {
			// if we're unable to flush the buffer, then something serious
			// has occurred and we should emit the error directly to the
			// console
			fmt.Printf("Failed to flush the log buffer: %v", err.Error())
		}
		// TODO: Wrap errors and return so they can be unpacked in main()
		return nil, fmt.Errorf("configuration validation failed after merging argsConfig: %s", err)

	}

	// Apply logging configuration. If error is encountered, pass it back to
	// caller to deal with.
	if err := baseConfig.SetLoggerConfig(); err != nil {
		return nil, err
	}

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("The config object that we are returning (raw format): %+v", baseConfig),
		Fields:  logrus.Fields{"line": logging.GetLineNumber()},
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprint("The config object that we are returning (string format): ", baseConfig.String()),
		Fields:  logrus.Fields{"line": logging.GetLineNumber()},
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: "Empty queued up log messages from log buffer using user-specified logging settings",
		Fields:  logrus.Fields{"line": logging.GetLineNumber()},
	})

	if err := logging.Buffer.Flush(baseConfig.GetLogger()); err != nil {
		// if we're unable to flush the buffer, then something serious
		// has occurred and we should emit the error directly to the
		// console
		fmt.Printf("Failed to flush the log buffer: %v", err.Error())
	}

	return &baseConfig, nil

}

// LoadConfigFile reads from an io.Reader and unmarshals a configuration file
// in TOML format into the associated Config struct.
func (c *Config) LoadConfigFile(fileHandle io.Reader) error {

	configFile, err := ioutil.ReadAll(fileHandle)
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

	return fmt.Sprintf("%s %s", c.GetAppName(), c.GetAppDescription())
}

// Version provides a version string that appears at the top of the
// application Help output
func (c Config) Version() string {

	versionString := fmt.Sprintf("%s %s\n%s",
		strings.ToTitle(c.GetAppName()), c.GetAppVersion(), c.GetAppURL())

	//divider := strings.Repeat("-", len(versionString))

	// versionBlock := fmt.Sprintf("\n%s\n%s\n%s\n",
	// 	divider, versionString, divider)

	//return versionBlock

	return "\n" + versionString + "\n"
}

// String satisfies the Stringer interface. This is intended for non-JSON
// formatting if using the TextFormatter logrus formatter.
func (c *Config) String() string {

	return fmt.Sprintf("AppName=%q, AppDescription=%q, AppVersion=%q, AppURL=%q, FilePattern=%q, FileExtensions=%q, Paths=%v, RecursiveSearch=%t, FileAge=%d, NumFilesToKeep=%d, KeepOldest=%t, Remove=%t, IgnoreErrors=%t, LogFormat=%q, LogFilePath=%q, ConfigFile=%q, ConsoleOutput=%q, LogLevel=%q, UseSyslog=%t, logger=%v, flagParser=%v,  logFileHandle=%v",

		c.GetAppName(),
		c.GetAppDescription(),
		c.GetAppVersion(),
		c.GetAppURL(),
		c.GetFilePattern(),
		c.GetFileExtensions(),
		c.GetPaths(),
		c.GetRecursiveSearch(),
		c.GetFileAge(),
		c.GetNumFilesToKeep(),
		c.GetKeepOldest(),
		c.GetRemove(),
		c.GetIgnoreErrors(),
		c.GetLogFormat(),
		c.GetLogFilePath(),
		c.GetConfigFile(),
		c.GetConsoleOutput(),
		c.GetLogLevel(),
		c.GetUseSyslog(),
		c.GetLogger(),
		c.GetFlagParser(),
		c.GetLogFileHandle(),
	)
}

// WriteDefaultHelpText is a helper function used to output Help text for
// situations where the Config object cannot be trusted to be in a usable
// state.
// TODO: Reconsider this; this feels fragile.
func WriteDefaultHelpText(appName string) {
	config := arg.Config{Program: appName}
	defaultConfig := Config{}
	parser, err := arg.NewParser(config, &defaultConfig)
	if err != nil {
		panic("failed to build go-arg parser for Help text generation")
	}
	parser.WriteUsage(os.Stdout)
}
