// +build !windows

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

import (

	"fmt"

	"github.com/atc0005/elbow/config"

	// Use `log` if we are going to override the default `log`, otherwise
	// import without an "override" if we want to use the `logrus` name.
	// https://godoc.org/github.com/sirupsen/logrus
	"github.com/sirupsen/logrus"

	// Official examples use either `lSyslog` or `logrus_syslog`
	// https://godoc.org/github.com/sirupsen/logrus/hooks/syslog
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"

	"log/syslog"

)

func enableSyslogLogging(config *config.Config, logger *logrus.Logger) error {

	// make sure that the user actually requested syslog logging as it is
	// currently supported on UNIX only.

	if !config.UseSyslog {
		return fmt.Errorf("syslog logging not requested, not enabling")
	}

	// Use roughly the same logging level as specified for the general logger
	// https://golang.org/pkg/log/syslog/#Priority
	// https://en.wikipedia.org/wiki/Syslog#Severity_level
	// https://pubs.opengroup.org/onlinepubs/009695399/functions/syslog.html
	var syslogLogLevel syslog.Priority

	switch config.LogLevel {
	case "emerg", "panic":
		// syslog: System is unusable; a panic condition.
		// logrus: calls panic
		syslogLogLevel = syslog.LOG_EMERG
	case "alert", "fatal":
		// syslog: A condition that should be corrected immediately, such as a
		// corrupted system database
		// logrus: calls os.Exit(1)
		syslogLogLevel = syslog.LOG_ALERT
	case "critical":
		// syslog: Critical conditions, such as hard device errors.
		syslogLogLevel = syslog.LOG_CRIT
	case "error":
		// syslog: Error conditions
		// logrus: Used for errors that should definitely be noted. Commonly
		// used for hooks to send errors to an error tracking service.
		syslogLogLevel = syslog.LOG_ERR
	case "warn":
		// syslog: Warning conditions
		// logrus: Non-critical entries that deserve eyes.
		syslogLogLevel = syslog.LOG_WARNING
	case "notice":
		// syslog: Normal but significant conditions; conditions that are not
		// error conditions, but that may require special handling.
		// logrus: N/A
		syslogLogLevel = syslog.LOG_NOTICE
	case "info":
		// syslog: Informational messages
		// logrus: General application operational entries
		syslogLogLevel = syslog.LOG_INFO
	case "debug":
		// syslog: Debug-level messages
		// logrus: Usually only enabled when debugging. Very verbose logging.
		syslogLogLevel = syslog.LOG_DEBUG
	case "trace":
		// syslog: N/A
		// logrus: Finer-grained informational events than debug.
		syslogLogLevel = syslog.LOG_DEBUG
	}

	logger.Debugf("syslog log level set to %v", syslogLogLevel)

	// Attempt to connect to local syslog
	//
	// FIXME: Is this setting a "cap" on the level of log messages that flow
	// through or is this an indication of the type of messages that ALL
	// logging calls will produce? I'm assuming we're setting a limiter here
	hook, err := lSyslog.NewSyslogHook("", "", syslogLogLevel, "")

	if err == nil {
		// https://github.com/sirupsen/logrus#hooks
		// https://github.com/sirupsen/logrus/blob/master/hooks/syslog/README.md
		// Seems to require `log.AddHook(hook)`` vs `log.Hooks.Add(hook)`
		logger.Info("Connected to syslog socket")
		logger.AddHook(hook)
		logger.Debug("AddHook() called for syslog logging")
	} else {
		return fmt.Errorf("unable to connect to syslog socket: %s", err)
	}

	return nil

}