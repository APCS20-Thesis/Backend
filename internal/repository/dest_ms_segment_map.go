package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type DestMasterSegmentMapRepository interface {
	CreateDestinationMasterSegmentMap(ctx context.Context, params *CreateDestinationMasterSegmentMapParams) (*model.DestMasterSegmentMap, error)
	ListDestinationMasterSegmentMaps(ctx context.Context, params *ListDestinationMasterSegmentMapsParams) ([]DestinationMasterSegmentMapItem, error)
}

type destMasterSegmentMapRepo struct {
	*gorm.DB
	TableName string
}

func NewDestMasterSegmentMapRepository(db *gorm.DB) DestMasterSegmentMapRepository {
	return &destMasterSegmentMapRepo{db, model.DestMasterSegmentMap{}.TableName()}
}

type CreateDestinationMasterSegmentMapParams struct {
	Tx              *gorm.DB
	MasterSegmentId int64
	DestinationId   int64
	MappingOptions  pqtype.NullRawMessage
}

func (r *destMasterSegmentMapRepo) CreateDestinationMasterSegmentMap(ctx context.Context, params *CreateDestinationMasterSegmentMapParams) (*model.DestMasterSegmentMap, error) {
	destMasterSegmentMap := &model.DestMasterSegmentMap{
		MasterSegmentId: params.MasterSegmentId,
		DestinationId:   params.DestinationId,
		MappingOptions:  params.MappingOptions,
	}

	var createErr error
	if params.Tx != nil {
		createErr = params.Tx.WithContext(ctx).Table(r.TableName).Create(destMasterSegmentMap).Error
	} else {
		createErr = r.WithContext(ctx).Table(r.TableName).Create(destMasterSegmentMap).Error
	}
	if createErr != nil {
		return nil, createErr
	}

	return destMasterSegmentMap, nil
}

type ListDestinationMasterSegmentMapsParams struct {
	DestinationId int64
	Ids           []int64
}

type DestinationMasterSegmentMapItem struct {
	ID                int64 `gorm:"primaryKey"`
	DestinationId     int64
	DestinationName   string
	MasterSegmentId   int64
	MasterSegmentName string
	MappingOptions    pqtype.NullRawMessage
	DataActionId      int64
}

func (r *destMasterSegmentMapRepo) ListDestinationMasterSegmentMaps(ctx context.Context, params *ListDestinationMasterSegmentMapsParams) ([]DestinationMasterSegmentMapItem, error) {
	var mappings []DestinationMasterSegmentMapItem
	query := r.WithContext(ctx).Table(r.TableName).
		Joins("LEFT JOIN master_segment ON dest_ms_segment_map.master_segment_id = master_segment.id " +
			"LEFT JOIN data_destination ON dest_ms_segment_map.destination_id = data_destination.id " +
			"LEFT JOIN data_action ON data_action.target_table = 'dest_ms_segment_map' AND data_action.object_id = dest_ms_segment_map.id").
		Select("dest_ms_segment_map.id AS id, " +
			"dest_ms_segment_map.destination_id AS destination_id, " +
			"data_destination.name AS destination_name, " +
			"dest_ms_segment_map.master_segment_id AS master_segment_id, " +
			"dest_ms_segment_map.mapping_options AS mapping_options, " +
			"master_segment.name AS master_segment_name, " +
			"data_action.id AS data_action_id")
	if params.DestinationId > 0 {
		query = query.Where("dest_ms_segment_map.destination_id = ?", params.DestinationId)
	}
	if len(params.Ids) > 0 {
		query = query.Where("dest_ms_segment_map.id IN ?", params.Ids)
	}

	err := query.Find(&mappings).Error
	if err != nil {
		return nil, err
	}

	return mappings, err
}
