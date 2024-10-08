package repository

import (
	"context"
	"errors"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type DataSourceRepository interface {
	CreateDataSource(ctx context.Context, params *CreateDataSourceParams) (*model.DataSource, error)
	GetDataSource(ctx context.Context, id int64) (*model.DataSource, error)
	UpdateDataSource(ctx context.Context, params *UpdateDataSourceParams) error
	ListDataSources(ctx context.Context, filter *ListDataSourcesFilters) (*ListDataSourcesResult, error)
}

type dataSourceRepo struct {
	*gorm.DB
	TableName string
}

func NewDataSourceRepository(db *gorm.DB) DataSourceRepository {
	return &dataSourceRepo{db, model.DataSource{}.TableName()}
}

type CreateDataSourceParams struct {
	Tx             *gorm.DB
	Name           string
	Description    string
	Type           model.DataSourceType
	Configurations pqtype.NullRawMessage
	AccountUuid    uuid.UUID
	ConnectionId   int64
}

func (r *dataSourceRepo) CreateDataSource(ctx context.Context, params *CreateDataSourceParams) (*model.DataSource, error) {
	dataSource := &model.DataSource{
		Name:           params.Name,
		Description:    params.Description,
		Type:           params.Type,
		Configurations: params.Configurations,
		ConnectionId:   params.ConnectionId,
		AccountUuid:    params.AccountUuid,
	}

	if params.Tx != nil {
		createErr := params.Tx.WithContext(ctx).Table(r.TableName).Create(&dataSource).Error
		return dataSource, createErr
	}

	createErr := r.WithContext(ctx).Table(r.TableName).Create(&dataSource).Error
	if createErr != nil {
		return nil, createErr
	}

	return dataSource, nil
}

func (r *dataSourceRepo) GetDataSource(ctx context.Context, id int64) (*model.DataSource, error) {
	var dataSource model.DataSource
	err := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).First(&dataSource).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &dataSource, nil
}

type UpdateDataSourceParams struct {
	ID             int64
	Name           string
	Type           model.DataSourceType
	Configurations pqtype.NullRawMessage
	MappingOptions pqtype.NullRawMessage
	AccountUuid    uuid.UUID
	Status         model.DataSourceStatus
}

func (r *dataSourceRepo) UpdateDataSource(ctx context.Context, params *UpdateDataSourceParams) error {
	dataSource := &model.DataSource{
		ID:             params.ID,
		Name:           params.Name,
		Type:           params.Type,
		Configurations: params.Configurations,
		AccountUuid:    params.AccountUuid,
	}

	updateErr := r.WithContext(ctx).Table(r.TableName).Where("id = ?", params.ID).Updates(&dataSource).Error
	if updateErr != nil {
		return updateErr
	}
	return nil
}

type ListDataSourcesFilters struct {
	Type        model.DataSourceType
	AccountUuid uuid.UUID
	Name        string
	Page        int
	PageSize    int
}

type ListDataSourcesResult struct {
	DataSource []model.DataSource
	Count      int64
}

func (r *dataSourceRepo) ListDataSources(ctx context.Context, filter *ListDataSourcesFilters) (*ListDataSourcesResult, error) {
	var (
		dataSources []model.DataSource
		count       int64
	)
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
	err := query.Count(&count).Scopes(Paginate(filter.Page, filter.PageSize)).Find(&dataSources).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &ListDataSourcesResult{
		DataSource: dataSources,
		Count:      count,
	}, nil
}
