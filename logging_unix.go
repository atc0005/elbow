// +build !windows

package main

import (

	// Use `log` if we are going to override the default `log`, otherwise
	// import without an "override" if we want to use the `logrus` name.
	// https://godoc.org/github.com/sirupsen/logrus
	"github.com/sirupsen/logrus"

	// Official examples use either `lSyslog` or `logrus_syslog`
	// https://godoc.org/github.com/sirupsen/logrus/hooks/syslog
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"

	"log/syslog"

)

func newLogger() (*logrus.Logger) {

	logger := logrus.New()

	// https://godoc.org/github.com/sirupsen/logrus#New
	// https://godoc.org/github.com/sirupsen/logrus#Logger

	// Attempt to connect to local syslog
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")

	if err == nil {
		// https://github.com/sirupsen/logrus#hooks
		// https://github.com/sirupsen/logrus/blob/master/hooks/syslog/README.md
		// Seems to require `log.AddHook(hook)`` vs `log.Hooks.Add(hook)`
		logger.AddHook(hook)
	}
	else {
		logger.Warn("Unable to connect to syslog socket:", err)
	}

	// FIXME: Is this what I should be returning here?
	return &logger
}
