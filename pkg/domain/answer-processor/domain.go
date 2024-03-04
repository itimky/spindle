package answerprocessor

import (
	"context"
	"fmt"
)

type AnswerProcessor struct {
	storage storage
}

func NewAnswerProcessor(
	repo storage,
) *AnswerProcessor {
	return &AnswerProcessor{
		storage: repo,
	}
}

func (p *AnswerProcessor) ProcessAnswer(
	ctx context.Context,
	params ProcessAnswerParams,
) error {
	matrixRes, err := p.storage.GetWeightMatrix(ctx, convertProcessAnswerParamsToGetWeightMatrixParams(params))
	if err != nil {
		return fmt.Errorf("get question weight matrix: %w", err)
	}

	answersRes, err := p.storage.GetOtherPersonAnswers(
		ctx,
		convertProcessAnswerParamsToGetOtherPersonAnswersParams(params),
	)
	if err != nil {
		return fmt.Errorf("get same level person answers: %w", err)
	}

	personsWeights := make([]RelatedPersonWeight, 0, len(answersRes.PersonAnswers))

	for i := range answersRes.PersonAnswers {
		personAnswer := &answersRes.PersonAnswers[i]
		weight := matrixRes.Matrix[params.AnswerID][personAnswer.AnswerID]
		personsWeights = append(personsWeights, RelatedPersonWeight{
			RelatedPersonID: personAnswer.PersonID,
			Weight:          weight,
		})
	}

	err = p.storage.UpdatePersonsWeights(ctx, convertPersonsWeightsToUpdatePersonsWeightsParams(params, personsWeights))
	if err != nil {
		return fmt.Errorf("update persons weights: %w", err)
	}

	return nil
}
