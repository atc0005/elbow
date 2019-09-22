package main

import (
	"os"

	// Use `log` if we are going to override the default `log`, otherwise
	// import without an "override" if we want to use the `logrus` name.
	// https://godoc.org/github.com/sirupsen/logrus
	"github.com/sirupsen/logrus"
)

func setLoggerConfig(config *Config, logger *logrus.Logger) {

	// switch on log format here

	// Log as JSON instead of the default ASCII formatter.
	// TODO: Use command-line option here
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// switch on log level here

	// TODO: Accept command-line parameter to determine this level
	// TODO: Setup mapping between command-line options and valid logrus levels
	// so that they can be referenced here
	logger.SetLevel(logrus.DebugLevel)

	// https://godoc.org/github.com/sirupsen/logrus#New
	// https://godoc.org/github.com/sirupsen/logrus#Logger

	if config.UseSyslog {
		log.Debug("Syslog logging requested, attempting to enable it")
		if err := enableSyslogLogging(config, logger); err != nil {

		}

	} else {
		log.Debug("Syslog logging not requested, not enabling")
	}

}
