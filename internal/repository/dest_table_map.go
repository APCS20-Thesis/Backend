package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type DestTableMapRepository interface {
	CreateDestinationTableMap(ctx context.Context, params *CreateDestinationTableMapParams) (*model.DestTableMap, error)
	ListDestinationTableMaps(ctx context.Context, params *ListDestinationTableMapsParams) ([]DestinationTableMapItem, error)
}

type destTableMapRepo struct {
	*gorm.DB
	TableName string
}

func NewDestTableMapRepository(db *gorm.DB) DestTableMapRepository {
	return &destTableMapRepo{db, model.DestTableMap{}.TableName()}
}

type CreateDestinationTableMapParams struct {
	Tx             *gorm.DB
	TableId        int64
	DestinationId  int64
	MappingOptions pqtype.NullRawMessage
}

func (r *destTableMapRepo) CreateDestinationTableMap(ctx context.Context, params *CreateDestinationTableMapParams) (*model.DestTableMap, error) {
	destTableMap := &model.DestTableMap{
		TableId:        params.TableId,
		DestinationId:  params.DestinationId,
		MappingOptions: params.MappingOptions,
	}

	var createErr error
	if params.Tx != nil {
		createErr = params.Tx.WithContext(ctx).Table(r.TableName).Create(destTableMap).Error
	} else {
		createErr = r.WithContext(ctx).Table(r.TableName).Create(destTableMap).Error
	}
	if createErr != nil {
		return nil, createErr
	}

	return destTableMap, nil
}

type ListDestinationTableMapsParams struct {
	DestinationId int64
}

type DestinationTableMapItem struct {
	ID             int64 `gorm:"primaryKey"`
	TableId        int64
	TableName      string
	MappingOptions pqtype.NullRawMessage
	DataActionId   int64
}

func (r *destTableMapRepo) ListDestinationTableMaps(ctx context.Context, params *ListDestinationTableMapsParams) ([]DestinationTableMapItem, error) {
	var mappings []DestinationTableMapItem
	query := r.WithContext(ctx).Table(r.TableName).Where("destination_id = ?", params.DestinationId).
		Joins("LEFT JOIN data_table ON dest_table_map.table_id = data_table.id " +
			"LEFT JOIN data_action ON data_action.target_table = 'dest_table_map' AND data_action.object_id = dest_table_map.id").
		Select("dest_table_map.id AS id, " +
			"dest_table_map.table_id AS table_id, " +
			"dest_table_map.mapping_options AS mapping_options, " +
			"data_table.name AS table_name, " +
			"data_action.id AS data_action_id")

	err := query.Find(&mappings).Error
	if err != nil {
		return nil, err
	}

	return mappings, err
}
