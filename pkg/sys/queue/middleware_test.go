package queue_test

import (
	"context"
	"testing"

	"github.com/itimky/spindle/pkg/sys/queue"
	"github.com/itimky/spindle/test"
	"github.com/stretchr/testify/suite"
)

type MiddlewareSuite struct {
	suite.Suite
}

func (s *MiddlewareSuite) Test_MiddlewareRecover_OK() {
	var msg queue.Message

	err := queue.MiddlewareRecover(
		func(ctx context.Context, event queue.Message) error {
			return nil
		})(context.Background(), msg)
	s.NoError(err)
}

func (s *MiddlewareSuite) Test_MiddlewareRecover_Error() {
	var msg queue.Message

	err := queue.MiddlewareRecover(
		func(ctx context.Context, event queue.Message) error {
			return test.Err
		})(context.Background(), msg)
	s.ErrorIs(err, test.Err)
}

func (s *MiddlewareSuite) Test_MiddlewareRecover_Panic() {
	var msg queue.Message

	err := queue.MiddlewareRecover(
		func(ctx context.Context, event queue.Message) error {
			panic("panic")
		})(context.Background(), msg)
	s.ErrorIs(err, queue.ErrPanicRecovered)
}

func TestMiddlewareSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(MiddlewareSuite))
}
