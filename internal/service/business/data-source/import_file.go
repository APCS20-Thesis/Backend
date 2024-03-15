package data_source

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/constants"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

func (b business) CreateDataActionImportFile(ctx context.Context, accountUuid string, dateTime string) (*model.DataAction, error) {
	dagId := accountUuid + "_" + dateTime
	dataAction, err := b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		ActionType:  model.ActionType_UploadDataFromFile,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
	})
	if err != nil {
		b.log.WithName("CreateDataActionImportFile").
			WithValues("Context", ctx).
			Error(err, "Cannot create data action import file")
		return nil, err
	}
	return dataAction, nil

}

func (b business) TriggerAirflowGenerateImportFile(ctx context.Context, request *api.ImportFileRequest, accountUuid string, dateTime string) error {
	_, err := b.airflowAdapter.TriggerGenerateDagImportFile(ctx, &airflow.TriggerGenerateDagImportFileRequest{
		Config: airflow.ImportFileRequestConfig{
			AccountUuid:            accountUuid,
			DeltaTableName:         request.DeltaTableName,
			BucketName:             constants.S3BucketName,
			Key:                    "data/" + accountUuid + "/" + dateTime + "_" + request.GetFileName(),
			WriteMode:              "overwrite",
			CsvReadOptionHeader:    request.Configurations.SkipRow > 0,
			CsvReadOptionMultiline: request.Configurations.Multiline,
			CsvReadOptionDelimiter: request.Configurations.Delimiter,
			CsvReadOptionSkipRow:   request.Configurations.SkipRow,
			DagId:                  accountUuid + "_" + dateTime,
		},
	}, request.FileType)

	return err
}

func (b business) ProcessImportFile(ctx context.Context, request *api.ImportFileRequest, accountUuid string, dateTime string) error {
	// Create DataAction
	dataAction, err := b.CreateDataActionImportFile(ctx, accountUuid, dateTime)
	if err != nil {
		return err
	}

	// Generate DagRun
	err = b.TriggerAirflowGenerateImportFile(ctx, request, accountUuid, dateTime)
	if err != nil {
		return err
	}

	// Create DataActionRun
	_, err = b.CreateDataActionRun(ctx, &repository.CreateDataActionRunParams{
		ActionId:    dataAction.ID,
		RunId:       0,
		AccountUuid: uuid.MustParse(accountUuid),
		Status:      model.DataActionRunStatus_Processing,
		DagRunId:    "",
	})
	if err != nil {
		return err
	}

	// Create Datasource
	configurations, _ := json.Marshal(model.FileConfigurations{
		FileName:      request.FileName,
		BucketName:    constants.S3BucketName,
		Key:           "data/" + accountUuid + "/" + dateTime + "_" + request.GetFileName(),
		CsvReadOption: request.Configurations,
	})

	mappingOptions, err := json.Marshal(request.MappingOptions)
	if err != nil {
		b.log.WithName("ProcessImportFile").
			WithValues("Mapping Options", mappingOptions).
			Error(err, "Cannot parse mappingOptions to JSON")
		return err
	}
	_, err = b.CreateDataSource(ctx, &repository.CreateDataSourceParams{
		Name:           request.Name,
		Description:    request.Description,
		Type:           model.DataSourceType_File,
		Configurations: pqtype.NullRawMessage{RawMessage: configurations, Valid: true},
		MappingOptions: pqtype.NullRawMessage{RawMessage: mappingOptions, Valid: true},
		AccountUuid:    uuid.MustParse(accountUuid),
	})
	if err != nil {
		return err
	}

	return nil
}
