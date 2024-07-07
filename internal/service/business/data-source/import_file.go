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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (b business) ProcessImportCsv(ctx context.Context, request *api.ImportCsvRequest, accountUuid string, dateTime string) error {

	s3Configurations := &airflow.S3Configurations{
		AccessKeyId:     b.config.S3StorageConfig.AccessKeyID,
		SecretAccessKey: b.config.S3StorageConfig.SecretAccessKey,
		BucketName:      constants.S3BucketName,
		Region:          b.config.S3StorageConfig.Region,
		Key:             "data/" + accountUuid + "/" + dateTime + "_" + request.GetFileName(),
	}
	actionType := model.ActionType_ImportDataFromFile

	var headers []string
	for _, mapping := range request.MappingOptions {
		headers = append(headers, mapping.DestinationFieldName)
	}

	configurations, err := json.Marshal(model.CsvConfigurations{
		FileName:      request.FileName,
		ConnectionId:  0,
		Key:           s3Configurations.Key,
		CsvReadOption: request.Configurations,
	})
	if err != nil {
		b.log.WithName("ProcessImportCsv").
			Error(err, "Cannot parse configurations to JSON")
		return err
	}

	mappingOptions, err := json.Marshal(request.MappingOptions)
	if err != nil {
		b.log.WithName("ProcessImportCsv").
			Error(err, "Cannot parse mappingOptions to JSON")
		return err
	}

	schema := make([]model.SchemaUnit, 0, len(request.MappingOptions))
	for _, mappingOption := range request.MappingOptions {
		schema = append(schema, model.SchemaUnit{
			ColumnName: mappingOption.DestinationFieldName,
		})
	}
	rawSchema, err := json.Marshal(schema)
	if err != nil {
		b.log.WithName("ProcessImportCsv").Error(err, "Cannot parse schema to JSON")
		return err
	}

	if request.TableId > 0 {
		dataTable, err := b.repository.DataTableRepository.GetDataTable(ctx, request.TableId)
		if err != nil {
			b.log.WithName("ProcessImportCsv").WithValues("tableId", request.TableId).Error(err, "Cannot get data table")
			return err
		}
		if dataTable.AccountUuid.String() != accountUuid {
			b.log.WithName("ProcessImportCsv").WithValues("tableId", request.TableId).Error(err, "No have permission with dataTable")
			return status.Error(codes.PermissionDenied, "No have permission with dataTable")
		}
	} else {
		err = b.repository.DataTableRepository.CheckExistsDataTableName(ctx, request.NewTableName, accountUuid)
		if err != nil {
			b.log.WithName("ProcessImportCsv").WithValues("newTableName", request.NewTableName).Error(err, "Check exist data table")
			return err
		}
	}

	err = b.repository.TransactionRepository.ImportCsvTransaction(ctx, &repository.ImportCsvTransactionParams{
		DataSourceName:           request.Name,
		DatSourceDescription:     request.Description,
		DataSourceType:           model.DataSourceType_FileCsv,
		DataSourceConfigurations: pqtype.NullRawMessage{RawMessage: configurations, Valid: true},
		AccountUuid:              uuid.MustParse(accountUuid),
		TableId:                  request.TableId,
		NewTableName:             request.NewTableName,
		MappingOptions:           pqtype.NullRawMessage{RawMessage: mappingOptions, Valid: true},
		DataActionType:           actionType,
		Schedule:                 "",
		DagId:                    "import_csv_" + accountUuid + "_" + dateTime,
		S3Configurations:         s3Configurations,
		WriteMode:                airflow.DeltaWriteMode(request.WriteMode),
		CsvReadOptions:           request.Configurations,
		Headers:                  headers,
		Schema:                   pqtype.NullRawMessage{RawMessage: rawSchema, Valid: true},
	}, b.airflowAdapter)

	if err != nil {
		b.log.WithName("ProcessImportCsv").
			WithValues("Request", request).
			Error(err, "Transaction failed")
		return err
	}

	return nil
}

func (b business) ProcessImportCsvFromS3(ctx context.Context, request *api.ImportCsvFromS3Request, accountUuid string, dateTime string) error {
	var s3Configurations *airflow.S3Configurations
	var actionType model.ActionType

	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.ConnectionId)
	if err != nil {
		return err
	}
	if connection.AccountUuid != uuid.MustParse(accountUuid) {
		return status.Error(codes.PermissionDenied, "No have permission with connection")
	}
	if connection.Type != model.ConnectionType_S3 {
		return status.Error(codes.InvalidArgument, "Invalid connection")
	}
	actionType = model.ActionType_ImportDataFromS3
	var configuration model.S3Configurations
	err = json.Unmarshal(connection.Configurations.RawMessage, &configuration)

	s3Configurations = &airflow.S3Configurations{
		AccessKeyId:     configuration.AccessKeyId,
		SecretAccessKey: configuration.SecretAccessKey,
		BucketName:      configuration.BucketName,
		Region:          configuration.Region,
		Key:             request.Key,
	}

	var headers []string
	for _, mapping := range request.MappingOptions {
		headers = append(headers, mapping.DestinationFieldName)
	}

	configurations, err := json.Marshal(model.CsvConfigurations{
		FileName:      request.FileName,
		ConnectionId:  request.ConnectionId,
		Key:           s3Configurations.Key,
		CsvReadOption: request.Configurations,
	})
	if err != nil {
		b.log.WithName("ProcessImportCsvFromS3").
			Error(err, "Cannot parse configurations to JSON")
		return err
	}

	mappingOptions, err := json.Marshal(request.MappingOptions)
	if err != nil {
		b.log.WithName("ProcessImportCsvFromS3").
			Error(err, "Cannot parse mappingOptions to JSON")
		return err
	}

	schema := make([]model.SchemaUnit, 0, len(request.MappingOptions))
	for _, mappingOption := range request.MappingOptions {
		schema = append(schema, model.SchemaUnit{
			ColumnName: mappingOption.DestinationFieldName,
		})
	}
	rawSchema, err := json.Marshal(schema)
	if err != nil {
		b.log.WithName("ProcessImportCsvFromS3").Error(err, "Cannot parse schema to JSON")
		return err
	}

	if request.TableId > 0 {
		dataTable, err := b.repository.DataTableRepository.GetDataTable(ctx, request.TableId)
		if err != nil {
			b.log.WithName("ProcessImportCsvFromS3").WithValues("tableId", request.TableId).Error(err, "Cannot get data table")
			return err
		}
		if dataTable.AccountUuid.String() != accountUuid {
			b.log.WithName("ProcessImportCsvFromS3").WithValues("tableId", request.TableId).Error(err, "No have permission with dataTable")
			return status.Error(codes.PermissionDenied, "No have permission with dataTable")
		}
	} else {
		err = b.repository.DataTableRepository.CheckExistsDataTableName(ctx, request.NewTableName, accountUuid)
		if err != nil {
			b.log.WithName("ProcessImportCsvFromS3").WithValues("newTableName", request.NewTableName).Error(err, "Check exist data table")
			return err
		}
	}

	err = b.repository.TransactionRepository.ImportCsvTransaction(ctx, &repository.ImportCsvTransactionParams{
		DataSourceName:           request.Name,
		DatSourceDescription:     request.Description,
		DataSourceType:           model.DataSourceType_FileCsv,
		DataSourceConfigurations: pqtype.NullRawMessage{RawMessage: configurations, Valid: true},
		AccountUuid:              uuid.MustParse(accountUuid),
		TableId:                  request.TableId,
		NewTableName:             request.NewTableName,
		MappingOptions:           pqtype.NullRawMessage{RawMessage: mappingOptions, Valid: true},
		DataActionType:           actionType,
		Schedule:                 "",
		DagId:                    "import_csv_s3_" + accountUuid + "_" + dateTime,
		S3Configurations:         s3Configurations,
		WriteMode:                airflow.DeltaWriteMode(request.WriteMode),
		CsvReadOptions:           request.Configurations,
		Headers:                  headers,
		Schema:                   pqtype.NullRawMessage{RawMessage: rawSchema, Valid: true},
		ConnectionId:             request.ConnectionId,
	}, b.airflowAdapter)

	if err != nil {
		b.log.WithName("ProcessImportCsvS3").
			WithValues("Request", request).
			Error(err, "Transaction failed")
		return err
	}

	return nil
}
