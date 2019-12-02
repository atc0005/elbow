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
		t.Log("Merged config file settings with base config successfully")
	}

	// Validate the base config settings after merging
	if err := baseConfig.Validate(); err != nil {
		t.Error("Unable to validate base configuration after merge:", err)
	} else {
		t.Log("Validation of base config settings after merge successful")
	}

	// TODO: We've confirmed that we end up with a config struct that passes
	// field validation, but we've not yet confirmed that field values in
	// the base config are overwritten for any non-nil field in the file
	// configuration. We should experiment by creating at least three
	// variations of this test:
	//
	// 1) in-memory complete file config (this test)
	// 2) template config file (complete)
	// 3) various partial in-memory file configurations
	//
	// This is test 1.
	// TODO: Create tests 2 & 3; test 3 can be composed of sub or table tests

	// This is where we compare the field values of the baseConfig struct
	// against the fileConfig struct to determine which are different

	// configReflect := reflect.TypeOf(baseConfig)

	// var field reflect.StructField
	// var ok bool

	// if field, ok = configReflect.FieldByName("FilePattern"); !ok {
	// 	t.Error("unable to reference struct field")
	// } else {
	// 	t.Log("FilePattern value is", field)
	// }

}
