package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	ImportCsvTransaction(ctx context.Context, params *ImportCsvTransactionParams, airflowAdapter airflow.AirflowAdapter) error
}

type transactionRepo struct {
	*gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepo{db}
}

type ImportCsvTransactionParams struct {
	DataSourceName           string
	DatSourceDescription     string
	DataSourceType           model.DataSourceType
	DataSourceConfigurations pqtype.NullRawMessage
	AccountUuid              uuid.UUID
	TableId                  int64
	NewTableName             string
	MappingOptions           pqtype.NullRawMessage
	TableNameInSource        string

	DataActionType   model.ActionType
	Schedule         string
	DagId            string
	S3Configurations *airflow.S3Configurations
	WriteMode        airflow.DeltaWriteMode
	CsvReadOptions   *api.ImportCsvRequest_ImportCsvConfigurations
	Headers          []string
}

func (r transactionRepo) ImportCsvTransaction(ctx context.Context, params *ImportCsvTransactionParams, airflowAdapter airflow.AirflowAdapter) error {
	dataSource := &model.DataSource{
		Name:           params.DataSourceName,
		Description:    params.DatSourceDescription,
		Type:           params.DataSourceType,
		Configurations: params.DataSourceConfigurations,
		AccountUuid:    params.AccountUuid,
		Status:         model.DataSourceStatus_Processing,
	}
	var dataTable model.DataTable

	tx := r.DB.Begin()
	err := tx.WithContext(ctx).Table("data_source").Create(&dataSource).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if params.TableId > 0 {
		dataTable = model.DataTable{Status: model.TableStatus_UPDATING}
		err = tx.WithContext(ctx).Table("data_table").
			Where("id = ?", params.TableId).
			Updates(&dataTable).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		dataTable = model.DataTable{
			Name:        params.NewTableName,
			Status:      model.TableStatus_DRAFT,
			AccountUuid: params.AccountUuid,
			Schema: pqtype.NullRawMessage{
				RawMessage: []byte("{}"),
				Valid:      false,
			},
		}
		err = tx.WithContext(ctx).Table("data_table").Create(&dataTable).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	sourceTableMap := &model.SourceTableMap{
		TableId:           dataTable.ID,
		SourceId:          dataSource.ID,
		MappingOptions:    params.MappingOptions,
		TableNameInSource: params.TableNameInSource,
	}

	err = tx.WithContext(ctx).Table("source_table_map").Create(&sourceTableMap).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = airflowAdapter.TriggerGenerateDagImportCsv(ctx, &airflow.TriggerGenerateDagImportCsvRequest{
		Config: airflow.ImportCsvRequestConfig{
			DagId:            params.DagId,
			AccountUuid:      params.AccountUuid.String(),
			DeltaTableName:   dataTable.Name,
			S3Configurations: params.S3Configurations,
			WriteMode:        params.WriteMode,
			CsvReadOptions:   params.CsvReadOptions,
			Headers:          params.Headers,
		},
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	dataAction := &model.DataAction{
		ActionType:       params.DataActionType,
		Status:           model.DataActionStatus_Pending,
		RunCount:         0,
		Schedule:         params.Schedule,
		DagId:            params.DagId,
		AccountUuid:      params.AccountUuid,
		SourceTableMapId: sourceTableMap.ID,
	}

	err = tx.WithContext(ctx).Table("data_action").Create(&dataAction).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	dataActionRun := &model.DataActionRun{
		ActionId:    dataAction.ID,
		RunId:       int64(dataAction.RunCount),
		Status:      model.DataActionRunStatus_Processing,
		AccountUuid: params.AccountUuid,
	}
	err = tx.WithContext(ctx).Table("data_action_run").Create(&dataActionRun).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
