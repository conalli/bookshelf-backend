package testutils

import "os"

// TestLogger is a dummy logger used only in testing.
type TestLogger struct{}

// NewLogger creates a new TestLogger.
func NewLogger() *TestLogger {
	return &TestLogger{}
}

// Info implements the logs.Logger interface but does nothing.
func (l *TestLogger) Info(args ...interface{}) {}

// Infof implements the logs.Logger interface but does nothing.
func (l *TestLogger) Infof(template string, args ...interface{}) {}

// Error implements the logs.Logger interface but does nothing.
func (l *TestLogger) Error(args ...interface{}) {}

// Errorf implements the logs.Logger interface but does nothing.
func (l *TestLogger) Errorf(template string, args ...interface{}) {}

// Fatal implements the logs.Logger interface but does nothing.
func (l *TestLogger) Fatal(args ...interface{}) {
	os.Exit(1)
}

// Fatalf implements the logs.Logger interface but does nothing.
func (l *TestLogger) Fatalf(template string, args ...interface{}) {
	os.Exit(1)
}
