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
	"testing"
)

// https://stackoverflow.com/questions/15311969/checking-the-equality-of-two-slices
// https://stackoverflow.com/a/15312097
//
// TODO: Not sure what file/package to place this.
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
func CompareConfig(got Config, wanted Config, t *testing.T) {

	// FIXME: How can we get all field names programatically so we don't have to
	// manually reference each field?

	t.Logf("got config (raw): %+v", got)
	t.Logf("got config (string): %s", got.String())

	t.Logf("wanted config (raw): %+v", wanted)
	t.Logf("wanted config (string): %s", wanted.String())

	// Validate both config structs first before attempting to do anything
	// useful with them.
	if err := got.Validate(); err != nil {
		t.Error("Validation failed for got config", err)
	}

	if err := wanted.Validate(); err != nil {
		t.Error("Validation failed for wanted config", err)
	}

	// Fail each field that fails to validate, but allow testing to continue
	// until all fields have been evaluated

	if got.AppName != wanted.AppName {
		t.Errorf("AppName: got (%v) does not equal wanted (%v)",
			got.AppName, wanted.AppName)
	} else {
		t.Logf("AppName: got (%v) == wanted (%v)",
			got.AppName, wanted.AppName)
	}

	if got.AppDescription != wanted.AppDescription {
		t.Errorf("AppDescription: got (%v) does not equal wanted (%v)",
			got.AppDescription, wanted.AppDescription)
	} else {
		t.Logf("AppDescription: got (%v) == wanted (%v)",
			got.AppDescription, wanted.AppDescription)
	}

	if got.AppURL != wanted.AppURL {
		t.Errorf("AppURL: got (%v) does not equal wanted (%v)",
			got.AppURL, wanted.AppURL)
	} else {
		t.Logf("AppURL: got (%v) == wanted (%v)",
			got.AppURL, wanted.AppURL)
	}

	if got.AppVersion != wanted.AppVersion {
		t.Errorf("AppVersion: got (%v) does not equal wanted (%v)",
			got.AppVersion, wanted.AppVersion)
	} else {
		t.Logf("AppVersion: got (%v) == wanted (%v)",
			got.AppVersion, wanted.AppVersion)
	}

	if !testStringSliceEqual(got.Paths, wanted.Paths) {
		t.Errorf("Paths: got (%q) does not equal wanted (%q)",
			got.Paths, wanted.Paths)
	} else {
		t.Logf("Paths: got (%q) == wanted (%q)",
			got.Paths, wanted.Paths)
	}

	if !testStringSliceEqual(got.Exclude, wanted.Exclude) {
		t.Errorf("Exclude: got (%q) does not equal wanted (%q)",
			got.Exclude, wanted.Exclude)
	} else {
		t.Logf("Exclude: got (%q) == wanted (%q)",
			got.Exclude, wanted.Exclude)
	}

	if !testStringSliceEqual(got.FileExtensions, wanted.FileExtensions) {
		t.Errorf("FileExtensions: got (%q) does not equal wanted (%q)",
			got.FileExtensions, wanted.FileExtensions)
	} else {
		t.Logf("FileExtensions: got (%q) == wanted (%q)",
			got.FileExtensions, wanted.FileExtensions)
	}

	if *got.FilePattern != *wanted.FilePattern {
		t.Errorf("FilePattern: got (%v) does not equal wanted (%v)",
			*got.FilePattern, *wanted.FilePattern)
	} else {
		t.Logf("FilePattern: got (%v) == wanted (%v)",
			*got.FilePattern, *wanted.FilePattern)
	}

	if *got.FileAge != *wanted.FileAge {
		t.Errorf("FileAge: got (%v) does not equal wanted (%v)",
			*got.FileAge, *wanted.FileAge)
	} else {
		t.Logf("FileAge: got (%v) == wanted (%v)",
			*got.FileAge, *wanted.FileAge)
	}

	if *got.NumFilesToKeep != *wanted.NumFilesToKeep {
		t.Errorf("NumFilesToKeep: got (%v) does not equal wanted (%v)",
			*got.NumFilesToKeep, *wanted.NumFilesToKeep)
	} else {
		t.Logf("NumFilesToKeep: got (%v) == wanted (%v)",
			*got.NumFilesToKeep, *wanted.NumFilesToKeep)
	}

	if *got.KeepOldest != *wanted.KeepOldest {
		t.Errorf("KeepOldest: got (%v) does not equal wanted (%v)",
			*got.KeepOldest, *wanted.KeepOldest)
	} else {
		t.Logf("KeepOldest: got (%v) == wanted (%v)",
			*got.KeepOldest, *wanted.KeepOldest)
	}

	if *got.Remove != *wanted.Remove {
		t.Errorf("Remove: got (%v) does not equal wanted (%v)",
			*got.Remove, *wanted.Remove)
	} else {
		t.Logf("Remove: got (%v) == wanted (%v)",
			*got.Remove, *wanted.Remove)
	}

	if *got.IgnoreErrors != *wanted.IgnoreErrors {
		t.Errorf("IgnoreErrors: got (%v) does not equal wanted (%v)",
			*got.IgnoreErrors, *wanted.IgnoreErrors)
	} else {
		t.Logf("IgnoreErrors: got (%v) == wanted (%v)",
			*got.IgnoreErrors, *wanted.IgnoreErrors)
	}

	if *got.RecursiveSearch != *wanted.RecursiveSearch {
		t.Errorf("RecursiveSearch: got (%v) does not equal wanted (%v)",
			*got.RecursiveSearch, *wanted.RecursiveSearch)
	} else {
		t.Logf("RecursiveSearch: got (%v) == wanted (%v)",
			*got.RecursiveSearch, *wanted.RecursiveSearch)
	}

	if *got.LogLevel != *wanted.LogLevel {
		t.Errorf("LogLevel: got (%v) does not equal wanted (%v)",
			*got.LogLevel, *wanted.LogLevel)
	} else {
		t.Logf("LogLevel: got (%v) == wanted (%v)",
			*got.LogLevel, *wanted.LogLevel)
	}

	if *got.LogFormat != *wanted.LogFormat {
		t.Errorf("LogFormat: got (%v) does not equal wanted (%v)",
			*got.LogFormat, *wanted.LogFormat)
	} else {
		t.Logf("LogFormat: got (%v) == wanted (%v)",
			*got.LogFormat, *wanted.LogFormat)
	}

	if *got.LogFilePath != *wanted.LogFilePath {
		t.Errorf("LogFilePath: got (%v) does not equal wanted (%v)",
			*got.LogFilePath, *wanted.LogFilePath)
	} else {
		t.Logf("LogFilePath: got (%v) == wanted (%v)",
			*got.LogFilePath, *wanted.LogFilePath)
	}

	if *got.ConsoleOutput != *wanted.ConsoleOutput {
		t.Errorf("ConsoleOutput: got (%v) does not equal wanted (%v)",
			*got.ConsoleOutput, *wanted.ConsoleOutput)
	} else {
		t.Logf("ConsoleOutput: got (%v) == wanted (%v)",
			*got.ConsoleOutput, *wanted.ConsoleOutput)
	}

	if *got.UseSyslog != *wanted.UseSyslog {
		t.Errorf("UseSyslog: got (%v) does not equal wanted (%v)",
			*got.UseSyslog, *wanted.UseSyslog)
	} else {
		t.Logf("UseSyslog: got (%v) == wanted (%v)",
			*got.UseSyslog, *wanted.UseSyslog)
	}

	if *got.ConfigFile != *wanted.ConfigFile {
		t.Errorf("ConfigFile: got (%v) does not equal wanted (%v)",
			*got.ConfigFile, *wanted.ConfigFile)
	} else {
		t.Logf("ConfigFile: got (%v) == wanted (%v)",
			*got.ConfigFile, *wanted.ConfigFile)
	}

	t.Log("CompareConfig complete")

}
