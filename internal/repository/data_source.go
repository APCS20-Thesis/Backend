package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type DataSourceRepository interface {
	CreateDataSource(ctx context.Context, params *CreateDataSourceParams) error
}

type dataSourceRepo struct {
	*gorm.DB
	TableName string
}

func NewDataSourceRepository(db *gorm.DB) DataSourceRepository {
	return &dataSourceRepo{db, model.DataSource{}.TableName()}
}

type CreateDataSourceParams struct {
	Name           string
	Description    string
	Type           model.DataSourceType
	Configuration  pqtype.NullRawMessage
	MappingOptions pqtype.NullRawMessage
	DeltaTableName string
	AccountUuid    uuid.UUID
}

type FileConfiguration struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
}

func (r *dataSourceRepo) CreateDataSource(ctx context.Context, params *CreateDataSourceParams) error {
	dataSource := &model.DataSource{
		Name:           params.Name,
		Description:    params.Description,
		Type:           params.Type,
		Configuration:  params.Configuration,
		AccountUuid:    params.AccountUuid,
		MappingOptions: params.MappingOptions,
		DeltaTableName: params.DeltaTableName,
	}

	createErr := r.WithContext(ctx).Table(r.TableName).Create(&dataSource).Error
	if createErr != nil {
		return createErr
	}

	return nil
}
