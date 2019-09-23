package main

import (
	"os"
	"strings"

	// Use `log` if we are going to override the default `log`, otherwise
	// import without an "override" if we want to use the `logrus` name.
	// https://godoc.org/github.com/sirupsen/logrus
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

	if strings.TrimSpace(config.LogFile) != "" {
		// If this is set, do not log to console unless writing to log file fails
		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			loggerOutput = file
		} else {
			log.Errorf("Failed to log to %s, will use %s instead.",
				config.LogFile, config.ConsoleOutput)
		}
	}

	// Apply chosen output based on earlier checks
	// Note: Can be any io.Writer
	logger.SetOutput(loggerOutput)

	switch config.LogLevel {
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "warn":
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
			log.Error("enabling syslog logging failed:", err)
		}
	} else {
		log.Debug("Syslog logging not requested, not enabling")
	}

}
