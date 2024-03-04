package answerprocessorstorage_test

import (
	"testing"

	"github.com/itimky/spindle/pkg/domain"
	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
	answerprocessorstorage "github.com/itimky/spindle/pkg/storage/answer-processor"
	"github.com/itimky/spindle/test"
	mocks "github.com/itimky/spindle/test/pkg/storage/answer-processor"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MatrixImMemSuite struct {
	suite.Suite

	readinessMock *mocks.Mockreadiness
	rwLockerMock  *mocks.MockrwLocker
}

func (s *MatrixImMemSuite) SetupTest() {
	s.rwLockerMock = mocks.NewMockrwLocker(s.T())
	s.readinessMock = mocks.NewMockreadiness(s.T())
}

func (s *MatrixImMemSuite) Test_Get() {
	testArgs := []struct {
		name           string
		mMap           map[domain.QuestionID]answerprocessor.WeightMatrix
		params         answerprocessor.GetWeightMatrixParams
		readyErr       error
		expectedResult answerprocessor.WeightMatrix
		expectedErr    error
	}{
		{
			name:        "err: wait ready error",
			mMap:        nil,
			expectedErr: test.Err,
			params: answerprocessor.GetWeightMatrixParams{
				QuestionID: "q1",
			},
			readyErr: test.Err,
		},
		{
			name: "err: not found: empty map",
			mMap: map[domain.QuestionID]answerprocessor.WeightMatrix{},
			params: answerprocessor.GetWeightMatrixParams{
				QuestionID: "q1",
			},
			expectedErr: answerprocessorstorage.ErrMatrixNotFound,
		},
		{
			name: "err: not found: non-empty map",
			mMap: map[domain.QuestionID]answerprocessor.WeightMatrix{
				"q2": {},
			},
			params: answerprocessor.GetWeightMatrixParams{
				QuestionID: "q1",
			},
			expectedErr: answerprocessorstorage.ErrMatrixNotFound,
		},
		{
			name: "ok: found",
			mMap: map[domain.QuestionID]answerprocessor.WeightMatrix{
				"q1": {
					"a1": {
						"a1": decimal.RequireFromString("0.1"),
					},
				},
			},
			params: answerprocessor.GetWeightMatrixParams{
				QuestionID: "q1",
			},
			expectedResult: answerprocessor.WeightMatrix{
				"a1": {
					"a1": decimal.RequireFromString("0.1"),
				},
			},
		},
	}

	for _, args := range testArgs {
		args := args

		s.Run(args.name, func() {
			s.rwLockerMock.EXPECT().RLock().Once()
			s.rwLockerMock.EXPECT().RUnlock().Once()

			s.readinessMock.EXPECT().WaitReady(mock.Anything).Return(args.readyErr).Once()

			matrixInMem := answerprocessorstorage.NewMatrixInMemFromSnapshot(
				s.readinessMock,
				s.rwLockerMock,
				args.mMap,
			)

			result, err := matrixInMem.Get(test.NewContext(s.T()), args.params)
			s.ErrorIs(err, args.expectedErr)
			s.EqualValues(args.expectedResult, result)
		})
	}
}

func (s *MatrixImMemSuite) Test_Bootstrap_Smoke() {
	matrixInMem := answerprocessorstorage.NewMatrixInMem(
		s.readinessMock,
		s.rwLockerMock,
	)

	s.readinessMock.EXPECT().MarkReady(mock.Anything).Once()

	err := matrixInMem.Bootstrap(test.NewContext(s.T()), map[domain.QuestionID]answerprocessor.WeightMatrix{
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
	})
	s.NoError(err)
}

func (s *MatrixImMemSuite) Test_Set() {
	testArgs := []struct {
		name        string
		mMap        map[domain.QuestionID]answerprocessor.WeightMatrix
		questionID  domain.QuestionID
		matrix      answerprocessor.WeightMatrix
		mMapAfter   map[domain.QuestionID]answerprocessor.WeightMatrix
		expectedErr error
		waitErr     error
	}{
		{
			name:        "err: wait ready error",
			mMap:        nil,
			expectedErr: test.Err,
			questionID:  "q1",
			matrix: answerprocessor.WeightMatrix{
				"a1": {
					"a1": decimal.RequireFromString("0.1"),
				},
			},
			waitErr: test.Err,
		},
		{
			name: "ok",
			mMap: map[domain.QuestionID]answerprocessor.WeightMatrix{
				"q1": {
					"a1": {
						"a1": decimal.RequireFromString("0.1"),
					},
				},
			},
			mMapAfter: map[domain.QuestionID]answerprocessor.WeightMatrix{
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
			questionID: "q2",
			matrix: answerprocessor.WeightMatrix{
				"a1": {
					"a1": decimal.RequireFromString("0.2"),
				},
			},
		},
	}

	for _, args := range testArgs {
		args := args

		s.Run(args.name, func() {
			s.rwLockerMock.EXPECT().Lock().Once()
			s.rwLockerMock.EXPECT().Unlock().Once()

			s.readinessMock.EXPECT().WaitReady(mock.Anything).Return(args.waitErr).Once()

			matrixInMem := answerprocessorstorage.NewMatrixInMemFromSnapshot(
				s.readinessMock,
				s.rwLockerMock,
				args.mMap,
			)

			err := matrixInMem.Set(test.NewContext(s.T()), args.questionID, args.matrix)
			s.ErrorIs(err, args.expectedErr)
			s.EqualValues(args.mMapAfter, args.mMap)
		})
	}
}

func TestMatrixImMemSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(MatrixImMemSuite))
}
