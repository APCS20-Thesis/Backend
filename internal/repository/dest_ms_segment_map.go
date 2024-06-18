package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type DestMasterSegmentMapRepository interface {
	CreateDestinationMasterSegmentMap(ctx context.Context, params *CreateDestinationMasterSegmentMapParams) (*model.DestMasterSegmentMap, error)
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
