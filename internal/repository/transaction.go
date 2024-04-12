package repository

import (
	"context"
	"errors"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
	"strconv"
)

type TransactionRepository interface {
	ImportCsvTransaction(ctx context.Context, params *ImportCsvTransactionParams, airflowAdapter airflow.AirflowAdapter) error
	ExportDataToCSVTransaction(ctx context.Context, params *ExportDataToCSVTransactionParams, airflowAdapter airflow.AirflowAdapter) error
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
	CsvReadOptions   *api.ImportCsvConfigurations
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
		dataTable = model.DataTable{ID: params.TableId, Status: model.DataTableStatus_UPDATING}
		err = tx.WithContext(ctx).Table("data_table").
			Where("id = ?", params.TableId).
			Updates(&dataTable).
			First(&dataTable).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		dataTable = model.DataTable{
			Name:        params.NewTableName,
			Status:      model.DataTableStatus_DRAFT,
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
		TableId:         dataTable.ID,
		SourceId:        dataSource.ID,
		MappingOptions:  params.MappingOptions,
		SourceTableName: params.TableNameInSource,
	}

	err = tx.WithContext(ctx).Table("source_table_map").Create(&sourceTableMap).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	//_, err = airflowAdapter.TriggerGenerateDagImportCsv(ctx, &airflow.TriggerGenerateDagImportCsvRequest{
	//	Config: airflow.ImportCsvRequestConfig{
	//		DagId:            params.DagId,
	//		AccountUuid:      params.AccountUuid.String(),
	//		DeltaTableName:   dataTable.Name,
	//		S3Configurations: params.S3Configurations,
	//		WriteMode:        params.WriteMode,
	//		CsvReadOptions:   params.CsvReadOptions,
	//		Headers:          params.Headers,
	//	},
	//})
	//if err != nil {
	//	tx.Rollback()
	//	return err
	//}
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

type ExportDataToCSVTransactionParams struct {
	AccountUuid uuid.UUID
	TableId     int64
	S3Key       string
}

func (r transactionRepo) ExportDataToCSVTransaction(ctx context.Context, params *ExportDataToCSVTransactionParams, airflowAdapter airflow.AirflowAdapter) error {

	var fileExportRecord model.FileExportRecord
	err := r.WithContext(ctx).Table("file_export_record").Where("data_table_id = ?", params.TableId).First(&fileExportRecord).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.NewExportFileCSVTransaction(ctx, params, airflowAdapter)
	}
	if err != nil {
		return err
	}

	return r.TriggerExportFileCSVTransaction(ctx, params, airflowAdapter, fileExportRecord.DataActionId)
}

func (r transactionRepo) TriggerExportFileCSVTransaction(ctx context.Context, params *ExportDataToCSVTransactionParams, airflowAdapter airflow.AirflowAdapter, dataActionId int64) error {
	var dataAction model.DataAction
	err := r.DB.WithContext(ctx).Table("data_action").First(&dataAction, "id = ?", dataActionId).Error
	if err != nil {
		return err
	}

	triggerDagRunResp, err := airflowAdapter.TriggerNewDagRun(ctx, dataAction.DagId, &airflow.TriggerNewDagRunRequest{})
	if err != nil {
		return err
	}

	tx := r.DB.Begin()

	dataActionRun := model.DataActionRun{
		ActionId:    dataAction.ID,
		RunId:       dataAction.RunCount + 1,
		DagRunId:    triggerDagRunResp.DagRunId,
		Status:      model.DataActionRunStatus_Processing,
		AccountUuid: params.AccountUuid,
	}
	err = tx.WithContext(ctx).Table("data_action_run").Create(&dataActionRun).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.WithContext(ctx).Table("data_action").Where("id = ?", dataAction.ID).Update("run_count", dataAction.RunCount+1).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.WithContext(ctx).Table("file_export_record").Create(&model.FileExportRecord{
		DataTableId:     params.TableId,
		Format:          model.FileType_CSV,
		AccountUuid:     params.AccountUuid,
		DataActionId:    dataAction.ID,
		DataActionRunId: dataActionRun.ID,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r transactionRepo) NewExportFileCSVTransaction(ctx context.Context, params *ExportDataToCSVTransactionParams, airflowAdapter airflow.AirflowAdapter) error {
	var dataTable model.DataTable
	err := r.DB.WithContext(ctx).Table("data_table").First(&dataTable, "id = ?", params.TableId).Error
	if err != nil {
		return err
	}

	tx := r.DB.Begin()
	dagId := "export_file_" + params.AccountUuid.String() + "_" + strconv.FormatInt(params.TableId, 10)

	// trigger new dag run
	_, err = airflowAdapter.TriggerGenerateDagExportFile(ctx, &airflow.TriggerGenerateDagExportFileRequest{
		Config: airflow.ExportFileRequestConfig{
			DagId:          dagId,
			AccountUuid:    params.AccountUuid.String(),
			DeltaTableName: dataTable.TableName(),
			SavedS3Path:    params.S3Key,
		}})
	if err != nil {
		tx.Rollback()
		return err
	}

	dataAction := model.DataAction{
		ActionType:  model.ActionType_ExportDataToCSV,
		Status:      model.DataActionStatus_Pending,
		RunCount:    1,
		DagId:       dagId,
		AccountUuid: params.AccountUuid,
	}
	err = tx.WithContext(ctx).Table("data_action").Create(&dataAction).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	fileExportRecord := model.FileExportRecord{
		DataTableId:  params.TableId,
		Format:       model.FileType_CSV,
		AccountUuid:  params.AccountUuid,
		DataActionId: dataAction.ID,
		S3Key:        params.S3Key,
	}
	err = tx.WithContext(ctx).Table("file_export_record").Create(&fileExportRecord).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
