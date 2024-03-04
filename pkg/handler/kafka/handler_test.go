package kafkahandler_test

import (
	"testing"

	kafkacontract "github.com/itimky/spindle/pkg/contract/kafka"
	matchfacade "github.com/itimky/spindle/pkg/facade/match"
	kafkahandler "github.com/itimky/spindle/pkg/handler/kafka"
	"github.com/itimky/spindle/pkg/sys"
	"github.com/itimky/spindle/pkg/sys/queue"
	"github.com/itimky/spindle/test"
	mocks "github.com/itimky/spindle/test/pkg/handler/kafka"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type KafkaHandlerSuite struct {
	suite.Suite

	matchFacadeMock *mocks.MockmatchFacade

	handler *kafkahandler.Handler
}

func (s *KafkaHandlerSuite) SetupTest() {
	s.matchFacadeMock = mocks.NewMockmatchFacade(s.T())

	s.handler = kafkahandler.NewHandler(
		s.matchFacadeMock,
	)
}

func (s *KafkaHandlerSuite) Test_handleAnswerV1() {
	testArgs := []struct {
		name         string
		params       queue.Message
		expectedErr  error
		facadeParams *matchfacade.ProcessAnswerParams
		facadeErr    error
	}{
		{
			name: "err: unmarshal error",
			params: queue.Message{
				Type: kafkacontract.MessageTypeAnswerV1,
				Data: []byte("invalid-json"),
			},
			expectedErr: sys.ErrInvalidJSON,
		},
		{
			name: "err: facade error",
			params: queue.Message{
				Type: kafkacontract.MessageTypeAnswerV1,
				Data: test.MustMarshalJSON(s.T(), kafkacontract.AnswerV1{
					PersonID:   "person-id",
					QuestionID: "question-id",
					AnswerID:   "answer-id",
				}),
			},
			expectedErr: test.Err,
			facadeParams: &matchfacade.ProcessAnswerParams{
				PersonID:   "person-id",
				QuestionID: "question-id",
				AnswerID:   "answer-id",
			},
			facadeErr: test.Err,
		},
		{
			name: "ok",
			params: queue.Message{
				Type: kafkacontract.MessageTypeAnswerV1,
				Data: test.MustMarshalJSON(s.T(), kafkacontract.AnswerV1{
					PersonID:   "person-id",
					QuestionID: "question-id",
					AnswerID:   "answer-id",
				}),
			},
			facadeParams: &matchfacade.ProcessAnswerParams{
				PersonID:   "person-id",
				QuestionID: "question-id",
				AnswerID:   "answer-id",
			},
		},
	}

	for _, testArg := range testArgs {
		testArg := testArg

		s.Run(testArg.name, func() {
			if testArg.facadeParams != nil {
				s.matchFacadeMock.
					EXPECT().
					ProcessAnswer(mock.Anything, *testArg.facadeParams).
					Return(testArg.facadeErr).
					Once()
			}

			err := queue.NewHandler(s.handler).Handle(test.NewContext(s.T()), testArg.params)
			s.ErrorIs(err, testArg.expectedErr)
		})
	}
}

func TestKafkaHandlerSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(KafkaHandlerSuite))
}
