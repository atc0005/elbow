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

func TestMergeConfigFileIntoBaseConfig(t *testing.T) {

	// Validation will fail if this is all we do since the default config
	// doesn't contain any Paths to process.
	baseConfig := NewDefaultConfig("x.y.z")

	testPaths := []string{"/tmp/elbow/path1"}
	baseConfig.Paths = testPaths
	baseConfig.logger = baseConfig.GetLogger()

	// TODO: Any reason to set this? The Validate() method does not currently
	// enforce that the FileExtensions field be set.
	//
	// testFileExtensions := []string{".tmp", ".war"}
	// baseConfig.FileExtensions = testFileExtensions

	// Validate the base config settings before making further changes that
	// could potentially break the configuration.
	if err := baseConfig.Validate(); err != nil {
		t.Error("Unable to validate base configuration before merge:", err)
	} else {
		t.Log("Validation of base config settings before merge successful")
	}

	defaultConfigFilePath := ""
	fileConfig := Config{
		ConfigFile: &defaultConfigFilePath,
		logger:     logrus.New(),
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
		t.Log("Merged config file settings with base config successfully")
	}

	// Validate the base config settings after merging
	if err := baseConfig.Validate(); err != nil {
		t.Error("Unable to validate base configuration after merge:", err)
	} else {
		t.Log("Validation of base config settings after merge successful")
	}
}
