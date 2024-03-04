package answerprocessor

import (
	"github.com/itimky/spindle/pkg/domain"
	"github.com/shopspring/decimal"
)

type WeightRow map[domain.AnswerID]decimal.Decimal

type WeightMatrix map[domain.AnswerID]WeightRow

type PersonAnswer struct {
	PersonID domain.PersonID
	AnswerID domain.AnswerID
}

type RelatedPersonWeight struct {
	RelatedPersonID domain.PersonID
	Weight          decimal.Decimal
}

type ProcessAnswerParams struct {
	PersonID   domain.PersonID
	QuestionID domain.QuestionID
	AnswerID   domain.AnswerID
}

type GetWeightMatrixParams struct {
	QuestionID domain.QuestionID
}

type GetWeightMatrixResult struct {
	Matrix WeightMatrix
}

type GetOtherPersonAnswersParams struct {
	QuestionID domain.QuestionID
}

type GetOtherPersonAnswersResult struct {
	PersonAnswers []PersonAnswer
}

type UpdatePersonsWeightsParams struct {
	PersonID       domain.PersonID
	QuestionID     domain.QuestionID
	RelatedWeights []RelatedPersonWeight
}
