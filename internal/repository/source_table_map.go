package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"gorm.io/gorm"
)

type SourceTableMapRepository interface {
	CreateSourceTableMap(ctx context.Context, params *CreateSourceTableMapParams) error
	GetSourceTableMapById(ctx context.Context, id int64) (*model.SourceTableMap, error)
}

type sourceTableMapRepo struct {
	*gorm.DB
	TableName string
}

func NewSourceTableMapRepository(db *gorm.DB) SourceTableMapRepository {
	return &sourceTableMapRepo{db, model.SourceTableMap{}.TableName()}
}

type CreateSourceTableMapParams struct {
	TableId  int64
	SourceId int64
}

func (r *sourceTableMapRepo) CreateSourceTableMap(ctx context.Context, params *CreateSourceTableMapParams) error {
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

func (r *sourceTableMapRepo) GetSourceTableMapById(ctx context.Context, id int64) (*model.SourceTableMap, error) {
	var sourceTableMap model.SourceTableMap

	err := r.WithContext(ctx).Table(r.TableName).First(&sourceTableMap, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &sourceTableMap, nil
}
