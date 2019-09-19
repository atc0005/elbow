// +build windows

package main

import (
	// Use `log` if we are going to override the default `log`, otherwise
	// import without an "override" if we want to use the `logrus` name.
	// https://godoc.org/github.com/sirupsen/logrus
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {

	// Stub to cover Windows builds
	return logrus.New()
}
