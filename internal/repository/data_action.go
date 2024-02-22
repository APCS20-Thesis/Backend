package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"

	"github.com/APCS20-Thesis/Backend/internal/model"
)

type DataActionRepository interface {
	CreateDataAction(ctx context.Context, params *CreateDataActionParams) (*model.DataAction, error)
	GetDataAction(ctx context.Context, id int64) (*model.DataAction, error)
	UpdateDataAction(ctx context.Context, params *UpdateDataActionParams) error
}

type dataActionRepo struct {
	*gorm.DB
	TableName string
}

func NewDataActionRepository(db *gorm.DB) DataActionRepository {
	return &dataActionRepo{db, model.DataAction{}.TableName()}
}

type CreateDataActionParams struct {
	ActionType  model.ActionType
	Schedule    string
	AccountUuid uuid.UUID
	DagId       string
	Status      string
}

func (r *dataActionRepo) CreateDataAction(ctx context.Context, params *CreateDataActionParams) (*model.DataAction, error) {
	dataAction := &model.DataAction{
		ActionType:  params.ActionType,
		Schedule:    params.Schedule,
		AccountUuid: params.AccountUuid,
		DagId:       params.DagId,
		RunCount:    0,
		Status:      params.Status,
	}

	createErr := r.WithContext(ctx).Table(r.TableName).Create(&dataAction).Error
	if createErr != nil {
		return nil, createErr
	}

	return dataAction, nil
}

func (r *dataActionRepo) GetDataAction(ctx context.Context, id int64) (*model.DataAction, error) {
	var dataAction model.DataAction
	err := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).First(&dataAction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &dataAction, nil
}

type UpdateDataActionParams struct {
	ID          int64
	ActionType  model.ActionType
	Payload     pqtype.NullRawMessage
	Schedule    string
	AccountUuid uuid.UUID
}

func (r *dataActionRepo) UpdateDataAction(ctx context.Context, params *UpdateDataActionParams) error {
	dataAction := &model.DataAction{
		ID:          params.ID,
		ActionType:  params.ActionType,
		Payload:     params.Payload,
		Schedule:    params.Schedule,
		AccountUuid: params.AccountUuid,
	}

	updateErr := r.WithContext(ctx).Table(r.TableName).Where("id = ?", params.ID).Updates(&dataAction).Error
	if updateErr != nil {
		return updateErr
	}

	return nil
}
