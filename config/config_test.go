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
	"runtime"
	"testing"
)

// TODO: Lots of variations here
func TestNewConfigFlagsOnly(t *testing.T) {

	// https://stackoverflow.com/questions/33723300/how-to-test-the-passing-of-arguments-in-golang

	// Save old command-line arguments so that we can restore them later
	oldArgs := os.Args

	// Defer restoring original command-line arguments
	defer func() { os.Args = oldArgs }()

	// TODO: A useful way to automate retrieving the app name?
	appName := "elbow"
	if runtime.GOOS == "windows" {
		appName += ".exe"
	}

	// Note to self: Don't add/escape double-quotes here. The shell strips
	// them away and the application never sees them.
	os.Args = []string{
		appName,
		"--paths", "/tmp/elbow/path1",
		"--keep", "1",
		"--recurse",
		"--keep-old",
		"--log-level", "info",
		"--use-syslog",
		"--log-format", "text",
		"--console-output", "stdout",
	}

	// TODO: Flesh this out
	_, err := NewConfig("x.y.z")
	if err != nil {
		t.Errorf("Error encountered when instantiating configuration: %s", err)
	} else {
		t.Log("No errors encountered when instantiating configuration")
	}

}

// FIXME: Needs a better name. This is *mostly* a default config struct with
// the addition of config.Paths and config.FileExtensions fields set to fill
// out the set.
func TestValidateFields(t *testing.T) {

	// Create struct
	c := NewDefaultConfig("x.y.z")

	testPaths := []string{"/tmp/elbow/path1"}
	testFileExtensions := []string{".tmp", ".war"}

	c.Paths = testPaths
	c.FileExtensions = testFileExtensions
	c.logger = c.GetLogger()

	// This should pass
	if err := c.Validate(); err != nil {
		t.Errorf("Validation failed for config created from NewDefaultConfig(): %s", err)
	} else {
		t.Log("No errors encountered when instantiating configuration")
	}

	// One at a time, set test-target fields to nil
	// Validate config struct
	// Set field back to good value

}

// func TestValidateDefaultConfig(t *testing.T) {

// 	// Create struct
// 	// Ensure all test-target fields are nil
// 	// Fail if any are missed by Validate() method

// 	c := NewDefaultConfig("x.y.z")

// 	if err := c.Validate(); err != nil {
// 		t.Errorf("Error encountered when instantiating configuration: %s", err)
// 	} else {
// 		t.Log("No errors encountered when instantiating configuration")
// 	}

// }
