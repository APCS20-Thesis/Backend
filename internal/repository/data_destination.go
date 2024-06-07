package repository

import (
	"context"
	"errors"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type DataDestinationRepository interface {
	CreateDataDestination(ctx context.Context, params *CreateDataDestinationParams) (*model.DataDestination, error)
	GetDataDestination(ctx context.Context, id int64) (*model.DataDestination, error)
	//ListDataDestinations(ctx context.Context, params *ListDataDestinationFilters) ([]model.DataDestination, error)
}

type dataDestinationRepo struct {
	*gorm.DB
	TableName string
}

func NewDataDestinationRepository(db *gorm.DB) DataDestinationRepository {
	return &dataDestinationRepo{db, model.DataDestination{}.TableName()}
}

type CreateDataDestinationParams struct {
	Name          string
	AccountUuid   uuid.UUID
	Type          model.DataDestinationType
	Configuration pqtype.NullRawMessage
	ConnectionId  int64
}

func (r *dataDestinationRepo) CreateDataDestination(ctx context.Context, params *CreateDataDestinationParams) (*model.DataDestination, error) {
	dest := &model.DataDestination{
		Name:           params.Name,
		Type:           params.Type,
		Configurations: params.Configuration,
		AccountUuid:    params.AccountUuid,
		ConnectionId:   params.ConnectionId,
	}

	createErr := r.WithContext(ctx).Table(r.TableName).Create(&dest).Error
	if createErr != nil {
		return nil, createErr
	}

	return dest, nil
}

func (r *dataDestinationRepo) GetDataDestination(ctx context.Context, id int64) (*model.DataDestination, error) {
	var destination model.DataDestination
	err := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).First(&destination).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Not found destination")
		}
		return nil, err
	}

	return &destination, nil
}
