package answerprocessorstorage

import (
	"context"

	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
)

type rwLocker interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

type readiness interface {
	WaitReady(ctx context.Context) error
	MarkReady(ctx context.Context)
}

type matrixStore interface {
	Get(
		ctx context.Context,
		questionID answerprocessor.GetWeightMatrixParams,
	) (answerprocessor.WeightMatrix, error)
}

type relationStore interface {
	GetOtherPersonAnswers(
		ctx context.Context,
		params answerprocessor.GetOtherPersonAnswersParams,
	) (*answerprocessor.GetOtherPersonAnswersResult, error)
	UpdatePersonsWeights(
		ctx context.Context,
		params answerprocessor.UpdatePersonsWeightsParams,
	) error
}
