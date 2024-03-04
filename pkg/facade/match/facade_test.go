package matchfacade_test

import (
	"testing"

	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
	matchfacade "github.com/itimky/spindle/pkg/facade/match"
	test "github.com/itimky/spindle/test"
	mocks "github.com/itimky/spindle/test/pkg/facade/match"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MatchFacadeSuite struct {
	suite.Suite

	answerProcessorMock *mocks.MockanswerProcessor

	facade *matchfacade.MatchFacade
}

func (s *MatchFacadeSuite) SetupTest() {
	s.answerProcessorMock = mocks.NewMockanswerProcessor(s.T())

	s.facade = matchfacade.NewMatchFacade(
		s.answerProcessorMock,
	)
}

func (s *MatchFacadeSuite) Test_ProcessAnswer() {
	testArgs := []struct {
		name          string
		params        matchfacade.ProcessAnswerParams
		expectedErr   error
		processParams *answerprocessor.ProcessAnswerParams
		processErr    error
	}{
		{
			name: "err: process answer error",
			params: matchfacade.ProcessAnswerParams{
				PersonID:   "person",
				QuestionID: "question",
				AnswerID:   "answer",
			},
			expectedErr: test.Err,
			processParams: &answerprocessor.ProcessAnswerParams{
				PersonID:   "person",
				QuestionID: "question",
				AnswerID:   "answer",
			},
			processErr: test.Err,
		},
		{
			name: "ok",
			params: matchfacade.ProcessAnswerParams{
				PersonID:   "person",
				QuestionID: "question",
				AnswerID:   "answer",
			},
			processParams: &answerprocessor.ProcessAnswerParams{
				PersonID:   "person",
				QuestionID: "question",
				AnswerID:   "answer",
			},
		},
	}

	for _, testArg := range testArgs {
		testArg := testArg

		s.Run(testArg.name, func() {
			if testArg.processParams != nil {
				s.answerProcessorMock.EXPECT().
					ProcessAnswer(mock.Anything, *testArg.processParams).
					Return(testArg.processErr).Once()
			}

			err := s.facade.ProcessAnswer(test.NewContext(s.T()), testArg.params)
			s.ErrorIs(err, testArg.expectedErr)
		})
	}
}

func TestMatchFacadeSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(MatchFacadeSuite))
}
