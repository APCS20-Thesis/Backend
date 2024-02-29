package data_source

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"strconv"
	"time"
)

func (b business) CreateDataActionImportFile(ctx context.Context, accountUuid string, dateTime string) (*model.DataAction, error) {
	dagId := accountUuid + "_" + dateTime
	dataAction, err := b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		ActionType:  model.ActionType_UploadDataFromFile,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      "Success",
	})
	if err != nil {
		b.log.WithName("CreateDataActionImportFile").
			WithValues("Context", ctx).
			Error(err, "Cannot create data action import file")
		return nil, err
	}
	return dataAction, nil

}

func (b business) TriggerAirflowImportFile(ctx context.Context, request *api.ImportFileRequest, accountUuid string, filePath string) error {
	deltaTableName := strconv.FormatInt(time.Now().Unix(), 10) + "customer"

	_, err := b.airflowAdapter.TriggerNewDagRunImportFile(ctx, &airflow.TriggerNewDagRunImportFileRequest{
		Config: airflow.ImportFileRequestConfig{
			AccountUuid:            accountUuid,
			DeltaTableName:         deltaTableName,
			CsvFilePath:            filePath,
			WriteMode:              "overwrite",
			CsvReadOptionHeader:    true,
			CsvReadOptionMultiline: true,
			CsvReadOptionDelimiter: ",",
		},
	}, "example_import_csv")

	return err
}

func (b business) ProcessImportFile(ctx context.Context, request *api.ImportFileRequest, accountUuid string, dateTime string, filePath string) error {
	// Create DataAction
	dataAction, err := b.CreateDataActionImportFile(ctx, accountUuid, dateTime)
	if err != nil {
		return err
	}

	// Create DataActionRun
	// TODO: implement call trigger new dag run
	dagRunId := "testDagRun"

	err = b.CreateDataActionRun(ctx, &repository.CreateDataActionRunParams{
		ActionId:    dataAction.ID,
		RunId:       0,
		AccountUuid: uuid.MustParse(accountUuid),
		Status:      model.DataActionRunStatus_Processing,
		DagRunId:    dagRunId,
	})
	if err != nil {
		return err
	}

	// Create Datasource
	configuration, _ := json.Marshal(model.FileConfiguration{
		FileName: request.FileName,
		FilePath: filePath,
	})

	mappingOptions, err := json.Marshal(request.MappingOptions)
	if err != nil {
		b.log.WithName("ProcessImportFile").
			WithValues("Mapping Options", mappingOptions).
			Error(err, "Cannot parse mappingOptions to JSON")
		return err
	}
	err = b.CreateDataSource(ctx, &repository.CreateDataSourceParams{
		Name:           request.Name,
		Description:    request.Description,
		Type:           model.DataSourceType_File,
		Configuration:  pqtype.NullRawMessage{RawMessage: configuration, Valid: true},
		MappingOptions: pqtype.NullRawMessage{RawMessage: mappingOptions, Valid: true},
		DeltaTableName: request.DeltaTableName,
		AccountUuid:    uuid.MustParse(accountUuid),
	})
	if err != nil {
		return err
	}

	err = b.TriggerAirflowImportFile(ctx, request, accountUuid, filePath)
	if err != nil {
		return err
	}

	return nil
}
