package systemfacade

import (
	"github.com/itimky/spindle/pkg/domain"
	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
)

func convertBootstrapWeightMatrixStorageParamsToQuestionWeightMatrixMap(
	params BootstrapWeightMatrixStorageParams,
) map[domain.QuestionID]answerprocessor.WeightMatrix {
	result := make(map[domain.QuestionID]answerprocessor.WeightMatrix)
	for _, questionMatrix := range params.Matrices {
		result[questionMatrix.QuestionID] = questionMatrix.Matrix
	}

	return result
}
