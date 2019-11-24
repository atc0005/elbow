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
)

// Validate verifies all struct fields have been provided acceptable
func (c Config) Validate() (bool, error) {

	// FilePattern is optional, but since has an underlying string type with a
	// default of empty string we can assert that the pointer isn't
	if c.FilePattern == nil {
		return false, fmt.Errorf("field FilePattern not configured")
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
		return false, fmt.Errorf("one or more paths not provided")
	}

	// RecursiveSearch is optional
	if c.RecursiveSearch == nil {
		return false, fmt.Errorf("field RecursiveSearch not configured")
	}

	// NumFilesToKeep is optional, but should be configured via
	// if specified we should make sure it is a non-negative number.
	switch {
	case c.NumFilesToKeep == nil:
		return false, fmt.Errorf("field NumFilesToKeep not configured")
	case *c.NumFilesToKeep < 0:
		return false, fmt.Errorf("invalid value provided for files to keep")
	}

	// We only want to work with positive file modification times 0 is
	// acceptable as it is the default value and indicates that the user has
	// not chosen to use the flag (or has chosen improperly and it will be
	// ignored).
	switch {
	case c.FileAge == nil:
		return false, fmt.Errorf("field FileAge not configured")
	case *c.FileAge < 0:
		return false, fmt.Errorf("negative number for file age not supported")
	}

	if c.KeepOldest == nil {
		return false, fmt.Errorf("field KeepOldest not configured")
	}

	if c.Remove == nil {
		return false, fmt.Errorf("field Remove not configured")
	}

	if c.IgnoreErrors == nil {
		return false, fmt.Errorf("field IgnoreErrors not configured")
	}

	switch {
	case c.LogFormat == nil:
		return false, fmt.Errorf("field LogFormat not configured")
	case *c.LogFormat == "text":
	case *c.LogFormat == "json":
	default:
		return false, fmt.Errorf("invalid option %q provided for log format", *c.LogFormat)
	}

	// logFilePath is optional
	// TODO: String validation if it is set?

	// Do nothing for valid choices, return false if invalid value specified
	switch {
	case c.ConsoleOutput == nil:
		return false, fmt.Errorf("field ConsoleOutput not configured")
	case *c.ConsoleOutput == "stdout":
	case *c.ConsoleOutput == "stderr":
	default:
		return false, fmt.Errorf("invalid option %q provided for console output destination", *c.ConsoleOutput)
	}

	switch {
	case c.LogLevel == nil:
		return false, fmt.Errorf("field LogLevel not configured")
	case *c.LogLevel == "emergency":
	case *c.LogLevel == "alert":
	case *c.LogLevel == "critical":
	case *c.LogLevel == "panic":
	case *c.LogLevel == "fatal":
	case *c.LogLevel == "error":
	case *c.LogLevel == "warn":
	case *c.LogLevel == "info":
	case *c.LogLevel == "notice":
	case *c.LogLevel == "debug":
	case *c.LogLevel == "trace":
	default:
		return false, fmt.Errorf("invalid option %q provided for log level", *c.LogLevel)
	}

	// UseSyslog is optional
	if c.UseSyslog == nil {
		return false, fmt.Errorf("field UseSyslog not configured")
	}

	// Make sure that a valid logger has been created
	if c.logger == nil {
		return false, fmt.Errorf("field logger not configured")
	}

	// Optimist
	return true, nil

}
