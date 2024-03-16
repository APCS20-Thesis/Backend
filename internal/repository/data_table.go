package repository

import (
	"context"
	"errors"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type DataTableRepository interface {
	CreateDataTable(ctx context.Context, params *CreateDataTableParams) (*model.DataTable, error)
	GetDataTable(ctx context.Context, id int64) (*model.DataTable, error)
	UpdateDataTable(ctx context.Context, params *UpdateDataTableParams) error
	ListDataTables(ctx context.Context, filter *ListDataTablesFilters) ([]model.DataTable, error)
}

type dataTableRepo struct {
	*gorm.DB
	TableName string
}

func NewDataTableRepository(db *gorm.DB) DataTableRepository {
	return &dataTableRepo{db, model.DataTable{}.TableName()}
}

type CreateDataTableParams struct {
	Name        string
	Schema      pqtype.NullRawMessage
	AccountUuid uuid.UUID
}

func (r *dataTableRepo) CreateDataTable(ctx context.Context, params *CreateDataTableParams) (*model.DataTable, error) {
	dataTable := &model.DataTable{
		Name:        params.Name,
		AccountUuid: params.AccountUuid,
		Schema:      params.Schema,
	}

	createErr := r.WithContext(ctx).Table(r.TableName).Create(&dataTable).Error
	if createErr != nil {
		return nil, createErr
	}

	return dataTable, nil
}

func (r *dataTableRepo) GetDataTable(ctx context.Context, id int64) (*model.DataTable, error) {
	var dataTable model.DataTable
	err := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).First(&dataTable).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &dataTable, nil
}

type UpdateDataTableParams struct {
	ID     int64
	Name   string
	Schema pqtype.NullRawMessage
}

func (r *dataTableRepo) UpdateDataTable(ctx context.Context, params *UpdateDataTableParams) error {
	dataTable := &model.DataTable{
		ID:     params.ID,
		Name:   params.Name,
		Schema: params.Schema,
	}

	updateErr := r.WithContext(ctx).Table(r.TableName).Where("id = ?", params.ID).Updates(&dataTable).Error
	if updateErr != nil {
		return updateErr
	}

	return nil
}

type ListDataTablesFilters struct {
	Name        string
	AccountUuid uuid.UUID
}

func (r *dataTableRepo) ListDataTables(ctx context.Context, filter *ListDataTablesFilters) ([]model.DataTable, error) {
	var dataTables []model.DataTable
	query := r.WithContext(ctx).Table(r.TableName)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.AccountUuid.String() != "" {
		query = query.Where("account_uuid = ?", filter.AccountUuid)
	}
	err := query.Find(&dataTables).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return dataTables, nil
}
