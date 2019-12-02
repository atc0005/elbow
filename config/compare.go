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
	"testing"
)

// https://stackoverflow.com/questions/15311969/checking-the-equality-of-two-slices
// https://stackoverflow.com/a/15312097
// Author: Stephen Weinberg
// https://creativecommons.org/licenses/by-sa/4.0/
//
// TODO: Not sure where to place this. If I keep it, I should add an entry
// to NOTICE.txt
func testStringSliceEqual(a []string, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// CompareConfig receives two Config objects and compares exported field
// values to determine equality.
func CompareConfig(actual Config, wanted Config, t *testing.T) {

	// FIXME: How can we get all field names programatically so we don't have to
	// manually reference each field?

	t.Logf("actual config (raw): %+v", actual)
	t.Logf("actual config (string): %s", actual.String())

	t.Logf("wanted config (raw): %+v", wanted)
	t.Logf("wanted config (string): %s", wanted.String())

	// Validate both config structs first before attempting to do anything
	// useful with them.
	if err := actual.Validate(); err != nil {
		t.Error("Validation failed for actual config", err)
	}

	if err := wanted.Validate(); err != nil {
		t.Error("Validation failed for wanted config", err)
	}

	// Fail each field that fails to validate, but allow testing to continue
	// until all fields have been evaluated

	if *actual.AppName != *wanted.AppName {
		t.Errorf("AppName: actual (%v) does not equal wanted (%v)",
			*actual.AppName, *wanted.AppName)
	} else {
		t.Logf("AppName: actual (%v) == wanted (%v)",
			*actual.AppName, *wanted.AppName)
	}

	if *actual.AppDescription != *wanted.AppDescription {
		t.Errorf("AppDescription: actual (%v) does not equal wanted (%v)",
			*actual.AppDescription, *wanted.AppDescription)
	} else {
		t.Logf("AppDescription: actual (%v) == wanted (%v)",
			*actual.AppDescription, *wanted.AppDescription)
	}

	if *actual.AppURL != *wanted.AppURL {
		t.Errorf("AppURL: actual (%v) does not equal wanted (%v)",
			*actual.AppURL, *wanted.AppURL)
	} else {
		t.Logf("AppURL: actual (%v) == wanted (%v)",
			*actual.AppURL, *wanted.AppURL)
	}

	if *actual.AppVersion != *wanted.AppVersion {
		t.Errorf("AppVersion: actual (%v) does not equal wanted (%v)",
			*actual.AppVersion, *wanted.AppVersion)
	} else {
		t.Logf("AppVersion: actual (%v) == wanted (%v)",
			*actual.AppVersion, *wanted.AppVersion)
	}

	if !testStringSliceEqual(actual.Paths, wanted.Paths) {
		t.Errorf("Paths: actual (%v) does not equal wanted (%v)",
			actual.Paths, wanted.Paths)
	} else {
		t.Logf("Paths: actual (%v) == wanted (%v)",
			actual.Paths, wanted.Paths)
	}

	if !testStringSliceEqual(actual.FileExtensions, wanted.FileExtensions) {
		t.Errorf("FileExtensions: actual (%v) does not equal wanted (%v)",
			actual.FileExtensions, wanted.FileExtensions)
	} else {
		t.Logf("FileExtensions: actual (%v) == wanted (%v)",
			actual.FileExtensions, wanted.FileExtensions)
	}

	if *actual.FilePattern != *wanted.FilePattern {
		t.Errorf("FilePattern: actual (%v) does not equal wanted (%v)",
			*actual.FilePattern, *wanted.FilePattern)
	} else {
		t.Logf("FilePattern: actual (%v) == wanted (%v)",
			*actual.FilePattern, *wanted.FilePattern)
	}

	if *actual.FileAge != *wanted.FileAge {
		t.Errorf("FileAge: actual (%v) does not equal wanted (%v)",
			*actual.FileAge, *wanted.FileAge)
	} else {
		t.Logf("FileAge: actual (%v) == wanted (%v)",
			*actual.FileAge, *wanted.FileAge)
	}

	if *actual.NumFilesToKeep != *wanted.NumFilesToKeep {
		t.Errorf("NumFilesToKeep: actual (%v) does not equal wanted (%v)",
			*actual.NumFilesToKeep, *wanted.NumFilesToKeep)
	} else {
		t.Logf("NumFilesToKeep: actual (%v) == wanted (%v)",
			*actual.NumFilesToKeep, *wanted.NumFilesToKeep)
	}

	if *actual.KeepOldest != *wanted.KeepOldest {
		t.Errorf("KeepOldest: actual (%v) does not equal wanted (%v)",
			*actual.KeepOldest, *wanted.KeepOldest)
	} else {
		t.Logf("KeepOldest: actual (%v) == wanted (%v)",
			*actual.KeepOldest, *wanted.KeepOldest)
	}

	if *actual.Remove != *wanted.Remove {
		t.Errorf("Remove: actual (%v) does not equal wanted (%v)",
			*actual.Remove, *wanted.Remove)
	} else {
		t.Logf("Remove: actual (%v) == wanted (%v)",
			*actual.Remove, *wanted.Remove)
	}

	if *actual.IgnoreErrors != *wanted.IgnoreErrors {
		t.Errorf("IgnoreErrors: actual (%v) does not equal wanted (%v)",
			*actual.IgnoreErrors, *wanted.IgnoreErrors)
	} else {
		t.Logf("IgnoreErrors: actual (%v) == wanted (%v)",
			*actual.IgnoreErrors, *wanted.IgnoreErrors)
	}

	if *actual.RecursiveSearch != *wanted.RecursiveSearch {
		t.Errorf("RecursiveSearch: actual (%v) does not equal wanted (%v)",
			*actual.RecursiveSearch, *wanted.RecursiveSearch)
	} else {
		t.Logf("RecursiveSearch: actual (%v) == wanted (%v)",
			*actual.RecursiveSearch, *wanted.RecursiveSearch)
	}

	if *actual.LogLevel != *wanted.LogLevel {
		t.Errorf("LogLevel: actual (%v) does not equal wanted (%v)",
			*actual.LogLevel, *wanted.LogLevel)
	} else {
		t.Logf("LogLevel: actual (%v) == wanted (%v)",
			*actual.LogLevel, *wanted.LogLevel)
	}

	if *actual.LogFormat != *wanted.LogFormat {
		t.Errorf("LogFormat: actual (%v) does not equal wanted (%v)",
			*actual.LogFormat, *wanted.LogFormat)
	} else {
		t.Logf("LogFormat: actual (%v) == wanted (%v)",
			*actual.LogFormat, *wanted.LogFormat)
	}

	if *actual.LogFilePath != *wanted.LogFilePath {
		t.Errorf("LogFilePath: actual (%v) does not equal wanted (%v)",
			*actual.LogFilePath, *wanted.LogFilePath)
	} else {
		t.Logf("LogFilePath: actual (%v) == wanted (%v)",
			*actual.LogFilePath, *wanted.LogFilePath)
	}

	if *actual.ConsoleOutput != *wanted.ConsoleOutput {
		t.Errorf("ConsoleOutput: actual (%v) does not equal wanted (%v)",
			*actual.ConsoleOutput, *wanted.ConsoleOutput)
	} else {
		t.Logf("ConsoleOutput: actual (%v) == wanted (%v)",
			*actual.ConsoleOutput, *wanted.ConsoleOutput)
	}

	if *actual.UseSyslog != *wanted.UseSyslog {
		t.Errorf("UseSyslog: actual (%v) does not equal wanted (%v)",
			*actual.UseSyslog, *wanted.UseSyslog)
	} else {
		t.Logf("UseSyslog: actual (%v) == wanted (%v)",
			*actual.UseSyslog, *wanted.UseSyslog)
	}

	t.Log("CompareConfig complete")

}
