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
	"testing"

	"github.com/sirupsen/logrus"
)

// TestMergeConfig creates multiple Config structs and merges them in
// sequence, verifying that after each MergeConfig operation that the initial
// config struct has been updated to reflect the new state.
func TestMergeConfig(t *testing.T) {

	// Validation will fail if this is all we do since the default config
	// doesn't contain any Paths to process.
	baseConfig := NewDefaultConfig("x.y.z")

	testPaths := []string{"/tmp/elbow/path1"}
	baseConfig.Paths = testPaths
	baseConfig.logger = baseConfig.GetLogger()

	// TODO: Any reason to set this? The Validate() method does not currently
	// enforce that the FileExtensions field be set.
	// Answer: Yes, because we want to ensure that the final MergeConfig()
	// result reflects our test case(s).
	//
	testFileExtensions := []string{".yaml", ".json"}
	baseConfig.FileExtensions = testFileExtensions

	// Validate the base config settings before making further changes that
	// could potentially break the configuration.
	if err := baseConfig.Validate(); err != nil {
		t.Error("Unable to validate base configuration before merge:", err)
	} else {
		t.Log("Validation of base config settings before merge successful")
	}

	defaultConfigFilePath := "/usr/local/etc/elbow/config.toml"
	fileConfig := Config{

		// This is required as well to pass validation checks
		ConfigFile: &defaultConfigFilePath,

		// Not going to merge this in, but we have to specify it in order to
		// pass validation checks.
		logger: logrus.New(),
	}

	// TODO: This currently mirrors the example config file. Replace with a read
	// against that file?
	var defaultConfigFile = []byte(`
		[appmetadata]

		app_name = "toml_app_name"
		app_description = "toml_app_description"
		app_version = "toml_app_version"
		app_url = "toml_app_url"


		[filehandling]

		file_age = 1
		files_to_keep = 2
		keep_oldest = true
		remove = true
		ignore_errors = true
		pattern = "reach-masterdev-"
		file_extensions = [
			".war",
			".tmp",
		]


		[search]

		recursive_search = true
		paths = [
			"/tmp/elbow/path1",
			"/tmp/elbow/path2",
		]


		[logging]

		log_level = "debug"
		log_format = "json"
		log_file_path = "/var/log/elbow.log"
		console_output = "stderr"
		use_syslog = true`)

	// Use our default in-memory config file settings
	r := bytes.NewReader(defaultConfigFile)

	if err := fileConfig.LoadConfigFile(r); err != nil {
		t.Error("Unable to load in-memory configuration:", err)
	} else {
		t.Log("Loaded in-memory configuration file")
	}

	// Validate the config file settings
	if err := fileConfig.Validate(); err != nil {
		t.Error("Unable to validate file config:", err)
	} else {
		t.Log("Validation of file config settings successful")
	}

	if err := MergeConfig(&baseConfig, fileConfig); err != nil {
		t.Errorf("Error merging config file settings with base config: %s", err)
	} else {
		t.Log("Merge of config file settings with base config successful")
	}

	// Validate the base config settings after merging
	if err := baseConfig.Validate(); err != nil {
		t.Error("Unable to validate base configuration after merge:", err)
	} else {
		t.Log("Validation of base config settings after merge successful")
	}

	// This is where we compare the field values of the baseConfig struct
	// against the fileConfig struct to determine if any are different. In
	// normal use of this application it is likely that the fields WOULD be
	// different, but in our test case we have explicitly defined most fields
	// of each config struct to have conflicting values. This allows us to
	// simply our test case(s) so that we can assume each field has a value
	// that should be compared and merged.

	CompareConfig(baseConfig, fileConfig, t)

	// TODO: Create Env var config by way of presetting env vars and then using
	// these lines to parse them and construct a struct object:
	// envArgsConfig := Config{}
	// arg.MustParse(&envArgsConfig)
	//
	// TODO: Merge that new config struct into the baseConfig struct
	// TODO: Compare the two structs
	//
	// TODO: Create an os.Args slice with all desired flags
	// TODO: Parse the flags
	// TODO: Merge the config structs
	// TODO: Compare the two structs

}
