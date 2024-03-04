package systemfacade

import (
	"github.com/itimky/spindle/pkg/domain"
	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
)

type QuestionWeightMatrix struct {
	QuestionID domain.QuestionID
	Matrix     answerprocessor.WeightMatrix
}

type BootstrapWeightMatrixStorageParams struct {
	Matrices []QuestionWeightMatrix
}

type UpdateWeightMatrixStorageParams QuestionWeightMatrix
