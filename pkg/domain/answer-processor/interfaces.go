package answerprocessor

import "context"

type storage interface {
	GetWeightMatrix(
		ctx context.Context,
		params GetWeightMatrixParams,
	) (*GetWeightMatrixResult, error)
	GetOtherPersonAnswers(
		ctx context.Context,
		params GetOtherPersonAnswersParams,
	) (*GetOtherPersonAnswersResult, error)
	UpdatePersonsWeights(
		ctx context.Context,
		params UpdatePersonsWeightsParams,
	) error
}
