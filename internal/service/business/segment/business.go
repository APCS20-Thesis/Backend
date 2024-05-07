package segment

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
)

type Business interface {
	CreateMasterSegment(ctx context.Context, request *api.CreateMasterSegmentRequest, accountUuid string) error
	ListMasterSegments(ctx context.Context, request *api.GetListMasterSegmentsRequest, accountUuid string) (int64, []*api.MasterSegment, error)
	CreateSegment(ctx context.Context, request *api.CreateSegmentRequest, accountUuid string) error
}

type business struct {
	log            logr.Logger
	repository     *repository.Repository
	airflowAdapter airflow.AirflowAdapter
}

func NewSegmentBusiness(log logr.Logger, repository *repository.Repository, airflowAdapter airflow.AirflowAdapter) Business {
	return &business{
		log:            log.WithName("SegmentBiz"),
		repository:     repository,
		airflowAdapter: airflowAdapter,
	}
}
