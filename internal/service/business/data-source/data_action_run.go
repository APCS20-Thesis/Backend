package data_source

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
)

func (b business) CreateDataActionRun(ctx context.Context, params *repository.CreateDataActionRunParams) (*model.DataActionRun, error) {
	dataActionRun, err := b.repository.DataActionRunRepository.CreateDataActionRun(ctx, params)
	if err != nil {
		b.log.WithName("CreateDataActionRun").
			WithValues("Context", ctx).
			Error(err, "Cannot create data action run")
		return nil, err
	}
	return dataActionRun, nil
}
