package answerprocessorstorage_test

import (
	"testing"

	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
	answerprocessorstorage "github.com/itimky/spindle/pkg/storage/answer-processor"
	"github.com/itimky/spindle/test"
	mocks "github.com/itimky/spindle/test/pkg/storage/answer-processor"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CompositeStorageSuite struct {
	suite.Suite

	matrixStoreMock   *mocks.MockmatrixStore
	relationStoreMock *mocks.MockrelationStore
	composite         *answerprocessorstorage.Composite
}

func (s *CompositeStorageSuite) SetupTest() {
	s.matrixStoreMock = mocks.NewMockmatrixStore(s.T())
	s.relationStoreMock = mocks.NewMockrelationStore(s.T())

	s.composite = answerprocessorstorage.NewComposite(
		s.matrixStoreMock,
		s.relationStoreMock,
	)
}

func (s *CompositeStorageSuite) Test_GetQuestionWeightMatrix() {
	testArgs := []struct {
		name           string
		params         answerprocessor.GetWeightMatrixParams
		expectedResult *answerprocessor.GetWeightMatrixResult
		expectedErr    error
		getParams      *answerprocessor.GetWeightMatrixParams
		getResult      answerprocessor.WeightMatrix
		getErr         error
	}{
		{
			name: "err: get err",
			params: answerprocessor.GetWeightMatrixParams{
				QuestionID: "q1",
			},
			expectedErr: test.Err,
			getParams: &answerprocessor.GetWeightMatrixParams{
				QuestionID: "q1",
			},
			getErr: test.Err,
		},
		{
			name: "ok",
			params: answerprocessor.GetWeightMatrixParams{
				QuestionID: "q1",
			},
			expectedResult: &answerprocessor.GetWeightMatrixResult{
				Matrix: answerprocessor.WeightMatrix{
					"a1": {
						"a1": decimal.RequireFromString("0.1"),
					},
				},
			},
			getParams: &answerprocessor.GetWeightMatrixParams{
				QuestionID: "q1",
			},
			getResult: answerprocessor.WeightMatrix{
				"a1": {
					"a1": decimal.RequireFromString("0.1"),
				},
			},
		},
	}

	for _, args := range testArgs {
		args := args

		s.Run(args.name, func() {
			if args.getParams != nil {
				s.matrixStoreMock.EXPECT().Get(mock.Anything, *args.getParams).Return(args.getResult, args.getErr).Once()
			}

			result, err := s.composite.GetWeightMatrix(
				test.NewContext(s.T()),
				args.params,
			)
			s.ErrorIs(err, args.expectedErr)
			s.Equal(args.expectedResult, result)
		})
	}
}

func (s *CompositeStorageSuite) Test_GetOtherPersonAnswers() {
	testArgs := []struct {
		name           string
		params         answerprocessor.GetOtherPersonAnswersParams
		expectedResult *answerprocessor.GetOtherPersonAnswersResult
		expectedErr    error
		getParams      *answerprocessor.GetOtherPersonAnswersParams
		getResult      *answerprocessor.GetOtherPersonAnswersResult
		getErr         error
	}{
		{
			name: "err: get err",
			params: answerprocessor.GetOtherPersonAnswersParams{
				QuestionID: "q1",
			},
			expectedErr: test.Err,
			getParams: &answerprocessor.GetOtherPersonAnswersParams{
				QuestionID: "q1",
			},
			getErr: test.Err,
		},
		{
			name: "ok",
			params: answerprocessor.GetOtherPersonAnswersParams{
				QuestionID: "q1",
			},
			expectedResult: &answerprocessor.GetOtherPersonAnswersResult{
				PersonAnswers: []answerprocessor.PersonAnswer{
					{
						PersonID: "p1",
						AnswerID: "a1",
					},
				},
			},
			getParams: &answerprocessor.GetOtherPersonAnswersParams{
				QuestionID: "q1",
			},
			getResult: &answerprocessor.GetOtherPersonAnswersResult{
				PersonAnswers: []answerprocessor.PersonAnswer{
					{
						PersonID: "p1",
						AnswerID: "a1",
					},
				},
			},
		},
	}

	for _, args := range testArgs {
		args := args

		s.Run(args.name, func() {
			if args.getParams != nil {
				s.relationStoreMock.EXPECT().
					GetOtherPersonAnswers(mock.Anything, *args.getParams).
					Return(args.getResult, args.getErr).
					Once()
			}

			result, err := s.composite.GetOtherPersonAnswers(test.NewContext(s.T()), args.params)
			s.ErrorIs(err, args.expectedErr)
			s.Equal(args.expectedResult, result)
		})
	}
}

func (s *CompositeStorageSuite) Test_UpdatePersonsWeights() {
	testArgs := []struct {
		name         string
		params       answerprocessor.UpdatePersonsWeightsParams
		expectedErr  error
		updateParams *answerprocessor.UpdatePersonsWeightsParams
		updateErr    error
	}{
		{
			name: "err: update err",
			params: answerprocessor.UpdatePersonsWeightsParams{
				QuestionID: "q1",
				PersonID:   "p1",
				RelatedWeights: []answerprocessor.RelatedPersonWeight{
					{
						RelatedPersonID: "p2",
						Weight:          decimal.RequireFromString("0.1"),
					},
				},
			},
			expectedErr: test.Err,
			updateParams: &answerprocessor.UpdatePersonsWeightsParams{
				QuestionID: "q1",
				PersonID:   "p1",
				RelatedWeights: []answerprocessor.RelatedPersonWeight{
					{
						RelatedPersonID: "p2",
						Weight:          decimal.RequireFromString("0.1"),
					},
				},
			},
			updateErr: test.Err,
		},
		{
			name: "ok",
			params: answerprocessor.UpdatePersonsWeightsParams{
				QuestionID: "q1",
				PersonID:   "p1",
				RelatedWeights: []answerprocessor.RelatedPersonWeight{
					{
						RelatedPersonID: "p2",
						Weight:          decimal.RequireFromString("0.1"),
					},
				},
			},
			updateParams: &answerprocessor.UpdatePersonsWeightsParams{
				QuestionID: "q1",
				PersonID:   "p1",
				RelatedWeights: []answerprocessor.RelatedPersonWeight{
					{
						RelatedPersonID: "p2",
						Weight:          decimal.RequireFromString("0.1"),
					},
				},
			},
		},
	}

	for _, args := range testArgs {
		args := args

		s.Run(args.name, func() {
			if args.updateParams != nil {
				s.relationStoreMock.EXPECT().UpdatePersonsWeights(mock.Anything, *args.updateParams).Return(args.updateErr).Once()
			}

			err := s.composite.UpdatePersonsWeights(test.NewContext(s.T()), args.params)
			s.ErrorIs(err, args.expectedErr)
		})
	}
}

func TestCompositeStorageSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(CompositeStorageSuite))
}
