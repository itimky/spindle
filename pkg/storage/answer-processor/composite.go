package answerprocessorstorage

import (
	"context"
	"fmt"

	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
)

type Composite struct {
	matrixStore   matrixStore
	relationStore relationStore
}

func NewComposite(
	matrixStore matrixStore,
	relationStore relationStore,
) *Composite {
	return &Composite{
		matrixStore:   matrixStore,
		relationStore: relationStore,
	}
}

func (s *Composite) GetWeightMatrix(
	ctx context.Context,
	params answerprocessor.GetWeightMatrixParams,
) (*answerprocessor.GetWeightMatrixResult, error) {
	matrix, err := s.matrixStore.Get(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("matrix store get: %w", err)
	}

	return convertWeightMatrixToGetQuestionWeightMatrixResult(matrix), nil
}

func (s *Composite) GetOtherPersonAnswers(
	ctx context.Context,
	params answerprocessor.GetOtherPersonAnswersParams,
) (*answerprocessor.GetOtherPersonAnswersResult, error) {
	answers, err := s.relationStore.GetOtherPersonAnswers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("get other person answers: %w", err)
	}

	return answers, nil
}

func (s *Composite) UpdatePersonsWeights(
	ctx context.Context,
	params answerprocessor.UpdatePersonsWeightsParams,
) error {
	if err := s.relationStore.UpdatePersonsWeights(ctx, params); err != nil {
		return fmt.Errorf("update persons weights: %w", err)
	}

	return nil
}
