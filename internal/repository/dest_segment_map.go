package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type DestSegmentMapRepository interface {
	CreateDestinationSegmentMap(ctx context.Context, params *CreateDestinationSegmentMapParams) (*model.DestSegmentMap, error)
	ListDestinationSegmentMaps(ctx context.Context, params *ListDestinationSegmentMapsParams) ([]DestinationSegmentMapItem, error)
}

type destSegmentMapRepo struct {
	*gorm.DB
	TableName string
}

func NewDestSegmentMapRepository(db *gorm.DB) DestSegmentMapRepository {
	return &destSegmentMapRepo{db, model.DestSegmentMap{}.TableName()}
}

type CreateDestinationSegmentMapParams struct {
	Tx             *gorm.DB
	SegmentId      int64
	DestinationId  int64
	MappingOptions pqtype.NullRawMessage
}

func (r *destSegmentMapRepo) CreateDestinationSegmentMap(ctx context.Context, params *CreateDestinationSegmentMapParams) (*model.DestSegmentMap, error) {
	destSegmentMap := &model.DestSegmentMap{
		SegmentId:      params.SegmentId,
		DestinationId:  params.DestinationId,
		MappingOptions: params.MappingOptions,
	}

	var createErr error
	if params.Tx != nil {
		createErr = params.Tx.WithContext(ctx).Table(r.TableName).Create(destSegmentMap).Error
	} else {
		createErr = r.WithContext(ctx).Table(r.TableName).Create(destSegmentMap).Error
	}
	if createErr != nil {
		return nil, createErr
	}

	return destSegmentMap, nil
}

type ListDestinationSegmentMapsParams struct {
	DestinationId int64
}

type DestinationSegmentMapItem struct {
	ID             int64 `gorm:"primaryKey"`
	SegmentId      int64
	SegmentName    string
	MappingOptions pqtype.NullRawMessage
	DataActionId   int64
}

func (r *destSegmentMapRepo) ListDestinationSegmentMaps(ctx context.Context, params *ListDestinationSegmentMapsParams) ([]DestinationSegmentMapItem, error) {
	var mappings []DestinationSegmentMapItem
	query := r.WithContext(ctx).Table(r.TableName).Where("destination_id = ?", params.DestinationId).
		Joins("LEFT JOIN segment ON dest_segment_map.segment_id = segment.id " +
			"LEFT JOIN data_action ON data_action.target_table = 'dest_segment_map' AND data_action.object_id = dest_segment_map.id").
		Select("dest_segment_map.id AS id, " +
			"dest_segment_map.segment_id AS segment_id, " +
			"dest_segment_map.mapping_options AS mapping_options, " +
			"segment.name AS segment_id, " +
			"data_action.id AS data_action_id")

	err := query.Find(&mappings).Error
	if err != nil {
		return nil, err
	}

	return mappings, err
}
