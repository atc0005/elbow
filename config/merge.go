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

	"github.com/atc0005/elbow/logging"
	"github.com/sirupsen/logrus"
)

// MergeConfig receives source, destination and default Config objects and
// merges select, non-nil field values from the source Config object to
// the destination config object, overwriting any field value already present.
//
// The goal is to respect the current documented configuration precedence for
// multiple configuration sources (e.g., config file and command-line flags).
func MergeConfig(destination *Config, source Config) error {

	// FIXME: How can we get all field names programatically so we don't have to
	// manually reference each field?

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: "MergeConfig starting",
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Source struct (raw): %+v", source),
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Source struct (string): %s", source.String()),
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Destination struct (raw): %+v", destination),
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Destination struct (string): %s", destination.String()),
	})

	if source.AppName != nil {
		*destination.AppName = *source.AppName
	}

	if source.AppDescription != nil {
		*destination.AppDescription = *source.AppDescription
	}

	if source.AppURL != nil {
		*destination.AppURL = *source.AppURL
	}

	if source.AppVersion != nil {
		*destination.AppVersion = *source.AppVersion
	}

	if source.Paths != nil {
		destination.Paths = source.Paths
	}

	if source.FileExtensions != nil {
		destination.FileExtensions = source.FileExtensions
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
		// TODO: Add debug output for each of these decisions, but enable it
		// only for debug/troubleshooting builds.
		// fmt.Printf("Using %q for ConsoleOutput\n", *source.ConsoleOutput)
		*destination.ConsoleOutput = *source.ConsoleOutput
	}

	if source.UseSyslog != nil {
		*destination.UseSyslog = *source.UseSyslog
	}

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: "MergeConfig ending",
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Source struct (raw): %+v", source),
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Source struct (string): %s", source.String()),
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Destination struct (raw): %+v", destination),
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Destination struct (string): %s", destination.String()),
	})

	// FIXME: Placeholder
	// FIXME: What useful error code would we return from this function?
	return nil
}
