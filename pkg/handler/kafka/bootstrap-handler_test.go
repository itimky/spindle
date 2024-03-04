package kafkahandler_test

import (
	"testing"

	kafkacontract "github.com/itimky/spindle/pkg/contract/kafka"
	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
	systemfacade "github.com/itimky/spindle/pkg/facade/system"
	kafkahandler "github.com/itimky/spindle/pkg/handler/kafka"
	"github.com/itimky/spindle/pkg/sys"
	"github.com/itimky/spindle/pkg/sys/queue"
	"github.com/itimky/spindle/test"
	mocks "github.com/itimky/spindle/test/pkg/handler/kafka"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WeightMatrixBootstrapHandlerSuite struct {
	suite.Suite

	systemFacadeMock *mocks.MocksystemFacade

	handler *kafkahandler.WeightMatrixBootstrapHandler
}

func (s *WeightMatrixBootstrapHandlerSuite) SetupTest() {
	s.systemFacadeMock = mocks.NewMocksystemFacade(s.T())

	s.handler = kafkahandler.NewWeightMatrixBootstrapHandler(
		s.systemFacadeMock,
	)
}

func (s *WeightMatrixBootstrapHandlerSuite) Test_Bootstrap() {
	testArgs := []struct {
		name         string
		params       []queue.Message
		expectedErr  error
		facadeParams *systemfacade.BootstrapWeightMatrixStorageParams
		facadeErr    error
	}{
		{
			name: "err: unmarshal error",
			params: []queue.Message{
				{
					Type: kafkacontract.MessageTypeWeightMatrixV1,
					Data: []byte("invalid-json"),
				},
			},
			expectedErr: sys.ErrInvalidJSON,
		},
		{
			name: "err: convert error",
			params: []queue.Message{
				{
					Type: kafkacontract.MessageTypeWeightMatrixV1,
					Data: test.MustMarshalJSON(s.T(), kafkacontract.WeightMatrixV1{
						QuestionID: "q1",
						Matrix: map[string]map[string]string{
							"a1": {"a1": "0.1", "a2": "invalid-decimal-string"},
						},
					}),
				},
			},
			expectedErr: sys.ErrInvalidDecimalString,
		},
		{
			name: "err: facade error",
			params: []queue.Message{
				{
					Type: kafkacontract.MessageTypeWeightMatrixV1,
					Data: test.MustMarshalJSON(s.T(), kafkacontract.WeightMatrixV1{
						QuestionID: "q1",
						Matrix: map[string]map[string]string{
							"a1": {"a1": "0.1", "a2": "0.2"},
							"a2": {"a1": "0.2", "a2": "0.4"},
						},
					}),
				},
			},
			expectedErr: test.Err,
			facadeParams: &systemfacade.BootstrapWeightMatrixStorageParams{
				Matrices: []systemfacade.QuestionWeightMatrix{
					{
						QuestionID: "q1",
						Matrix: answerprocessor.WeightMatrix{
							"a1": {"a1": decimal.RequireFromString("0.1"), "a2": decimal.RequireFromString("0.2")},
							"a2": {"a1": decimal.RequireFromString("0.2"), "a2": decimal.RequireFromString("0.4")},
						},
					},
				},
			},
			facadeErr: test.Err,
		},
		{
			name: "ok",
			params: []queue.Message{
				{
					Type: kafkacontract.MessageTypeWeightMatrixV1,
					Data: test.MustMarshalJSON(s.T(), kafkacontract.WeightMatrixV1{
						QuestionID: "q1",
						Matrix: map[string]map[string]string{
							"a1": {"a1": "0.1", "a2": "0.2"},
							"a2": {"a1": "0.2", "a2": "0.4"},
						},
					}),
				},
				{
					Type: kafkacontract.MessageTypeWeightMatrixV1,
					Data: test.MustMarshalJSON(s.T(), kafkacontract.WeightMatrixV1{
						QuestionID: "q2",
						Matrix: map[string]map[string]string{
							"a1": {"a1": "0.2", "a2": "0.4"},
							"a2": {"a1": "0.4", "a2": "0.8"},
						},
					}),
				},
			},
			facadeParams: &systemfacade.BootstrapWeightMatrixStorageParams{
				Matrices: []systemfacade.QuestionWeightMatrix{
					{
						QuestionID: "q1",
						Matrix: answerprocessor.WeightMatrix{
							"a1": {"a1": decimal.RequireFromString("0.1"), "a2": decimal.RequireFromString("0.2")},
							"a2": {"a1": decimal.RequireFromString("0.2"), "a2": decimal.RequireFromString("0.4")},
						},
					},
					{
						QuestionID: "q2",
						Matrix: answerprocessor.WeightMatrix{
							"a1": {"a1": decimal.RequireFromString("0.2"), "a2": decimal.RequireFromString("0.4")},
							"a2": {"a1": decimal.RequireFromString("0.4"), "a2": decimal.RequireFromString("0.8")},
						},
					},
				},
			},
		},
	}

	for _, testArg := range testArgs {
		testArg := testArg

		s.Run(testArg.name, func() {
			if testArg.facadeParams != nil {
				s.systemFacadeMock.EXPECT().
					BootstrapWeightMatrixStorage(mock.Anything, *testArg.facadeParams).
					Return(testArg.facadeErr).
					Once()
			}

			err := s.handler.Bootstrap(test.NewContext(s.T()), testArg.params)
			s.ErrorIs(err, testArg.expectedErr)
		})
	}
}

func (s *WeightMatrixBootstrapHandlerSuite) Test_Handle() {
	testArgs := []struct {
		name         string
		params       queue.Message
		expectedErr  error
		facadeParams *systemfacade.UpdateWeightMatrixStorageParams
		facadeErr    error
	}{
		{
			name: "err: unmarshal error",
			params: queue.Message{
				Type: kafkacontract.MessageTypeWeightMatrixV1,
				Data: []byte("invalid-json"),
			},
			expectedErr: sys.ErrInvalidJSON,
		},
		{
			name: "err: convert error",
			params: queue.Message{
				Type: kafkacontract.MessageTypeWeightMatrixV1,
				Data: test.MustMarshalJSON(s.T(), kafkacontract.WeightMatrixV1{
					QuestionID: "q1",
					Matrix: map[string]map[string]string{
						"a1": {"a1": "0.1", "a2": "invalid-decimal-string"},
					},
				}),
			},
			expectedErr: sys.ErrInvalidDecimalString,
		},
		{
			name: "err: facade error",
			params: queue.Message{
				Type: kafkacontract.MessageTypeWeightMatrixV1,
				Data: test.MustMarshalJSON(s.T(), kafkacontract.WeightMatrixV1{
					QuestionID: "q1",
					Matrix: map[string]map[string]string{
						"a1": {"a1": "0.1", "a2": "0.2"},
						"a2": {"a1": "0.2", "a2": "0.4"},
					},
				}),
			},
			expectedErr: test.Err,
			facadeParams: &systemfacade.UpdateWeightMatrixStorageParams{
				QuestionID: "q1",
				Matrix: answerprocessor.WeightMatrix{
					"a1": {"a1": decimal.RequireFromString("0.1"), "a2": decimal.RequireFromString("0.2")},
					"a2": {"a1": decimal.RequireFromString("0.2"), "a2": decimal.RequireFromString("0.4")},
				},
			},
			facadeErr: test.Err,
		},
		{
			name: "ok",
			params: queue.Message{
				Type: kafkacontract.MessageTypeWeightMatrixV1,
				Data: test.MustMarshalJSON(s.T(), kafkacontract.WeightMatrixV1{
					QuestionID: "q1",
					Matrix: map[string]map[string]string{
						"a1": {"a1": "0.1", "a2": "0.2"},
						"a2": {"a1": "0.2", "a2": "0.4"},
					},
				}),
			},
			facadeParams: &systemfacade.UpdateWeightMatrixStorageParams{
				QuestionID: "q1",
				Matrix: answerprocessor.WeightMatrix{
					"a1": {"a1": decimal.RequireFromString("0.1"), "a2": decimal.RequireFromString("0.2")},
					"a2": {"a1": decimal.RequireFromString("0.2"), "a2": decimal.RequireFromString("0.4")},
				},
			},
		},
	}

	for _, testArg := range testArgs {
		testArg := testArg

		s.Run(testArg.name, func() {
			if testArg.facadeParams != nil {
				s.systemFacadeMock.EXPECT().
					UpdateWeightMatrixStorage(mock.Anything, *testArg.facadeParams).
					Return(testArg.facadeErr).
					Once()
			}

			err := s.handler.Handle(test.NewContext(s.T()), testArg.params)
			s.ErrorIs(err, testArg.expectedErr)
		})
	}
}

func TestWeightMatrixBootstrapHandlerSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(WeightMatrixBootstrapHandlerSuite))
}
