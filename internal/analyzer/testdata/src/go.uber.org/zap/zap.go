package zap

type Logger struct{}

func NewNop() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string, fields ...any) {}
func (l *Logger) Error(msg string, fields ...any) {}
