package data_action

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type Business interface {
	ProcessGetListDataActionRuns(ctx context.Context, request *api.GetListDataActionRunsRequest, accountUuid string) (*api.GetListDataActionRunsResponse, error)
	ProcessNewDataActionRun(ctx context.Context, request *api.TriggerDataActionRunRequest, accountUuid string) error
}

type business struct {
	db             *gorm.DB
	log            logr.Logger
	repository     *repository.Repository
	airflowAdapter airflow.AirflowAdapter
}

func NewDataActionBusiness(db *gorm.DB, log logr.Logger, repository *repository.Repository, airflowAdapter airflow.AirflowAdapter) Business {
	return &business{
		db:             db,
		log:            log.WithName("DataActionBiz"),
		repository:     repository,
		airflowAdapter: airflowAdapter,
	}
}
