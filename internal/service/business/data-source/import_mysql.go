package data_source

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

func (b business) ProcessImportFromMySQLSource(ctx context.Context, request *api.ImportFromMySQLSourceRequest, accountUuid uuid.UUID) error {
	logger := b.log.WithName("ProcessImportFromMySQLSource").WithValues("request", request)

	mySQLConnection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.ConnectionId)
	if err != nil {
		logger.Error(err, "cannot get connection")
		return err
	}
	var dbConfiguration model.MySQLConfiguration
	err = json.Unmarshal(mySQLConnection.Configurations.RawMessage, &dbConfiguration)
	if err != nil {
		logger.Error(err, "cannot parse db configuration")
		return err
	}

	config, err := json.Marshal(request)
	if err != nil {
		logger.Error(err, "cannot parse db configuration")
		return err
	}

	tx := b.db.Begin()

	dataSource, err := b.repository.DataSourceRepository.CreateDataSource(ctx, &repository.CreateDataSourceParams{
		Tx:             tx,
		Name:           request.Name,
		Description:    request.Description,
		Type:           model.DataSourceType_MySQL,
		Configurations: pqtype.NullRawMessage{RawMessage: config, Valid: config != nil},
		AccountUuid:    accountUuid,
		ConnectionId:   request.ConnectionId,
	})
	if err != nil {
		logger.Error(err, "cannot create data source")
		tx.Rollback()
		return err
	}

	var dataTable *model.DataTable
	if request.DeltaTableId > 0 {
		dataTable, err = b.repository.UpdateDataTable(ctx, &repository.UpdateDataTableParams{
			Tx:     tx,
			ID:     request.DeltaTableId,
			Status: model.DataTableStatus_UPDATING,
		})
		if err != nil {
			logger.Error(err, "cannot update data table")
			tx.Rollback()
			return err
		}
	} else {
		err := b.repository.DataTableRepository.CheckExistsDataTableName(ctx, request.DeltaTableName, accountUuid.String())
		if err != nil {
			logger.Error(err, "check table name exist")
			tx.Rollback()
			return err
		}
		dataTable, err = b.repository.DataTableRepository.CreateDataTable(ctx, &repository.CreateDataTableParams{
			Tx:          tx,
			Name:        request.DeltaTableName,
			Schema:      pqtype.NullRawMessage{},
			AccountUuid: accountUuid,
		})
		if err != nil {
			logger.Error(err, "cannot create data table")
			tx.Rollback()
			return err
		}
	}

	sourceTableMap, err := b.repository.SourceTableMapRepository.CreateSourceTableMap(ctx, &repository.CreateSourceTableMapParams{
		Tx:       tx,
		TableId:  dataTable.ID,
		SourceId: dataSource.ID,
	})
	if err != nil {
		logger.Error(err, "cannot create source table map")
		tx.Rollback()
		return err
	}

	dagId := utils.GenerateDagId(accountUuid.String(), model.ActionType_ImportDataFromMySQL)

	_, err = b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		TargetTable: model.TargetTable_SourceTableMap,
		ActionType:  model.ActionType_ImportDataFromMySQL,
		Schedule:    "",
		AccountUuid: accountUuid,
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    sourceTableMap.ID,
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		tx.Rollback()
		return err
	}

	_, err = b.airflowAdapter.TriggerGenerateDagImportMySQL(ctx, &airflow.TriggerGenerateDagImportMySQLRequest{
		Conf: airflow.DagImportMySQLConfig{
			DagId:          dagId,
			AccountUuid:    accountUuid.String(),
			DeltaTableName: request.DeltaTableName,
			Headers:        request.Headers,
			DatabaseConfiguration: airflow.DagImportMySQLDatabaseConfiguration{
				Host:     dbConfiguration.Host,
				Port:     dbConfiguration.Port,
				Database: dbConfiguration.Database,
				User:     dbConfiguration.User,
				Password: dbConfiguration.Password,
				Table:    request.SourceTableName,
			},
			WriteMode: airflow.DeltaWriteMode(request.WriteMode),
		},
	})
	if err != nil {
		logger.Error(err, "cannot trigger generate dag import mysql in airflow")
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
