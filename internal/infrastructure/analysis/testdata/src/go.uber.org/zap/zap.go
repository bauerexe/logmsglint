package zap

type Logger struct{}

type SugaredLogger struct{}

func NewNop() *Logger {
	return &Logger{}
}

func (l *Logger) Info(string, ...any)  {}
func (l *Logger) Warn(string, ...any)  {}
func (l *Logger) Error(string, ...any) {}
func (l *Logger) Debug(string, ...any) {}

func (l *Logger) Sugar() *SugaredLogger {
	return &SugaredLogger{}
}

func (l *SugaredLogger) Info(string, ...any)  {}
func (l *SugaredLogger) Warn(string, ...any)  {}
func (l *SugaredLogger) Error(string, ...any) {}
func (l *SugaredLogger) Debug(string, ...any) {}
