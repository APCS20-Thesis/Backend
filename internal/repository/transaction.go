package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strconv"
)

const (
	tableMasterSegment = "master_segment"
	tableSegmentTable  = "segment"
	tableAudienceTable = "audience_table"
	tableBehaviorTable = "behavior_table"
)

type TransactionRepository interface {
	ImportCsvTransaction(ctx context.Context, params *ImportCsvTransactionParams, airflowAdapter airflow.AirflowAdapter) error
	ExportDataToCSVTransaction(ctx context.Context, params *ExportDataToCSVTransactionParams, airflowAdapter airflow.AirflowAdapter) error
	CreateMasterSegmentTransaction(ctx context.Context, params *CreateMasterSegmentTransactionParams, airflowAdapter airflow.AirflowAdapter) error
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
	DataActionType           model.ActionType
	Schedule                 string
	DagId                    string
	S3Configurations         *airflow.S3Configurations
	WriteMode                airflow.DeltaWriteMode
	CsvReadOptions           *api.ImportCsvConfigurations
	Headers                  []string
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
		TableId:        dataTable.ID,
		SourceId:       dataSource.ID,
		MappingOptions: params.MappingOptions,
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
		TargetTable: model.TargetTable_SourceTableMap,
		ActionType:  params.DataActionType,
		Status:      model.DataActionStatus_Pending,
		RunCount:    0,
		Schedule:    params.Schedule,
		DagId:       params.DagId,
		AccountUuid: params.AccountUuid,
		ObjectId:    sourceTableMap.ID,
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

	var dataAction model.DataAction
	err := r.WithContext(ctx).Table(string(model.TargetTable_DataTable)).
		Where("id = ? AND action_type = ?", params.TableId, string(model.ActionType_ExportDataToCSV)).
		First(&dataAction).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.NewExportFileCSVTransaction(ctx, params, airflowAdapter)
	}
	if err != nil {
		return err
	}

	return r.TriggerExportFileCSVTransaction(ctx, params, airflowAdapter, dataAction.ID)
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
		TargetTable: model.TargetTable_DataTable,
		ActionType:  model.ActionType_ExportDataToCSV,
		Status:      model.DataActionStatus_Pending,
		RunCount:    1,
		DagId:       dagId,
		AccountUuid: params.AccountUuid,
		ObjectId:    dataTable.ID,
	}
	err = tx.WithContext(ctx).Table("data_action").Create(&dataAction).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	fileExportRecord := model.FileExportRecord{
		DataTableId: params.TableId,
		Format:      model.FileType_CSV,
		AccountUuid: params.AccountUuid,
		S3Key:       params.S3Key,
	}
	err = tx.WithContext(ctx).Table("file_export_record").Create(&fileExportRecord).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

type CreateMasterSegmentTransactionParams struct {
	MasterSegmentName  string
	Description        string
	AccountUuid        uuid.UUID
	AudienceName       string
	BuildConfiguration AudienceBuildConfiguration
	BehaviorTables     []*CreateBehaviorTableParams
}

func (r transactionRepo) CreateMasterSegmentTransaction(ctx context.Context, params *CreateMasterSegmentTransactionParams, airflowAdapter airflow.AirflowAdapter) error {
	// prepare audience table
	buildConfiguration, err := json.Marshal(params.BuildConfiguration)
	if err != nil {
		return err
	}

	err = r.Transaction(func(tx *gorm.DB) error {
		// Create master segment
		var masterSegment = model.MasterSegment{
			Description: params.Description,
			Name:        params.MasterSegmentName,
			AccountUuid: params.AccountUuid,
			Status:      model.MasterSegmentStatus_DRAFT,
		}
		txErr := tx.WithContext(ctx).Table(tableMasterSegment).Create(&masterSegment).Error
		if txErr != nil {
			return txErr
		}

		// Create audience table
		txErr = tx.WithContext(ctx).Table(tableAudienceTable).Create(&model.AudienceTable{
			MasterSegmentId:    masterSegment.ID,
			BuildConfiguration: pqtype.NullRawMessage{RawMessage: buildConfiguration, Valid: buildConfiguration != nil},
			Name:               params.AudienceName,
		}).Error
		if txErr != nil {
			return txErr
		}

		// Create behavior tables
		var (
			modelBehaviorTables []*model.BehaviorTable
		)
		for _, table := range params.BehaviorTables {
			modelBehaviorTables = append(modelBehaviorTables, &model.BehaviorTable{
				MasterSegmentId: masterSegment.ID,
				DataTableId:     table.TableId,
				ForeignKey:      table.ForeignKey,
				JoinKey:         table.JoinKey,
				Name:            table.Name,
			})
		}
		txErr = tx.WithContext(ctx).Table(tableBehaviorTable).Create(modelBehaviorTables).Error
		if txErr != nil {
			return txErr
		}

		// Airflow triggers generating dag create segment
		txErr = r.TriggerAirflowCreateMasterSegment(ctx, tx, masterSegment.ID, params, airflowAdapter)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r transactionRepo) TriggerAirflowCreateMasterSegment(ctx context.Context, tx *gorm.DB, masterSegmentId int64, params *CreateMasterSegmentTransactionParams, airflowAdapter airflow.AirflowAdapter) error {
	var mainTable model.DataTable
	err := r.WithContext(ctx).Table(model.DataTable{}.TableName()).Where("id = ? AND account_uuid = ?", params.BuildConfiguration.MainTableId, params.AccountUuid).First(&mainTable).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return status.Error(codes.PermissionDenied, "user not have access to main table")
		}
		return err
	}

	var behaviorTables []airflow.CreateMasterSegmentConfig_BehaviorTable
	for _, table := range params.BehaviorTables {
		var behaviorTable model.DataTable
		err := r.DB.Table(model.DataTable{}.TableName()).Where("id = ? AND account_uuid = ?", table.TableId, params.AccountUuid).First(&behaviorTable).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return status.Error(codes.PermissionDenied, "user not have access to one behavior table")
			}
			return err
		}
		var columns []airflow.CreateMasterSegmentConfig_TableColumns
		for _, column := range table.SelectedColumns {
			columns = append(columns, airflow.CreateMasterSegmentConfig_TableColumns{
				TableColumnName:    column.TableColumnName,
				AudienceColumnName: column.NewTableColumnName,
			})
		}
		behaviorTables = append(behaviorTables, airflow.CreateMasterSegmentConfig_BehaviorTable{
			TableName:         behaviorTable.Name,
			BehaviorTableName: table.Name,
			JoinKey:           table.JoinKey,
			ForeignKey:        table.ForeignKey,
			Columns:           columns,
		})
	}
	jsonBehaviorTables, err := json.Marshal(behaviorTables)
	if err != nil {
		return err
	}

	var attributeTables []airflow.CreateMasterSegmentConfig_AttributeTable
	for _, table := range params.BuildConfiguration.AttributeTables {
		var attributeTable model.DataTable
		err = r.DB.Table("data_table").Where("id = ? AND account_uuid = ?", table.TableId, params.AccountUuid).First(&attributeTable).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return status.Error(codes.PermissionDenied, "user not have access to attribute table")
			}
			return err
		}
		var columns []airflow.CreateMasterSegmentConfig_TableColumns
		for _, column := range table.SelectedColumns {
			columns = append(columns, airflow.CreateMasterSegmentConfig_TableColumns{
				TableColumnName:    column.TableColumnName,
				AudienceColumnName: column.NewTableColumnName,
			})
		}
		attributeTables = append(attributeTables, airflow.CreateMasterSegmentConfig_AttributeTable{
			TableName:  attributeTable.Name,
			JoinKey:    table.JoinKey,
			ForeignKey: table.ForeignKey,
			Columns:    columns,
		})
	}
	jsonAttributeTables, err := json.Marshal(attributeTables)
	if err != nil {
		return err
	}

	var mainAttributes []airflow.CreateMasterSegmentConfig_TableColumns
	for _, each := range params.BuildConfiguration.SelectedColumns {
		mainAttributes = append(mainAttributes, airflow.CreateMasterSegmentConfig_TableColumns{
			TableColumnName:    each.TableColumnName,
			AudienceColumnName: each.NewTableColumnName,
		})
	}
	jsonMainAttributes, err := json.Marshal(mainAttributes)
	if err != nil {
		return err
	}

	dagId := utils.GenerateDagIdForCreateMasterSegment(params.AccountUuid.String(), masterSegmentId)
	request := airflow.TriggerGenerateDagCreateMasterSegmentRequest{
		Config: airflow.CreateMasterSegmentConfig{
			DagId:           dagId,
			AccountUuid:     params.AccountUuid.String(),
			MasterSegmentId: masterSegmentId,
			MainTableName:   mainTable.Name,
			MainAttributes:  string(jsonMainAttributes),
			AttributeTables: string(jsonAttributeTables),
			BehaviorTables:  string(jsonBehaviorTables),
		},
	}
	err = airflowAdapter.TriggerGenerateDagCreateMasterSegment(ctx, &request)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	tx.WithContext(ctx).Table("data_action").Create(&model.DataAction{
		ActionType:  model.ActionType_CreateMasterSegment,
		Payload:     pqtype.NullRawMessage{RawMessage: payload, Valid: payload != nil},
		Status:      model.DataActionStatus_Pending,
		RunCount:    1,
		DagId:       dagId,
		TargetTable: tableMasterSegment,
		ObjectId:    masterSegmentId,
		AccountUuid: params.AccountUuid,
	})

	return nil
}
