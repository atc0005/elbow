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

// NOTE: This file was created as a way of hotfixing the goconst linting errors
// as part of resolving GH-176. This file is also intended to help
// start a more thorough exploration of whether Go has a similar concept to
// collecting constants in header files like C/C++ does.
//
// TODO: Evaluate use of constants throughout the entire codebase

// Workarounds for golantci-lint errors:
// string `STRING` has N occurrences, make it a constant (goconst)
const (
	FakeValue        string = "fakeValue"
	FieldValueNotSet string = "NotSet"
	WindowsOSName    string = "windows"
	WindowsAppSuffix string = ".exe"
)

const (

	// DefaultAppName is the default name for this application
	DefaultAppName string = "elbow"

	// DefaultAppDescription is the description for this application shown in
	// HelpText output.
	DefaultAppDescription string = "prunes content matching specific patterns, either in a single directory or recursively through a directory tree."

	// DefaultAppVersion is a placeholder that is used when the application is
	// compiled without the use of the Makefile (which handles setting the
	// value via a build tag)
	//
	// NOTE: We use an unexported variable that is set via Makefile builds
	// instead of using a constant.
	//
	// DefaultAppVersion string = "dev build"

	// DefaultAppURL is the website where users can learn more about the
	// application, submit problem reports, etc.
	DefaultAppURL string = "https://github.com/atc0005/elbow"
)
