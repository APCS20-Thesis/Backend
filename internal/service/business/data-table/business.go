package data_table

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
)

type Business interface {
	CreateDataTable(ctx context.Context, params *repository.CreateDataTableParams) (*model.DataTable, error)
	UpdateDataTable(ctx context.Context, params *repository.UpdateDataTableParams) error
	GetDataTable(ctx context.Context, request *api.GetDataTableRequest, accountUuid string) (*api.GetDataTableResponse, error)
	GetListDataTables(ctx context.Context, request *api.GetListDataTablesRequest, accountUuid string) ([]*api.GetListDataTablesResponse_DataTable, int64, error)
	ExportDataTableToFile(ctx context.Context, request *api.ExportDataTableToFileRequest, accountUuid string) (*api.ExportDataTableToFileResponse, error)
	GetListFileExportRecords(ctx context.Context, request *api.GetListFileExportRecordsRequest, accountUuid string) ([]*api.GetListFileExportRecordsResponse_FileExportRecord, error)
	GetQueryDataTable(ctx context.Context, request *api.GetQueryDataTableRequest, accountUuid string) (*api.GetQueryDataTableResponse, error)
}

type business struct {
	log            logr.Logger
	repository     *repository.Repository
	airflowAdapter airflow.AirflowAdapter
	queryAdapter   query.QueryAdapter
}

func NewDataTableBusiness(log logr.Logger, repository *repository.Repository, airflowAdapter airflow.AirflowAdapter, queryAdapter query.QueryAdapter) Business {
	return &business{
		log:            log.WithName("DataTableBiz"),
		repository:     repository,
		airflowAdapter: airflowAdapter,
		queryAdapter:   queryAdapter,
	}
}
