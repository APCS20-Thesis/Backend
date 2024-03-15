package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DataActionRunRepository interface {
	CreateDataActionRun(ctx context.Context, params *CreateDataActionRunParams) (*model.DataActionRun, error)
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
	DagRunId    string
	Status      model.DataActionRunStatus
	AccountUuid uuid.UUID
}

func (r *dataActionRunRepo) CreateDataActionRun(ctx context.Context, params *CreateDataActionRunParams) (*model.DataActionRun, error) {
	dataActionRun := &model.DataActionRun{
		ActionId:    params.ActionId,
		RunId:       params.RunId,
		Status:      params.Status,
		AccountUuid: params.AccountUuid,
		DagRunId:    params.DagRunId,
	}
	if err := r.WithContext(ctx).Table(r.TableName).Create(&dataActionRun).Error; err != nil {
		return nil, err
	}
	return dataActionRun, nil
}

func (r *dataActionRunRepo) UpdateDataActionRunStatus(ctx context.Context, id int64, status model.DataActionRunStatus) error {
	err := r.WithContext(ctx).Table(r.TableName).Save(&model.DataActionRun{
		ID:     id,
		Status: status,
	}).Error

	return err
}
