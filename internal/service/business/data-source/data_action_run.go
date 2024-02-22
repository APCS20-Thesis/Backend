package data_source

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/repository"
)

func (b business) CreateDataActionRun(ctx context.Context, params *repository.CreateDataActionRunParams) error {
	err := b.repository.DataActionRunRepository.CreateDataActionRun(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
