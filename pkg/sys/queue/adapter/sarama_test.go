package queueadapter_test

import (
	"context"
	"testing"

	"github.com/IBM/sarama"
	"github.com/itimky/spindle/pkg/sys/log"
	"github.com/itimky/spindle/pkg/sys/queue"
	queueadapter "github.com/itimky/spindle/pkg/sys/queue/adapter"
	"github.com/itimky/spindle/test"
	saramamocks "github.com/itimky/spindle/test/github.com/IBM/sarama"
	mocks "github.com/itimky/spindle/test/pkg/sys/queue/adapter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ConsumerSaramaSuite struct {
	suite.Suite

	sessionMock *saramamocks.MockConsumerGroupSession
	claimMock   *saramamocks.MockConsumerGroupClaim
	handlerMock *mocks.Mockhandler

	consumer *queueadapter.ConsumerSarama
}

func (s *ConsumerSaramaSuite) SetupSuite() {
	s.sessionMock = saramamocks.NewMockConsumerGroupSession(s.T())
	s.claimMock = saramamocks.NewMockConsumerGroupClaim(s.T())
	s.handlerMock = mocks.NewMockhandler(s.T())

	s.consumer = queueadapter.NewConsumerSarama(s.handlerMock)
}

func (s *ConsumerSaramaSuite) TestConsumer_ConsumeClaim() {
	testArgs := []struct {
		name        string
		expectedErr error
		ctxFn       func() context.Context
		msg         *sarama.ConsumerMessage
		handleMsg   *queue.Message
		handleErr   error
		markMsg     *sarama.ConsumerMessage
	}{
		{
			name:        "err: no logger",
			expectedErr: log.ErrNoLogger,
			ctxFn:       context.Background,
		},
		{
			name:        "err: handler err",
			expectedErr: test.Err,
			ctxFn: func() context.Context {
				return test.NewContext(s.T())
			},
			msg: &sarama.ConsumerMessage{
				Headers: []*sarama.RecordHeader{
					{
						Key:   []byte("type"),
						Value: []byte("test"),
					},
				},
				Value: []byte("value"),
			},
			handleMsg: &queue.Message{
				Type: "test",
				Data: []byte("value"),
			},
			handleErr: test.Err,
		},
		{
			name: "ok",
			ctxFn: func() context.Context {
				return test.NewContext(s.T())
			},
			msg: &sarama.ConsumerMessage{
				Headers: []*sarama.RecordHeader{
					{
						Key:   []byte("type"),
						Value: []byte("test"),
					},
				},
				Value: []byte("value"),
			},
			handleMsg: &queue.Message{
				Type: "test",
				Data: []byte("value"),
			},
			markMsg: &sarama.ConsumerMessage{
				Headers: []*sarama.RecordHeader{
					{
						Key:   []byte("type"),
						Value: []byte("test"),
					},
				},
				Value: []byte("value"),
			},
		},
		{
			name: "ok: ctx done",
			ctxFn: func() context.Context {
				ctx, cancel := context.WithCancel(test.NewContext(s.T()))
				cancel()

				return ctx
			},
		},
	}

	for _, testArg := range testArgs {
		testArg := testArg

		s.Run(testArg.name, func() {
			msgChan := make(chan *sarama.ConsumerMessage, 1)
			s.claimMock.EXPECT().Messages().Return(msgChan).Once()
			if testArg.msg != nil {
				msgChan <- testArg.msg
				close(msgChan)
			}

			s.sessionMock.EXPECT().Context().Return(testArg.ctxFn()).Once()

			if testArg.handleMsg != nil {
				s.handlerMock.EXPECT().Handle(mock.Anything, *testArg.handleMsg).Return(testArg.handleErr).Once()
			}

			if testArg.markMsg != nil {
				s.sessionMock.EXPECT().MarkMessage(testArg.markMsg, "").Once()
			}

			err := s.consumer.ConsumeClaim(s.sessionMock, s.claimMock)
			s.ErrorIs(err, testArg.expectedErr)
		})
	}
}

func (s *ConsumerSaramaSuite) Test_Setup() {
	err := s.consumer.Setup(s.sessionMock)
	s.NoError(err)
}

func (s *ConsumerSaramaSuite) Test_Cleanup() {
	err := s.consumer.Cleanup(s.sessionMock)
	s.NoError(err)
}

func TestConsumerSaramaSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ConsumerSaramaSuite))
}

type BootstrapConsumerSaramaSuite struct {
	suite.Suite

	consumerMock     *mocks.MockpartitionConsumerSarama
	bootstrapperMock *mocks.Mockbootstrapper
	handlerMock      *mocks.Mockhandler

	adapter *queueadapter.BootstrapConsumerSarama
}

func (s *BootstrapConsumerSaramaSuite) SetupSuite() {
	s.consumerMock = mocks.NewMockpartitionConsumerSarama(s.T())
	s.bootstrapperMock = mocks.NewMockbootstrapper(s.T())
	s.handlerMock = mocks.NewMockhandler(s.T())

	s.adapter = queueadapter.NewBootstrapConsumerSarama(
		s.consumerMock,
		s.bootstrapperMock,
		s.handlerMock,
	)
}

func (s *BootstrapConsumerSaramaSuite) Test_Bootstrap() {
	testArgs := []struct {
		name        string
		expectedErr error
		ctxFn       func() context.Context
		watermark   int64
		chanMsgs    []*sarama.ConsumerMessage
		chanErr     *sarama.ConsumerError
		bootParams  []queue.Message
		bootErr     error
	}{
		{
			name:        "err: context done",
			expectedErr: context.Canceled,
			ctxFn: func() context.Context {
				ctx, cancel := context.WithCancel(test.NewContext(s.T()))
				cancel()

				return ctx
			},
		},
		{
			name: "ok: nil msg",
			ctxFn: func() context.Context {
				return test.NewContext(s.T())
			},
			chanMsgs: []*sarama.ConsumerMessage{
				nil,
			},
			bootParams: []queue.Message{},
		},
		{
			name:        "err: chan err",
			expectedErr: test.Err,
			ctxFn: func() context.Context {
				return test.NewContext(s.T())
			},
			chanErr: &sarama.ConsumerError{
				Err: test.Err,
			},
		},
		{
			name:        "err: bootstrapper err",
			expectedErr: test.Err,
			watermark:   1,
			ctxFn: func() context.Context {
				return test.NewContext(s.T())
			},
			chanMsgs: []*sarama.ConsumerMessage{
				{
					Offset: 0,
					Headers: []*sarama.RecordHeader{
						{
							Key:   []byte("type"),
							Value: []byte("test"),
						},
					},
					Value: []byte("payload"),
				},
			},
			bootParams: []queue.Message{
				{
					Type: "test",
					Data: []byte("payload"),
				},
			},
			bootErr: test.Err,
		},
		{
			name:      "ok",
			watermark: 2,
			ctxFn: func() context.Context {
				return test.NewContext(s.T())
			},
			chanMsgs: []*sarama.ConsumerMessage{
				{
					Offset: 0,
					Headers: []*sarama.RecordHeader{
						{
							Key:   []byte("type"),
							Value: []byte("test-0"),
						},
					},
					Value: []byte("payload-0"),
				},
				{
					Offset: 1,
					Headers: []*sarama.RecordHeader{
						{
							Key:   []byte("type"),
							Value: []byte("test-1"),
						},
					},
					Value: []byte("payload-1"),
				},
			},
			bootParams: []queue.Message{
				{
					Type: "test-0",
					Data: []byte("payload-0"),
				},
				{
					Type: "test-1",
					Data: []byte("payload-1"),
				},
			},
		},
	}

	for _, testArg := range testArgs {
		testArg := testArg

		s.Run(testArg.name, func() {
			s.consumerMock.EXPECT().HighWaterMarkOffset().Return(testArg.watermark).Once()

			msgChan := make(chan *sarama.ConsumerMessage, len(testArg.chanMsgs))
			for i := range testArg.chanMsgs {
				msgChan <- testArg.chanMsgs[i]
			}
			s.consumerMock.EXPECT().Messages().Return(msgChan).Once()

			errChan := make(chan *sarama.ConsumerError, 1)
			if testArg.chanErr != nil {
				errChan <- &sarama.ConsumerError{
					Err: testArg.chanErr,
				}
			}
			s.consumerMock.EXPECT().Errors().Return(errChan).Once()

			if testArg.bootParams != nil {
				s.bootstrapperMock.EXPECT().Bootstrap(mock.Anything, testArg.bootParams).Return(testArg.bootErr).Once()
			}

			err := s.adapter.Bootstrap(testArg.ctxFn())
			s.ErrorIs(err, testArg.expectedErr)
		})
	}
}

func (s *BootstrapConsumerSaramaSuite) Test_Consume() {
	testArgs := []struct {
		name        string
		expectedErr error
		ctxFn       func() context.Context
		chanMsgs    []*sarama.ConsumerMessage
		chanErr     *sarama.ConsumerError
		handleMsg   *queue.Message
		handleErr   error
	}{
		{
			name:        "err: context done",
			expectedErr: context.Canceled,
			ctxFn: func() context.Context {
				ctx, cancel := context.WithCancel(test.NewContext(s.T()))
				cancel()

				return ctx
			},
		},
		{
			name: "ok: nil msg",
			ctxFn: func() context.Context {
				return test.NewContext(s.T())
			},
			chanMsgs: []*sarama.ConsumerMessage{
				nil,
			},
		},
		{
			name:        "err: chan err",
			expectedErr: test.Err,
			ctxFn: func() context.Context {
				return test.NewContext(s.T())
			},
			chanErr: &sarama.ConsumerError{
				Err: test.Err,
			},
		},
		{
			name:        "err: handler err",
			expectedErr: test.Err,
			ctxFn: func() context.Context {
				return test.NewContext(s.T())
			},
			chanMsgs: []*sarama.ConsumerMessage{
				{
					Offset: 0,
					Headers: []*sarama.RecordHeader{
						{
							Key:   []byte("type"),
							Value: []byte("test"),
						},
					},
					Value: []byte("payload"),
				},
			},
			handleMsg: &queue.Message{
				Type: "test",
				Data: []byte("payload"),
			},
			handleErr: test.Err,
		},
		{
			name: "ok",
			ctxFn: func() context.Context {
				return test.NewContext(s.T())
			},
			chanMsgs: []*sarama.ConsumerMessage{
				{
					Offset: 0,
					Headers: []*sarama.RecordHeader{
						{
							Key:   []byte("type"),
							Value: []byte("test"),
						},
					},
					Value: []byte("payload"),
				},
				nil, // break loop
			},
			handleMsg: &queue.Message{
				Type: "test",
				Data: []byte("payload"),
			},
		},
	}

	for _, testArg := range testArgs {
		testArg := testArg

		s.Run(testArg.name, func() {
			msgChan := make(chan *sarama.ConsumerMessage, len(testArg.chanMsgs))
			for i := range testArg.chanMsgs {
				msgChan <- testArg.chanMsgs[i]
			}
			s.consumerMock.EXPECT().Messages().Return(msgChan).Once()

			errChan := make(chan *sarama.ConsumerError, 1)
			if testArg.chanErr != nil {
				errChan <- &sarama.ConsumerError{
					Err: testArg.chanErr,
				}
			}
			s.consumerMock.EXPECT().Errors().Return(errChan).Once()

			if testArg.handleMsg != nil {
				s.handlerMock.EXPECT().Handle(mock.Anything, *testArg.handleMsg).Return(testArg.handleErr).Once()
			}

			err := s.adapter.Consume(testArg.ctxFn())
			s.ErrorIs(err, testArg.expectedErr)
		})
	}
}

func TestBootstrapConsumerSaramaSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(BootstrapConsumerSaramaSuite))
}
