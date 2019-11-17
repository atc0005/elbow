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

// NOTE: This file was created as a way of hotfixing the goconst linting errors
// as part of resolving GH-176. This file is also intended to help
// start a more thorough exploration of whether Go has a similar concept to
// collecting constants in header files like C/C++ does.
//
// TODO: Evaluate use of constants throughout the entire codebase

// Workarounds for golantci-lint errors:
// string `STRING` has N occurrences, make it a constant (goconst)
const fakeValue = "fakeValue"
const defaultAppName string = "Elbow"
const fieldValueNotSet string = "NotSet"
const windowsOSName string = "windows"
const logFormatText string = "text"
const logFormatJSON string = "json"
const windowsAppSuffix string = ".exe"
