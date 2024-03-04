package answerprocessorstorage

import answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"

func convertWeightMatrixToGetQuestionWeightMatrixResult(
	matrix answerprocessor.WeightMatrix,
) *answerprocessor.GetWeightMatrixResult {
	return &answerprocessor.GetWeightMatrixResult{
		Matrix: matrix,
	}
}
