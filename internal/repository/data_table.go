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
	UpdateStatusDataTable(ctx context.Context, id int64, status model.DataTableStatus) error
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
		Status:      model.DataTableStatus_DRAFT,
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
		return nil, err
	}

	return &dataTable, nil
}

type UpdateDataTableParams struct {
	ID     int64
	Name   string
	Schema pqtype.NullRawMessage
	Status model.DataTableStatus
}

func (r *dataTableRepo) UpdateDataTable(ctx context.Context, params *UpdateDataTableParams) error {
	dataTable := &model.DataTable{
		ID:     params.ID,
		Name:   params.Name,
		Schema: params.Schema,
		Status: params.Status,
	}

	updateErr := r.WithContext(ctx).Table(r.TableName).Where("id = ?", params.ID).Updates(&dataTable).Error
	if updateErr != nil {
		return updateErr
	}

	return nil
}

type ListDataTablesFilters struct {
	Name         string
	AccountUuid  string
	DataTableIds []int64
}

func (r *dataTableRepo) ListDataTables(ctx context.Context, filter *ListDataTablesFilters) ([]model.DataTable, error) {
	var dataTables []model.DataTable
	query := r.WithContext(ctx).Table(r.TableName)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.AccountUuid != "" {
		query = query.Where("account_uuid = ?", filter.AccountUuid)
	}
	if len(filter.DataTableIds) > 0 {
		query = query.Where("id IN ?", filter.DataTableIds)
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

func (r *dataTableRepo) UpdateStatusDataTable(ctx context.Context, id int64, status model.DataTableStatus) error {
	return r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).Update("status", status).Error
}
