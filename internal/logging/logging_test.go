// Copyright 2020 Adam Chalkley
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

package logging

import (
	// "io"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogBufferFlushNilLoggerShouldFail(t *testing.T) {

	var nilLogger *logrus.Logger

	var logBuffer LogBuffer

	if err := logBuffer.Flush(nilLogger); err == nil {
		t.Error("passed nil *logrus.Logger without error")
	} else {
		t.Log("received error as expected:", err)
	}
}

func TestLogBufferFlushShouldSucceed(t *testing.T) {

	var testLogBuffer LogBuffer

	logger := logrus.New()
	// Configure logger to throw everything away
	// logger.SetOutput(io.Discard)
	logger.SetLevel(logrus.TraceLevel)

	type test struct {
		entryLevel logrus.Level
		// potentially used for dealing with PanicLevel and FatalLevel?
		// result     error
	}

	tests := []test{
		// TODO: Need to add coverage for messages at these log levels:
		//
		// {entryLevel: logrus.PanicLevel, result: nil},
		// {entryLevel: logrus.FatalLevel, result: nil},
		//
		// Problem: Flushing either of these types results in that immediate
		// action; FatalLevel forces an exit, PanicLevel forces a panic.

		{entryLevel: logrus.ErrorLevel},
		{entryLevel: logrus.WarnLevel},
		{entryLevel: logrus.InfoLevel},
		{entryLevel: logrus.DebugLevel},
		{entryLevel: logrus.TraceLevel},
	}

	// Create test log buffer entries
	for _, v := range tests {

		testLogBuffer.Add(LogRecord{
			Level:   v.entryLevel,
			Message: fmt.Sprintf("This is a message at level %v.", v.entryLevel),
		})
	}

	// Verify that the number of entries matches up with the same number of
	// active test entries
	if len(testLogBuffer) != len(tests) {
		t.Errorf("Expected %d log buffer entries, Got %d",
			len(testLogBuffer), len(tests))
	} else {
		t.Log("Number of log buffer entries matches test entries")
	}

	if err := testLogBuffer.Flush(logger); err != nil {
		t.Error("Failed to flush log entries:", err)
	} else {
		t.Log("Flushed log buffer entry as expected")
	}
}

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

	// TODO: Evaluate replacing bare strings with constants (see constants.go)
	tests := []test{
		{logLevel: LogLevelEmergency, loggerLevel: logrus.PanicLevel},
		{logLevel: LogLevelPanic, loggerLevel: logrus.PanicLevel},
		{logLevel: LogLevelAlert, loggerLevel: logrus.FatalLevel},
		{logLevel: LogLevelCritical, loggerLevel: logrus.FatalLevel},
		{logLevel: LogLevelFatal, loggerLevel: logrus.FatalLevel},
		{logLevel: LogLevelError, loggerLevel: logrus.ErrorLevel},
		{logLevel: LogLevelWarn, loggerLevel: logrus.WarnLevel},
		{logLevel: LogLevelNotice, loggerLevel: logrus.WarnLevel},
		{logLevel: LogLevelInfo, loggerLevel: logrus.InfoLevel},
		{logLevel: LogLevelDebug, loggerLevel: logrus.DebugLevel},
		{logLevel: LogLevelTrace, loggerLevel: logrus.TraceLevel},
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
		{format: LogFormatText, result: nil},
		{format: LogFormatJSON, result: nil},
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

	// TODO: Evaluate replacing bare strings with constants (see constants.go)
	tests := []test{
		{consoleOutput: ConsoleOutputStdout, result: nil},
		{consoleOutput: ConsoleOutputStderr, result: nil},
	}

	for _, give := range tests {
		got := SetLoggerConsoleOutput(logger, give.consoleOutput)
		if got != give.result {
			t.Error("Expected", give.result, "Got", got)
		}
	}

}

func TestEnableSyslogLogging(t *testing.T) {
	// TODO: Need to implement this

	t.Log("TODO: Need to implement this test.")
}
