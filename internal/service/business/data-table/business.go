package data_table

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
)

type Business interface {
	CreateDataTable(ctx context.Context, params *repository.CreateDataTableParams) (*model.DataTable, error)
	UpdateDataTable(ctx context.Context, params *repository.UpdateDataTableParams) error
	GetDataTable(ctx context.Context, request *api.GetDataTableRequest, accountUuid string) (*api.GetDataTableResponse, error)
	GetListDataTables(ctx context.Context, request *api.GetListDataTablesRequest, accountUuid string) ([]*api.GetListDataTablesResponse_DataTable, error)
}

type business struct {
	log        logr.Logger
	repository *repository.Repository
}

func NewDataTableBusiness(log logr.Logger, repository *repository.Repository) Business {
	return &business{
		log:        log.WithName("DataTableBiz"),
		repository: repository,
	}
}
