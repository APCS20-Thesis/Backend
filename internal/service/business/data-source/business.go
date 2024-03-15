package data_source

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
)

type Business interface {
	ProcessImportFile(ctx context.Context, request *api.ImportFileRequest, accountUuid string, dateTime string) error

	TriggerAirflowGenerateImportFile(ctx context.Context, request *api.ImportFileRequest, accountUuid string, dateTime string) error
	CreateDataActionImportFile(ctx context.Context, accountUuid string, dateTime string) (*model.DataAction, error)

	CreateDataActionRun(ctx context.Context, params *repository.CreateDataActionRunParams) (*model.DataActionRun, error)

	CreateDataSource(ctx context.Context, params *repository.CreateDataSourceParams) (*model.DataSource, error)
	GetDataSource(ctx context.Context, request *api.GetDataSourceRequest, accountUuid string) (*api.GetDataSourceResponse, error)
	GetListDataSources(ctx context.Context, request *api.GetListDataSourcesRequest, accountUuid string) ([]*api.DataSourceBase, error)

	CreateSourceConnection(ctx context.Context, params *repository.CreateSourceConnectionParams) (*model.SourceConnection, error)
	UpdateSourceConnection(ctx context.Context, params *repository.UpdateSourceConnectionParams) error
	GetSourceConnection(ctx context.Context, request *api.GetSourceConnectionRequest, accountUuid string) (*api.GetSourceConnectionResponse, error)
	GetListSourceConnections(ctx context.Context, request *api.GetListSourceConnectionsRequest, accountUuid string) ([]*api.SourceConnectionBase, error)
}

type business struct {
	log            logr.Logger
	repository     *repository.Repository
	airflowAdapter airflow.AirflowAdapter
}

func NewDataSourceBusiness(log logr.Logger, repository *repository.Repository, airflowAdapter airflow.AirflowAdapter) Business {
	return &business{
		log:            log.WithName("DataSourceBiz"),
		repository:     repository,
		airflowAdapter: airflowAdapter,
	}
}
