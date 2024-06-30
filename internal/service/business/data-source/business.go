package data_source

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Business interface {
	ProcessImportCsv(ctx context.Context, request *api.ImportCsvRequest, accountUuid string, dateTime string) error

	CreateDataActionRun(ctx context.Context, params *repository.CreateDataActionRunParams) (*model.DataActionRun, error)

	CreateDataSource(ctx context.Context, params *repository.CreateDataSourceParams) (*model.DataSource, error)
	GetDataSource(ctx context.Context, request *api.GetDataSourceRequest, accountUuid string) (*api.GetDataSourceResponse, error)
	GetListDataSources(ctx context.Context, request *api.GetListDataSourcesRequest, accountUuid string) (*api.GetListDataSourcesResponse, error)

	ProcessImportCsvFromS3(ctx context.Context, request *api.ImportCsvFromS3Request, accountUuid string, dateTime string) error
	ProcessImportFromMySQLSource(ctx context.Context, request *api.ImportFromMySQLSourceRequest, accountUuid uuid.UUID) error

	GetListSourceTableMappings(ctx context.Context, request *api.GetListSourceTableMapRequest) (*api.GetListSourceTableMapResponse, error)
}

type business struct {
	db             *gorm.DB
	log            logr.Logger
	repository     *repository.Repository
	airflowAdapter airflow.AirflowAdapter
	config         *config.Config
}

func NewDataSourceBusiness(db *gorm.DB, log logr.Logger, repository *repository.Repository, airflowAdapter airflow.AirflowAdapter, config *config.Config) Business {

	return &business{
		db:             db,
		log:            log.WithName("DataSourceBiz"),
		repository:     repository,
		airflowAdapter: airflowAdapter,
		config:         config,
	}
}
