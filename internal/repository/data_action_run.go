package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DataActionRunRepository interface {
	CreateDataActionRun(params CreateDataActionRunParams) error
	UpdateDataActionRunStatus(ctx context.Context, id int64, status model.DataActionRunStatus) error
}

type dataActionRunRepo struct {
	*gorm.DB
	TableName string
}

func NewDataActionRunRepository(db *gorm.DB) DataActionRunRepository {
	return &dataActionRunRepo{db, model.DataActionRun{}.TableName()}
}

type CreateDataActionRunParams struct {
	ActionId    int64
	RunId       int64
	Status      model.DataActionRunStatus
	AccountUuid uuid.UUID
}

func (r *dataActionRunRepo) CreateDataActionRun(params CreateDataActionRunParams) error {
	if err := r.DB.Create(&model.DataActionRun{
		ActionId:    params.ActionId,
		RunId:       params.RunId,
		Status:      params.Status,
		AccountUuid: params.AccountUuid,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *dataActionRunRepo) UpdateDataActionRunStatus(ctx context.Context, id int64, status model.DataActionRunStatus) error {
	err := r.WithContext(ctx).Table(r.TableName).Save(&model.DataActionRun{
		ID:     id,
		Status: status,
	}).Error

	return err
}
