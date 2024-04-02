package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type DataActionRunRepository interface {
	CreateDataActionRun(ctx context.Context, params *CreateDataActionRunParams) (*model.DataActionRun, error)
	UpdateDataActionRunStatus(ctx context.Context, id int64, status model.DataActionRunStatus) error
	GetListDataActionRuns(ctx context.Context, params *GetListDataActionRunsParams) ([]model.DataActionRun, error)
	GetListDataActionRunsWithExtraInfo(ctx context.Context, params *GetListDataActionRunsWithExtraInfoParams) ([]DataActionRunWithExtraInfo, error)
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
	err := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).Update("status", status).Error

	return err
}

type GetListDataActionRunsParams struct {
	Ids         []int64
	Statuses    []model.DataActionRunStatus
	AccountUuid uuid.UUID
}

func (r *dataActionRunRepo) GetListDataActionRuns(ctx context.Context, params *GetListDataActionRunsParams) ([]model.DataActionRun, error) {
	query := r.WithContext(ctx).Table(r.TableName)
	if len(params.Ids) > 0 {
		query.Where("id IN ?", params.Ids)
	}
	if len(params.Statuses) > 0 {
		query.Where("status IN ?", params.Statuses)
	}

	emptyUuid, _ := uuid.Parse("")
	if params.AccountUuid != emptyUuid {
		query = query.Where("account_uuid = ?", params.AccountUuid)
	}

	var dataActionRuns []model.DataActionRun
	err := query.Find(&dataActionRuns).Error
	if err != nil {
		return nil, err
	}

	return dataActionRuns, nil
}

type GetListDataActionRunsWithExtraInfoParams struct {
	Ids         []int64
	Statuses    []model.DataActionRunStatus
	AccountUuid uuid.UUID
}

type DataActionRunWithExtraInfo struct {
	ID          int64                     `gorm:"column:id"`
	ActionId    int64                     `gorm:"column:action_id"`
	RunId       int64                     `gorm:"column:run_id"`
	DagRunId    string                    `gorm:"column:dag_run_id"`
	Status      model.DataActionRunStatus `gorm:"column:status"`
	Error       string                    `gorm:"column:error"`
	AccountUuid uuid.UUID                 `gorm:"column:account_uuid"`
	CreatedAt   time.Time                 `gorm:"column:created_at"`
	UpdatedAt   time.Time                 `gorm:"column:updated_at"`
	// extra
	DagId string `gorm:"column:dag_id"`
}

func (r *dataActionRunRepo) GetListDataActionRunsWithExtraInfo(ctx context.Context, params *GetListDataActionRunsWithExtraInfoParams) ([]DataActionRunWithExtraInfo, error) {
	query := r.WithContext(ctx).Table(r.TableName)
	if len(params.Ids) > 0 {
		query.Where("data_action_run.id IN ?", params.Ids)
	}
	if len(params.Statuses) > 0 {
		query.Where("data_action_run.status IN ?", params.Statuses)
	}
	emptyUuid, _ := uuid.Parse("")
	if params.AccountUuid != emptyUuid {
		query = query.Where("data_action_run.account_uuid = ?", params.AccountUuid)
	}

	query = query.Joins("LEFT JOIN data_action ON data_action.id = data_action_run.action_id").Select(
		"data_action_run.id AS id, " +
			"data_action_run.action_id AS action_id, " +
			"data_action_run.run_id AS run_id, " +
			"data_action_run.dag_run_id AS dag_run_id, " +
			"data_action_run.status AS status, " +
			"data_action_run.error AS error, " +
			"data_action_run.account_uuid AS account_uuid, " +
			"data_action.dag_id AS dag_id ",
	)

	var dataActionRuns []DataActionRunWithExtraInfo
	err := query.Scan(&dataActionRuns).Error
	if err != nil {
		return nil, err
	}

	return dataActionRuns, nil
}
