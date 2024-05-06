package repository

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type SegmentRepository interface {
	CreateMasterSegment(ctx context.Context, params *model.MasterSegment) error
	ListMasterSegments(ctx context.Context, params *ListMasterSegmentsParams) ([]model.MasterSegment, error)
	CreateAudienceTable(ctx context.Context, params *CreateAudienceTableParams) error
	CreateBehaviorTable(ctx context.Context, params *CreateBehaviorTableParams) error
	CreateSegment(ctx context.Context, params *CreateSegmentParams) error
}

type segmentRepo struct {
	*gorm.DB
	MasterSegmentTableName string
	SegmentTableName       string
	AudienceTableName      string
	BehaviorTableName      string
}

func NewSegmentRepository(db *gorm.DB) SegmentRepository {
	return &segmentRepo{db,
		model.MasterSegment{}.TableName(),
		model.Segment{}.TableName(),
		model.AudienceTable{}.TableName(),
		model.BehaviorTable{}.TableName(),
	}
}

type CreateMasterSegmentParams struct {
	Name        string
	Description string
	AccountUuid uuid.UUID
}

func (r *segmentRepo) CreateMasterSegment(ctx context.Context, params *model.MasterSegment) error {
	err := r.WithContext(ctx).Table(r.MasterSegmentTableName).Create(&model.MasterSegment{
		Description: params.Description,
		Name:        params.Name,
		AccountUuid: params.AccountUuid,
		Status:      model.MasterSegmentStatus_DRAFT,
	}).Error

	return err
}

type CreateAudienceTableParams struct {
	MasterSegmentId    int64
	Name               string
	BuildConfiguration AudienceBuildConfiguration
}

type AudienceBuildConfiguration struct {
	MainTableId     int64                    `json:"mainTableId"`
	SelectedColumns []*api.TransferredColumn `json:"selectedColumns"`
	AttributeTables []*AttributeTableInfo    `json:"attributeTables"`
}

type AttributeTableInfo struct {
	TableId         int64                    `json:"tableId"`
	ForeignKey      string                   `json:"foreignKey"`
	JoinKey         string                   `json:"joinKey"`
	SelectedColumns []*api.TransferredColumn `json:"selectedColumns"`
}

func (r *segmentRepo) CreateAudienceTable(ctx context.Context, params *CreateAudienceTableParams) error {
	buildConfiguration, err := json.Marshal(params.BuildConfiguration)
	if err != nil {
		return err
	}

	err = r.WithContext(ctx).Table(r.AudienceTableName).Create(&model.AudienceTable{
		MasterSegmentId:    params.MasterSegmentId,
		BuildConfiguration: pqtype.NullRawMessage{RawMessage: buildConfiguration, Valid: true},
		Name:               params.Name,
	}).Error

	return err
}

type CreateBehaviorTableParams struct {
	MasterSegmentId int64
	Name            string
	TableId         int64
	ForeignKey      string
	JoinKey         string
	SelectedColumns []*api.TransferredColumn
}

func (r *segmentRepo) CreateBehaviorTable(ctx context.Context, params *CreateBehaviorTableParams) error {
	err := r.WithContext(ctx).Table(r.BehaviorTableName).Create(&model.BehaviorTable{
		MasterSegmentId: params.MasterSegmentId,
		DataTableId:     params.TableId,
		ForeignKey:      params.ForeignKey,
		JoinKey:         params.JoinKey,
		Name:            params.Name,
	}).Error

	return err
}

type CreateSegmentParams struct {
	Name            string
	Description     string
	MasterSegmentId int64
	Condition       pqtype.NullRawMessage
	AccountUuid     uuid.UUID
}

func (r *segmentRepo) CreateSegment(ctx context.Context, params *CreateSegmentParams) error {
	err := r.WithContext(ctx).Table(r.SegmentTableName).Create(&model.Segment{
		MasterSegmentId: params.MasterSegmentId,
		Condition:       params.Condition,
		Description:     params.Description,
		Name:            params.Name,
		AccountUuid:     params.AccountUuid,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

type ListMasterSegmentsParams struct {
	AccountUuid uuid.UUID
}

func (r *segmentRepo) ListMasterSegments(ctx context.Context, params *ListMasterSegmentsParams) ([]model.MasterSegment, error) {
	var masterSegments []model.MasterSegment
	err := r.WithContext(ctx).Table(r.MasterSegmentTableName).
		Where("account_uuid = ?", params.AccountUuid).
		Find(&masterSegments).Error
	if err != nil {
		return nil, err
	}

	return masterSegments, nil
}
