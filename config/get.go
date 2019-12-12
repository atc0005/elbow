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

package config

import (
	"os"

	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
)

// GetAppName returns the appName field if it's non-nil, app default value
// otherwise
func (c *Config) GetAppName() string {
	if c == nil || c.AppName == "" {
		return defaultAppName
	}
	return c.AppName
}

// GetAppDescription returns the appDescription field if it's non-nil, app
// default value otherwise
func (c *Config) GetAppDescription() string {
	if c == nil || c.AppDescription == "" {
		return "prunes content matching specific patterns, either in a single directory or recursively through a directory tree."
	}
	return c.AppDescription

}

// GetAppVersion returns the appVersion field if it's non-nil, app default
// value otherwise
func (c *Config) GetAppVersion() string {
	if c == nil || c.AppVersion == "" {
		return "dev"
	}
	return c.AppVersion
}

// GetAppURL returns the appURL field if it's non-nil, app default value
// otherwise
func (c *Config) GetAppURL() string {
	if c == nil || c.AppURL == "" {
		return "https://github.com/atc0005/elbow"
	}
	return c.AppURL
}

// GetFilePattern returns the filePattern field if it's non-nil, app default
// value otherwise
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
	return c.FileExtensions
}

// GetFileAge returns the fileAge field if it's non-nil, app default value
// otherwise
func (c *Config) GetFileAge() int {
	if c == nil || c.FileAge == nil {
		return 0
	}
	return *c.FileAge
}

// GetNumFilesToKeep returns the numFilesToKeep field if it's non-nil, zero
// value otherwise.
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

// GetRemove returns the remove field if it's non-nil, app default value
// otherwise
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

// GetPaths returns the paths field if it's non-nil, app default value
// otherwise
func (c *Config) GetPaths() []string {
	if c == nil || c.Paths == nil {
		return nil
	}
	return c.Paths
}

// GetRecursiveSearch returns the recursiveSearch field if it's non-nil, zero
// value otherwise.
func (c *Config) GetRecursiveSearch() bool {
	if c == nil || c.RecursiveSearch == nil {
		return false
	}
	return *c.RecursiveSearch
}

// GetLogLevel returns the logLevel field if it's non-nil, app default value
// otherwise
func (c *Config) GetLogLevel() string {
	if c == nil || c.LogLevel == nil {
		return "info"
	}
	return *c.LogLevel
}

// GetLogFormat returns the logFormat field if it's non-nil, app default value
// otherwise
func (c *Config) GetLogFormat() string {
	if c == nil || c.LogFormat == nil {
		return logFormatText
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

// GetConsoleOutput returns the consoleOutput field if it's non-nil, zero
// value otherwise.
func (c *Config) GetConsoleOutput() string {
	if c == nil || c.ConsoleOutput == nil {
		return "stdout"
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

// GetLogger returns the logger field if it's non-nil, app default value
// otherwise
func (c *Config) GetLogger() *logrus.Logger {
	if c == nil || c.logger == nil {
		//return nil

		// FIXME: Is this the best logic?
		c.logger = logrus.New()
		//c.logger.Out = os.Stderr

	}
	return c.logger
}

// GetFlagParser returns the flagParser field if it's non-nil, app default
// value otherwise
func (c *Config) GetFlagParser() *arg.Parser {
	if c == nil || c.flagParser == nil {
		return nil
	}
	return c.flagParser
}

// GetLogFileHandle returns the logFileHandle field if it's non-nil, app
// default value otherwise
func (c *Config) GetLogFileHandle() *os.File {
	if c == nil || c.logFileHandle == nil {
		return nil
	}
	return c.logFileHandle
}
