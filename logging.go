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
	"strings"

	"github.com/sirupsen/logrus"
)

func setLoggerConfig(config *Config, logger *logrus.Logger) {

	switch config.LogFormat {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	case "json":
		// Log as JSON instead of the default ASCII formatter.
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	// NOTE: If config.LogFile is set, console output is muted
	var loggerOutput *os.File
	switch {
	case config.ConsoleOutput == "stdout":
		loggerOutput = os.Stdout
	case config.ConsoleOutput == "stderr":
		loggerOutput = os.Stderr
	}

	if strings.TrimSpace(config.LogFilePath) != "" {
		// If this is set, do not log to console unless writing to log file fails
		// FIXME: How do we defer the file close without killing the file handle?
		// https://github.com/sirupsen/logrus/blob/de736cf91b921d56253b4010270681d33fdf7cb5/logger.go#L332
		file, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			loggerOutput = file

			// This is what we'll use to close the file handle from main()
			// https://kgrz.io/reading-files-in-go-an-overview.html
			config.LogFileHandle = file
		} else {
			log.Errorf("Failed to log to %s, will use %s instead.",
				config.LogFilePath, config.ConsoleOutput)
		}
	}

	// Apply chosen output based on earlier checks
	// Note: Can be any io.Writer
	logger.SetOutput(loggerOutput)

	// https://godoc.org/github.com/sirupsen/logrus#Level
	// https://golang.org/pkg/log/syslog/#Priority
	// https://en.wikipedia.org/wiki/Syslog#Severity_level
	switch config.LogLevel {
	case "emerg", "panic":
		logger.SetLevel(logrus.PanicLevel)
	case "alert", "critical", "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "warn", "notice":
		logger.SetLevel(logrus.WarnLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	}

	// https://godoc.org/github.com/sirupsen/logrus#New
	// https://godoc.org/github.com/sirupsen/logrus#Logger

	if config.UseSyslog {
		log.Debug("Syslog logging requested, attempting to enable it")
		if err := enableSyslogLogging(config, logger); err != nil {
			// TODO: Is this sufficient cause for failing? Perhaps if a local
			// log file is not also set consider it a failure?
			log.Errorf("Failed to enable syslog logging: %s", err)
			log.Warn("Proceeding without syslog logging")
		}
	} else {
		log.Debug("Syslog logging not requested, not enabling")
	}

}
