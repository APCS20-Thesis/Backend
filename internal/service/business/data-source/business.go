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
	CreateDataActionRun(ctx context.Context, params *repository.CreateDataActionRunParams) error
	CreateDataSource(ctx context.Context, params *repository.CreateDataSourceParams) error
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
