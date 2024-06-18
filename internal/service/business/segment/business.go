package segment

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type Business interface {
	CreateMasterSegment(ctx context.Context, request *api.CreateMasterSegmentRequest, accountUuid string) error
	ListMasterSegments(ctx context.Context, request *api.GetListMasterSegmentsRequest, accountUuid string) (int64, []*api.MasterSegment, error)

	CreateSegment(ctx context.Context, request *api.CreateSegmentRequest, accountUuid string) error
	ListSegments(ctx context.Context, request *api.GetListSegmentsRequest, accountUuid string) ([]*api.Segment, error)
	GetSegmentDetail(ctx context.Context, request *api.GetSegmentDetailRequest, accountUuid string) (*api.GetSegmentDetailResponse, error)

	GetMasterSegmentDetail(ctx context.Context, request *api.GetMasterSegmentDetailRequest, accountUuid string) (*api.MasterSegmentDetail, error)
}

type business struct {
	db             *gorm.DB
	log            logr.Logger
	repository     *repository.Repository
	airflowAdapter airflow.AirflowAdapter
}

func NewSegmentBusiness(db *gorm.DB, log logr.Logger, repository *repository.Repository, airflowAdapter airflow.AirflowAdapter) Business {
	return &business{
		db:             db,
		log:            log.WithName("SegmentBiz"),
		repository:     repository,
		airflowAdapter: airflowAdapter,
	}
}
