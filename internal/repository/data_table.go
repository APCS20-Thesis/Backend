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

type DataTableRepository interface {
	CreateDataTable(ctx context.Context, params *CreateDataTableParams) (*model.DataTable, error)
	GetDataTable(ctx context.Context, id int64) (*model.DataTable, error)
	UpdateDataTable(ctx context.Context, params *UpdateDataTableParams) (*model.DataTable, error)
	ListDataTables(ctx context.Context, filter *ListDataTablesFilters) (*ListDataTablesResult, error)
	UpdateStatusDataTable(ctx context.Context, id int64, status model.DataTableStatus) error
	GetDataTableDeltaPath(ctx context.Context, id int64) (string, error)
	GetSourcesOfDataTables(ctx context.Context, tableIds []int64) (map[int64][]model.DataSource, error)
	GetDestinationsOfDataTables(ctx context.Context, tableIds []int64) (map[int64][]model.DataDestination, error)
	CheckExistsDataTable(ctx context.Context, tableName string, accountUuid string) error
}

type dataTableRepo struct {
	*gorm.DB
	TableName string
}

func NewDataTableRepository(db *gorm.DB) DataTableRepository {
	return &dataTableRepo{db, model.DataTable{}.TableName()}
}

type CreateDataTableParams struct {
	Tx          *gorm.DB
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

	var createErr error
	if params.Tx != nil {
		createErr = params.Tx.WithContext(ctx).Table(r.TableName).Create(&dataTable).Error
	} else {
		createErr = r.WithContext(ctx).Table(r.TableName).Create(&dataTable).Error
	}
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
	Tx     *gorm.DB
	ID     int64
	Name   string
	Schema pqtype.NullRawMessage
	Status model.DataTableStatus
}

func (r *dataTableRepo) UpdateDataTable(ctx context.Context, params *UpdateDataTableParams) (*model.DataTable, error) {
	dataTable := &model.DataTable{
		ID:     params.ID,
		Name:   params.Name,
		Schema: params.Schema,
		Status: params.Status,
	}
	var updateErr error
	if params.Tx != nil {
		updateErr = params.Tx.WithContext(ctx).Table(r.TableName).
			Where("id = ?", params.ID).
			Updates(&dataTable).
			First(&dataTable).Error
	} else {
		updateErr = r.WithContext(ctx).Table(r.TableName).
			Where("id = ?", params.ID).
			Updates(&dataTable).First(&dataTable).
			Error
	}
	if updateErr != nil {
		return nil, updateErr
	}

	return dataTable, nil
}

type ListDataTablesFilters struct {
	Name         string
	AccountUuid  string
	DataTableIds []int64
	Page         int
	PageSize     int
}

type ListDataTablesResult struct {
	Count      int64
	DataTables []model.DataTable
}

func (r *dataTableRepo) ListDataTables(ctx context.Context, filter *ListDataTablesFilters) (*ListDataTablesResult, error) {
	var (
		dataTables []model.DataTable
		count      int64
	)
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
	err := query.Count(&count).Scopes(Paginate(filter.Page, filter.PageSize)).Find(&dataTables).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &ListDataTablesResult{
		Count:      count,
		DataTables: dataTables,
	}, nil
}

func (r *dataTableRepo) UpdateStatusDataTable(ctx context.Context, id int64, status model.DataTableStatus) error {
	return r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).Update("status", status).Error
}

func (r *dataTableRepo) GetDataTableDeltaPath(ctx context.Context, id int64) (string, error) {
	var dataTable model.DataTable
	err := r.WithContext(ctx).Table(r.TableName).Where("id = ?", id).First(&dataTable).Error
	if err != nil {
		return "", err
	}

	return "/data/bronze/" + dataTable.AccountUuid.String() + "/" + dataTable.Name, err
}

type DataTableDataSource struct {
	TableId    int64
	SourceId   int64
	SourceName string
	SourceType model.DataSourceType
}

func (r *dataTableRepo) GetSourcesOfDataTables(ctx context.Context, tableIds []int64) (map[int64][]model.DataSource, error) {
	if len(tableIds) <= 0 {
		return nil, nil
	}

	query := r.WithContext(ctx).Table(r.TableName).
		Where("data_table.id IN ?", tableIds).
		Joins("JOIN source_table_map ON data_table.id = source_table_map.table_id").
		Joins("LEFT JOIN data_source ON source_table_map.source_id = data_source.id").
		Select("data_table.id AS table_id, " +
			"data_source.id AS source_id, " +
			"data_source.name AS source_name, " +
			"data_source.type AS source_type")

	var results []DataTableDataSource
	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	tableSourceMap := make(map[int64][]model.DataSource)
	for _, each := range results {
		if each.SourceId > 0 {
			tableSourceMap[each.TableId] = append(tableSourceMap[each.TableId], model.DataSource{
				ID:   each.SourceId,
				Name: each.SourceName,
				Type: each.SourceType,
			})
		}
	}

	return tableSourceMap, nil
}

type DataTableDataDestination struct {
	TableId         int64
	DestinationId   int64
	DestinationName string
	DestinationType model.DataDestinationType
}

func (r *dataTableRepo) GetDestinationsOfDataTables(ctx context.Context, tableIds []int64) (map[int64][]model.DataDestination, error) {
	if len(tableIds) <= 0 {
		return nil, nil
	}

	query := r.WithContext(ctx).Table(r.TableName).
		Where("data_table.id IN ?", tableIds).
		Joins("JOIN dest_table_map ON data_table.id = dest_table_map.table_id").
		Joins("LEFT JOIN data_destination ON dest_table_map.destination_id = data_destination.id").
		Select("data_table.id AS table_id, " +
			"data_destination.id AS destination_id, " +
			"data_destination.name AS destination_name, " +
			"data_destination.type AS destination_type")

	var results []DataTableDataDestination
	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	tableDestinationMap := make(map[int64][]model.DataDestination)
	for _, each := range results {
		if each.DestinationId > 0 {
			tableDestinationMap[each.TableId] = append(tableDestinationMap[each.TableId], model.DataDestination{
				ID:   each.DestinationId,
				Name: each.DestinationName,
				Type: each.DestinationType,
			})
		}
	}

	return tableDestinationMap, nil
}

func (r *dataTableRepo) CheckExistsDataTable(ctx context.Context, tableName string, accountUuid string) error {
	err := r.WithContext(ctx).Table(r.TableName).
		Where("name = ? AND account_uuid = ?", tableName, accountUuid).
		First(&model.DataTable{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else if err != nil {
		return err
	}
	return status.Error(codes.AlreadyExists, "Table name is already used")
}
