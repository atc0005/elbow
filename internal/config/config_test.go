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
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/atc0005/elbow/internal/units"
	"github.com/sirupsen/logrus"
)

func GetBaseProjectDir(t *testing.T) string {
	t.Helper()

	// https://stackoverflow.com/questions/48570228/get-the-parent-path
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Our project root is two directories up from internal/config.
	dir := filepath.Join(wd, "../..")

	return dir
}

// TODO: Lots of variations here
func TestNewConfigFlagsOnly(t *testing.T) {

	// https://stackoverflow.com/questions/33723300/how-to-test-the-passing-of-arguments-in-golang

	// Save old command-line arguments so that we can restore them later
	oldArgs := os.Args

	// Defer restoring original command-line arguments
	defer func() { os.Args = oldArgs }()

	// TODO: A useful way to automate retrieving the app name?
	appName := strings.ToLower(DefaultAppName)
	if runtime.GOOS == WindowsOSName {
		appName += WindowsAppSuffix
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
	_, err := NewConfig()
	if err != nil {
		t.Errorf("Error encountered when instantiating configuration: %s", err)
	} else {
		t.Log("No errors encountered when instantiating configuration")
	}

}

func TestLoadConfigFileBaseline(t *testing.T) {

	// TODO: This currently mirrors the example config file. Replace with a read
	// against that file?
	var defaultConfigFile = []byte(`
		[filehandling]

		pattern = "reach-masterdev-"
		file_extensions = [
			".war",
			".tmp",
		]

		file_age = 1
		files_to_keep = 2
		keep_oldest = false
		remove = false
		ignore_errors = true

		[search]

		paths = [
			"/tmp/elbow/path1",
			"/tmp/elbow/path2",
		]

		recursive_search = true


		[logging]

		log_level = "debug"
		log_format = "text"

		# If set, all output to the console will be muted and sent here instead
		log_file_path = "/tmp/log.json"

		console_output = "stdout"
		use_syslog = false`)

	// Construct a mostly empty config struct to load our config settings into.
	// We only define values for settings that we have no intention of using
	// from the config file, such as a logger object and an empty path to
	// the log file that we already know the path to.
	defaultConfigFilePath := ""
	c := Config{
		ConfigFile: &defaultConfigFilePath,
		logger:     logrus.New(),
	}

	// Use our default in-memory config file settings
	r := bytes.NewReader(defaultConfigFile)

	if err := c.LoadConfigFile(r); err != nil {
		t.Error("Unable to load in-memory configuration:", err)
	} else {
		t.Log("Loaded in-memory configuration file")
	}

	t.Log("LoadConfigFile wiped the existing struct, reconstructing AppMetadata fields to pass validation checks")
	c.AppName = c.GetAppName()
	c.AppVersion = c.GetAppVersion()
	c.AppURL = c.GetAppURL()
	c.AppDescription = c.GetAppDescription()

	t.Logf("Current config settings: %s", c.String())

	if err := c.Validate(); err != nil {
		t.Error("Unable to validate configuration:", err)
	} else {
		t.Log("Validation successful")
	}

}

// This function is intended to test the example config file included in the
// repo. That example config file is a template of valid settings, so we
// should run tests against it to verify everything checks out.
func TestLoadConfigFileTemplate(t *testing.T) {

	// Construct a mostly empty config struct to load our config settings into.
	// We only define values for settings that we have no intention of using
	// from the config file, such as a logger object and an empty path to
	// the log file that we already know the path to.
	defaultConfigFilePath := ""
	c := Config{
		ConfigFile: &defaultConfigFilePath,
		logger:     logrus.New(),
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("Current working directory:", cwd)

	// this file is located in the base of the repo
	exampleConfigFile := "config.example.toml"

	// Get path above cwd in order to load config file (from base path of repo)
	baseDir := GetBaseProjectDir(t)
	if err := os.Chdir(baseDir); err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Log("New working directory:", cwd)
	}

	if fileDetails, err := os.Stat(exampleConfigFile); os.IsNotExist(err) {
		t.Errorf("requested config file not found: %v", err)
	} else {
		t.Log("config file found")
		t.Log("name:", fileDetails.Name())
		t.Log("size:", units.ByteCountSI(fileDetails.Size()))
		t.Log("date/time stamp:", fileDetails.ModTime())

		fh, err := os.Open(exampleConfigFile)
		if err != nil {
			t.Errorf("Unable to open config file: %v", err)
		} else {
			t.Log("Successfully opened config file", exampleConfigFile)
		}

		// #nosec G307
		// Believed to be a false-positive from recent gosec release
		// https://github.com/securego/gosec/issues/714
		defer func() {
			if err := fh.Close(); err != nil {
				// Ignore "file already closed" errors
				if !errors.Is(err, os.ErrClosed) {
					t.Errorf(
						"failed to close file %q: %s",
						exampleConfigFile,
						err.Error(),
					)
				}
			}
		}()

		if err := c.LoadConfigFile(fh); err != nil {
			t.Error("Unable to load configuration file:", err)
		} else {
			t.Log("Loaded configuration file")
		}

		t.Log("LoadConfigFile wiped the existing struct, reconstructing AppMetadata fields to pass validation checks")
		c.AppName = c.GetAppName()
		c.AppVersion = c.GetAppVersion()
		c.AppURL = c.GetAppURL()
		c.AppDescription = c.GetAppDescription()

		t.Logf("Current config settings: %s", c.String())

		if err := c.Validate(); err != nil {
			t.Error("Unable to validate configuration:", err)
		} else {
			t.Log("Validation successful")
		}

	}
}
