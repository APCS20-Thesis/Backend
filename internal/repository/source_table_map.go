package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
	"time"
)

type SourceTableMapRepository interface {
	CreateSourceTableMap(ctx context.Context, params *CreateSourceTableMapParams) (*model.SourceTableMap, error)
	GetSourceTableMapById(ctx context.Context, id int64) (*model.SourceTableMap, error)
	ListSourceTableMap(ctx context.Context, params *ListSourceTableMapParams) (*ListSourceTableMapResult, error)
}

type sourceTableMapRepo struct {
	*gorm.DB
	TableName string
}

func NewSourceTableMapRepository(db *gorm.DB) SourceTableMapRepository {
	return &sourceTableMapRepo{db, model.SourceTableMap{}.TableName()}
}

type CreateSourceTableMapParams struct {
	Tx             *gorm.DB
	TableId        int64
	SourceId       int64
	MappingOptions pqtype.NullRawMessage
}

func (r *sourceTableMapRepo) CreateSourceTableMap(ctx context.Context, params *CreateSourceTableMapParams) (*model.SourceTableMap, error) {
	sourceTableMap := &model.SourceTableMap{
		TableId:        params.TableId,
		SourceId:       params.SourceId,
		MappingOptions: params.MappingOptions,
	}

	var createErr error
	if params.Tx != nil {
		createErr = params.Tx.WithContext(ctx).Table(r.TableName).Create(sourceTableMap).Error
	} else {
		createErr = r.WithContext(ctx).Table(r.TableName).Create(sourceTableMap).Error
	}
	if createErr != nil {
		return nil, createErr
	}

	return sourceTableMap, nil
}

func (r *sourceTableMapRepo) GetSourceTableMapById(ctx context.Context, id int64) (*model.SourceTableMap, error) {
	var sourceTableMap model.SourceTableMap

	err := r.WithContext(ctx).Table(r.TableName).First(&sourceTableMap, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &sourceTableMap, nil
}

type (
	ListSourceTableMapParams struct {
		TableId  int64
		SourceId int64
		Ids      []int64
	}
	TableSourceMapWithExtraInfo struct {
		ID             int64 `gorm:"primaryKey"`
		TableId        int64
		TableName      string
		SourceId       int64
		SourceName     string
		SourceType     model.DataSourceType
		MappingOptions pqtype.NullRawMessage
		CreatedAt      time.Time `gorm:"autoCreateTime"`
		UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	}
	ListSourceTableMapResult struct {
		TableSourceMaps []TableSourceMapWithExtraInfo
		Count           int64
	}
)

func (r *sourceTableMapRepo) ListSourceTableMap(ctx context.Context, params *ListSourceTableMapParams) (*ListSourceTableMapResult, error) {
	query := r.WithContext(ctx).Table(r.TableName)
	if params.SourceId > 0 {
		query.Where("source_table_map.source_id = ?", params.SourceId)
	}
	if params.TableId > 0 {
		query.Where("source_table_map.table_id = ?", params.TableId)
	}
	if len(params.Ids) > 0 {
		query.Where("source_table_map.id IN ?", params.Ids)
	}
	query = query.
		Joins("LEFT JOIN data_table ON source_table_map.table_id = data_table.id").
		Joins("LEFT JOIN data_source ON source_table_map.source_id = data_source.id").
		Select(
			"source_table_map.id AS id, " +
				"source_table_map.table_id AS table_id, " +
				"source_table_map.source_id AS source_id, " +
				"source_table_map.mapping_options AS mapping_options, " +
				"source_table_map.created_at AS created_at, " +
				"source_table_map.updated_at AS updated_at, " +
				"data_table.name AS table_name, " +
				"data_source.name AS source_name, " +
				"data_source.type AS source_type ",
		)

	var (
		maps  []TableSourceMapWithExtraInfo
		count int64
	)
	err := query.Order("source_table_map.id desc").Count(&count).Scan(&maps).Error
	if err != nil {
		return nil, err
	}

	return &ListSourceTableMapResult{
		TableSourceMaps: maps,
		Count:           count,
	}, nil
}
