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
	var schema []*api.SchemaColumn
	if dataTable.Schema.Valid && dataTable.Schema.RawMessage != nil {
		err := json.Unmarshal(dataTable.Schema.RawMessage, &schema)
		if err != nil {
			logger.Error(err, "cannot unmarshal schema")
			schema = []*api.SchemaColumn{}
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

	mySQLDesConfig := model.MySQLDestinationConfiguration{
		TableName: request.DestinationTableName,
	}
	jsonMySQLDesConfig, err := json.Marshal(mySQLDesConfig)
	if err != nil {
		logger.Error(err, "cannot marshal config ")
	}

	tx := b.db.Begin()
	// 1. Create destination
	destination, err := b.repository.DataDestinationRepository.CreateDataDestination(ctx, &repository.CreateDataDestinationParams{
		Name:          "MySQL - " + request.DestinationTableName,
		AccountUuid:   uuid.MustParse(accountUuid),
		Type:          model.DataDestinationType_MYSQL,
		Configuration: pqtype.NullRawMessage{RawMessage: jsonMySQLDesConfig, Valid: jsonMySQLDesConfig != nil},
		ConnectionId:  request.ConnectionId,
	})
	if err != nil {
		logger.Error(err, "cannot create data destination")
		tx.Rollback()
		return err
	}

	// 2. Create dest table map
	dstTableMap, err := b.repository.DestTableMapRepository.CreateDestinationTableMap(ctx, &repository.CreateDestinationTableMapParams{
		Tx:             tx,
		TableId:        request.DataTableId,
		DestinationId:  destination.ID,
		MappingOptions: pqtype.NullRawMessage{RawMessage: jsonMappingOptions, Valid: jsonMappingOptions != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create destination table mapping")
		tx.Rollback()
		return err
	}

	// 3. Create data action
	dataAction, err := b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		Tx:          tx,
		TargetTable: model.TargetTable_DestTableMap,
		ActionType:  model.ActionType_ExportToMySQL,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    dstTableMap.ID,
		Payload:     pqtype.NullRawMessage{RawMessage: jsonPayload, Valid: jsonPayload != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		tx.Rollback()
		return err
	}

	// 4. Create data action run
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

	mySQLDesConfig := model.MySQLDestinationConfiguration{
		TableName: request.DestinationTableName,
	}
	jsonMySQLDesConfig, err := json.Marshal(mySQLDesConfig)
	if err != nil {
		logger.Error(err, "cannot marshal config ")
	}

	tx := b.db.Begin()
	// 1. Create destination
	destination, err := b.repository.DataDestinationRepository.CreateDataDestination(ctx, &repository.CreateDataDestinationParams{
		Name:          "MySQL - " + request.DestinationTableName,
		AccountUuid:   uuid.MustParse(accountUuid),
		Type:          model.DataDestinationType_MYSQL,
		Configuration: pqtype.NullRawMessage{RawMessage: jsonMySQLDesConfig, Valid: jsonMySQLDesConfig != nil},
		ConnectionId:  request.ConnectionId,
	})
	if err != nil {
		logger.Error(err, "cannot create data destination")
		tx.Rollback()
		return err
	}

	// 2. Create dest ms segment map
	dstMsSegmentMap, err := b.repository.DestMasterSegmentMapRepository.CreateDestinationMasterSegmentMap(ctx, &repository.CreateDestinationMasterSegmentMapParams{
		Tx:              tx,
		MasterSegmentId: request.MasterSegmentId,
		DestinationId:   destination.ID,
		//MappingOptions:  pqtype.NullRawMessage{},
	})
	if err != nil {
		logger.Error(err, "cannot create destination master segment mapping")
		tx.Rollback()
		return err
	}

	// 3. Create data action
	dataAction, err := b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		Tx:          tx,
		TargetTable: model.TargetTable_DestMasterSegmentMap,
		ActionType:  model.ActionType_ExportToMySQL,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    dstMsSegmentMap.ID,
		Payload:     pqtype.NullRawMessage{RawMessage: jsonPayload, Valid: jsonPayload != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		tx.Rollback()
		return err
	}

	// 4. Create data action run
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

	segment, err := b.repository.SegmentRepository.GetSegment(ctx, request.SegmentId)
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

	mySQLDesConfig := model.MySQLDestinationConfiguration{
		TableName: request.DestinationTableName,
	}
	jsonMySQLDesConfig, err := json.Marshal(mySQLDesConfig)
	if err != nil {
		logger.Error(err, "cannot marshal config ")
	}

	tx := b.db.Begin()
	// 1. Create destination
	destination, err := b.repository.DataDestinationRepository.CreateDataDestination(ctx, &repository.CreateDataDestinationParams{
		Name:          "MySQL - " + request.DestinationTableName,
		AccountUuid:   uuid.MustParse(accountUuid),
		Type:          model.DataDestinationType_MYSQL,
		Configuration: pqtype.NullRawMessage{RawMessage: jsonMySQLDesConfig, Valid: jsonMySQLDesConfig != nil},
		ConnectionId:  request.ConnectionId,
	})
	if err != nil {
		logger.Error(err, "cannot create data destination")
		tx.Rollback()
		return err
	}

	// 2. Create dest segment map
	dstSegmentMap, err := b.repository.DestSegmentMapRepository.CreateDestinationSegmentMap(ctx, &repository.CreateDestinationSegmentMapParams{
		Tx:            tx,
		SegmentId:     request.SegmentId,
		DestinationId: destination.ID,
		//MappingOptions:  pqtype.NullRawMessage{},
	})
	if err != nil {
		logger.Error(err, "cannot create destination segment mapping")
		tx.Rollback()
		return err
	}

	// 3. Create data action
	dataAction, err := b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		Tx:          tx,
		TargetTable: model.TargetTable_DestSegmentMap,
		ActionType:  model.ActionType_ExportToMySQL,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    dstSegmentMap.ID,
		Payload:     pqtype.NullRawMessage{RawMessage: jsonPayload, Valid: jsonPayload != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		tx.Rollback()
		return err
	}

	// 4. Create data action run
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

	return nil
}
