package kafkahandler

import (
	"context"

	matchfacade "github.com/itimky/spindle/pkg/facade/match"
	systemfacade "github.com/itimky/spindle/pkg/facade/system"
)

type matchFacade interface {
	ProcessAnswer(ctx context.Context, params matchfacade.ProcessAnswerParams) error
}

type systemFacade interface {
	BootstrapWeightMatrixStorage(
		ctx context.Context,
		params systemfacade.BootstrapWeightMatrixStorageParams,
	) error
	UpdateWeightMatrixStorage(
		ctx context.Context,
		params systemfacade.UpdateWeightMatrixStorageParams,
	) error
}
