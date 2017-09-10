/*
Package console implements simple hierarchical logging functionality.
*/
package console

import (
	"log"
	"os"
	"strings"
)

// Levels for implementing the debug and trace message functionality.
const (
	LevelTrace uint8 = iota
	LevelDebug
	LevelInfo
	LevelStatus
	LevelWarn
	LevelSilent
)

// These variables are initialized in init()
var logLevel = LevelInfo
var logger *log.Logger
var logLevelStrings = [...]string{
	"trace", "debug", "info", "status", "warn", "silent",
}

//===========================================================================
// Interact with debug output
//===========================================================================

// Init the console logger with the previx and log options
func Init(prefix string, flag int) {
	logger = log.New(os.Stdout, prefix, flag)
}

// LogLevel returns a string representation of the current level.
func LogLevel() string {
	return logLevelStrings[logLevel]
}

// SetLogLevel modifies the log level for messages at runtime. Ensures that
// the highest level that can be set is the trace level. This function is
// often called from outside of the package in an init() function to define
// how logging is handled in the console.
func SetLogLevel(level uint8) {
	if level > LevelSilent {
		level = LevelSilent
	}

	logLevel = level
}

//===========================================================================
// Debugging output functions
//===========================================================================

// Print to the standard logger at the specified level. Arguments are handled
// in the manner of log.Printf, but a newline is appended.
func print(level uint8, msg string, a ...interface{}) {
	if level >= logLevel {
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}

		logger.Printf(msg, a...)
	}
}

// Warn prints to the standard logger if level is warn or greater; arguments
// are handled in the manner of log.Printf, but a newline is appended.
func Warn(msg string, a ...interface{}) {
	print(LevelWarn, msg, a...)
}

// Warne is a helper function to simply warn about an error received.
func Warne(err error) {
	Warn(err.Error())
}

// Status prints to the standard logger if level is status or greater;
// arguments are handled in the manner of log.Printf, but a newline is
// appended.
func Status(msg string, a ...interface{}) {
	print(LevelStatus, msg, a...)
}

// Info prints to the standard logger if level is info or greater; arguments
// are handled in the manner of log.Printf, but a newline is appended.
func Info(msg string, a ...interface{}) {
	print(LevelInfo, msg, a...)
}

// Debug prints to the standard logger if level is debug or greater;
// arguments are handled in the manner of log.Printf, but a newline is
// appended.
func Debug(msg string, a ...interface{}) {
	print(LevelDebug, msg, a...)
}

// Trace prints to the standard logger if level is trace or greater;
// arguments are handled in the manner of log.Printf, but a newline is
// appended.
func Trace(msg string, a ...interface{}) {
	print(LevelTrace, msg, a...)
}
