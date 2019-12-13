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

package logging

// Fix linting error
// string `fakeValue` has 3 occurrences, make it a constant (goconst)
const fakeValue = "fakeValue"

// Log levels
const (
	// https://godoc.org/github.com/sirupsen/logrus#Level
	// https://golang.org/pkg/log/syslog/#Priority
	// https://en.wikipedia.org/wiki/Syslog#Severity_level

	// LogLevelEmergency represents messages at Emergency level. System is
	// unusable; a panic condition. This level is mapped to logrus.PanicLevel
	// and syslog.LOG_EMERG. logrus calls panic.
	LogLevelEmergency string = "emergency"

	// LogLevelPanic represents messages at Emergency level. System is
	// unusable; a panic condition. This level is mapped to logrus.PanicLevel
	// and syslog.LOG_EMERG. logrus calls panic.
	LogLevelPanic string = "panic"

	// LogLevelAlert represents a condition that should be corrected
	// immediately, such as a corrupted system database. This level is mapped
	// to logrus.FatalLevel and syslog.LOG_ALERT. logrus calls os.Exit(1)
	LogLevelAlert string = "alert"

	// LogLevelFatal is used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	// This level is mapped to logrus.FatalLevel and syslog.LOG_ALERT. logrus
	// calls os.Exit(1)
	LogLevelFatal string = "fatal"

	// LogLevelCritical is for critical conditions, such as hard device
	// errors. This level is mapepd to logrus.FatalLevel and syslog.LOG_CRIT.
	// logrus calls os.Exit(1)
	LogLevelCritical string = "critical"

	// LogLevelError is for errors that should definitely be noted. Commonly
	// used for hooks to send errors to an error tracking service. This level
	// maps to logrus.ErrorLevel and syslog.LOG_ERR.
	LogLevelError string = "error"

	// LogLevelWarn is for non-critical entries that deserve eyes. This level
	// maps to logrus.WarnLevel and syslog.LOG_WARNING.
	LogLevelWarn string = "warn"

	// LogLevelNotice is for normal but significant conditions; conditions
	// that are not error conditions, but that may require special handling.
	// This level maps to logrus.WarnLevel and syslog.LOG_NOTICE.
	LogLevelNotice string = "notice"

	// LogLevelInfo is for general application operational entries. This level
	// maps to logrus.InfoLevel and syslog.LOG_INFO.
	LogLevelInfo string = "info"

	// LogLevelDebug is for debug-level messages and is usually on enabled
	// when debugging. Very verbose logging. This level is mapped to
	// logrus.DebugLevel and syslog.LOG_DEBUG.
	LogLevelDebug string = "debug"

	// LogLevelTrace is for finer-grained informational events than debug.
	// This level maps to logrus.TraceLevel and syslog.LOG_DEBUG.
	LogLevelTrace string = "trace"
)

// Log formats used by logrus
const (

	// LogFormatText represents the logrus text formatter.
	LogFormatText string = "text"

	// LogFormatJSON represents the logrus JSON formatter.
	LogFormatJSON string = "json"
)

const (

	// ConsoleOutputStdout represents os.Stdout
	ConsoleOutputStdout string = "stdout"

	// ConsoleOutputStderr represents os.Stderr
	ConsoleOutputStderr string = "stderr"
)
