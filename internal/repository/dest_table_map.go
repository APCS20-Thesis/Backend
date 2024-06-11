package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"gorm.io/gorm"
)

type DestTableMapRepository interface {
	CreateDestinationTableMap(ctx context.Context, params *CreateDestinationTableMapParams) (*model.DestTableMap, error)
}

type destTableMapRepo struct {
	*gorm.DB
	TableName string
}

func NewDestTableMapRepository(db *gorm.DB) DestTableMapRepository {
	return &destTableMapRepo{db, model.DestTableMap{}.TableName()}
}

type CreateDestinationTableMapParams struct {
	Tx            *gorm.DB
	TableId       int64
	DestinationId int64
}

func (r *destTableMapRepo) CreateDestinationTableMap(ctx context.Context, params *CreateDestinationTableMapParams) (*model.DestTableMap, error) {
	destTableMap := &model.DestTableMap{
		TableId:       params.TableId,
		DestinationId: params.DestinationId,
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
