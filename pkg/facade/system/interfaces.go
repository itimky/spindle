package systemfacade

import (
	"context"

	"github.com/itimky/spindle/pkg/domain"
	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
)

type weightMatrixStore interface {
	Bootstrap(
		ctx context.Context,
		matrices map[domain.QuestionID]answerprocessor.WeightMatrix,
	) error
	Set(
		ctx context.Context,
		questionID domain.QuestionID,
		matrix answerprocessor.WeightMatrix,
	) error
}
