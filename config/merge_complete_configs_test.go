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
	"os"
	"runtime"
	"testing"

	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
)

// TestMergeConfigUsingCompletConfigObjects creates multiple Config structs
// and merges them in sequence, verifying that after each MergeConfig
// operation that the initial config struct has been updated to reflect the
// new state.
func TestMergeConfigUsingCompletConfigObjects(t *testing.T) {

	//
	// Base Config testing
	//

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

	//
	// File Config testing
	//

	// ConfigFilePath for use with fileConfig struct tests
	fcConfigFilePath := "/usr/local/etc/elbow/config.toml"
	fileConfig := Config{

		// This is required as well to pass validation checks
		ConfigFile: &fcConfigFilePath,

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

	//
	// Environment variables config testing
	//

	// Setup environment variables for parsing with alexflint/go-arg package

	evConfigFilePath := ""

	// Bolt these on directly as we're likely going to abandon support for
	// overriding these values anyway (haven't been able to come up with a
	// legitimate reason why others would need or want to do so)
	evAppName := "ElbowEnvVar"
	evAppDescription := "something nifty here"
	evAppURL := "https://example.com"
	evAppVersion := "x.y.z"

	envConfig := Config{

		// See earlier notes
		AppMetadata: AppMetadata{
			AppName:        &evAppName,
			AppDescription: &evAppDescription,
			AppURL:         &evAppURL,
			AppVersion:     &evAppVersion,
		},

		// This is required as well to pass validation checks
		ConfigFile: &evConfigFilePath,

		// Not going to merge this in, but we have to specify it in order to
		// pass validation checks.
		logger: logrus.New(),
	}

	envVarTables := []struct {
		envVar string
		value  string
	}{
		{"ELBOW_FILE_PATTERN", "reach-masterqa-"},
		{"ELBOW_FILE_AGE", "3"},
		{"ELBOW_KEEP", "4"},
		{"ELBOW_KEEP_OLD", "false"},
		{"ELBOW_REMOVE", "false"},
		{"ELBOW_IGNORE_ERRORS", "false"},
		{"ELBOW_RECURSE", "false"},
		{"ELBOW_LOG_LEVEL", "warn"},
		{"ELBOW_LOG_FORMAT", "text"},
		{"ELBOW_LOG_FILE", "/var/log/elbow/env.log"},
		{"ELBOW_CONSOLE_OUTPUT", "stdout"},
		{"ELBOW_USE_SYSLOG", "false"},
		{"ELBOW_CONFIG_FILE", "/tmp/config.toml"},
		{"ELBOW_PATHS", "/tmp/elbow/path3"},
		{"ELBOW_EXTENSIONS", ".docx,.pptx"},
	}

	for _, table := range envVarTables {
		t.Logf("Setting %q to %q", table.envVar, table.value)
		os.Setenv(table.envVar, table.value)
	}

	// https://stackoverflow.com/questions/33723300/how-to-test-the-passing-of-arguments-in-golang
	// Save old command-line arguments so that we can restore them later
	oldArgs := os.Args

	// Defer restoring original command-line arguments
	defer func() { os.Args = oldArgs }()

	// Wipe existing command-line arguments so that the unexpected testing
	// package flags don't trip alexflint/go-arg package logic for invalid
	// flags.
	// https://github.com/alexflint/go-arg/issues/97#issuecomment-561995206
	os.Args = nil

	t.Log("Parsing environment variables")
	arg.MustParse(&envConfig)
	t.Logf("Results of parsing environment variables: %v", envConfig.String())

	// Validate the config file settings
	if err := envConfig.Validate(); err != nil {
		t.Error("Unable to validate environment vars config:", err)
	} else {
		t.Log("Validation of environment vars config settings successful")
	}

	if err := MergeConfig(&baseConfig, envConfig); err != nil {
		t.Errorf("Error merging environment vars config settings with base config: %s", err)
	} else {
		t.Log("Merge of environment vars config settings with base config successful")
	}

	// Validate the base config settings after merging
	if err := baseConfig.Validate(); err != nil {
		t.Error("Unable to validate base configuration after merge:", err)
	} else {
		t.Log("Validation of base config settings after merge successful")
	}

	CompareConfig(baseConfig, envConfig, t)

	// Unset environment variables that we just set
	for _, table := range envVarTables {
		t.Logf("Unsetting %q\n", table.envVar)
		os.Unsetenv(table.envVar)
	}

	//
	// Flags Config testing
	//

	// Bolt these on directly as we're likely going to abandon support for
	// overriding these values anyway (haven't been able to come up with a
	// legitimate reason why others would need or want to do so)
	flagsAppName := "ElbowFlagVar"
	flagsAppDescription := "nothing fancy"
	flagsAppURL := "https://example.org"
	flagsAppVersion := "0.1.2"

	flagsConfigFilePath := ""

	flagsConfig := Config{

		// See earlier notes
		AppMetadata: AppMetadata{
			AppName:        &flagsAppName,
			AppDescription: &flagsAppDescription,
			AppURL:         &flagsAppURL,
			AppVersion:     &flagsAppVersion,
		},

		// This is required as well to pass validation checks
		ConfigFile: &flagsConfigFilePath,

		// Not going to merge this in, but we have to specify it in order to
		// pass validation checks.
		logger: logrus.New(),
	}

	// TODO: A useful way to automate retrieving the app name?
	appName := "elbow"
	if runtime.GOOS == "windows" {
		appName += ".exe"
	}

	// Note to self: Don't add/escape double-quotes here. The shell strips
	// them away and the application never sees them.
	os.Args = []string{
		appName,
		"--paths", "/tmp/elbow/path4",
		"--pattern", "reach-master-",
		"--age", "5",
		"--keep", "6",
		"--remove",
		"--ignore-errors",
		"--recurse",
		"--keep-old",
		"--log-level", "info",
		"--use-syslog",
		"--log-format", "json",
		"--console-output", "stderr",
		"--log-file", "/var/log/elbow/flags.log",
		"--config-file", "/tmp/configfile.toml",
		"--extensions", ".java", ".class",
	}

	t.Log("Parsing command-line flags")
	arg.MustParse(&flagsConfig)
	t.Logf("Results of parsing flags: %v", flagsConfig.String())

	// Validate the config file settings
	if err := flagsConfig.Validate(); err != nil {
		t.Error("Unable to validate flags config:", err)
	} else {
		t.Log("Validation of flags config settings successful")
	}

	if err := MergeConfig(&baseConfig, flagsConfig); err != nil {
		t.Errorf("Error merging flags config settings with base config: %s", err)
	} else {
		t.Log("Merge of flags config settings with base config successful")
	}

	// Validate the base config settings after merging
	if err := baseConfig.Validate(); err != nil {
		t.Error("Unable to validate base configuration after merge:", err)
	} else {
		t.Log("Validation of base config settings after merge successful")
	}

	CompareConfig(baseConfig, flagsConfig, t)

}
