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

package main

import (
	"os"
	"runtime"
	"testing"

	"github.com/atc0005/elbow/internal/config"
	"github.com/atc0005/elbow/internal/logging"
)

func TestMain(t *testing.T) {

	// https://stackoverflow.com/questions/33723300/how-to-test-the-passing-of-arguments-in-golang

	// Save old command-line arguments so that we can restore them later
	oldArgs := os.Args

	// Defer restoring original command-line arguments
	defer func() { os.Args = oldArgs }()

	// TODO: A useful way to automate retrieving the app name?
	appName := config.DefaultAppName
	if runtime.GOOS == config.WindowsOSName {
		appName += config.WindowsAppSuffix
	}

	// Note to self: Don't add/escape double-quotes here. The shell strips
	// them away and the application never sees them.
	// TODO: Evaluate replacing bare strings with constants (see constants.go)
	os.Args = []string{
		appName,
		"--paths", "/tmp/elbow/path1",
		"--keep", "1",
		"--recurse",
		"--keep-old",
		"--log-level", logging.LogLevelInfo,
		"--use-syslog",
		"--log-format", logging.LogFormatText,
		"--console-output", logging.ConsoleOutputStdout,
	}

	// TODO: Flesh this out
	_, err := config.NewConfig()
	if err != nil {
		t.Errorf("Error encountered when instantiating configuration: %s", err)
	} else {
		t.Log("No errors encountered when instantiating configuration")
	}

}
