package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"

	"github.com/APCS20-Thesis/Backend/internal/model"
)

type DataActionRepository interface {
	CreateDataAction(ctx context.Context, params *CreateDataAction) error
	GetDataAction(ctx context.Context, id int64) (*model.DataAction, error)
	UpdateDataAction(ctx context.Context, params *UpdateDataActionParams) error
	DeleteDataAction(ctx context.Context, id int64) error
}

type dataActionRepo struct {
	*gorm.DB
	TableName string
}

func NewDataActionRepository(db *gorm.DB) DataActionRepository {
	return &dataActionRepo{db, model.DataAction{}.TableName()}
}

type CreateDataAction struct {
	ActionType  model.ActionType
	Payload     pqtype.NullRawMessage
	Schedule    string
	AccountUuid uuid.UUID
}

func (r *dataActionRepo) CreateDataAction(ctx context.Context, params *CreateDataAction) error {
	dataAction := &model.DataAction{
		ActionType:  params.ActionType,
		Payload:     params.Payload,
		Schedule:    params.Schedule,
		AccountUuid: params.AccountUuid,
	}

	createErr := r.WithContext(ctx).Table(r.TableName).Create(&dataAction).Error
	if createErr != nil {
		return createErr
	}

	return nil
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

func (r *dataActionRepo) DeleteDataAction(ctx context.Context, id int64) error {
	deleteErr := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).Delete(&model.DataAction{}).Error
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
