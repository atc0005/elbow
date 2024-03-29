//go:build windows
// +build windows

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

package logging

import (
	"github.com/sirupsen/logrus"
)

// EnableSyslogLogging attempts to enable local syslog logging for non-Windows
// systems. For Windows systems the attempt is skipped.
func EnableSyslogLogging(_ *logrus.Logger, logBuffer *LogBuffer, _ string) error {

	logBuffer.Add(LogRecord{
		// TODO: Not sure what log level is appropriate here. We are already
		// reporting failures enabling syslog logging from the caller, but are
		// not noting elsewhere that Windows syslog support is not available.
		Level:   logrus.WarnLevel,
		Message: "This is a Windows build. Syslog support is not currently supported for this platform.",
	})

	return nil

}
