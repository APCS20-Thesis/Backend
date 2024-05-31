package data_source

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
)

type Business interface {
	ProcessImportCsv(ctx context.Context, request *api.ImportCsvRequest, accountUuid string, dateTime string) error

	CreateDataActionRun(ctx context.Context, params *repository.CreateDataActionRunParams) (*model.DataActionRun, error)

	CreateDataSource(ctx context.Context, params *repository.CreateDataSourceParams) (*model.DataSource, error)
	GetDataSource(ctx context.Context, request *api.GetDataSourceRequest, accountUuid string) (*api.GetDataSourceResponse, error)
	GetListDataSources(ctx context.Context, request *api.GetListDataSourcesRequest, accountUuid string) ([]*api.GetListDataSourcesResponse_DataSource, error)

	CreateConnection(ctx context.Context, request *api.CreateConnectionRequest, accountUuid string) (*model.Connection, error)
	UpdateConnection(ctx context.Context, params *repository.UpdateConnectionParams) error
	GetConnection(ctx context.Context, request *api.GetConnectionRequest, accountUuid string) (*api.GetConnectionResponse, error)
	GetListConnections(ctx context.Context, request *api.GetListConnectionsRequest, accountUuid string) ([]*api.GetListConnectionsResponse_Connection, error)
	DeleteConnection(ctx context.Context, request *api.DeleteConnectionRequest, accountUuid string) error

	ProcessImportCsvFromS3(ctx context.Context, request *api.ImportCsvFromS3Request, accountUuid string, dateTime string) error
}

type business struct {
	log            logr.Logger
	repository     *repository.Repository
	airflowAdapter airflow.AirflowAdapter
	config         *config.Config
}

func NewDataSourceBusiness(log logr.Logger, repository *repository.Repository, airflowAdapter airflow.AirflowAdapter, config *config.Config) Business {

	return &business{
		log:            log.WithName("DataSourceBiz"),
		repository:     repository,
		airflowAdapter: airflowAdapter,
		config:         config,
	}
}
