package logger

import "context"

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)

	With(args ...any) Logger
}

func New(level string) Logger {
	return newSlogLogger(level)
}

func NewNoop() Logger {
	return noopLogger{}
}
