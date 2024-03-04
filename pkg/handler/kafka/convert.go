package kafkahandler

import (
	"fmt"

	kafkacontract "github.com/itimky/spindle/pkg/contract/kafka"
	"github.com/itimky/spindle/pkg/domain"
	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
	matchfacade "github.com/itimky/spindle/pkg/facade/match"
	systemfacade "github.com/itimky/spindle/pkg/facade/system"
	"github.com/itimky/spindle/pkg/sys"
	"github.com/shopspring/decimal"
)

func convertAnswerV1ToProcessAnswerParams(
	event kafkacontract.AnswerV1,
) matchfacade.ProcessAnswerParams {
	return matchfacade.ProcessAnswerParams{
		PersonID:   domain.PersonID(event.PersonID),
		QuestionID: domain.QuestionID(event.QuestionID),
		AnswerID:   domain.AnswerID(event.AnswerID),
	}
}

func convertWeightMatrixV1ToQuestionWeightMatrix(
	event kafkacontract.WeightMatrixV1,
) (*systemfacade.QuestionWeightMatrix, error) {
	matrix := make(answerprocessor.WeightMatrix, len(event.Matrix))

	for iQID := range event.Matrix {
		row := make(map[domain.AnswerID]decimal.Decimal)

		for jQID := range event.Matrix[iQID] {
			weight, err := decimal.NewFromString(event.Matrix[iQID][jQID])
			if err != nil {
				return nil, fmt.Errorf("%w: %s", sys.ErrInvalidDecimalString, err.Error())
			}

			row[domain.AnswerID(jQID)] = weight
		}

		matrix[domain.AnswerID(iQID)] = row
	}

	return &systemfacade.QuestionWeightMatrix{
		QuestionID: domain.QuestionID(event.QuestionID),
		Matrix:     matrix,
	}, nil
}

func convertQuestionWeightMatricesToBootstrapWeightMatrixStorageParams(
	matrices []systemfacade.QuestionWeightMatrix,
) systemfacade.BootstrapWeightMatrixStorageParams {
	return systemfacade.BootstrapWeightMatrixStorageParams{
		Matrices: matrices,
	}
}
