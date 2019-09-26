// +build windows

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
	// Use `log` if we are going to override the default `log`, otherwise
	// import without an "override" if we want to use the `logrus` name.
	// https://godoc.org/github.com/sirupsen/logrus
	"fmt"

	// TODO: Is this needed here with a global `log` objection already created
	// from this package?
	"github.com/sirupsen/logrus"
)

func enableSyslogLogging(config *Config, logger *logrus.Logger) error {

	// make sure that the user actually requested syslog logging as it is
	// currently supported on UNIX only.

	if !config.UseSyslog {
		return fmt.Errorf("syslog logging not requested, not enabling")
	}

	log.Debug("This is a Windows build. Syslog support is not currently supported for this platform.")

	return nil

}
