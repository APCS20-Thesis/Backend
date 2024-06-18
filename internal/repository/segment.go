package repository

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
	"time"
)

type SegmentRepository interface {
	CreateMasterSegment(ctx context.Context, params *model.MasterSegment) error
	ListMasterSegments(ctx context.Context, params *ListMasterSegmentsParams) ([]model.MasterSegment, error)
	GetMasterSegment(ctx context.Context, masterSegmentId int64, accountUuid string) (model.MasterSegment, error)

	CreateAudienceTable(ctx context.Context, params *CreateAudienceTableParams) error
	GetAudienceTable(ctx context.Context, params GetAudienceTableParams) (model.AudienceTable, error)

	CreateBehaviorTable(ctx context.Context, params *CreateBehaviorTableParams) error
	ListBehaviorTables(ctx context.Context, params ListBehaviorTablesParams) ([]model.BehaviorTable, error)

	CreateSegment(ctx context.Context, params *CreateSegmentParams) error
	ListSegments(ctx context.Context, filter *ListSegmentFilter) ([]SegmentListItem, error)
	GetSegment(ctx context.Context, segmentId int64, accountUuid string) (model.Segment, error)
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
	SqlCondition    string
	AccountUuid     uuid.UUID
}

func (r *segmentRepo) CreateSegment(ctx context.Context, params *CreateSegmentParams) error {
	err := r.WithContext(ctx).Table(r.SegmentTableName).Create(&model.Segment{
		MasterSegmentId: params.MasterSegmentId,
		Condition:       params.Condition,
		SqlCondition:    params.SqlCondition,
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

type ListBehaviorTablesParams struct {
	MasterSegmentId int64
}

func (r *segmentRepo) ListBehaviorTables(ctx context.Context, params ListBehaviorTablesParams) ([]model.BehaviorTable, error) {
	var behaviorTables []model.BehaviorTable
	err := r.WithContext(ctx).Table(r.BehaviorTableName).
		Where("master_segment_id = ?", params.MasterSegmentId).
		Find(&behaviorTables).Error
	if err != nil {
		return nil, err
	}

	return behaviorTables, nil
}

type GetAudienceTableParams struct {
	MasterSegmentId int64
}

func (r *segmentRepo) GetAudienceTable(ctx context.Context, params GetAudienceTableParams) (model.AudienceTable, error) {
	var audienceTable model.AudienceTable
	err := r.WithContext(ctx).Table(r.AudienceTableName).
		Where("master_segment_id = ?", params.MasterSegmentId).
		Find(&audienceTable).Error
	if err != nil {
		return model.AudienceTable{}, err
	}

	return audienceTable, nil
}

func (r *segmentRepo) GetMasterSegment(ctx context.Context, masterSegmentId int64, accountUuid string) (model.MasterSegment, error) {
	var masterSegment model.MasterSegment
	err := r.WithContext(ctx).Table(r.MasterSegmentTableName).
		Where("account_uuid = ? AND id = ?", accountUuid, masterSegmentId).
		First(&masterSegment).Error
	if err != nil {
		return model.MasterSegment{}, err
	}

	return masterSegment, nil
}

type SegmentListItem struct {
	ID                int64
	Name              string
	MasterSegmentId   int64
	MasterSegmentName string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ListSegmentFilter struct {
	AccountUuid      string
	MasterSegmentIds []int64
}

func (r *segmentRepo) ListSegments(ctx context.Context, filter *ListSegmentFilter) ([]SegmentListItem, error) {
	var segments []SegmentListItem
	query := r.WithContext(ctx).Table(r.SegmentTableName)
	if len(filter.MasterSegmentIds) > 0 {
		query = query.Where("segment.master_segment_id IN ?", filter.MasterSegmentIds)
	}
	if filter.AccountUuid != "" {
		query = query.Where("segment.account_uuid = ?", filter.AccountUuid)
	}
	err := query.
		Joins("LEFT JOIN master_segment ON segment.master_segment_id = master_segment.id").
		Select("segment.id AS id," +
			"segment.name AS name, " +
			"segment.master_segment_id as master_segment_id, " +
			"master_segment.name as master_segment_name, " +
			"segment.created_at AS created_at, " +
			"segment.updated_at AS updated_at").
		Scan(&segments).Error
	if err != nil {
		return nil, err
	}

	return segments, nil
}

func (r *segmentRepo) GetSegment(ctx context.Context, segmentId int64, accountUuid string) (model.Segment, error) {
	var segment model.Segment
	err := r.WithContext(ctx).Table(r.SegmentTableName).
		Where("account_uuid = ? AND id = ?", accountUuid, segmentId).
		First(&segment).Error
	if err != nil {
		return model.Segment{}, err
	}
	return segment, nil
}
