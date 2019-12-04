package config

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
)

// MergeConfigTest creates multiple Config structs and merges them in
// sequence, verifying that after each MergeConfig operation that the initial
// config struct has been updated to reflect the new state.
func MergeConfigTest(t *testing.T) {

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
		fmt.Println("Unable to validate base configuration before merge:", err)
	} else {
		fmt.Printf("Validation of base config settings before merge successful\n")
	}

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
		fmt.Println("Unable to load in-memory configuration:", err)
	} else {
		fmt.Printf("Loaded in-memory configuration file\n")
	}

	// Validate the config file settings
	if err := fileConfig.Validate(); err != nil {
		fmt.Println("Unable to validate file config:", err)
	} else {
		fmt.Printf("Validation of file config settings successful\n")
	}

	if err := MergeConfig(&baseConfig, fileConfig); err != nil {
		fmt.Printf("Error merging config file settings with base config: %s\n", err)
	} else {
		fmt.Printf("Merge of config file settings with base config successful\n")
	}

	// Validate the base config settings after merging
	if err := baseConfig.Validate(); err != nil {
		fmt.Println("Unable to validate base configuration after merge:", err)
	} else {
		fmt.Printf("Validation of base config settings after merge successful\n")
	}

	// This is where we compare the field values of the baseConfig struct
	// against the fileConfig struct to determine if any are different. In
	// normal use of this application it is likely that the fields WOULD be
	// different, but in our test case we have explicitly defined most fields
	// of each config struct to have conflicting values. This allows us to
	// simply our test case(s) so that we can assume each field has a value
	// that should be compared and merged.

	CompareConfig(baseConfig, fileConfig, t)

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
		{"ELBOW_EXTENSIONS", ".docx, .pptx"},
	}

	for _, table := range envVarTables {
		fmt.Printf("Setting %q to %q\n", table.envVar, table.value)
		os.Setenv(table.envVar, table.value)
	}

	fmt.Println("Parsing environment variables")
	arg.MustParse(&envConfig)
	fmt.Printf("Results of parsing environment variables: %v\n", envConfig.String())

	// Validate the config file settings
	if err := envConfig.Validate(); err != nil {
		fmt.Println("Unable to validate environment vars config:", err)
	} else {
		fmt.Print("Validation of environment vars config settings successful")
	}

	if err := MergeConfig(&baseConfig, envConfig); err != nil {
		fmt.Printf("Error merging environment vars config settings with base config: %s\n", err)
	} else {
		fmt.Println("Merge of environment vars config settings with base config successful")
	}

	// Validate the base config settings after merging
	if err := baseConfig.Validate(); err != nil {
		fmt.Println("Unable to validate base configuration after merge:", err)
	} else {
		fmt.Println("Validation of base config settings after merge successful")
	}

	CompareConfig(baseConfig, envConfig, t)

	// TODO: Create an os.Args slice with all desired flags
	// TODO: Parse the flags
	// TODO: Merge the config structs
	// TODO: Compare the two structs

}
