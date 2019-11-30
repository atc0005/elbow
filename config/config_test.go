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
	"bytes"
	"github.com/atc0005/elbow/units"
	"os"
	"path"
	"runtime"
	"testing"
)

var defaultConfigFile = []byte(`
[appmetadata]

app_name = "toml_app_name"
app_description = "toml_app_description"
app_version = "toml_app_version"
app_url = "toml_app_url"

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
#log_file_path = "/tmp/log.json"

console_output = "stdout"
use_syslog = false`)

func GetBaseProjectDir(t *testing.T) string {

	// https://stackoverflow.com/questions/23847003/golang-tests-and-working-directory
	_, filename, _, _ := runtime.Caller(1)
	// The ".." reflects the path above the current working directory
	dir := path.Join(path.Dir(filename), "..")
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

// TODO: Need to ensure that the configuration was loaded properly.
func TestLoadConfigOnTopOfBaseConfig(t *testing.T) {

	c := NewDefaultConfig("x.y.z")

	// testPaths := []string{"/tmp/elbow/path1"}
	// testFileExtensions := []string{".tmp", ".war"}

	// c.Paths = testPaths
	// c.FileExtensions = testFileExtensions
	c.logger = c.GetLogger()

	// Use stock configuration
	r := bytes.NewReader(defaultConfigFile)

	if err := c.LoadConfigFile(r); err != nil {
		t.Error("Unable to load in-memory configuration:", err)
	} else {
		t.Log("Loaded in-memory configuration file")
	}

	if err := c.Validate(); err != nil {
		t.Error("Unable to validate configuration:", err)
	} else {
		t.Log("Validation successful")
	}

}

// This function is intended to test the example config file included in the
// repo. Since it is intended to reflect a template of valid settings, we
// should run tests against it to verify everything checks out.
func TestLoadConfigFileExampleConfigInRepo(t *testing.T) {

	c := NewDefaultConfig("x.y.z")

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
	}

	// this isn't handled by the config file settings and is ordinarily taken
	// care of by the time the config file is pulled in
	c.logger = c.GetLogger()

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
		defer fh.Close()

		if err := c.LoadConfigFile(fh); err != nil {
			t.Error("Unable to load configuration file:", err)
		} else {
			t.Log("Loaded configuration file")
		}

		if err := c.Validate(); err != nil {
			t.Error("Unable to validate configuration:", err)
		} else {
			t.Log("Validation successful")
		}

	}
}
