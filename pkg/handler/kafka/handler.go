package kafkahandler

import (
	"context"
	"encoding/json"
	"fmt"

	kafkacontract "github.com/itimky/spindle/pkg/contract/kafka"
	"github.com/itimky/spindle/pkg/sys"
	"github.com/itimky/spindle/pkg/sys/queue"
)

type Handler struct {
	matchFacade matchFacade
}

func NewHandler(
	matchFacade matchFacade,
) *Handler {
	return &Handler{
		matchFacade: matchFacade,
	}
}

func (h *Handler) handleAnswerV1(ctx context.Context, msg queue.Message) error {
	var event kafkacontract.AnswerV1

	err := json.Unmarshal(msg.Data, &event)
	if err != nil {
		return fmt.Errorf("%w: %s", sys.ErrInvalidJSON, err.Error())
	}

	err = h.matchFacade.ProcessAnswer(ctx, convertAnswerV1ToProcessAnswerParams(event))
	if err != nil {
		return fmt.Errorf("process answer: %w", err)
	}

	return nil
}

func (h *Handler) Routes() map[queue.MessageType]queue.HandlerFunc {
	return map[queue.MessageType]queue.HandlerFunc{
		kafkacontract.MessageTypeAnswerV1: h.handleAnswerV1,
	}
}
