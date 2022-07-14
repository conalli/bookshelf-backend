package logstest

// TestLogger is a dummy logger used only in testing.
type TestLogger struct{}

// New creates a new TestLogger.
func New() *TestLogger {
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
