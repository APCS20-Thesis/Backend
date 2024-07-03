package predict_model

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type Business interface {
	ProcessTrainPredictModel(ctx context.Context, request *api.TrainPredictModelRequest, accountUuid string) error
	ProcessGetListPredictModels(ctx context.Context, request *api.GetListPredictModelsRequest, accountUuid string) (*api.GetListPredictModelsResponse, error)
	ProcessGetPredictModelDetail(ctx context.Context, request *api.GetPredictModelDetailRequest, accountUuid string) (*api.GetPredictModelDetailResponse, error)
}

type business struct {
	db             *gorm.DB
	log            logr.Logger
	repository     *repository.Repository
	airflowAdapter airflow.AirflowAdapter
}

func NewPredictModelBusiness(db *gorm.DB, log logr.Logger, repository *repository.Repository, airflowAdapter airflow.AirflowAdapter) Business {
	return &business{
		db:             db,
		log:            log.WithName("PredictModelBiz"),
		repository:     repository,
		airflowAdapter: airflowAdapter,
	}
}
