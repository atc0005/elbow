// +build windows

package main

import (
	// Use `log` if we are going to override the default `log`, otherwise
	// import without an "override" if we want to use the `logrus` name.
	// https://godoc.org/github.com/sirupsen/logrus
	"fmt"

	"github.com/sirupsen/logrus"
)

func enableSyslogLogging(config *Config, logger *logrus.Logger) error {

	// make sure that the user actually requested syslog logging as it is
	// currently supported on UNIX only.

	if !config.UseSyslog {
		return fmt.Errorf("Syslog logging not requested, not enabling")
	}

	log.Debug("This is a Windows build. Syslog support is not currently supported for Windows builds.")

	return nil

}
