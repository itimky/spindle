package test

import (
	"context"
	"os/signal"
	"syscall"
	"testing"

	"github.com/itimky/spindle/pkg/sys/log"
	logadapter "github.com/itimky/spindle/pkg/sys/log/adapter"
)

func NewContext(t *testing.T) context.Context {
	t.Helper()

	ctx, _ := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGTERM,
	)

	return log.ToContext(ctx, logadapter.NewNop())
}
