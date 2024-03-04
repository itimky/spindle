package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/itimky/spindle/pkg/sys/log"
	logadapter "github.com/itimky/spindle/pkg/sys/log/adapter"
	"github.com/itimky/spindle/pkg/sys/run"
)

func main() {
	ctx, cancelCtx := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGTERM,
	)
	defer cancelCtx()

	logger := logadapter.NewZeroLog()
	ctx = log.ToContext(ctx, logger)

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		logger.Errorf("env parse: %s", err)

		return
	}

	answerProcessorComponents := NewAnswerProcessor(cfg, logger)
	systemComponents := NewSystem(cfg, logger, answerProcessorComponents.matrixInMem)

	supervisor := make(run.Supervisor)
	supervisor["answerProcessorConsumer"] = answerProcessorComponents.RunConsumer
	supervisor["systemBootstrapConsumer"] = systemComponents.RunBootstrapConsumer

	err := supervisor.RunIdle(ctx)
	if err != nil {
		logger.Errorf("run: %s", err)
	}
}
