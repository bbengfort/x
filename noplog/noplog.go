/*
Package noplog implements a noop logger with a null writer.

The writer implements the io.Write interface but does not do anything.
All logging function definitions are empty noops.

Common use:

    // Initialize the package.
    func init() {
        // Set the random seed to something different each time.
        // This is very common in init
        rand.Seed(time.Now().Unix())

        // Stop the grpc verbose logging
        grpclog.SetLogger(noplog.New())
    }

This functionality exists because some third party packages (looking at you,
grpc) have internal logging that may interfere with primary application
logging. The noplog can be passed to those libraries, allowing you to stifle
those log messages in favor of your application logs.
*/
package noplog

import "log"

// New returns a NopLogger that sets the log module to write to a null
// writer, thereby ensuring that logging is a noop for any package that uses
// the logger.
func New() *NopLogger {
	return &NopLogger{
		log.New(NullWriter(1), "", log.LstdFlags),
	}
}

// NullWriter implements the io.Write interface but doesn't do anything.
type NullWriter int

// Write implements the io.Write interface but is a noop.
func (NullWriter) Write([]byte) (int, error) { return 0, nil }

// NopLogger is a noop logger for passing to grpclog to minimize spew.
type NopLogger struct {
	*log.Logger
}

// Fatal is a noop
func (l *NopLogger) Fatal(args ...interface{}) {}

// Fatalf is a noop
func (l *NopLogger) Fatalf(format string, args ...interface{}) {}

// Fatalln is a noop
func (l *NopLogger) Fatalln(args ...interface{}) {}

// Print is a noop
func (l *NopLogger) Print(args ...interface{}) {}

// Printf is a noop
func (l *NopLogger) Printf(format string, args ...interface{}) {}

// Println is a noop
func (l *NopLogger) Println(v ...interface{}) {}
