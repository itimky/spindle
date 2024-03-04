package answerprocessor_test

import (
	"testing"

	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
	"github.com/itimky/spindle/test"
	mocks "github.com/itimky/spindle/test/pkg/domain/answer-processor"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AnswerProcessorSuite struct {
	suite.Suite

	mockStorage *mocks.Mockstorage
	processor   *answerprocessor.AnswerProcessor
}

func (s *AnswerProcessorSuite) SetupTest() {
	s.mockStorage = mocks.NewMockstorage(s.T())
	s.processor = answerprocessor.NewAnswerProcessor(s.mockStorage)
}

func (s *AnswerProcessorSuite) Test_ProcessAnswer() {
	testArgs := []struct {
		name             string
		params           answerprocessor.ProcessAnswerParams
		expectedErr      error
		getMatrixParams  *answerprocessor.GetWeightMatrixParams
		getMatrixResult  *answerprocessor.GetWeightMatrixResult
		getMatrixErr     error
		getAnswersParams *answerprocessor.GetOtherPersonAnswersParams
		getAnswersResult *answerprocessor.GetOtherPersonAnswersResult
		getAnswersErr    error
		updateParams     *answerprocessor.UpdatePersonsWeightsParams
		updateErr        error
	}{
		{
			name: "err: get matrix error",
			params: answerprocessor.ProcessAnswerParams{
				PersonID:   "person-one",
				QuestionID: "question",
				AnswerID:   "answer-one",
			},
			expectedErr: test.Err,
			getMatrixParams: &answerprocessor.GetWeightMatrixParams{
				QuestionID: "question",
			},
			getMatrixErr: test.Err,
		},
		{
			name: "err: get answers error",
			params: answerprocessor.ProcessAnswerParams{
				PersonID:   "person-one",
				QuestionID: "question",
				AnswerID:   "answer-one",
			},
			expectedErr: test.Err,
			getMatrixParams: &answerprocessor.GetWeightMatrixParams{
				QuestionID: "question",
			},
			getMatrixResult: &answerprocessor.GetWeightMatrixResult{
				Matrix: answerprocessor.WeightMatrix{},
			},
			getAnswersParams: &answerprocessor.GetOtherPersonAnswersParams{
				QuestionID: "question",
			},
			getAnswersErr: test.Err,
		},
		{
			name: "err: update error",
			params: answerprocessor.ProcessAnswerParams{
				PersonID:   "person-one",
				QuestionID: "question",
				AnswerID:   "answer-one",
			},
			expectedErr: test.Err,
			getMatrixParams: &answerprocessor.GetWeightMatrixParams{
				QuestionID: "question",
			},
			getMatrixResult: &answerprocessor.GetWeightMatrixResult{
				Matrix: answerprocessor.WeightMatrix{
					"answer-one": {
						"answer-one": decimal.RequireFromString("1"),
						"answer-two": decimal.RequireFromString("2"),
					},
					"answer-two": {
						"answer-one": decimal.RequireFromString("2"),
						"answer-two": decimal.RequireFromString("4"),
					},
				},
			},
			getAnswersParams: &answerprocessor.GetOtherPersonAnswersParams{
				QuestionID: "question",
			},
			getAnswersResult: &answerprocessor.GetOtherPersonAnswersResult{
				PersonAnswers: []answerprocessor.PersonAnswer{
					{
						PersonID: "person-two",
						AnswerID: "answer-one",
					},
					{
						PersonID: "person-three",
						AnswerID: "answer-two",
					},
				},
			},
			updateParams: &answerprocessor.UpdatePersonsWeightsParams{
				PersonID:   "person-one",
				QuestionID: "question",
				RelatedWeights: []answerprocessor.RelatedPersonWeight{
					{
						RelatedPersonID: "person-two",
						Weight:          decimal.RequireFromString("1"),
					},
					{
						RelatedPersonID: "person-three",
						Weight:          decimal.RequireFromString("2"),
					},
				},
			},
			updateErr: test.Err,
		},
		{
			name: "ok",
			params: answerprocessor.ProcessAnswerParams{
				PersonID:   "person-one",
				QuestionID: "question",
				AnswerID:   "answer-one",
			},
			getMatrixParams: &answerprocessor.GetWeightMatrixParams{
				QuestionID: "question",
			},
			getMatrixResult: &answerprocessor.GetWeightMatrixResult{
				Matrix: answerprocessor.WeightMatrix{
					"answer-one": {
						"answer-one": decimal.RequireFromString("1"),
						"answer-two": decimal.RequireFromString("2"),
					},
					"answer-two": {
						"answer-one": decimal.RequireFromString("2"),
						"answer-two": decimal.RequireFromString("4"),
					},
				},
			},
			getAnswersParams: &answerprocessor.GetOtherPersonAnswersParams{
				QuestionID: "question",
			},
			getAnswersResult: &answerprocessor.GetOtherPersonAnswersResult{
				PersonAnswers: []answerprocessor.PersonAnswer{
					{
						PersonID: "person-two",
						AnswerID: "answer-one",
					},
					{
						PersonID: "person-three",
						AnswerID: "answer-two",
					},
				},
			},
			updateParams: &answerprocessor.UpdatePersonsWeightsParams{
				PersonID:   "person-one",
				QuestionID: "question",
				RelatedWeights: []answerprocessor.RelatedPersonWeight{
					{
						RelatedPersonID: "person-two",
						Weight:          decimal.RequireFromString("1"),
					},
					{
						RelatedPersonID: "person-three",
						Weight:          decimal.RequireFromString("2"),
					},
				},
			},
		},
	}

	for _, testArg := range testArgs {
		testArg := testArg
		s.Run(testArg.name, func() {
			if testArg.getMatrixParams != nil {
				s.mockStorage.EXPECT().GetWeightMatrix(
					mock.Anything,
					*testArg.getMatrixParams,
				).Return(testArg.getMatrixResult, testArg.getMatrixErr).Once()
			}

			if testArg.getAnswersParams != nil {
				s.mockStorage.EXPECT().GetOtherPersonAnswers(
					mock.Anything,
					*testArg.getAnswersParams,
				).Return(testArg.getAnswersResult, testArg.getAnswersErr).Once()
			}

			if testArg.updateParams != nil {
				s.mockStorage.EXPECT().UpdatePersonsWeights(
					mock.Anything,
					*testArg.updateParams,
				).Return(testArg.updateErr).Once()
			}

			err := s.processor.ProcessAnswer(test.NewContext(s.T()), testArg.params)
			s.ErrorIs(err, testArg.expectedErr)
		})
	}
}

func TestAnswerProcessorSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(AnswerProcessorSuite))
}
