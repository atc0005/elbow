// +build !windows

package main

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

func enableSyslogLogging(config *Config, logger *logrus.Logger) error {

	// make sure that the user actually requested syslog logging as it is
	// currently supported on UNIX only.

	if !config.UseSyslog {
		return fmt.Errorf("Syslog logging not requested, not enabling")
	}

	// Attempt to connect to local syslog
	// TODO: We need to decide whether we're using the same logging level
	// as specified for the general logger and update this reference
	// accordingly
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")

	if err == nil {
		// https://github.com/sirupsen/logrus#hooks
		// https://github.com/sirupsen/logrus/blob/master/hooks/syslog/README.md
		// Seems to require `log.AddHook(hook)`` vs `log.Hooks.Add(hook)`
		logger.AddHook(hook)
	} else {
		return fmt.Errorf("unable to connect to syslog socket:", err)
	}

	return nil

}
