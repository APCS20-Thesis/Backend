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
	ListDataDestinations(ctx context.Context, params *ListDataDestinationsParams) (*ListDataDestinationsResult, error)
}

type dataDestinationRepo struct {
	*gorm.DB
	TableName string
}

func NewDataDestinationRepository(db *gorm.DB) DataDestinationRepository {
	return &dataDestinationRepo{db, model.DataDestination{}.TableName()}
}

type CreateDataDestinationParams struct {
	Tx            *gorm.DB
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

	var createErr error
	if params.Tx != nil {
		createErr = params.Tx.WithContext(ctx).Table(r.TableName).Create(&dest).Error
	} else {
		createErr = r.WithContext(ctx).Table(r.TableName).Create(&dest).Error
	}
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

type ListDataDestinationsParams struct {
	Page        int
	PageSize    int
	Type        model.DataDestinationType
	AccountUuid string
}

type ListDataDestinationsResult struct {
	Destinations []model.DataDestination
	Count        int64
}

func (r *dataDestinationRepo) ListDataDestinations(ctx context.Context, params *ListDataDestinationsParams) (*ListDataDestinationsResult, error) {
	var (
		destinations []model.DataDestination
		count        int64
	)

	query := r.WithContext(ctx).Table(r.TableName).Where("account_uuid = ?", params.AccountUuid)
	if params.Type != "" {
		query.Where("type = ?", params.Type)
	}
	err := query.Count(&count).Scopes(Paginate(params.Page, params.PageSize)).Find(&destinations).Error
	if err != nil {
		return nil, err
	}

	return &ListDataDestinationsResult{
		Destinations: destinations,
		Count:        count,
	}, nil
}
