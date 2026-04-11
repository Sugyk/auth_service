package logger

import "context"

type noopLogger struct{}

func (l noopLogger) Debug(ctx context.Context, msg string, args ...any) {}
func (l noopLogger) Info(ctx context.Context, msg string, args ...any)  {}
func (l noopLogger) Warn(ctx context.Context, msg string, args ...any)  {}
func (l noopLogger) Error(ctx context.Context, msg string, args ...any) {}

func (l noopLogger) With(args ...any) Logger {
	return l
}
