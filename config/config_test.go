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
func TestValidate(t *testing.T) {

	// Create and initialize struct
	// One at a time, set test-target fields to nil
	// Validate config struct
	// Set field back to good value

	c := NewDefaultConfig("x.y.z")

	testPaths := []string{"/tmp/elbow/path1"}
	testFileExtensions := []string{".tmp", ".war"}

	c.Paths = testPaths
	c.FileExtensions = testFileExtensions
	c.logger = c.GetLogger()

	t.Run("Validating mostly default config", func(t *testing.T) {
		// This should pass
		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config created from NewDefaultConfig(): %s", err)
		} else {
			t.Log("No errors encountered when instantiating configuration")
		}
	})

	t.Run("AppName set to nil", func(t *testing.T) {
		tmpAppName := c.AppName
		c.AppName = nil
		if err := c.Validate(); err == nil {
			t.Errorf("Config passed, but should have failed on nil AppName: %s", err)
		} else {
			t.Logf("Config failed as expected after setting AppName to nil: %s", err)
		}
		// Set back to prior value
		c.AppName = tmpAppName

		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config after restoring AppName: %s", err)
		} else {
			t.Log("Validation successful after restoring AppName field")
		}
	})

	t.Run("AppDescription set to nil", func(t *testing.T) {
		tmpAppDescription := c.AppDescription
		c.AppDescription = nil
		if err := c.Validate(); err == nil {
			t.Errorf("Config passed, but should have failed on nil AppDescription: %s", err)
		} else {
			t.Logf("Config failed as expected after setting AppDescription to nil: %s", err)
		}
		// Set back to prior value
		c.AppDescription = tmpAppDescription

		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config after restoring AppDescription: %s", err)
		} else {
			t.Log("Validation successful after restoring AppDescription field")
		}
	})

	t.Run("AppVersion set to nil", func(t *testing.T) {
		tmpAppVersion := c.AppVersion
		c.AppVersion = nil
		if err := c.Validate(); err == nil {
			t.Errorf("Config passed, but should have failed on nil AppVersion: %s", err)
		} else {
			t.Logf("Config failed as expected after setting AppVersion to nil: %s", err)
		}
		// Set back to prior value
		c.AppVersion = tmpAppVersion

		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config after restoring AppVersion: %s", err)
		} else {
			t.Log("Validation successful after restoring AppVersion field")
		}
	})

	t.Run("AppURL set to nil", func(t *testing.T) {
		tmpAppURL := c.AppURL
		c.AppURL = nil
		if err := c.Validate(); err == nil {
			t.Errorf("Config passed, but should have failed on nil AppURL: %s", err)
		} else {
			t.Logf("Config failed as expected after setting AppURL to nil: %s", err)
		}
		// Set back to prior value
		c.AppURL = tmpAppURL

		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config after restoring AppURL: %s", err)
		} else {
			t.Log("Validation successful after restoring AppURL field")
		}
	})

	t.Run("FilePattern set to nil", func(t *testing.T) {
		tmpFilePattern := c.FilePattern
		//t.Logf("c.FilePattern before setting to nil: %p", c.FilePattern)
		c.FilePattern = nil
		if err := c.Validate(); err == nil {
			t.Errorf("Config passed, but should have failed on nil FilePattern: %s", err)
		} else {
			t.Logf("Config failed as expected after setting FilePattern to nil: %s", err)
		}
		// Set back to prior value
		c.FilePattern = tmpFilePattern
		//t.Logf("c.FilePattern address after resetting back to original value: %p", c.FilePattern)

		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config after restoring FilePattern: %s", err)
		} else {
			t.Log("Validation successful after restoring FilePattern field")
		}
	})

	t.Run("Paths set to nil", func(t *testing.T) {
		tmpPaths := c.Paths
		c.Paths = nil
		if err := c.Validate(); err == nil {
			t.Errorf("Config passed, but should have failed on nil Paths: %s", err)
		} else {
			t.Logf("Config failed as expected after setting Paths to nil: %s", err)
		}
		// Set back to prior value
		c.Paths = tmpPaths

		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config after restoring Paths: %s", err)
		} else {
			t.Log("Validation successful after restoring Paths field")
		}
	})

	t.Run("RecursiveSearch set to nil", func(t *testing.T) {
		tmpRecursiveSearch := c.RecursiveSearch
		c.RecursiveSearch = nil
		if err := c.Validate(); err == nil {
			t.Errorf("Config passed, but should have failed on nil RecursiveSearch: %s", err)
		} else {
			t.Logf("Config failed as expected after setting RecursiveSearch to nil: %s", err)
		}
		// Set back to prior value
		c.RecursiveSearch = tmpRecursiveSearch

		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config after restoring RecursiveSearch: %s", err)
		} else {
			t.Log("Validation successful after restoring RecursiveSearch field")
		}
	})

	t.Run("NumFilesToKeep set to nil", func(t *testing.T) {
		tmpNumFilesToKeep := c.NumFilesToKeep
		c.NumFilesToKeep = nil
		if err := c.Validate(); err == nil {
			t.Errorf("Config passed, but should have failed on nil NumFilesToKeep: %s", err)
		} else {
			t.Logf("Config failed as expected after setting NumFilesToKeep to nil: %s", err)
		}
		// Set back to prior value
		c.NumFilesToKeep = tmpNumFilesToKeep

		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config after restoring NumFilesToKeep: %s", err)
		} else {
			t.Log("Validation successful after restoring NumFilesToKeep field")
		}
	})

	t.Run("NumFilesToKeep set to invalid value", func(t *testing.T) {
		tmpNumFilesToKeep := *c.NumFilesToKeep
		*c.NumFilesToKeep = -1
		if err := c.Validate(); err == nil {
			t.Errorf("Config passed, but should have failed on invalid value %d for NumFilesToKeep: %s", *c.NumFilesToKeep, err)
		} else {
			t.Logf("Config failed as expected after setting NumFilesToKeep to %d: %s", *c.NumFilesToKeep, err)
		}
		// Set back to prior value
		*c.NumFilesToKeep = tmpNumFilesToKeep

		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config after restoring NumFilesToKeep: %s", err)
		} else {
			t.Log("Validation successful after restoring NumFilesToKeep field")
		}
	})

	t.Run("FileAge set to nil", func(t *testing.T) {
		tmpFileAge := c.FileAge
		c.FileAge = nil
		if err := c.Validate(); err == nil {
			t.Errorf("Config passed, but should have failed on nil FileAge: %s", err)
		} else {
			t.Logf("Config failed as expected after setting FileAge to nil: %s", err)
		}
		// Set back to prior value
		c.FileAge = tmpFileAge

		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config after restoring FileAge: %s", err)
		} else {
			t.Log("Validation successful after restoring FileAge field")
		}
	})

	t.Run("FileAge set to invalid value", func(t *testing.T) {
		tmpFileAge := *c.FileAge
		*c.FileAge = -1
		if err := c.Validate(); err == nil {
			t.Errorf("Config passed, but should have failed on invalid value %d for FileAge: %s", *c.FileAge, err)
		} else {
			t.Logf("Config failed as expected after setting FileAge to %d: %s", *c.FileAge, err)
		}
		// Set back to prior value
		*c.FileAge = tmpFileAge

		if err := c.Validate(); err != nil {
			t.Errorf("Validation failed for config after restoring FileAge: %s", err)
		} else {
			t.Log("Validation successful after restoring FileAge field")
		}
	})

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
