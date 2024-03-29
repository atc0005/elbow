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

package config

import (
	"fmt"

	"github.com/atc0005/elbow/internal/logging"
)

// Validate verifies all struct fields have been provided acceptable values
func (c Config) Validate() error {

	if c.AppName == "" {
		return fmt.Errorf("field AppName not configured")
	}

	if c.AppDescription == "" {
		return fmt.Errorf("field AppDescription not configured")
	}

	if c.AppVersion == "" {
		return fmt.Errorf("field AppVersion not configured")
	}

	if c.AppURL == "" {
		return fmt.Errorf("field AppURL not configured")
	}

	// FilePattern is optional, but since has an underlying string type with a
	// default of empty string we can assert that the pointer isn't
	if c.FilePattern == nil {
		return fmt.Errorf("field FilePattern not configured")
	}

	// FileExtensions is optional
	// Discovered files are checked against FileExtensions later
	// This isn't a pointer, but rather a string slice. The user may opt to
	// not configure this setting, so having a `nil` state for this setting is
	// normal?
	//
	// if c.FileExtensions == nil {
	// 	return false, fmt.Errorf("file extensions option not configured")
	// }

	if c.Paths == nil {
		return fmt.Errorf("one or more paths not provided")
	}

	// RecursiveSearch is optional
	if c.RecursiveSearch == nil {
		return fmt.Errorf("field RecursiveSearch not configured")
	}

	// NumFilesToKeep is optional, but should be configured via
	// if specified we should make sure it is a non-negative number.
	switch {
	case c.NumFilesToKeep == nil:
		return fmt.Errorf("field NumFilesToKeep not configured")
	case *c.NumFilesToKeep < 0:
		return fmt.Errorf("negative number for files to keep not supported")
	}

	// We only want to work with positive file modification times 0 is
	// acceptable as it is the default value and indicates that the user has
	// not chosen to use the flag (or has chosen improperly and it will be
	// ignored).
	switch {
	case c.FileAge == nil:
		return fmt.Errorf("field FileAge not configured")
	case *c.FileAge < 0:
		return fmt.Errorf("negative number for file age not supported")
	}

	if c.KeepOldest == nil {
		return fmt.Errorf("field KeepOldest not configured")
	}

	if c.Remove == nil {
		return fmt.Errorf("field Remove not configured")
	}

	if c.IgnoreErrors == nil {
		return fmt.Errorf("field IgnoreErrors not configured")
	}

	switch {
	case c.LogFormat == nil:
		return fmt.Errorf("field LogFormat not configured")
	case *c.LogFormat == logging.LogFormatText:
	case *c.LogFormat == logging.LogFormatJSON:
	default:
		return fmt.Errorf("invalid option %q provided for log format", *c.LogFormat)
	}

	// LogFilePath is optional, but should still have a non-nil value
	if c.LogFilePath == nil {
		return fmt.Errorf("field LogFilePath not configured")
	}

	// Do nothing for valid choices, return false if invalid value specified
	switch {
	case c.ConsoleOutput == nil:
		return fmt.Errorf("field ConsoleOutput not configured")
	// TODO: Evaluate replacing bare strings with constants (see constants.go)
	case *c.ConsoleOutput == logging.ConsoleOutputStdout:
	case *c.ConsoleOutput == logging.ConsoleOutputStderr:
	default:
		return fmt.Errorf("invalid option %q provided for console output destination", *c.ConsoleOutput)
	}

	switch {
	case c.LogLevel == nil:
		return fmt.Errorf("field LogLevel not configured")
	case *c.LogLevel == logging.LogLevelEmergency:
	case *c.LogLevel == logging.LogLevelAlert:
	case *c.LogLevel == logging.LogLevelCritical:
	case *c.LogLevel == logging.LogLevelPanic:
	case *c.LogLevel == logging.LogLevelFatal:
	case *c.LogLevel == logging.LogLevelError:
	case *c.LogLevel == logging.LogLevelWarn:
	case *c.LogLevel == logging.LogLevelInfo:
	case *c.LogLevel == logging.LogLevelNotice:
	case *c.LogLevel == logging.LogLevelDebug:
	case *c.LogLevel == logging.LogLevelTrace:
	default:
		return fmt.Errorf("invalid option %q provided for log level", *c.LogLevel)
	}

	// UseSyslog is optional, but should still be initialized
	if c.UseSyslog == nil {
		return fmt.Errorf("field UseSyslog not configured")
	}

	// Make sure that a valid logger has been created
	if c.logger == nil {
		return fmt.Errorf("field logger not configured")
	}

	// Using a config file is optional, but should still be initialized so
	// that user values can be stored later if specified.
	if c.ConfigFile == nil {
		return fmt.Errorf("field ConfigFile not configured")
	}

	// Optimist
	return nil

}
