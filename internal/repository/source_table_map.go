package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"gorm.io/gorm"
)

type SourceTableMapRepository interface {
	CreateSourceTableMap(ctx context.Context, params *CreateSourceTableMapParams) (*model.SourceTableMap, error)
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
	Tx       *gorm.DB
	TableId  int64
	SourceId int64
}

func (r *sourceTableMapRepo) CreateSourceTableMap(ctx context.Context, params *CreateSourceTableMapParams) (*model.SourceTableMap, error) {
	sourceTableMap := &model.SourceTableMap{
		TableId:  params.TableId,
		SourceId: params.SourceId,
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
