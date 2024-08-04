package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
	"time"
)

type SegmentRepository interface {
	CreateMasterSegment(ctx context.Context, params *model.MasterSegment) error
	ListMasterSegments(ctx context.Context, filter *ListMasterSegmentsFilter) (*ListMasterSegmentsResult, error)
	GetMasterSegment(ctx context.Context, masterSegmentId int64) (model.MasterSegment, error)
	UpdateMasterSegment(ctx context.Context, params *UpdateMasterSegmentParams) error

	CreateAudienceTable(ctx context.Context, params *CreateAudienceTableParams) error
	GetAudienceTable(ctx context.Context, params GetAudienceTableParams) (model.AudienceTable, error)
	UpdateAudienceTable(ctx context.Context, params *UpdateAudienceTableParams) error

	CreateBehaviorTable(ctx context.Context, params *CreateBehaviorTableParams) error
	ListBehaviorTables(ctx context.Context, params ListBehaviorTablesParams) ([]model.BehaviorTable, error)
	UpdateBehaviorTable(ctx context.Context, params *UpdateBehaviorTableParams) error

	CreateSegment(ctx context.Context, params *CreateSegmentParams) (*model.Segment, error)
	ListSegments(ctx context.Context, filter *ListSegmentFilter) ([]SegmentListItem, error)
	GetSegment(ctx context.Context, segmentId int64, accountUuid string) (model.Segment, error)
	UpdateSegment(ctx context.Context, params *UpdateSegmentParams) error
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
	Tx              *gorm.DB
	Name            string
	Description     string
	MasterSegmentId int64
	Condition       pqtype.NullRawMessage
	SqlCondition    string
	AccountUuid     uuid.UUID
}

func (r *segmentRepo) CreateSegment(ctx context.Context, params *CreateSegmentParams) (*model.Segment, error) {
	segment := &model.Segment{
		MasterSegmentId: params.MasterSegmentId,
		Condition:       params.Condition,
		SqlCondition:    params.SqlCondition,
		Description:     params.Description,
		Name:            params.Name,
		AccountUuid:     params.AccountUuid,
		Status:          model.SegmentStatus_DRAFT,
	}

	var createErr error
	if params.Tx != nil {
		createErr = params.Tx.WithContext(ctx).Table(r.SegmentTableName).Create(segment).Error
	} else {
		createErr = r.WithContext(ctx).Table(r.SegmentTableName).Create(segment).Error
	}
	if createErr != nil {
		return nil, createErr
	}

	return segment, nil
}

type ListMasterSegmentsFilter struct {
	AccountUuid uuid.UUID
	Statuses    []string
	Name        string
	Page        int
	PageSize    int
}

type ListMasterSegmentsResult struct {
	MasterSegments []model.MasterSegment
	Count          int64
}

func (r *segmentRepo) ListMasterSegments(ctx context.Context, filter *ListMasterSegmentsFilter) (*ListMasterSegmentsResult, error) {
	var (
		masterSegments []model.MasterSegment
		count          int64
	)
	query := r.WithContext(ctx).Table(r.MasterSegmentTableName)
	if filter.AccountUuid.String() != "" {
		query = query.Where("account_uuid = ?", filter.AccountUuid)
	}
	if len(filter.Statuses) > 0 {
		query = query.Where("status IN ?", filter.Statuses)
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	err := query.Count(&count).Scopes(Paginate(filter.Page, filter.PageSize)).Find(&masterSegments).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &ListMasterSegmentsResult{
		MasterSegments: masterSegments,
		Count:          count,
	}, nil
}

type UpdateMasterSegmentParams struct {
	Tx     *gorm.DB
	Id     int64
	Status model.MasterSegmentStatus
}

func (r *segmentRepo) UpdateMasterSegment(ctx context.Context, params *UpdateMasterSegmentParams) error {
	masterSegment := model.MasterSegment{
		ID:     params.Id,
		Status: params.Status,
	}
	var updateErr error
	if params.Tx != nil {
		updateErr = params.Tx.WithContext(ctx).Table(r.MasterSegmentTableName).Where("id = ?", params.Id).Updates(&masterSegment).Error
	} else {
		updateErr = r.WithContext(ctx).Table(r.MasterSegmentTableName).Where("id = ?", params.Id).Updates(&masterSegment).Error
	}
	if updateErr != nil {
		return updateErr
	}
	return nil
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

type UpdateBehaviorTableParams struct {
	Tx     *gorm.DB
	Id     int64
	Schema pqtype.NullRawMessage
}

func (r *segmentRepo) UpdateBehaviorTable(ctx context.Context, params *UpdateBehaviorTableParams) error {
	behaviorTable := model.BehaviorTable{
		ID:     params.Id,
		Schema: params.Schema,
	}
	var updateErr error
	if params.Tx != nil {
		updateErr = params.Tx.WithContext(ctx).Table(r.BehaviorTableName).Where("id = ?", params.Id).Updates(&behaviorTable).Error
	} else {
		updateErr = r.WithContext(ctx).Table(r.BehaviorTableName).Where("id = ?", params.Id).Updates(&behaviorTable).Error
	}
	if updateErr != nil {
		return updateErr
	}
	return nil
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

type UpdateAudienceTableParams struct {
	Tx     *gorm.DB
	Id     int64
	Schema pqtype.NullRawMessage
}

func (r *segmentRepo) UpdateAudienceTable(ctx context.Context, params *UpdateAudienceTableParams) error {
	audienceTable := model.AudienceTable{
		ID:     params.Id,
		Schema: params.Schema,
	}
	var updateErr error
	if params.Tx != nil {
		updateErr = params.Tx.WithContext(ctx).Table(r.AudienceTableName).Where("id = ?", params.Id).Updates(&audienceTable).Error
	} else {
		updateErr = r.WithContext(ctx).Table(r.AudienceTableName).Where("id = ?", params.Id).Updates(&audienceTable).Error
	}
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func (r *segmentRepo) GetMasterSegment(ctx context.Context, masterSegmentId int64) (model.MasterSegment, error) {
	var masterSegment model.MasterSegment
	err := r.WithContext(ctx).Table(r.MasterSegmentTableName).
		Where("id = ?", masterSegmentId).
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
	Status            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ListSegmentFilter struct {
	AccountUuid      string
	MasterSegmentIds []int64
	Statuses         []string
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
	if len(filter.Statuses) > 0 {
		query = query.Where("segment.status IN ?", filter.Statuses)
	}
	err := query.
		Joins("LEFT JOIN master_segment ON segment.master_segment_id = master_segment.id").
		Select("segment.id AS id," +
			"segment.name AS name, " +
			"segment.master_segment_id as master_segment_id, " +
			"master_segment.name as master_segment_name, " +
			"segment.status AS status, " +
			"segment.created_at AS created_at, " +
			"segment.updated_at AS updated_at").
		Order("segment.updated_at desc").Scan(&segments).Error
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

type UpdateSegmentParams struct {
	Tx     *gorm.DB
	Id     int64
	Status string
}

func (r *segmentRepo) UpdateSegment(ctx context.Context, params *UpdateSegmentParams) error {
	segment := model.Segment{
		ID:     params.Id,
		Status: model.SegmentStatus(params.Status),
	}
	var updateErr error
	if params.Tx != nil {
		updateErr = params.Tx.WithContext(ctx).Table(r.SegmentTableName).Where("id = ?", params.Id).Updates(&segment).Error
	} else {
		updateErr = r.WithContext(ctx).Table(r.SegmentTableName).Where("id = ?", params.Id).Updates(&segment).Error
	}
	if updateErr != nil {
		return updateErr
	}
	return nil
}
