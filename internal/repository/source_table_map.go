package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"gorm.io/gorm"
)

type SourceTableMapRepository interface {
	CreateSourceTableMap(ctx context.Context, params *CreateSourceTableMapParams) error
}

type SourceTableMapRepo struct {
	*gorm.DB
	TableName string
}

func NewSourceTableMapRepository(db *gorm.DB) SourceTableMapRepository {
	return &SourceTableMapRepo{db, model.SourceTableMap{}.TableName()}
}

type CreateSourceTableMapParams struct {
	TableId  int64
	SourceId int64
}

func (r *SourceTableMapRepo) CreateSourceTableMap(ctx context.Context, params *CreateSourceTableMapParams) error {
	SourceTableMap := &model.SourceTableMap{
		TableId:  params.TableId,
		SourceId: params.SourceId,
	}

	createErr := r.WithContext(ctx).Table(r.TableName).Create(&SourceTableMap).Error
	if createErr != nil {
		return createErr
	}

	return nil
}
