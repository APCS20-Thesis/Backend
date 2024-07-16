package repository

import (
	"context"

	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type PredictModelRepository interface {
	CreatePredictModel(ctx context.Context, params *CreatePredictModelParams) (*model.PredictModel, error)
	ListPredictModels(ctx context.Context, params *ListPredictModelsParams) (*ListPredictModelsResult, error)
	GetPredictModel(ctx context.Context, id int64) (*model.PredictModel, error)
	UpdatePredictModel(ctx context.Context, params *UpdatePredictModelParams) error
}

type predictModelRepo struct {
	*gorm.DB
	TableName string
}

func NewPredictModelRepository(db *gorm.DB) PredictModelRepository {
	return &predictModelRepo{
		DB:        db,
		TableName: model.PredictModel{}.TableName(),
	}
}

type CreatePredictModelParams struct {
	Tx                  *gorm.DB
	Name                string
	MasterSegmentId     int64
	TrainConfigurations pqtype.NullRawMessage
}

func (r *predictModelRepo) CreatePredictModel(ctx context.Context, params *CreatePredictModelParams) (*model.PredictModel, error) {
	predictModel := model.PredictModel{
		Name:                params.Name,
		MasterSegmentId:     params.MasterSegmentId,
		TrainConfigurations: params.TrainConfigurations,
		Status:              model.PredictModelStatus_DRAFT,
	}

	var createErr error
	if params.Tx != nil {
		createErr = params.Tx.WithContext(ctx).Table(r.TableName).Create(&predictModel).Error
	} else {
		createErr = r.WithContext(ctx).Table(r.TableName).Create(&predictModel).Error
	}
	if createErr != nil {
		return nil, createErr
	}

	return &predictModel, nil
}

type ListPredictModelsParams struct {
	Page            int
	PageSize        int
	MasterSegmentId int64
	Statuses        []model.PredictModelStatus
	Ids             []int64
}

type ListPredictModelsResult struct {
	PredictModels []model.PredictModel
	Count         int64
}

func (r *predictModelRepo) ListPredictModels(ctx context.Context, params *ListPredictModelsParams) (*ListPredictModelsResult, error) {
	var (
		models []model.PredictModel
		count  int64
	)

	query := r.WithContext(ctx).Table(r.TableName)
	if params.MasterSegmentId > 0 {
		query.Where("master_segment_id = ?", params.MasterSegmentId)
	}
	if len(params.Statuses) != 0 {
		query.Where("status IN ?", params.Statuses)
	}
	if len(params.Ids) > 0 {
		query.Where("id IN ?", params.Ids)
	}
	err := query.Count(&count).Scopes(Paginate(params.Page, params.PageSize)).Find(&models).Error
	if err != nil {
		return nil, err
	}

	return &ListPredictModelsResult{
		PredictModels: models,
		Count:         count,
	}, nil
}

func (r *predictModelRepo) GetPredictModel(ctx context.Context, id int64) (*model.PredictModel, error) {
	var md model.PredictModel

	err := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).First(&md).Error
	if err != nil {
		return nil, err
	}

	return &md, nil
}

type UpdatePredictModelParams struct {
	Tx     *gorm.DB
	Id     int64
	Status model.PredictModelStatus
}

func (r *predictModelRepo) UpdatePredictModel(ctx context.Context, params *UpdatePredictModelParams) error {
	predictModel := model.PredictModel{
		ID:     params.Id,
		Status: params.Status,
	}
	var updateErr error
	if params.Tx != nil {
		updateErr = params.Tx.WithContext(ctx).Table(r.TableName).Where("id = ?", params.Id).Updates(&predictModel).Error
	} else {
		updateErr = r.WithContext(ctx).Table(r.TableName).Where("id = ?", params.Id).Updates(&predictModel).Error
	}
	if updateErr != nil {
		return updateErr
	}
	return nil
}
