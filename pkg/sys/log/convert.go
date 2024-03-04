package log

import (
	"context"
)

type loggerKey struct{}

func ToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) (Logger, error) { //nolint: ireturn
	logger, ok := ctx.Value(loggerKey{}).(Logger)
	if !ok {
		return nil, ErrNoLogger
	}

	return logger, nil
}
