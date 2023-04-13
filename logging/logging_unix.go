//go:build !windows
// +build !windows

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
	"fmt"

	// Use `log` if we are going to override the default `log`, otherwise
	// import without an "override" if we want to use the `logrus` name.
	// https://godoc.org/github.com/sirupsen/logrus
	"github.com/sirupsen/logrus"

	// Official examples use either `lSyslog` or `logrus_syslog`
	// https://godoc.org/github.com/sirupsen/logrus/hooks/syslog
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"

	"log/syslog"
)

// EnableSyslogLogging attempts to enable local syslog logging for non-Windows
// systems. For Windows systems the attempt is skipped.
func EnableSyslogLogging(logger *logrus.Logger, logBuffer *LogBuffer, logLevel string) error {

	// Use roughly the same logging level as specified for the general logger
	// https://golang.org/pkg/log/syslog/#Priority
	// https://en.wikipedia.org/wiki/Syslog#Severity_level
	// https://pubs.opengroup.org/onlinepubs/009695399/functions/syslog.html
	var syslogLogLevel syslog.Priority

	switch logLevel {
	case LogLevelEmergency, LogLevelPanic:
		syslogLogLevel = syslog.LOG_EMERG
	case LogLevelAlert, LogLevelFatal:
		syslogLogLevel = syslog.LOG_ALERT
	case LogLevelCritical:
		syslogLogLevel = syslog.LOG_CRIT
	case LogLevelError:
		syslogLogLevel = syslog.LOG_ERR
	case LogLevelWarn:
		syslogLogLevel = syslog.LOG_WARNING
	case LogLevelNotice:
		syslogLogLevel = syslog.LOG_NOTICE
	case LogLevelInfo:
		syslogLogLevel = syslog.LOG_INFO
	case LogLevelDebug:
		syslogLogLevel = syslog.LOG_DEBUG
	case LogLevelTrace:
		syslogLogLevel = syslog.LOG_DEBUG
	default:
		return fmt.Errorf("invalid syslog log level: %q", logLevel)
	}

	logBuffer.Add(LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("syslog log level set to %v", syslogLogLevel),
		Fields:  logrus.Fields{"log_level": logLevel},
	})

	// Attempt to connect to local syslog
	//
	// FIXME: Is this setting a "cap" on the level of log messages that flow
	// through or is this an indication of the type of messages that ALL
	// logging calls will produce? I'm assuming we're setting a limiter here
	hook, err := lSyslog.NewSyslogHook("", "", syslogLogLevel, "")

	if err == nil {

		logBuffer.Add(LogRecord{
			Level:   logrus.InfoLevel,
			Message: "Connected to syslog socket",
		})

		logger.AddHook(hook)

		logBuffer.Add(LogRecord{
			Level:   logrus.DebugLevel,
			Message: "Finished using AddHook() to enable syslog logging",
		})

	} else {
		return fmt.Errorf("unable to connect to syslog socket: %w", err)
	}

	return nil

}
