package data_destination

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (b business) ProcessExportDataToFile(ctx context.Context, request *api.ExportDataToFileRequest, accountUuid string) (*api.ExportDataToFileResponse, error) {

	var err error
	switch true {
	case request.TableId > 0:
		err = b.ExportDataTableToS3File(ctx, request, accountUuid)
	case request.SegmentId > 0:
		err = b.ExportSegmentToS3File(ctx, request, accountUuid)
	case request.MasterSegmentId > 0:
		err = b.ExportMasterSegmentToS3File(ctx, request, accountUuid)
	default:
		err = status.Error(codes.InvalidArgument, "require table_id, segment_id, master_segment_id")
	}
	if err != nil {
		return nil, err
	}

	return &api.ExportDataToFileResponse{
		Code:    int32(code.Code_OK),
		Message: "Success",
	}, nil
}

func (b business) ExportDataTableToS3File(ctx context.Context, request *api.ExportDataToFileRequest, accountUuid string) error {
	logger := b.log.WithName("ExportDataTableToS3File").WithValues("request", request)

	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.ConnectionId)
	if err != nil {
		logger.Error(err, "cannot get connection")
		return err
	}
	var s3Configuration model.S3Configurations
	err = json.Unmarshal(connection.Configurations.RawMessage, &s3Configuration)
	if err != nil {
		logger.Error(err, "cannot unmarshal s3 configuration")
		return err
	}

	dataTable, err := b.repository.DataTableRepository.GetDataTable(ctx, request.TableId)
	if err != nil {
		logger.Error(err, "cannot get data table")
		return err
	}

	file := request.FileName + "." + strings.ToLower(request.FileType)
	dstKey := request.FilePath + "/" + file

	dagId := utils.GenerateDagId(accountUuid, model.ActionType_ExportDataToS3CSV)
	tablePathKey := utils.GenerateDeltaTablePath(accountUuid, dataTable.Name)

	dstConfiguration, err := json.Marshal(model.S3FileDestinationConfiguration{
		FileName: request.FileName,
		FileType: request.FileType,
		FilePath: request.FilePath,
	})
	if err != nil {
		logger.Error(err, "cannot marshal destination configuration")
		return err
	}

	// BEGIN TRANSACTION
	tx := b.db.Begin()

	destination, err := b.repository.DataDestinationRepository.CreateDataDestination(ctx, &repository.CreateDataDestinationParams{
		Tx:            tx,
		Name:          connection.Name + " " + file,
		AccountUuid:   uuid.MustParse(accountUuid),
		Type:          model.DataDestinationType_S3FileCSV,
		Configuration: pqtype.NullRawMessage{RawMessage: dstConfiguration, Valid: dstConfiguration != nil},
		ConnectionId:  request.ConnectionId,
	})
	if err != nil {
		logger.Error(err, "cannot create data destination")
		tx.Rollback()
		return err
	}

	destTableMap, err := b.repository.DestTableMapRepository.CreateDestinationTableMap(ctx, &repository.CreateDestinationTableMapParams{
		Tx:            tx,
		TableId:       request.TableId,
		DestinationId: destination.ID,
		//MappingOptions: pqtype.NullRawMessage{},
	})
	if err != nil {
		logger.Error(err, "cannot create destination table map")
		tx.Rollback()
		return err
	}

	payload := airflow.ExportFileRequestConfig{
		DagId: dagId,
		Key:   tablePathKey,
		S3Configurations: airflow.S3Configurations{
			AccessKeyId:     s3Configuration.AccessKeyId,
			SecretAccessKey: s3Configuration.SecretAccessKey,
			BucketName:      s3Configuration.BucketName,
			Region:          s3Configuration.Region,
			Key:             dstKey,
		},
	}
	_, err = b.airflowAdapter.TriggerGenerateDagExportFile(ctx, &airflow.TriggerGenerateDagExportFileRequest{
		Config: payload,
	})
	if err != nil {
		logger.Error(err, "cannot trigger generate dag export file")
		tx.Rollback()
		return err
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err, "cannot marshal trigger generate dag payload")
		tx.Rollback()
		return err
	}
	dataAction, err := b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		Tx:          tx,
		TargetTable: model.TargetTable_DestTableMap,
		ActionType:  model.ActionType_ExportDataToS3CSV,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    destTableMap.ID,
		Payload:     pqtype.NullRawMessage{RawMessage: jsonPayload, Valid: jsonPayload != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		tx.Rollback()
		return err
	}

	_, err = b.repository.DataActionRunRepository.CreateDataActionRun(ctx, &repository.CreateDataActionRunParams{
		Tx:          tx,
		ActionId:    dataAction.ID,
		RunId:       0,
		Status:      model.DataActionRunStatus_Creating,
		AccountUuid: uuid.MustParse(accountUuid),
	})
	if err != nil {
		logger.Error(err, "cannot create data action run")
		tx.Rollback()
		return err
	}

	tx.Commit()
	// END TRANSACTION

	return nil
}

func (b business) ExportSegmentToS3File(ctx context.Context, request *api.ExportDataToFileRequest, accountUuid string) error {
	logger := b.log.WithName("ExportSegmentToS3File").WithValues("request", request)

	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.ConnectionId)
	if err != nil {
		logger.Error(err, "cannot get connection")
		return err
	}
	var s3Configuration model.S3Configurations
	err = json.Unmarshal(connection.Configurations.RawMessage, &s3Configuration)
	if err != nil {
		logger.Error(err, "cannot unmarshal s3 configuration")
		return err
	}

	segment, err := b.repository.SegmentRepository.GetSegment(ctx, request.SegmentId)
	if err != nil {
		logger.Error(err, "cannot get segment")
		return err
	}

	dstConfiguration, err := json.Marshal(model.S3FileDestinationConfiguration{
		FileName: request.FileName,
		FileType: request.FileType,
		FilePath: request.FilePath,
	})
	if err != nil {
		logger.Error(err, "cannot marshal destination configuration")
		return err
	}

	file := request.FileName + "." + strings.ToLower(request.FileType)
	dstKey := request.FilePath + "/" + file

	dagId := utils.GenerateDagId(accountUuid, model.ActionType_ExportDataToS3CSV)
	audiencePathKey := utils.GenerateDeltaAudiencePath(segment.MasterSegmentId)

	// BEGIN TRANSACTION
	tx := b.db.Begin()

	destination, err := b.repository.DataDestinationRepository.CreateDataDestination(ctx, &repository.CreateDataDestinationParams{
		Tx:            tx,
		Name:          connection.Name + " " + file,
		AccountUuid:   uuid.MustParse(accountUuid),
		Type:          model.DataDestinationType_S3FileCSV,
		Configuration: pqtype.NullRawMessage{RawMessage: dstConfiguration, Valid: dstConfiguration != nil},
		ConnectionId:  request.ConnectionId,
	})
	if err != nil {
		logger.Error(err, "cannot create data destination")
		tx.Rollback()
		return err
	}

	destSegmentMap, err := b.repository.DestSegmentMapRepository.CreateDestinationSegmentMap(ctx, &repository.CreateDestinationSegmentMapParams{
		Tx:            tx,
		SegmentId:     request.SegmentId,
		DestinationId: destination.ID,
		//MappingOptions: pqtype.NullRawMessage{},
	})
	if err != nil {
		logger.Error(err, "cannot create destination segment map")
		tx.Rollback()
		return err
	}

	payload := airflow.ExportFileRequestConfig{
		DagId:     dagId,
		Key:       audiencePathKey,
		Condition: segment.SqlCondition,
		S3Configurations: airflow.S3Configurations{
			AccessKeyId:     s3Configuration.AccessKeyId,
			SecretAccessKey: s3Configuration.SecretAccessKey,
			BucketName:      s3Configuration.BucketName,
			Region:          s3Configuration.Region,
			Key:             dstKey,
		},
	}
	_, err = b.airflowAdapter.TriggerGenerateDagExportFile(ctx, &airflow.TriggerGenerateDagExportFileRequest{
		Config: payload,
	})
	if err != nil {
		logger.Error(err, "cannot trigger generate dag export file")
		tx.Rollback()
		return err
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err, "cannot marshal trigger generate dag payload")
		tx.Rollback()
		return err
	}
	dataAction, err := b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		Tx:          tx,
		TargetTable: model.TargetTable_DestSegmentMap,
		ActionType:  model.ActionType_ExportDataToS3CSV,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    destSegmentMap.ID,
		Payload:     pqtype.NullRawMessage{RawMessage: jsonPayload, Valid: jsonPayload != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		tx.Rollback()
		return err
	}

	_, err = b.repository.DataActionRunRepository.CreateDataActionRun(ctx, &repository.CreateDataActionRunParams{
		Tx:          tx,
		ActionId:    dataAction.ID,
		RunId:       0,
		Status:      model.DataActionRunStatus_Creating,
		AccountUuid: uuid.MustParse(accountUuid),
	})
	if err != nil {
		logger.Error(err, "cannot create data action run")
		tx.Rollback()
		return err
	}

	tx.Commit()
	// END TRANSACTION

	return nil
}

func (b business) ExportMasterSegmentToS3File(ctx context.Context, request *api.ExportDataToFileRequest, accountUuid string) error {
	logger := b.log.WithName("ExportSegmentToS3File").WithValues("request", request)

	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.ConnectionId)
	if err != nil {
		logger.Error(err, "cannot get connection")
		return err
	}
	var s3Configuration model.S3Configurations
	err = json.Unmarshal(connection.Configurations.RawMessage, &s3Configuration)
	if err != nil {
		logger.Error(err, "cannot unmarshal s3 configuration")
		return err
	}

	dstConfiguration, err := json.Marshal(model.S3FileDestinationConfiguration{
		FileName: request.FileName,
		FileType: request.FileType,
		FilePath: request.FilePath,
	})
	if err != nil {
		logger.Error(err, "cannot marshal destination configuration")
		return err
	}

	file := request.FileName + "." + strings.ToLower(request.FileType)
	dstKey := request.FilePath + "/" + file

	dagId := utils.GenerateDagId(accountUuid, model.ActionType_ExportDataToS3CSV)
	audiencePathKey := utils.GenerateDeltaAudiencePath(request.MasterSegmentId)

	audience, err := b.repository.SegmentRepository.GetAudienceTable(ctx, repository.GetAudienceTableParams{MasterSegmentId: request.MasterSegmentId})
	if err != nil {
		logger.Error(err, "cannot get audience")
		return err
	}
	var schema []*api.SchemaColumn
	if audience.Schema.Valid && audience.Schema.RawMessage != nil {
		err = json.Unmarshal(audience.Schema.RawMessage, &schema)
		if err != nil {
			logger.Error(err, "cannot unmarshal audience schema")
			return err
		}
	}
	mappingOptions := utils.Map(schema, func(col *api.SchemaColumn) *api.MappingOptionItem {
		return &api.MappingOptionItem{
			SourceFieldName:      col.ColumnName,
			DestinationFieldName: col.ColumnName,
		}
	})
	jsonMappingOptions, err := json.Marshal(mappingOptions)
	if err != nil {
		logger.Error(err, "cannot marshal mapping options")
		return err
	}

	// BEGIN TRANSACTION
	tx := b.db.Begin()

	destination, err := b.repository.DataDestinationRepository.CreateDataDestination(ctx, &repository.CreateDataDestinationParams{
		Tx:            tx,
		Name:          connection.Name + " " + file,
		AccountUuid:   uuid.MustParse(accountUuid),
		Type:          model.DataDestinationType_S3FileCSV,
		Configuration: pqtype.NullRawMessage{RawMessage: dstConfiguration, Valid: dstConfiguration != nil},
		ConnectionId:  request.ConnectionId,
	})
	if err != nil {
		logger.Error(err, "cannot create data destination")
		tx.Rollback()
		return err
	}

	destMsSegmentMap, err := b.repository.DestMasterSegmentMapRepository.CreateDestinationMasterSegmentMap(ctx, &repository.CreateDestinationMasterSegmentMapParams{
		Tx:              tx,
		MasterSegmentId: request.MasterSegmentId,
		DestinationId:   destination.ID,
		MappingOptions:  pqtype.NullRawMessage{RawMessage: jsonMappingOptions, Valid: jsonMappingOptions != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create destination segment map")
		tx.Rollback()
		return err
	}

	payload := airflow.ExportFileRequestConfig{
		DagId:     dagId,
		Key:       audiencePathKey,
		Condition: "",
		S3Configurations: airflow.S3Configurations{
			AccessKeyId:     s3Configuration.AccessKeyId,
			SecretAccessKey: s3Configuration.SecretAccessKey,
			BucketName:      s3Configuration.BucketName,
			Region:          s3Configuration.Region,
			Key:             dstKey,
		},
	}
	_, err = b.airflowAdapter.TriggerGenerateDagExportFile(ctx, &airflow.TriggerGenerateDagExportFileRequest{
		Config: payload,
	})
	if err != nil {
		logger.Error(err, "cannot trigger generate dag export file")
		tx.Rollback()
		return err
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err, "cannot marshal trigger generate dag payload")
		tx.Rollback()
		return err
	}
	dataAction, err := b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		Tx:          tx,
		TargetTable: model.TargetTable_DestMasterSegmentMap,
		ActionType:  model.ActionType_ExportDataToS3CSV,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    destMsSegmentMap.ID,
		Payload:     pqtype.NullRawMessage{RawMessage: jsonPayload, Valid: jsonPayload != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		tx.Rollback()
		return err
	}

	_, err = b.repository.DataActionRunRepository.CreateDataActionRun(ctx, &repository.CreateDataActionRunParams{
		Tx:          tx,
		ActionId:    dataAction.ID,
		RunId:       0,
		Status:      model.DataActionRunStatus_Creating,
		AccountUuid: uuid.MustParse(accountUuid),
	})
	if err != nil {
		logger.Error(err, "cannot create data action run")
		tx.Rollback()
		return err
	}

	tx.Commit()
	// END TRANSACTION

	return nil
}
