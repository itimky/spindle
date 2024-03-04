package answerprocessorstorage

import (
	"context"
	"fmt"

	"github.com/itimky/spindle/pkg/domain"
	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
)

type MatrixInMem struct {
	readiness readiness
	rwLocker  rwLocker
	matrices  map[domain.QuestionID]answerprocessor.WeightMatrix
}

func NewMatrixInMem(
	readiness readiness,
	rwLocker rwLocker,
) *MatrixInMem {
	return &MatrixInMem{
		readiness: readiness,
		rwLocker:  rwLocker,
		matrices:  make(map[domain.QuestionID]answerprocessor.WeightMatrix),
	}
}

func NewMatrixInMemFromSnapshot(
	readiness readiness,
	rwLocker rwLocker,
	matrices map[domain.QuestionID]answerprocessor.WeightMatrix,
) *MatrixInMem {
	return &MatrixInMem{
		readiness: readiness,
		rwLocker:  rwLocker,
		matrices:  matrices,
	}
}

func (s *MatrixInMem) Get(
	ctx context.Context,
	params answerprocessor.GetWeightMatrixParams,
) (answerprocessor.WeightMatrix, error) {
	s.rwLocker.RLock()
	defer s.rwLocker.RUnlock()

	if err := s.readiness.WaitReady(ctx); err != nil {
		return nil, fmt.Errorf("wait ready: %w", err)
	}

	matrix, ok := s.matrices[params.QuestionID]
	if !ok {
		return nil, ErrMatrixNotFound
	}

	return matrix, nil
}

func (s *MatrixInMem) Bootstrap(
	ctx context.Context,
	matrices map[domain.QuestionID]answerprocessor.WeightMatrix,
) error {
	s.matrices = matrices

	s.readiness.MarkReady(ctx)

	return nil
}

func (s *MatrixInMem) Set(
	ctx context.Context,
	questionID domain.QuestionID,
	matrix answerprocessor.WeightMatrix,
) error {
	s.rwLocker.Lock()
	defer s.rwLocker.Unlock()

	if err := s.readiness.WaitReady(ctx); err != nil {
		return fmt.Errorf("wait ready: %w", err)
	}

	s.matrices[questionID] = matrix

	return nil
}
