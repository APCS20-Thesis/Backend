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
)

func (b business) ProcessExportToMySQLDestination(ctx context.Context, request *api.ExportToMySQLDestinationRequest, accountUuid string) error {
	var err error

	if request.DataTableId != 0 {
		err = b.ExportDataTableToMySQL(ctx, request, accountUuid)
	} else if request.MasterSegmentId != 0 {
		err = b.ExportMasterSegmentAudienceToMySQL(ctx, request, accountUuid)
	} else {
		err = b.ExportSegmentToMySQL(ctx, request, accountUuid)
	}
	if err != nil {
		return err
	}

	return nil
}

func (b business) ExportDataTableToMySQL(ctx context.Context, request *api.ExportToMySQLDestinationRequest, accountUuid string) error {
	logger := b.log.WithName("ExportDataTableToMySQL").WithValues("request", request)

	dagId := utils.GenerateDagId(accountUuid, model.ActionType_ExportToMySQL)

	dataTable, err := b.repository.DataTableRepository.GetDataTable(ctx, request.DataTableId)
	if err != nil {
		logger.Error(err, "cannot get data table")
		return err
	}

	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.ConnectionId)
	if err != nil {
		logger.Error(err, "cannot get connection")
		return err
	}
	var mySQLConnection model.MySQLConfiguration
	err = json.Unmarshal(connection.Configurations.RawMessage, &mySQLConnection)
	if err != nil {
		logger.Error(err, "cannot unmarshal connection configuration")
		return err
	}

	payload := &airflow.TriggerGenerateDagExportMySQLRequest{
		Conf: airflow.DagExportMySQLConfig{
			DagId:           dagId,
			AccountUuid:     accountUuid,
			DeltaTableName:  dataTable.Name,
			MasterSegmentId: 0,
			Condition:       "",
			DatabaseConfiguration: airflow.DagExportMySQLDatabaseConfiguration{
				Host:     mySQLConnection.Host,
				Port:     mySQLConnection.Port,
				User:     mySQLConnection.User,
				Password: mySQLConnection.Password,
				Database: mySQLConnection.Database,
			},
			DestinationTableName: request.DestinationTableName,
		},
	}
	_, err = b.airflowAdapter.TriggerGenerateDagExportMySQL(ctx, payload)
	if err != nil {
		logger.Error(err, "cannot trigger generate dag export mysql")
		return err
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err, "cannot marshal payload")
		return err
	}
	_, err = b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		TargetTable: model.TargetTable_DataTable,
		ActionType:  model.ActionType_ExportToMySQL,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    request.DataTableId,
		Payload:     pqtype.NullRawMessage{RawMessage: jsonPayload, Valid: jsonPayload != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		return err
	}

	config, err := json.Marshal(request)
	_, err = b.repository.DataDestinationRepository.CreateDataDestination(ctx, &repository.CreateDataDestinationParams{
		Name:          "MySQL - " + request.DestinationTableName,
		AccountUuid:   uuid.MustParse(accountUuid),
		Type:          model.DataDestinationType_MYSQL,
		Configuration: pqtype.NullRawMessage{RawMessage: config, Valid: config != nil},
		ConnectionId:  request.ConnectionId,
	})
	if err != nil {
		logger.Error(err, "cannot create data destination")
		return err
	}

	return nil
}

func (b business) ExportMasterSegmentAudienceToMySQL(ctx context.Context, request *api.ExportToMySQLDestinationRequest, accountUuid string) error {
	logger := b.log.WithName("ExportDataTableToMySQL").WithValues("request", request)

	dagId := utils.GenerateDagId(accountUuid, model.ActionType_ExportToMySQL)

	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.ConnectionId)
	if err != nil {
		logger.Error(err, "cannot get connection")
		return err
	}
	var mySQLConnection model.MySQLConfiguration
	err = json.Unmarshal(connection.Configurations.RawMessage, &mySQLConnection)
	if err != nil {
		logger.Error(err, "cannot unmarshal connection configuration")
		return err
	}

	payload := &airflow.TriggerGenerateDagExportMySQLRequest{
		Conf: airflow.DagExportMySQLConfig{
			DagId:           dagId,
			AccountUuid:     accountUuid,
			DeltaTableName:  "",
			MasterSegmentId: request.MasterSegmentId,
			Condition:       "",
			DatabaseConfiguration: airflow.DagExportMySQLDatabaseConfiguration{
				Host:     mySQLConnection.Host,
				Port:     mySQLConnection.Port,
				User:     mySQLConnection.User,
				Password: mySQLConnection.Password,
				Database: mySQLConnection.Database,
			},
			DestinationTableName: request.DestinationTableName,
		},
	}
	_, err = b.airflowAdapter.TriggerGenerateDagExportMySQL(ctx, payload)
	if err != nil {
		logger.Error(err, "cannot trigger generate dag export mysql")
		return err
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err, "cannot marshal payload")
		return err
	}
	_, err = b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		TargetTable: model.TargetTable_DataTable,
		ActionType:  model.ActionType_ExportToMySQL,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    request.DataTableId,
		Payload:     pqtype.NullRawMessage{RawMessage: jsonPayload, Valid: jsonPayload != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		return err
	}

	config, err := json.Marshal(request)
	_, err = b.repository.DataDestinationRepository.CreateDataDestination(ctx, &repository.CreateDataDestinationParams{
		Name:          "MySQL - " + request.DestinationTableName,
		AccountUuid:   uuid.MustParse(accountUuid),
		Type:          model.DataDestinationType_MYSQL,
		Configuration: pqtype.NullRawMessage{RawMessage: config, Valid: config != nil},
		ConnectionId:  request.ConnectionId,
	})
	if err != nil {
		logger.Error(err, "cannot create data destination")
		return err
	}

	return nil
}

func (b business) ExportSegmentToMySQL(ctx context.Context, request *api.ExportToMySQLDestinationRequest, accountUuid string) error {
	logger := b.log.WithName("ExportSegmentToMySQL").WithValues("request", request)

	dagId := utils.GenerateDagId(accountUuid, model.ActionType_ExportToMySQL)

	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.ConnectionId)
	if err != nil {
		logger.Error(err, "cannot get connection")
		return err
	}
	var mySQLConnection model.MySQLConfiguration
	err = json.Unmarshal(connection.Configurations.RawMessage, &mySQLConnection)
	if err != nil {
		logger.Error(err, "cannot unmarshal connection configuration")
		return err
	}

	segment, err := b.repository.SegmentRepository.GetSegment(ctx, request.SegmentId, accountUuid)
	if err != nil {
		logger.Error(err, "cannot get segment")
		return err
	}

	payload := &airflow.TriggerGenerateDagExportMySQLRequest{
		Conf: airflow.DagExportMySQLConfig{
			DagId:           dagId,
			AccountUuid:     accountUuid,
			DeltaTableName:  "",
			MasterSegmentId: segment.MasterSegmentId,
			Condition:       segment.SqlCondition,
			DatabaseConfiguration: airflow.DagExportMySQLDatabaseConfiguration{
				Host:     mySQLConnection.Host,
				Port:     mySQLConnection.Port,
				User:     mySQLConnection.User,
				Password: mySQLConnection.Password,
				Database: mySQLConnection.Database,
			},
			DestinationTableName: request.DestinationTableName,
		},
	}
	_, err = b.airflowAdapter.TriggerGenerateDagExportMySQL(ctx, payload)
	if err != nil {
		logger.Error(err, "cannot trigger generate dag export mysql")
		return err
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err, "cannot marshal payload")
		return err
	}
	_, err = b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		TargetTable: model.TargetTable_DataTable,
		ActionType:  model.ActionType_ExportToMySQL,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    request.DataTableId,
		Payload:     pqtype.NullRawMessage{RawMessage: jsonPayload, Valid: jsonPayload != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		return err
	}

	config, err := json.Marshal(request)
	_, err = b.repository.DataDestinationRepository.CreateDataDestination(ctx, &repository.CreateDataDestinationParams{
		Name:          "MySQL - " + request.DestinationTableName,
		AccountUuid:   uuid.MustParse(accountUuid),
		Type:          model.DataDestinationType_MYSQL,
		Configuration: pqtype.NullRawMessage{RawMessage: config, Valid: config != nil},
		ConnectionId:  request.ConnectionId,
	})
	if err != nil {
		logger.Error(err, "cannot create data destination")
		return err
	}

	return nil
}
