package repository

import (
	"context"
	"errors"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type DataActionRepository interface {
	CreateDataAction(ctx context.Context, params *CreateDataActionParams) (*model.DataAction, error)
	GetDataAction(ctx context.Context, id int64) (*model.DataAction, error)
	UpdateDataAction(ctx context.Context, params *UpdateDataActionParams) error
	GetListDataActions(ctx context.Context, params *GetListDataActionsParams) (*GetListDataActionsResult, error)
}

type dataActionRepo struct {
	*gorm.DB
	TableName string
}

func NewDataActionRepository(db *gorm.DB) DataActionRepository {
	return &dataActionRepo{db, model.DataAction{}.TableName()}
}

type CreateDataActionParams struct {
	TargetTable model.DataActionTargetTable
	ActionType  model.ActionType
	Schedule    string
	AccountUuid uuid.UUID
	DagId       string
	Status      model.DataActionStatus
	ObjectId    int64
	Payload     pqtype.NullRawMessage
	RunCount    int64
}

func (r *dataActionRepo) CreateDataAction(ctx context.Context, params *CreateDataActionParams) (*model.DataAction, error) {
	dataAction := &model.DataAction{
		ActionType:  params.ActionType,
		Payload:     params.Payload,
		Status:      params.Status,
		Schedule:    params.Schedule,
		DagId:       params.DagId,
		TargetTable: params.TargetTable,
		ObjectId:    params.ObjectId,
		AccountUuid: params.AccountUuid,
		RunCount:    params.RunCount,
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &dataAction, nil
}

type UpdateDataActionParams struct {
	Tx          *gorm.DB
	ID          int64
	ActionType  model.ActionType
	Payload     pqtype.NullRawMessage
	Schedule    string
	AccountUuid uuid.UUID
	Status      model.DataActionStatus
	RunCount    int64
}

func (r *dataActionRepo) UpdateDataAction(ctx context.Context, params *UpdateDataActionParams) error {
	dataAction := &model.DataAction{
		ID:          params.ID,
		ActionType:  params.ActionType,
		Payload:     params.Payload,
		Schedule:    params.Schedule,
		AccountUuid: params.AccountUuid,
		Status:      params.Status,
		RunCount:    params.RunCount,
	}

	var updateErr error
	if params.Tx != nil {
		updateErr = params.Tx.WithContext(ctx).Table(r.TableName).Where("id = ?", params.ID).Updates(&dataAction).Error
	} else {
		updateErr = r.WithContext(ctx).Table(r.TableName).Where("id = ?", params.ID).Updates(&dataAction).Error
	}
	if updateErr != nil {
		return updateErr
	}

	return nil
}

type GetListDataActionsParams struct {
	Ids         []int64
	ActionTypes []string
	Statuses    []model.DataActionStatus
	AccountUuid uuid.UUID
	DagId       string
	Page        int
	PageSize    int
	TargetTable model.DataActionTargetTable
	ObjectId    int64
}

type GetListDataActionsResult struct {
	DataActions []model.DataAction
	Count       int64
}

func (r *dataActionRepo) GetListDataActions(ctx context.Context, params *GetListDataActionsParams) (*GetListDataActionsResult, error) {
	var count int64

	query := r.WithContext(ctx).Table(r.TableName)
	if params.DagId != "" {
		query = query.Where("dag_id = ?", params.DagId)
	}
	if len(params.Ids) > 0 {
		query.Where("id IN ?", params.Ids)
	}
	if len(params.ActionTypes) > 0 {
		query.Where("action_type IN ?", params.ActionTypes)
	}
	if len(params.Statuses) > 0 {
		query.Where("status IN ?", params.Statuses)
	}
	if params.TargetTable != "" {
		query.Where("target_table = ?", params.TargetTable)
	}
	if params.ObjectId != 0 {
		query.Where("object_id = ?", params.ObjectId)
	}

	emptyUuid, _ := uuid.Parse("")
	if params.AccountUuid != emptyUuid {
		query = query.Where("account_uuid = ?", params.AccountUuid)
	}

	var dataActions []model.DataAction
	err := query.Order("updated_at desc").Count(&count).Scopes(Paginate(params.Page, params.PageSize)).Find(&dataActions).Error
	if err != nil {
		return nil, err
	}

	return &GetListDataActionsResult{
		DataActions: dataActions,
		Count:       count,
	}, nil
}
