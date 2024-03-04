package systemfacade

import (
	"context"
	"fmt"
)

type SystemFacade struct {
	weightMatrixStore weightMatrixStore
}

func NewSystemFacade(
	weightMatrixStore weightMatrixStore,
) *SystemFacade {
	return &SystemFacade{
		weightMatrixStore: weightMatrixStore,
	}
}

func (p *SystemFacade) BootstrapWeightMatrixStorage(
	ctx context.Context,
	params BootstrapWeightMatrixStorageParams,
) error {
	err := p.weightMatrixStore.Bootstrap(ctx, convertBootstrapWeightMatrixStorageParamsToQuestionWeightMatrixMap(params))
	if err != nil {
		return fmt.Errorf("bootstrap: %w", err)
	}

	return nil
}

func (p *SystemFacade) UpdateWeightMatrixStorage(
	ctx context.Context,
	params UpdateWeightMatrixStorageParams,
) error {
	err := p.weightMatrixStore.Set(ctx, params.QuestionID, params.Matrix)
	if err != nil {
		return fmt.Errorf("set: %w", err)
	}

	return nil
}
