package systemfacade_test

import (
	"testing"

	"github.com/itimky/spindle/pkg/domain"
	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
	systemfacade "github.com/itimky/spindle/pkg/facade/system"
	"github.com/itimky/spindle/test"
	mocks "github.com/itimky/spindle/test/pkg/facade/system"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SystemFacadeSuite struct {
	suite.Suite

	weightMatrixStoreMock *mocks.MockweightMatrixStore

	facade *systemfacade.SystemFacade
}

func (s *SystemFacadeSuite) SetupTest() {
	s.weightMatrixStoreMock = mocks.NewMockweightMatrixStore(s.T())

	s.facade = systemfacade.NewSystemFacade(
		s.weightMatrixStoreMock,
	)
}

func (s *SystemFacadeSuite) Test_BootstrapWeightMatrixStorage() {
	testArgs := []struct {
		name            string
		params          systemfacade.BootstrapWeightMatrixStorageParams
		expectedErr     error
		bootstrapParams map[domain.QuestionID]answerprocessor.WeightMatrix
		bootstrapErr    error
	}{
		{
			name: "err: bootstrap error",
			params: systemfacade.BootstrapWeightMatrixStorageParams{
				Matrices: []systemfacade.QuestionWeightMatrix{
					{
						QuestionID: "q1",
						Matrix: answerprocessor.WeightMatrix{
							"a1": {
								"a1": decimal.RequireFromString("0.1"),
							},
						},
					},
					{
						QuestionID: "q2",
						Matrix: answerprocessor.WeightMatrix{
							"a1": {
								"a1": decimal.RequireFromString("0.2"),
							},
						},
					},
				},
			},
			expectedErr: test.Err,
			bootstrapParams: map[domain.QuestionID]answerprocessor.WeightMatrix{
				"q1": {
					"a1": {
						"a1": decimal.RequireFromString("0.1"),
					},
				},
				"q2": {
					"a1": {
						"a1": decimal.RequireFromString("0.2"),
					},
				},
			},
			bootstrapErr: test.Err,
		},
		{
			name: "ok",
			params: systemfacade.BootstrapWeightMatrixStorageParams{
				Matrices: []systemfacade.QuestionWeightMatrix{
					{
						QuestionID: "q1",
						Matrix: answerprocessor.WeightMatrix{
							"a1": {
								"a1": decimal.RequireFromString("0.1"),
							},
						},
					},
					{
						QuestionID: "q2",
						Matrix: answerprocessor.WeightMatrix{
							"a1": {
								"a1": decimal.RequireFromString("0.2"),
							},
						},
					},
				},
			},
			bootstrapParams: map[domain.QuestionID]answerprocessor.WeightMatrix{
				"q1": {
					"a1": {
						"a1": decimal.RequireFromString("0.1"),
					},
				},
				"q2": {
					"a1": {
						"a1": decimal.RequireFromString("0.2"),
					},
				},
			},
		},
	}

	for _, args := range testArgs {
		args := args

		s.Run(args.name, func() {
			if args.bootstrapParams != nil {
				s.weightMatrixStoreMock.EXPECT().
					Bootstrap(mock.Anything, args.bootstrapParams).
					Return(args.bootstrapErr).Once()
			}

			err := s.facade.BootstrapWeightMatrixStorage(test.NewContext(s.T()), args.params)
			s.ErrorIs(err, args.expectedErr)
		})
	}
}

func (s *SystemFacadeSuite) Test_UpdateWeightMatrixStorage() {
	testArgs := []struct {
		name        string
		params      systemfacade.UpdateWeightMatrixStorageParams
		expectedErr error
		setParams   domain.QuestionID
		setMatrix   answerprocessor.WeightMatrix
		setErr      error
	}{
		{
			name: "err: set error",
			params: systemfacade.UpdateWeightMatrixStorageParams{
				QuestionID: "q1",
				Matrix: answerprocessor.WeightMatrix{
					"a1": {
						"a1": decimal.RequireFromString("0.1"),
					},
				},
			},
			expectedErr: test.Err,
			setParams:   "q1",
			setMatrix: answerprocessor.WeightMatrix{
				"a1": {
					"a1": decimal.RequireFromString("0.1"),
				},
			},
			setErr: test.Err,
		},
		{
			name: "ok",
			params: systemfacade.UpdateWeightMatrixStorageParams{
				QuestionID: "q1",
				Matrix: answerprocessor.WeightMatrix{
					"a1": {
						"a1": decimal.RequireFromString("0.1"),
					},
				},
			},
			setParams: "q1",
			setMatrix: answerprocessor.WeightMatrix{
				"a1": {
					"a1": decimal.RequireFromString("0.1"),
				},
			},
		},
	}

	for _, args := range testArgs {
		args := args

		s.Run(args.name, func() {
			if args.setParams != "" {
				s.weightMatrixStoreMock.EXPECT().
					Set(mock.Anything, args.setParams, args.setMatrix).
					Return(args.setErr).Once()
			}

			err := s.facade.UpdateWeightMatrixStorage(test.NewContext(s.T()), args.params)
			s.ErrorIs(err, args.expectedErr)
		})
	}
}

func TestSystemFacadeSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(SystemFacadeSuite))
}
