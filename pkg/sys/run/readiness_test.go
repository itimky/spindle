package run_test

import (
	"context"
	"testing"
	"time"

	"github.com/itimky/spindle/pkg/sys/run"
	"github.com/itimky/spindle/test"
	"github.com/stretchr/testify/suite"
)

type ReadinessSuite struct {
	suite.Suite

	readiness *run.Readiness
}

func (s *ReadinessSuite) SetupTest() {
	s.readiness = run.NewReadiness()
}

func (s *ReadinessSuite) TestReadiness_Wait_Success() {
	ctx := test.NewContext(s.T())

	go func() {
		s.readiness.MarkReady(ctx)
	}()

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err := s.readiness.WaitReady(ctx)
	s.NoError(err)
}

func (s *ReadinessSuite) TestReadiness_Wait_ContextCancelled() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()

	err := s.readiness.WaitReady(ctx)
	s.ErrorIs(err, context.Canceled)
}

func (s *ReadinessSuite) TestReadiness_MarkReady_MultipleTimes() {
	ctx := test.NewContext(s.T())

	s.readiness.MarkReady(ctx)
	s.readiness.MarkReady(ctx)
}

func TestReadinessSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ReadinessSuite))
}
