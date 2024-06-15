package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type DestSegmentMapRepository interface {
	CreateDestinationSegmentMap(ctx context.Context, params *CreateDestinationSegmentMapParams) (*model.DestSegmentMap, error)
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
