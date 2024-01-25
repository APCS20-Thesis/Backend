package data_source

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
)

type Business interface {
	ProcessCreateDataSourceMySQL(ctx context.Context, request *api.CreateDataSourceMySQLRequest) (*api.CreateDataSourceMySQLResponse, error)
}

type business struct {
	log        logr.Logger
	repository *repository.Repository
}

func NewDataSourceBusiness(log logr.Logger, repository *repository.Repository) Business {
	return &business{
		log:        log.WithName("DataSourceBiz"),
		repository: repository,
	}
}
