package run

import (
	"context"
	"fmt"
	"sync/atomic"
)

type Readiness struct {
	// readyCh is used to WaitReady until MarkReady is called
	readyCh chan struct{}
	// readyFlag is used to not close readyCh multiple times
	readyFlag atomic.Bool
}

func NewReadiness() *Readiness {
	return &Readiness{
		readyCh:   make(chan struct{}),
		readyFlag: atomic.Bool{},
	}
}

func (r *Readiness) MarkReady(_ context.Context) {
	wasReady := r.readyFlag.Swap(true)
	if wasReady {
		return
	}

	close(r.readyCh)
}

func (r *Readiness) WaitReady(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context done: %w", ctx.Err())
	case <-r.readyCh:
		return nil
	}
}
