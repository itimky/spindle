package matchfacade

import (
	"context"

	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
)

type answerProcessor interface {
	ProcessAnswer(ctx context.Context, params answerprocessor.ProcessAnswerParams) error
}
