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
	GetListDataActionRuns(ctx context.Context, params *GetListDataActionRunsParams) (*GetListDataActionRunsResult, error)
	GetListDataActionRunsWithExtraInfo(ctx context.Context, params *GetListDataActionRunsWithExtraInfoParams) ([]DataActionRunWithExtraInfo, error)
	GetTotalRunsPerDay(ctx context.Context, params *GetTotalRunsPerDayParams) ([]TotalRunsPerDay, error)
	GetTotalRunsPerType(ctx context.Context, params *GetTotalRunsPerTypeParams) ([]TotalRunsPerType, error)
}

type dataActionRunRepo struct {
	*gorm.DB
	TableName string
}

func NewDataActionRunRepository(db *gorm.DB) DataActionRunRepository {
	return &dataActionRunRepo{db, model.DataActionRun{}.TableName()}
}

type CreateDataActionRunParams struct {
	Tx          *gorm.DB
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
	var createErr error
	if params.Tx != nil {
		createErr = params.Tx.WithContext(ctx).Table(r.TableName).Create(&dataActionRun).Error
	} else {
		createErr = r.WithContext(ctx).Table(r.TableName).Create(&dataActionRun).Error
	}
	if createErr != nil {
		return nil, createErr
	}

	return dataActionRun, nil
}

func (r *dataActionRunRepo) UpdateDataActionRunStatus(ctx context.Context, id int64, status model.DataActionRunStatus) error {
	err := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).Update("status", status).Error

	return err
}

type GetListDataActionRunsParams struct {
	Ids         []int64
	ActionTypes []string
	Statuses    []model.DataActionRunStatus
	AccountUuid uuid.UUID
	Page        int
	PageSize    int
}
type GetListDataActionRunsResult struct {
	DataActionRuns []DataActionRunWithExtraInfo
	Count          int64
}

func (r *dataActionRunRepo) GetListDataActionRuns(ctx context.Context, params *GetListDataActionRunsParams) (*GetListDataActionRunsResult, error) {
	query := r.WithContext(ctx).Table(r.TableName)
	var count int64

	if len(params.Ids) > 0 {
		query.Where("data_action_run.id IN ?", params.Ids)
	}
	if len(params.ActionTypes) > 0 {
		query.Where("data_action_run.action_type IN ?", params.ActionTypes)
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
			"data_action_run.created_at AS created_at, " +
			"data_action_run.updated_at AS updated_at, " +
			"data_action.action_type AS action_type, " +
			"data_action.target_table AS target_table," +
			"data_action.object_id AS object_id ",
	)

	var dataActionRuns []DataActionRunWithExtraInfo
	err := query.Order("data_action_run.updated_at desc").Count(&count).Scopes(Paginate(params.Page, params.PageSize)).Find(&dataActionRuns).Error
	if err != nil {
		return nil, err
	}

	return &GetListDataActionRunsResult{
		DataActionRuns: dataActionRuns,
		Count:          count,
	}, nil
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
	TargetTable string                    `gorm:"column:target_table"`
	ObjectId    int64                     `gorm:"column:object_id"`
	// extra
	DagId      string           `gorm:"column:dag_id"`
	ActionType model.ActionType `gorm:"column:action_type"`
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
			"data_action.dag_id AS dag_id," +
			"data_action.action_type AS action_type," +
			"data_action.target_table AS target_table," +
			"data_action.object_id AS object_id ",
	)

	var dataActionRuns []DataActionRunWithExtraInfo
	err := query.Scan(&dataActionRuns).Error
	if err != nil {
		return nil, err
	}

	return dataActionRuns, nil
}

type TotalRunsPerDay struct {
	Date  time.Time `gorm:"column:date"`
	Total int       `gorm:"column:total"`
}

type GetTotalRunsPerDayParams struct {
	AccountUuid string
}

func (r *dataActionRunRepo) GetTotalRunsPerDay(ctx context.Context, params *GetTotalRunsPerDayParams) ([]TotalRunsPerDay, error) {

	var result []TotalRunsPerDay
	err := r.WithContext(ctx).Table(r.TableName).Select("DATE_TRUNC('day', created_at) AS date, count(*) AS total").
		Where("created_at >= ? AND account_uuid = ?", time.Now().AddDate(0, 0, -15), params.AccountUuid).
		Group("DATE_TRUNC('day', created_at)").
		Order("date ASC").
		Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

type TotalRunsPerType struct {
	Type  model.ActionType
	Total int
}
type GetTotalRunsPerTypeParams struct {
	AccountUuid string
}

func (r *dataActionRunRepo) GetTotalRunsPerType(ctx context.Context, params *GetTotalRunsPerTypeParams) ([]TotalRunsPerType, error) {
	var result []TotalRunsPerType
	err := r.WithContext(ctx).Table(r.TableName).
		Joins("LEFT JOIN data_action ON data_action.id = data_action_run.action_id").
		Where("data_action_run.account_uuid = ?", params.AccountUuid).
		Select("data_action.action_type AS type, COUNT(*) AS total").
		Group("data_action.action_type").
		Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
