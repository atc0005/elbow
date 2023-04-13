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

// NOTE: This file is a "box on the shelf" for potential later use, but should
// be considered earmarked for retirement. Content in this file is subject to
// removal at any time.

package config

import (
	"reflect"
)

// GetStructTag returns the requested struct tag value, if set, and an error
// value indicating whether any problems were encountered.
func GetStructTag(c Config, fieldname string, tagName string) (string, bool) {

	t := reflect.TypeOf(c)

	var field reflect.StructField
	var ok bool
	var tagValue string

	if field, ok = t.FieldByName(fieldname); !ok {
		return "", false
	}

	if tagValue, ok = field.Tag.Lookup(tagName); !ok {
		return "", false
	}

	return tagValue, true

}

// SetDefaultConfig applies application default values to Config object fields
// TODO: Is this still needed? NewDefaultConfig() is handling this now?
func (c *Config) SetDefaultConfig() {

	// These fields are intentionally ignored
	// FileExtensions
	// Paths

	// TODO: Create default logger object?

	c.AppName = c.GetAppName()
	c.AppDescription = c.GetAppDescription()
	c.AppURL = c.GetAppURL()
	c.AppVersion = c.GetAppVersion()
	*c.FilePattern = c.GetFilePattern()
	*c.FileAge = c.GetFileAge()
	*c.NumFilesToKeep = c.GetNumFilesToKeep()
	*c.KeepOldest = c.GetKeepOldest()
	*c.Remove = c.GetRemove()
	*c.IgnoreErrors = c.GetIgnoreErrors()
	*c.RecursiveSearch = c.GetRecursiveSearch()
	*c.LogLevel = c.GetLogLevel()
	*c.LogFormat = c.GetLogFormat()
	*c.LogFilePath = c.GetLogFilePath()
	*c.ConsoleOutput = c.GetConsoleOutput()
	*c.UseSyslog = c.GetUseSyslog()
	*c.ConfigFile = c.GetConfigFile()
}
