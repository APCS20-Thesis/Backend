package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
)

func (s *Service) CreateDataSourceMySQL(ctx context.Context, request *api.CreateDataSourceMySQLRequest) (*api.CreateDataSourceMySQLResponse, error) {
	return s.business.DataSourceBusiness.ProcessCreateDataSourceMySQL(ctx, request)
}
