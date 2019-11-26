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
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGetLineNumber(t *testing.T) {
	got := GetLineNumber()
	if got < 1 {
		t.Errorf("Line number is incorrect, got: %d, want: greater than 0.", got)
	}

}

func TestSetLoggerLevelShouldPanic(t *testing.T) {

	// https://stackoverflow.com/questions/31595791/how-to-test-panics
	defer func() {

		r := recover()
		t.Logf("Panic message: %q", r)

		if r == nil {
			t.Errorf("SetLoggerLevel accepted an invalid level without panicing.")
		}
	}()

	logger := logrus.New()
	badLogLevel := "tacos"

	SetLoggerLevel(logger, badLogLevel)
}

func TestSetLoggerLevelShouldSucceed(t *testing.T) {

	// TODO: Pass in a valid logLevel string, call logger.GetLevel()
	// and compare against the expected value

	type test struct {
		logLevel    string
		loggerLevel logrus.Level
	}

	tests := []test{
		test{logLevel: "emerg", loggerLevel: logrus.PanicLevel},
		test{logLevel: "panic", loggerLevel: logrus.PanicLevel},
		test{logLevel: "alert", loggerLevel: logrus.FatalLevel},
		test{logLevel: "critical", loggerLevel: logrus.FatalLevel},
		test{logLevel: "fatal", loggerLevel: logrus.FatalLevel},
		test{logLevel: "error", loggerLevel: logrus.ErrorLevel},
		test{logLevel: "warn", loggerLevel: logrus.WarnLevel},
		test{logLevel: "notice", loggerLevel: logrus.WarnLevel},
		test{logLevel: "info", loggerLevel: logrus.InfoLevel},
		test{logLevel: "debug", loggerLevel: logrus.DebugLevel},
		test{logLevel: "trace", loggerLevel: logrus.TraceLevel},
	}

	logger := logrus.New()

	for _, v := range tests {
		give := v.logLevel
		SetLoggerLevel(logger, give)
		want := v.loggerLevel
		got := logger.GetLevel()

		if got != v.loggerLevel {
			t.Error("Expected", want, "Got", got)
			t.FailNow()
		} else {
			t.Log("Got", got, "as expected for requested level of", give)
		}
	}

}
