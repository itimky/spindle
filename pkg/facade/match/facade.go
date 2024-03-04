package matchfacade

import (
	"context"
	"fmt"

	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
)

type MatchFacade struct {
	answerProcessor answerProcessor
}

func NewMatchFacade(
	answerProcessor answerProcessor,
) *MatchFacade {
	return &MatchFacade{
		answerProcessor: answerProcessor,
	}
}

func (f *MatchFacade) ProcessAnswer(
	ctx context.Context,
	params ProcessAnswerParams,
) error {
	err := f.answerProcessor.ProcessAnswer(ctx, answerprocessor.ProcessAnswerParams(params))
	if err != nil {
		return fmt.Errorf("process answer: %w", err)
	}

	return nil
}
