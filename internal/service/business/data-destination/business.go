package data_destination

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/adapter/gophish"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type Business interface {
	CreateGophishUserGroupFromSegment(ctx context.Context, accountUuid string, request *api.CreateGophishUserGroupFromSegmentRequest) error
	ProcessExportToMySQLDestination(ctx context.Context, request *api.ExportToMySQLDestinationRequest, accountUuid string) error
	ProcessExportDataToFile(ctx context.Context, request *api.ExportDataToFileRequest, accountUuid string) (*api.ExportDataToFileResponse, error)

	ProcessGetListDataDestinations(ctx context.Context, request *api.GetListDataDestinationsRequest, accountUuid string) (*api.GetListDataDestinationsResponse, error)
	ProcessGetDataDestinationDetail(ctx context.Context, request *api.GetDataDestinationDetailRequest, accountUuid string) (*api.GetDataDestinationDetailResponse, error)

	ProcessGetListDestinationMap(ctx context.Context, request *api.GetListDestinationMapRequest, accountUuid string) (*api.GetListDestinationMapResponse, error)
}

type business struct {
	db             *gorm.DB
	log            logr.Logger
	repository     *repository.Repository
	gophishAdapter gophish.GophishAdapter
	queryAdapter   query.QueryAdapter
	airflowAdapter airflow.AirflowAdapter
	config         *config.Config
}

func NewDataDestinationBusiness(db *gorm.DB, log logr.Logger, repository *repository.Repository, gophishAdapter gophish.GophishAdapter, queryAdapter query.QueryAdapter, airflowAdapter airflow.AirflowAdapter, config *config.Config) Business {
	return &business{
		db:             db,
		log:            log.WithName("DataDestinationBiz"),
		repository:     repository,
		gophishAdapter: gophishAdapter,
		queryAdapter:   queryAdapter,
		airflowAdapter: airflowAdapter,
		config:         config,
	}
}
