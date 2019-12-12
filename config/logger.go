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

import (
	"fmt"

	"github.com/atc0005/elbow/logging"
	"github.com/sirupsen/logrus"
)

// SetLoggerConfig applies chosen configuration settings that control logging
// output.
func (c *Config) SetLoggerConfig() error {

	var err error

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: "Calling SetLoggerConfig()",
		Fields:  logrus.Fields{"line": logging.GetLineNumber()},
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("Current state of config object: %+v\n", c),
		Fields:  logrus.Fields{"line": logging.GetLineNumber()},
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("The address of the logger SetLoggerConfig received: %p\n", c.GetLogger()),
		Fields:  logrus.Fields{"line": logging.GetLineNumber()},
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: "Current state of individual logging related fields",
		Fields: logrus.Fields{
			"line":           logging.GetLineNumber(),
			"log_format":     c.GetLogFormat(),
			"console_output": c.GetConsoleOutput(),
			"use_syslog":     c.GetUseSyslog(),
		},
	})

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("logger.Out field at start of SetLoggerFormatter(): %p\n", c.GetLogger().Out),
		Fields: logrus.Fields{
			"line": logging.GetLineNumber(),
		},
	})

	if err = logging.SetLoggerFormatter(c.logger, c.GetLogFormat()); err != nil {
		return err
	}

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("logger.Out field after SetLoggerFormatter: %p\n", c.GetLogger().Out),
		Fields: logrus.Fields{
			"line": logging.GetLineNumber(),
		},
	})

	if err = logging.SetLoggerConsoleOutput(c.logger, c.GetConsoleOutput()); err != nil {
		return err
	}

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: fmt.Sprintf("logger.Out field after SetLoggerConsoleOutput(): %p\n", c.GetLogger().Out),
		Fields: logrus.Fields{
			"line": logging.GetLineNumber(),
		},
	})

	// FIXME: This seems like a pretty ugly tradeoff just to avoid golint
	// complaining about the use of an else block; we now have `err` declared
	// at function scope instead of per block scope in order to directly
	// assign to struct field.
	c.logFileHandle, err = logging.SetLoggerLogFile(c.logger, c.GetLogFilePath())
	if err != nil {

		// Need to collect the error for display later
		logging.Buffer.Add(logging.LogRecord{
			Level:   logrus.ErrorLevel,
			Message: fmt.Sprintf("%s", err),
			Fields: logrus.Fields{
				"log_file":        c.GetLogFilePath(),
				"line":            logging.GetLineNumber(),
				"log_file_handle": c.GetLogFileHandle(),
			},
		})

		return err
	}

	if err = logging.SetLoggerLevel(c.logger, c.GetLogLevel()); err != nil {
		return err
	}

	// https://godoc.org/github.com/sirupsen/logrus#New
	// https://godoc.org/github.com/sirupsen/logrus#Logger

	// make sure that the user actually requested syslog logging as it is
	// currently supported on UNIX only.
	if c.GetUseSyslog() {
		logging.Buffer.Add(logging.LogRecord{
			Level:   logrus.InfoLevel,
			Message: "Syslog logging requested, attempting to enable it",
			Fields: logrus.Fields{
				"use_syslog": c.GetUseSyslog(),
				"line":       logging.GetLineNumber(),
			},
		})

		if err := logging.EnableSyslogLogging(c.logger, &logging.Buffer, c.GetLogLevel()); err != nil {
			// TODO: Is this sufficient cause for failing? Perhaps if a local
			// log file is not also set consider it a failure?

			logging.Buffer.Add(logging.LogRecord{
				Level:   logrus.ErrorLevel,
				Message: fmt.Sprintf("Failed to enable syslog logging: %s", err),
				Fields: logrus.Fields{
					"use_syslog": c.GetUseSyslog(),
					"line":       logging.GetLineNumber(),
				},
			})

			logging.Buffer.Add(logging.LogRecord{
				Level:   logrus.WarnLevel,
				Message: "Proceeding without syslog logging",
				Fields: logrus.Fields{
					"use_syslog": c.GetUseSyslog(),
					"line":       logging.GetLineNumber(),
				},
			})
		}
	} else {
		logging.Buffer.Add(logging.LogRecord{
			Level:   logrus.DebugLevel,
			Message: "Syslog logging not requested, not enabling",
			Fields: logrus.Fields{
				"use_syslog": c.GetUseSyslog(),
				"line":       logging.GetLineNumber(),
			},
		})
	}

	logging.Buffer.Add(logging.LogRecord{
		Level:   logrus.DebugLevel,
		Message: "logging object details at end of SetLoggerConfig()",
		Fields: logrus.Fields{
			"logger": fmt.Sprintf("%+v\n", c.GetLogger()),
			// TODO: Re-evaluate potential for field ref on nil pointer
			"logger_out": fmt.Sprintf("%+v\n", c.GetLogger().Out),
			"line":       logging.GetLineNumber(),
		},
	})

	// FIXME: Placeholder for now

	return nil

}
