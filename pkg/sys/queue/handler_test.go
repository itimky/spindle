package queue_test

import (
	"context"
	"testing"

	"github.com/itimky/spindle/pkg/sys/queue"
	"github.com/itimky/spindle/test"
	mocks "github.com/itimky/spindle/test/pkg/sys/queue"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite

	routerMock *mocks.Mockrouter
}

func (s *HandlerSuite) SetupTest() {
	s.routerMock = mocks.NewMockrouter(s.T())
}

func (s *HandlerSuite) Test_Handle_ErrNoRoute() {
	s.routerMock.EXPECT().Routes().Return(map[queue.MessageType]queue.HandlerFunc{}).Once()

	handler := queue.NewHandler(
		s.routerMock,
	)

	err := handler.Handle(test.NewContext(s.T()), queue.Message{
		Type: "test",
		Data: []byte("value"),
	})
	s.ErrorIs(err, queue.ErrNoRoute)
}

func (s *HandlerSuite) Test_Handle_Middlewares() {
	var mw1Called, mw2Called, hCalled bool

	mw1 := func(h queue.HandlerFunc) queue.HandlerFunc {
		return func(ctx context.Context, msg queue.Message) error {
			mw1Called = true

			return h(ctx, msg)
		}
	}
	mw2 := func(h queue.HandlerFunc) queue.HandlerFunc {
		return func(ctx context.Context, msg queue.Message) error {
			mw2Called = true

			return h(ctx, msg)
		}
	}

	s.routerMock.EXPECT().Routes().Return(map[queue.MessageType]queue.HandlerFunc{
		"test": func(ctx context.Context, msg queue.Message) error {
			hCalled = true

			return nil
		},
	}).Once()

	handler := queue.NewHandler(
		s.routerMock,
		mw1,
		mw2,
	)

	err := handler.Handle(test.NewContext(s.T()), queue.Message{
		Type: "test",
		Data: []byte("value"),
	})
	s.Require().NoError(err)
	s.True(mw1Called)
	s.True(mw2Called)
	s.True(hCalled)
}

func TestHandlerSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(HandlerSuite))
}
