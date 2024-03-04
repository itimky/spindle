package run_test

import (
	"context"
	"testing"
	"time"

	"github.com/itimky/spindle/pkg/sys/log"
	"github.com/itimky/spindle/pkg/sys/run"
	"github.com/itimky/spindle/test"
	"github.com/stretchr/testify/suite"
)

type SupervisorSuite struct {
	suite.Suite
}

func (s *SupervisorSuite) Test_RunUntilAnyExit_ErrNoLogger() {
	supervisor := make(run.Supervisor)
	err := supervisor.RunUntilAnyExit(context.Background())
	s.ErrorIs(err, log.ErrNoLogger)
}

func (s *SupervisorSuite) Test_RunUntilAnyExit_ParentContextCancelled() {
	rootCtx, cancel := context.WithCancel(test.NewContext(s.T()))

	called := false

	supervisor := run.Supervisor{
		"goroutine": func(ctx context.Context) error {
			<-ctx.Done()
			called = true

			return nil
		},
	}

	go func() {
		time.Sleep(time.Millisecond)
		cancel()
	}()

	err := supervisor.RunUntilAnyExit(rootCtx)
	s.Require().NoError(err)

	s.True(called)
}

func (s *SupervisorSuite) Test_RunUntilAnyExit_CancelOnAnyGoroutineCompletion() {
	rootCtx := test.NewContext(s.T())

	called := false

	supervisor := run.Supervisor{
		"goroutine-1": func(ctx context.Context) error {
			return nil
		},
		"goroutine-2": func(ctx context.Context) error {
			<-ctx.Done()
			called = true

			return nil
		},
	}

	err := supervisor.RunUntilAnyExit(rootCtx)
	s.Require().NoError(err)

	s.True(called)
}

func (s *SupervisorSuite) Test_RunUntilAnyExit_CancelOnAnyGoroutineError() {
	rootCtx := test.NewContext(s.T())

	called := false

	supervisor := run.Supervisor{
		"goroutine-1": func(ctx context.Context) error {
			<-ctx.Done()

			called = true

			return nil
		},
		"goroutine-2": func(ctx context.Context) error {
			return test.Err
		},
	}

	err := supervisor.RunUntilAnyExit(rootCtx)
	s.Require().NoError(err)

	s.True(called)
}

func (s *SupervisorSuite) Test_RunUntilAnyExit_CancelOnAnyGoroutinePanic() {
	rootCtx := test.NewContext(s.T())

	called := false

	supervisor := run.Supervisor{
		"goroutine-1": func(ctx context.Context) error {
			<-ctx.Done()

			called = true

			return nil
		},
		"goroutine-2": func(ctx context.Context) error {
			panic("test")
		},
	}

	err := supervisor.RunUntilAnyExit(rootCtx)
	s.Require().NoError(err)

	s.True(called)
}

func (s *SupervisorSuite) Test_RunIdle() {
	supervisor := make(run.Supervisor)

	ctx, cancel := context.WithCancel(test.NewContext(s.T()))
	cancel()

	err := supervisor.RunIdle(ctx)
	s.Require().NoError(err)
}

func (s *SupervisorSuite) Test_RunIdle_ErrNoLogger() {
	supervisor := make(run.Supervisor)
	err := supervisor.RunIdle(context.Background())
	s.ErrorIs(err, log.ErrNoLogger)
}

func TestSupervisorSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(SupervisorSuite))
}
