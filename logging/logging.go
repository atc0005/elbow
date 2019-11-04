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

// Package logging is intended mostly as a set of helper functions around
// configuring and using a common logger to provide structured, leveled
// logging.
package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// LogRecord holds logrus.Field values along with additional metadata that can be
// used later to complete the log message submission process.
type LogRecord struct {
	Level   logrus.Level
	Message string
	Fields  logrus.Fields
}

// LogBuffer represents a slice of LogRecord objects
type LogBuffer []LogRecord

// Add passed LogRecord type to slice of LogRecord objects
func (lb *LogBuffer) Add(r LogRecord) {
	*lb = append(*lb, r)
}

// Flush LogRecord entries after applying user-provided logging settings
func (lb LogBuffer) Flush(logger *logrus.Logger) {

	for _, entry := range lb {

		switch {

		case entry.Level == logrus.PanicLevel:
			logger.WithFields(entry.Fields).Panic(entry.Message)

		case entry.Level == logrus.FatalLevel:
			logger.WithFields(entry.Fields).Fatal(entry.Message)

		case entry.Level == logrus.ErrorLevel:
			logger.WithFields(entry.Fields).Error(entry.Message)

		case entry.Level == logrus.WarnLevel:
			logger.WithFields(entry.Fields).Warn(entry.Message)

		case entry.Level == logrus.InfoLevel:
			logger.WithFields(entry.Fields).Info(entry.Message)

		case entry.Level == logrus.DebugLevel:
			logger.WithFields(entry.Fields).Debug(entry.Message)

		case entry.Level == logrus.TraceLevel:
			logger.WithFields(entry.Fields).Trace(entry.Message)

		default:
			panic("This should not have been reachable")

		}

	}

	// TODO: Empty slice now that we're done processing all items

}

// SetLoggerFormatter sets a user-specified logging format for the provided
// logger object.
func SetLoggerFormatter(logger *logrus.Logger, format string) {
	switch format {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	case "json":
		// Log as JSON instead of the default ASCII formatter.
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
}

// SetLoggerConsoleOutput configures the chosen console output to one of
// stdout or stderr.
func SetLoggerConsoleOutput(logger *logrus.Logger, consoleOutput string) {
	var loggerOutput *os.File
	switch {
	case consoleOutput == "stdout":
		loggerOutput = os.Stdout
	case consoleOutput == "stderr":
		loggerOutput = os.Stderr
	}

	// Apply chosen output based on earlier checks
	// Note: Can be any io.Writer
	logger.SetOutput(loggerOutput)
}

// SetLoggerLogFile configures a log file as the destination for all log
// messages. NOTE: If successfully set, console output is muted.
func SetLoggerLogFile(logger *logrus.Logger, logFilePath string) (*os.File, error) {

	var file *os.File
	var err error

	if strings.TrimSpace(logFilePath) != "" {
		// If this is set, do not log to console unless writing to log file fails
		// FIXME: How do we defer the file close without killing the file handle?
		// https://github.com/sirupsen/logrus/blob/de736cf91b921d56253b4010270681d33fdf7cb5/logger.go#L332
		file, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to log to %s, will leave configuration as is",
				logFilePath)
		}
		// The `file` handle is what we'll use to close the file handle from main()
		// https://kgrz.io/reading-files-in-go-an-overview.html
		logger.SetOutput(file)
	}

	return file, nil
}

// SetLoggerLevel applies the requested logger level to filter out messages
// with a lower level than the one configured.
func SetLoggerLevel(logger *logrus.Logger, logLevel string) {

	// https://godoc.org/github.com/sirupsen/logrus#Level
	// https://golang.org/pkg/log/syslog/#Priority
	// https://en.wikipedia.org/wiki/Syslog#Severity_level
	switch logLevel {
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

}
