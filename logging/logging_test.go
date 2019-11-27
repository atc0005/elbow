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

// Fix linting error
// string `fakeValue` has 3 occurrences, make it a constant (goconst)
const fakeValue = "fakeValue"

func TestGetLineNumber(t *testing.T) {
	got := GetLineNumber()
	if got < 1 {
		t.Errorf("Line number is incorrect, got: %d, want: greater than 0.", got)
	}

}

func TestSetLoggerLevelShouldFail(t *testing.T) {

	logger := logrus.New()

	give := fakeValue
	got := SetLoggerLevel(logger, give)
	if got == nil {
		t.Error("Expected error for", give, "Got", got)
	} else {
		t.Logf("Got error as expected for %v: %v", give, got)
	}

}

// Pass in a valid logLevel string, call logger.GetLevel()
// and compare against the expected value
func TestSetLoggerLevelShouldSucceed(t *testing.T) {

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
		if err := SetLoggerLevel(logger, give); err != nil {
			t.Error("Error when calling SetLoggerLevel(): ", err)
		} else {
			t.Log("No error when calling SetLoggerLevel()")
		}
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

func TestSetLoggerFormatterShouldFail(t *testing.T) {

	logger := logrus.New()

	give := fakeValue
	got := SetLoggerFormatter(logger, give)
	if got == nil {
		t.Error("Expected error for", give, "Got", got)
	} else {
		t.Logf("Got error as expected for %v: %v", give, got)
	}
}

func TestSetLoggerFormatterShouldSucceed(t *testing.T) {

	type test struct {
		format string
		result error
	}

	logger := logrus.New()

	tests := []test{
		test{format: "text", result: nil},
		test{format: "json", result: nil},
	}

	for _, give := range tests {
		got := SetLoggerFormatter(logger, give.format)
		if got != give.result {
			t.Error("Expected", give.result, "Got", got)
		}
	}

}

func TestSetLoggerConsoleOutputShouldFail(t *testing.T) {

	logger := logrus.New()

	give := fakeValue
	got := SetLoggerConsoleOutput(logger, give)
	if got == nil {
		t.Error("Expected error for", give, "Got", got)
	} else {
		t.Logf("Got error as expected for %v: %v", give, got)
	}
}

func TestSetLoggerConsoleOutputShouldSucceed(t *testing.T) {

	type test struct {
		consoleOutput string
		result        error
	}

	logger := logrus.New()

	tests := []test{
		test{consoleOutput: "stdout", result: nil},
		test{consoleOutput: "stderr", result: nil},
	}

	for _, give := range tests {
		got := SetLoggerConsoleOutput(logger, give.consoleOutput)
		if got != give.result {
			t.Error("Expected", give.result, "Got", got)
		}
	}

}
