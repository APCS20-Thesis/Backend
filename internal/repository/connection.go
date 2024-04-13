package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ConnectionRepository interface {
	CreateConnection(ctx context.Context, params *CreateConnectionParams) (*model.Connection, error)
	GetConnection(ctx context.Context, id int64) (*model.Connection, error)
	UpdateConnection(ctx context.Context, params *UpdateConnectionParams) error
	ListConnections(ctx context.Context, filter *FilterConnection) ([]model.Connection, error)
	DeleteConnection(ctx context.Context, id int64) error
}

type ConnectionRepo struct {
	*gorm.DB
	TableName string
}

func NewConnectionRepository(db *gorm.DB) ConnectionRepository {
	return &ConnectionRepo{db, model.Connection{}.TableName()}
}

type CreateConnectionParams struct {
	Name           string
	Type           model.ConnectionType
	Configurations pqtype.NullRawMessage
	AccountUuid    uuid.UUID
}

func (r *ConnectionRepo) CreateConnection(ctx context.Context, params *CreateConnectionParams) (*model.Connection, error) {
	var existConnectionName int64
	err := r.WithContext(ctx).Table(r.TableName).Where("name = ? AND account_uuid = ?", params.Name, params.AccountUuid).Count(&existConnectionName).Error
	if err != nil {
		return nil, err
	}
	if existConnectionName > 0 {
		return nil, status.Error(codes.AlreadyExists, "This name already exists")
	}
	connection := &model.Connection{
		Name:           params.Name,
		AccountUuid:    params.AccountUuid,
		Configurations: params.Configurations,
		Type:           params.Type,
	}

	createErr := r.WithContext(ctx).Table(r.TableName).Create(&connection).Error
	if createErr != nil {
		return nil, createErr
	}

	return connection, nil
}

func (r *ConnectionRepo) GetConnection(ctx context.Context, id int64) (*model.Connection, error) {
	var connection model.Connection
	err := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).First(&connection).Error
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

type UpdateConnectionParams struct {
	ID             int64
	Name           string
	Configurations pqtype.NullRawMessage
	AccountUuid    uuid.UUID
}

func (r *ConnectionRepo) UpdateConnection(ctx context.Context, params *UpdateConnectionParams) error {
	var existConnectionName int64
	err := r.WithContext(ctx).Table(r.TableName).
		Where("name = ? AND account_uuid = ? AND id <> ?", params.Name, params.AccountUuid, params.ID).
		Count(&existConnectionName).Error
	if err != nil {
		return err
	}
	if existConnectionName > 0 {
		return status.Error(codes.AlreadyExists, "This name already exists")
	}
	updateErr := r.WithContext(ctx).Table(r.TableName).Where("id = ?", params.ID).
		Updates(model.Connection{
			Name:           params.Name,
			Configurations: params.Configurations,
		}).
		Error
	if updateErr != nil {
		return updateErr
	}
	return nil
}

type FilterConnection struct {
	Name        string
	Type        model.ConnectionType
	AccountUuid uuid.UUID
}

func (r *ConnectionRepo) ListConnections(ctx context.Context, filter *FilterConnection) ([]model.Connection, error) {
	var connections []model.Connection
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
	err := query.Find(&connections).Error
	if err != nil {
		return nil, err
	}
	return connections, nil
}

func (r *ConnectionRepo) DeleteConnection(ctx context.Context, id int64) error {
	deleteErr := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).Delete(&model.Connection{}).Error
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}
