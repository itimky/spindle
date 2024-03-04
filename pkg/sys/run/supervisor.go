package run

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/itimky/spindle/pkg/sys/log"
)

type goroutine func(ctx context.Context) error

type Supervisor map[string]goroutine

// RunUntilAnyExit runs all goroutines until any of them exits.
func (s *Supervisor) RunUntilAnyExit(ctx context.Context) error {
	logger, err := log.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("logger from context: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var group sync.WaitGroup

	for name, routine := range *s {
		group.Add(1)

		go func(name string, routine goroutine) {
			logger := logger.WithField("goroutine", name)

			defer group.Done()
			defer cancel()
			defer func() {
				if r := recover(); r != nil {
					logger.Errorf("%s: %s", name, r)
				}
			}()

			err := routine(ctx)
			if err != nil && !errors.Is(err, context.Canceled) {
				logger.Errorf("%s: %s", name, err)
			}
		}(name, routine)
	}

	logger.Info("running goroutines")

	group.Wait()

	logger.Info("all goroutines have been stopped")

	return nil
}

// RunIdle do nothing until the context is cancelled.
func (s *Supervisor) RunIdle(ctx context.Context) error {
	logger, err := log.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("from context: %w", err)
	}

	logger.Info("running idle")

	<-ctx.Done()

	logger.Info("idle has been stopped")

	return nil
}
