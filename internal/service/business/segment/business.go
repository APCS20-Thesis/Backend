package segment

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type Business interface {
	CreateMasterSegment(ctx context.Context, request *api.CreateMasterSegmentRequest, accountUuid string) error
	ListMasterSegments(ctx context.Context, request *api.GetListMasterSegmentsRequest, accountUuid string) (*api.GetListMasterSegmentsResponse, error)

	CreateSegment(ctx context.Context, request *api.CreateSegmentRequest, accountUuid string) error
	ListSegments(ctx context.Context, request *api.GetListSegmentsRequest, accountUuid string) ([]*api.Segment, error)
	GetSegmentDetail(ctx context.Context, request *api.GetSegmentDetailRequest, accountUuid string) (*api.GetSegmentDetailResponse, error)
	ProcessApplyPredictModel(ctx context.Context, request *api.ApplyPredictModelRequest, accountUuid string) (*api.ApplyPredictModelResponse, error)
	ProcessGetListPredictionActions(ctx context.Context, request *api.GetListPredictionActionsRequest, accountUuid string) (*api.GetListPredictionActionsResponse, error)
	ProcessGetResultPredictionActions(ctx context.Context, request *api.GetResultPredictionActionsRequest, accountUuid string) (*api.GetResultPredictionActionsResponse, error)

	GetMasterSegmentDetail(ctx context.Context, request *api.GetMasterSegmentDetailRequest, accountUuid string) (*api.MasterSegmentDetail, error)
	ListMasterSegmentProfiles(ctx context.Context, request *api.GetListMasterSegmentProfilesRequest, accountUuid string) (int64, []string, error)
	GetMasterSegmentProfile(ctx context.Context, request *api.GetMasterSegmentProfileRequest, accountUuid string) (string, error)
	GetBehaviorProfile(ctx context.Context, request *api.GetBehaviorProfileRequest, accountUuid string) (*BehaviorProfileRecords, error)
	TotalProfilesMasterSegment(ctx context.Context, request *api.TotalProfilesMasterSegmentRequest, accountUuid string) (int64, error)
	SyncOnCreateMasterSegment(ctx context.Context, masterSegmentId int64, actionStatus model.DataActionStatus) error
}

type business struct {
	db             *gorm.DB
	log            logr.Logger
	repository     *repository.Repository
	airflowAdapter airflow.AirflowAdapter
	queryAdapter   query.QueryAdapter
	config         *config.Config
}

func NewSegmentBusiness(db *gorm.DB, log logr.Logger, repository *repository.Repository, airflowAdapter airflow.AirflowAdapter, queryAdapter query.QueryAdapter, config *config.Config) Business {
	return &business{
		db:             db,
		log:            log.WithName("SegmentBiz"),
		repository:     repository,
		airflowAdapter: airflowAdapter,
		queryAdapter:   queryAdapter,
		config:         config,
	}
}
