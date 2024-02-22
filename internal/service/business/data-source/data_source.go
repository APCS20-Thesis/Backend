package data_source

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/repository"
)

func (b business) CreateDataSource(ctx context.Context, params *repository.CreateDataSourceParams) error {
	err := b.repository.DataSourceRepository.CreateDataSource(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
