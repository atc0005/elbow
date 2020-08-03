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
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/alexflint/go-arg"
	"github.com/atc0005/elbow/logging"
)

// TODO: Evaluate replacing bare strings with constants (see constants.go)

// TestMergeConfigUsingIncompleteConfigObjects creates multiple Config structs
// and merges them in sequence, verifying that after each MergeConfig
// operation that the initial config struct has been updated to reflect the
// new state.
func TestMergeConfigUsingIncompleteConfigObjects(t *testing.T) {

	//
	// Base Config testing
	//

	// Validation will fail if this is all we do since the default config
	// doesn't contain any Paths to process.
	baseConfig := NewDefaultConfig("x.y.z")

	testBaseConfigPaths := []string{"/tmp/elbow/path1"}
	baseConfig.Paths = testBaseConfigPaths
	baseConfig.logger = baseConfig.GetLogger()

	// Question: Any reason to set this? The Validate() method does not
	// currently enforce that the FileExtensions field be set.
	// Answer: Yes, because we want to ensure that the final MergeConfig()
	// result reflects our test case(s).
	//
	testBaseConfigFileExtensions := []string{".yaml", ".json"}
	baseConfig.FileExtensions = testBaseConfigFileExtensions

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

	// TODO: This currently mirrors the example config file. Replace with a read
	// against that file?
	var defaultConfigFile = []byte(`

		[filehandling]

		file_age = 90
		files_to_keep = 1
		file_extensions = [
			".war",
		]

		[search]

		recursive_search = true
		paths = [
			"/tmp/elbow/path1",
		]

		[logging]

		log_level = "notice"
		log_format = "text"
		use_syslog = true`)

	// Use our default in-memory config file settings
	r := bytes.NewReader(defaultConfigFile)

	fileConfig := Config{}
	if err := fileConfig.LoadConfigFile(r); err != nil {
		t.Error("Unable to load in-memory configuration:", err)
	} else {
		t.Log("Loaded in-memory configuration file")
	}

	t.Log("LoadConfigFile wiped the existing struct, reconstructing AppMetadata fields to pass validation checks")
	fileConfig.AppName = baseConfig.GetAppName()
	fileConfig.AppVersion = baseConfig.GetAppVersion()
	fileConfig.AppURL = baseConfig.GetAppURL()
	fileConfig.AppDescription = baseConfig.GetAppDescription()

	// NOTE: We cannot validate the fileConfig object at this point because it
	// is a partial object, missing the rest of the settings that are required
	// for a full config validation check

	// Build EXPECTED baseConfig after merge so we can use Compare() against
	// it and the actual baseConfig

	expectedAppNameAfterFileMerge := baseConfig.GetAppName()
	expectedAppDescriptionAfterFileMerge := baseConfig.GetAppDescription()
	expectedAppURLAfterFileMerge := baseConfig.GetAppURL()
	expectedAppVersionAfterFileMerge := baseConfig.GetAppVersion()
	expectedFilePatternAfterFileMerge := baseConfig.GetFilePattern()

	// Explicitly set these; we want to ensure the final merged config has
	// the values we provided (incomplete fileConfig) and the prior baseConfig
	// settings that we are not overriding
	// NOTE: Paths and FileExtensions are set below after config struct is
	// instantiated
	// TODO: Evaluate replacing bare strings with constants (see constants.go)
	expectedPathsAfterFileMerge := []string{"/tmp/elbow/path1"}
	expectedFileExtensionsAfterFileMerge := []string{".war"}
	expectedFileAgeAfterFileMerge := 90
	expectedNumFilesToKeepAfterFileMerge := 1
	expectedRecursiveSearchAfterFileMerge := true
	expectedLogLevelAfterFileMerge := logging.LogLevelNotice
	expectedLogFormatAfterFileMerge := logging.LogFormatText
	expectedUseSyslogAfterFileMerge := true

	expectedKeepOldestAfterFileMerge := baseConfig.GetKeepOldest()
	expectedRemoveAfterFileMerge := baseConfig.GetRemove()
	expectedIgnoreErrorsAfterFileMerge := baseConfig.GetIgnoreErrors()
	expectedLogFilePathAfterFileMerge := baseConfig.GetLogFilePath()
	expectedConsoleOutputAfterFileMerge := baseConfig.GetConsoleOutput()
	expectedConfigFileAfterFileMerge := baseConfig.GetConfigFile()

	expectedBaseConfigAfterFileConfigMerge := Config{
		AppMetadata: AppMetadata{
			AppName:        expectedAppNameAfterFileMerge,
			AppDescription: expectedAppDescriptionAfterFileMerge,
			AppURL:         expectedAppURLAfterFileMerge,
			AppVersion:     expectedAppVersionAfterFileMerge,
		},
		FileHandling: FileHandling{
			FilePattern:    &expectedFilePatternAfterFileMerge,
			FileAge:        &expectedFileAgeAfterFileMerge,
			NumFilesToKeep: &expectedNumFilesToKeepAfterFileMerge,
			KeepOldest:     &expectedKeepOldestAfterFileMerge,
			Remove:         &expectedRemoveAfterFileMerge,
			IgnoreErrors:   &expectedIgnoreErrorsAfterFileMerge,
		},
		Logging: Logging{
			LogLevel:      &expectedLogLevelAfterFileMerge,
			LogFormat:     &expectedLogFormatAfterFileMerge,
			LogFilePath:   &expectedLogFilePathAfterFileMerge,
			ConsoleOutput: &expectedConsoleOutputAfterFileMerge,
			UseSyslog:     &expectedUseSyslogAfterFileMerge,
		},
		Search: Search{
			//Paths: ,
			RecursiveSearch: &expectedRecursiveSearchAfterFileMerge,
		},
		ConfigFile: &expectedConfigFileAfterFileMerge,
		logger:     baseConfig.GetLogger(),
	}

	expectedBaseConfigAfterFileConfigMerge.Paths = expectedPathsAfterFileMerge
	expectedBaseConfigAfterFileConfigMerge.FileExtensions = expectedFileExtensionsAfterFileMerge

	// Validate the expectedBaseConfigAfterFileConfigMerge config settings
	// to ensure that we're not going to compare against a broken configuration
	if err := expectedBaseConfigAfterFileConfigMerge.Validate(); err != nil {
		t.Error("Unable to validate expectedBaseConfigAfterFileConfigMerge before merge:", err)
	} else {
		t.Log("Validation of expectedBaseConfigAfterFileConfigMerge settings before merge successful")
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
	// against the expectedBaseConfigAfterFileConfigMerge struct to determine
	// if the merge results are as expected.

	CompareConfig(baseConfig, expectedBaseConfigAfterFileConfigMerge, t)

	//
	// Environment variables config testing
	//

	// Setup subset of total environment variables for parsing with
	// alexflint/go-arg package. These values should override baseConfig
	// settings
	// TODO: Evaluate replacing bare strings with constants (see constants.go)
	envVarTables := []struct {
		envVar string
		value  string
	}{
		{"ELBOW_FILE_PATTERN", "reach-masterqa-"},
		{"ELBOW_FILE_AGE", "3"},
		{"ELBOW_KEEP", "4"},
		{"ELBOW_KEEP_OLD", "false"},
		{"ELBOW_REMOVE", "true"},
		{"ELBOW_LOG_FORMAT", logging.LogFormatText},
		{"ELBOW_LOG_FILE", "/var/log/elbow/env.log"},
		{"ELBOW_PATHS", "/tmp/elbow/path3"},
		{"ELBOW_EXTENSIONS", ".docx,.pptx"},
	}

	envConfig := Config{}

	t.Log("Explicitly setting AppMetadata fields to pass validation checks")
	envConfig.AppName = baseConfig.GetAppName()
	envConfig.AppVersion = baseConfig.GetAppVersion()
	envConfig.AppURL = baseConfig.GetAppURL()
	envConfig.AppDescription = baseConfig.GetAppDescription()

	for _, table := range envVarTables {
		t.Logf("Setting %q to %q", table.envVar, table.value)
		if err := os.Setenv(table.envVar, table.value); err != nil {
			t.Errorf("Unable to set environment variable: %v", err)
		}
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

	// NOTE: We cannot validate envConfig here since the set of options is
	// incomplete.

	// Build EXPECTED baseConfig after env vars merge so we can use Compare()
	// against it and the actual baseConfig

	expectedAppNameAfterEnvVarsMerge := baseConfig.GetAppName()
	expectedAppDescriptionAfterEnvVarsMerge := baseConfig.GetAppDescription()
	expectedAppURLAfterEnvVarsMerge := baseConfig.GetAppURL()
	expectedAppVersionAfterEnvVarsMerge := baseConfig.GetAppVersion()

	// Explicitly set these; we want to ensure the final merged config has
	// the values we provided (incomplete fileConfig) and the prior baseConfig
	// settings that we are not overriding
	// NOTE: Paths and FileExtensions are set below after config struct is
	// instantiated
	// TODO: Evaluate replacing bare strings with constants (see constants.go)
	expectedPathsAfterEnvVarsMerge := []string{"/tmp/elbow/path3"}
	expectedFileExtensionsAfterEnvVarsMerge := []string{".docx", ".pptx"}
	expectedFilePatternAfterEnvVarsMerge := "reach-masterqa-"
	expectedFileAgeAfterEnvVarsMerge := 3
	expectedNumFilesToKeepAfterEnvVarsMerge := 4
	expectedKeepOldestAfterEnvVarsMerge := false
	expectedRemoveAfterEnvVarsMerge := true
	expectedLogFormatAfterEnvVarsMerge := logging.LogFormatText
	expectedLogFilePathAfterEnvVarsMerge := "/var/log/elbow/env.log"

	expectedRecursiveSearchAfterEnvVarsMerge := baseConfig.GetRecursiveSearch()
	expectedLogLevelAfterEnvVarsMerge := baseConfig.GetLogLevel()
	expectedUseSyslogAfterEnvVarsMerge := baseConfig.GetUseSyslog()
	expectedIgnoreErrorsAfterEnvVarsMerge := baseConfig.GetIgnoreErrors()
	expectedConsoleOutputAfterEnvVarsMerge := baseConfig.GetConsoleOutput()
	expectedConfigFileAfterEnvVarsMerge := baseConfig.GetConfigFile()

	expectedBaseConfigAfterEnvVarsMerge := Config{
		AppMetadata: AppMetadata{
			AppName:        expectedAppNameAfterEnvVarsMerge,
			AppDescription: expectedAppDescriptionAfterEnvVarsMerge,
			AppURL:         expectedAppURLAfterEnvVarsMerge,
			AppVersion:     expectedAppVersionAfterEnvVarsMerge,
		},
		FileHandling: FileHandling{
			FilePattern:    &expectedFilePatternAfterEnvVarsMerge,
			FileAge:        &expectedFileAgeAfterEnvVarsMerge,
			NumFilesToKeep: &expectedNumFilesToKeepAfterEnvVarsMerge,
			KeepOldest:     &expectedKeepOldestAfterEnvVarsMerge,
			Remove:         &expectedRemoveAfterEnvVarsMerge,
			IgnoreErrors:   &expectedIgnoreErrorsAfterEnvVarsMerge,
		},
		Logging: Logging{
			LogLevel:      &expectedLogLevelAfterEnvVarsMerge,
			LogFormat:     &expectedLogFormatAfterEnvVarsMerge,
			LogFilePath:   &expectedLogFilePathAfterEnvVarsMerge,
			ConsoleOutput: &expectedConsoleOutputAfterEnvVarsMerge,
			UseSyslog:     &expectedUseSyslogAfterEnvVarsMerge,
		},
		Search: Search{
			RecursiveSearch: &expectedRecursiveSearchAfterEnvVarsMerge,
		},
		ConfigFile: &expectedConfigFileAfterEnvVarsMerge,
		logger:     baseConfig.GetLogger(),
	}

	expectedBaseConfigAfterEnvVarsMerge.Paths = expectedPathsAfterEnvVarsMerge
	expectedBaseConfigAfterEnvVarsMerge.FileExtensions = expectedFileExtensionsAfterEnvVarsMerge

	// Validate the env vars settings
	if err := expectedBaseConfigAfterEnvVarsMerge.Validate(); err != nil {
		t.Error("Unable to validate expectedBaseConfigAfterEnvVarsMerge before merge:", err)
	} else {
		t.Log("Validation of expectedBaseConfigAfterEnvVarsMerge before merge successful")
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

	CompareConfig(baseConfig, expectedBaseConfigAfterEnvVarsMerge, t)

	// Unset environment variables that we just set
	for _, table := range envVarTables {
		t.Logf("Unsetting %q\n", table.envVar)
		if err := os.Unsetenv(table.envVar); err != nil {
			t.Errorf("Unable to unset environment variable: %v", err)

		}
	}

	//
	// Flags Config testing
	//

	flagsConfig := Config{}

	t.Log("Explicitly setting AppMetadata fields to pass validation checks")
	flagsConfig.AppName = baseConfig.GetAppName()
	flagsConfig.AppVersion = baseConfig.GetAppVersion()
	flagsConfig.AppURL = baseConfig.GetAppURL()
	flagsConfig.AppDescription = baseConfig.GetAppDescription()

	// TODO: A useful way to automate retrieving the app name?
	appName := strings.ToLower(DefaultAppName)
	if runtime.GOOS == WindowsOSName {
		appName += WindowsAppSuffix
	}

	// Note to self: Don't add/escape double-quotes here. The shell strips
	// them away and the application never sees them.
	// TODO: Evaluate replacing bare strings with constants (see constants.go)
	os.Args = []string{
		appName,
		"--pattern", "reach-master-",
		"--age", "5",
		"--keep", "6",
		"--remove",
		"--log-level", logging.LogLevelPanic,
		"--log-format", logging.LogFormatJSON,
		"--console-output", logging.ConsoleOutputStderr,
		"--extensions", ".java", ".class",
	}

	t.Log("Parsing command-line flags")
	arg.MustParse(&flagsConfig)
	t.Logf("Results of parsing flags: %v", flagsConfig.String())

	// NOTE: We cannot validate flagsConfig here since the set of options is
	// incomplete.

	// Build EXPECTED baseConfig after flags merge so we can use Compare()
	// against it and the actual baseConfig

	expectedAppNameAfterFlagsMerge := baseConfig.GetAppName()
	expectedAppDescriptionAfterFlagsMerge := baseConfig.GetAppDescription()
	expectedAppURLAfterFlagsMerge := baseConfig.GetAppURL()
	expectedAppVersionAfterFlagsMerge := baseConfig.GetAppVersion()
	expectedPathsAfterFlagsMerge := baseConfig.GetPaths()
	expectedKeepOldestAfterFlagsMerge := baseConfig.GetKeepOldest()
	expectedLogFilePathAfterFlagsMerge := baseConfig.GetLogFilePath()
	expectedRecursiveSearchAfterFlagsMerge := baseConfig.GetRecursiveSearch()
	expectedUseSyslogAfterFlagsMerge := baseConfig.GetUseSyslog()
	expectedIgnoreErrorsAfterFlagsMerge := baseConfig.GetIgnoreErrors()
	expectedConfigFileAfterFlagsMerge := baseConfig.GetConfigFile()

	// Explicitly set these; we want to ensure the final merged config has
	// the values we provided (incomplete fileConfig) and the prior baseConfig
	// settings that we are not overriding
	// NOTE: Paths and FileExtensions are set below after config struct is
	// instantiated
	// TODO: Evaluate replacing bare strings with constants (see constants.go)
	expectedFileExtensionsAfterFlagsMerge := []string{".java", ".class"}
	expectedFilePatternAfterFlagsMerge := "reach-master-"
	expectedFileAgeAfterFlagsMerge := 5
	expectedNumFilesToKeepAfterFlagsMerge := 6
	expectedRemoveAfterFlagsMerge := true
	expectedLogFormatAfterFlagsMerge := logging.LogFormatJSON
	expectedLogLevelAfterFlagsMerge := logging.LogLevelPanic
	expectedConsoleOutputAfterFlagsMerge := logging.ConsoleOutputStderr

	expectedBaseConfigAfterFlagsMerge := Config{
		AppMetadata: AppMetadata{
			AppName:        expectedAppNameAfterFlagsMerge,
			AppDescription: expectedAppDescriptionAfterFlagsMerge,
			AppURL:         expectedAppURLAfterFlagsMerge,
			AppVersion:     expectedAppVersionAfterFlagsMerge,
		},
		FileHandling: FileHandling{
			FilePattern:    &expectedFilePatternAfterFlagsMerge,
			FileAge:        &expectedFileAgeAfterFlagsMerge,
			NumFilesToKeep: &expectedNumFilesToKeepAfterFlagsMerge,
			KeepOldest:     &expectedKeepOldestAfterFlagsMerge,
			Remove:         &expectedRemoveAfterFlagsMerge,
			IgnoreErrors:   &expectedIgnoreErrorsAfterFlagsMerge,
		},
		Logging: Logging{
			LogLevel:      &expectedLogLevelAfterFlagsMerge,
			LogFormat:     &expectedLogFormatAfterFlagsMerge,
			LogFilePath:   &expectedLogFilePathAfterFlagsMerge,
			ConsoleOutput: &expectedConsoleOutputAfterFlagsMerge,
			UseSyslog:     &expectedUseSyslogAfterFlagsMerge,
		},
		Search: Search{
			RecursiveSearch: &expectedRecursiveSearchAfterFlagsMerge,
		},
		ConfigFile: &expectedConfigFileAfterFlagsMerge,
		logger:     baseConfig.GetLogger(),
	}

	expectedBaseConfigAfterFlagsMerge.Paths = expectedPathsAfterFlagsMerge
	expectedBaseConfigAfterFlagsMerge.FileExtensions = expectedFileExtensionsAfterFlagsMerge

	// Validate the config file settings
	if err := expectedBaseConfigAfterFlagsMerge.Validate(); err != nil {
		t.Error("Unable to validate expectedBaseConfigAfterFlagsMerge before merging:", err)
	} else {
		t.Log("Validation of expectedBaseConfigAfterFlagsMerge before merging successful")
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

	CompareConfig(baseConfig, expectedBaseConfigAfterFlagsMerge, t)

}
