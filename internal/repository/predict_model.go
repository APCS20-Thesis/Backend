package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type PredictModelRepository interface {
	CreatePredictModel(ctx context.Context, params *CreatePredictModelParams) (*model.PredictModel, error)
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
