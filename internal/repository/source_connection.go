package repository

import (
	"context"
	"errors"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type SourceConnectionRepository interface {
	CreateSourceConnection(ctx context.Context, params *CreateSourceConnectionParams) (*model.SourceConnection, error)
	GetSourceConnection(ctx context.Context, id int64) (*model.SourceConnection, error)
	UpdateSourceConnection(ctx context.Context, params *UpdateSourceConnectionParams) error
	ListSourceConnections(ctx context.Context, filter *FilterSourceConnection) ([]model.SourceConnection, error)
}

type SourceConnectionRepo struct {
	*gorm.DB
	TableName string
}

func NewSourceConnectionRepository(db *gorm.DB) SourceConnectionRepository {
	return &SourceConnectionRepo{db, model.SourceConnection{}.TableName()}
}

type CreateSourceConnectionParams struct {
	Name           string
	Type           model.ConnectionType
	Configurations pqtype.NullRawMessage
	AccountUuid    uuid.UUID
}

func (r *SourceConnectionRepo) CreateSourceConnection(ctx context.Context, params *CreateSourceConnectionParams) (*model.SourceConnection, error) {
	SourceConnection := &model.SourceConnection{
		Name:           params.Name,
		AccountUuid:    params.AccountUuid,
		Configurations: params.Configurations,
		Type:           params.Type,
	}

	createErr := r.WithContext(ctx).Table(r.TableName).Create(&SourceConnection).Error
	if createErr != nil {
		return nil, createErr
	}

	return SourceConnection, nil
}

func (r *SourceConnectionRepo) GetSourceConnection(ctx context.Context, id int64) (*model.SourceConnection, error) {
	var SourceConnection model.SourceConnection
	err := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).First(&SourceConnection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &SourceConnection, nil
}

type UpdateSourceConnectionParams struct {
	ID             int64
	Name           string
	Configurations pqtype.NullRawMessage
}

func (r *SourceConnectionRepo) UpdateSourceConnection(ctx context.Context, params *UpdateSourceConnectionParams) error {
	SourceConnection := &model.SourceConnection{
		ID:             params.ID,
		Name:           params.Name,
		Configurations: params.Configurations,
	}

	updateErr := r.WithContext(ctx).Table(r.TableName).Where("id = ?", params.ID).Updates(&SourceConnection).Error
	if updateErr != nil {
		return updateErr
	}

	return nil
}

type FilterSourceConnection struct {
	Name        string
	Type        model.ConnectionType
	AccountUuid uuid.UUID
}

func (r *SourceConnectionRepo) ListSourceConnections(ctx context.Context, filter *FilterSourceConnection) ([]model.SourceConnection, error) {
	var SourceConnections []model.SourceConnection
	query := r.WithContext(ctx).Table(r.TableName)
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.AccountUuid.String() != "" {
		query = query.Where("account_uuid = ?", filter.AccountUuid)
	}
	err := query.Find(&SourceConnections).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return SourceConnections, nil
}
