package kafkahandler

import (
	"context"
	"encoding/json"
	"fmt"

	kafkacontract "github.com/itimky/spindle/pkg/contract/kafka"
	systemfacade "github.com/itimky/spindle/pkg/facade/system"
	"github.com/itimky/spindle/pkg/sys"
	"github.com/itimky/spindle/pkg/sys/queue"
)

type WeightMatrixBootstrapHandler struct {
	systemFacade systemFacade
}

func NewWeightMatrixBootstrapHandler(
	systemFacade systemFacade,
) *WeightMatrixBootstrapHandler {
	return &WeightMatrixBootstrapHandler{
		systemFacade: systemFacade,
	}
}

func (h *WeightMatrixBootstrapHandler) Bootstrap(
	ctx context.Context,
	msgs []queue.Message,
) error {
	weightMatrices := make([]systemfacade.QuestionWeightMatrix, 0, len(msgs))

	for _, msg := range msgs {
		var event kafkacontract.WeightMatrixV1

		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			return fmt.Errorf("%w: %s", sys.ErrInvalidJSON, err.Error())
		}

		weightMatrix, err := convertWeightMatrixV1ToQuestionWeightMatrix(event)
		if err != nil {
			return fmt.Errorf("convert weight matrix v1 to update weight matrix storage params: %w", err)
		}

		weightMatrices = append(weightMatrices, *weightMatrix)
	}

	err := h.systemFacade.BootstrapWeightMatrixStorage(
		ctx,
		convertQuestionWeightMatricesToBootstrapWeightMatrixStorageParams(weightMatrices),
	)
	if err != nil {
		return fmt.Errorf("bootstrap weight matrix storage: %w", err)
	}

	return nil
}

func (h *WeightMatrixBootstrapHandler) Handle(
	ctx context.Context,
	msg queue.Message,
) error {
	var event kafkacontract.WeightMatrixV1

	err := json.Unmarshal(msg.Data, &event)
	if err != nil {
		return fmt.Errorf("%w: %s", sys.ErrInvalidJSON, err.Error())
	}

	params, err := convertWeightMatrixV1ToQuestionWeightMatrix(event)
	if err != nil {
		return fmt.Errorf("convert weight matrix v1 to update weight matrix storage params: %w", err)
	}

	err = h.systemFacade.UpdateWeightMatrixStorage(ctx, systemfacade.UpdateWeightMatrixStorageParams(*params))
	if err != nil {
		return fmt.Errorf("process question matrix: %w", err)
	}

	return nil
}
