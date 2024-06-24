package job

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/sqlc-dev/pqtype"
)

func (j *job) SyncDagRunStatus(ctx context.Context) {
	jobLog := j.logger.WithName("SyncDagRunStatus")
	// get data action run that are in status PROCESSING
	dataActionRuns, err := j.repository.DataActionRunRepository.GetListDataActionRunsWithExtraInfo(ctx, &repository.GetListDataActionRunsWithExtraInfoParams{
		Statuses: []model.DataActionRunStatus{model.DataActionRunStatus_Processing},
	})
	if err != nil {
		jobLog.Error(err, "fail to get list data action_runs")
		return
	}
	jobLog.Info("dataActionRuns", "data", dataActionRuns)

	for _, dataActionRun := range dataActionRuns {
		if dataActionRun.DagRunId == "" {
			err = j.repository.DataActionRunRepository.UpdateDataActionRunStatus(ctx, dataActionRun.ID, model.DataActionRunStatus_Failed)
			continue
		}

		response, err := j.airflowAdapter.GetDagRun(ctx, dataActionRun.DagId, dataActionRun.DagRunId)
		if err != nil {
			jobLog.Error(err, "cannot get dag run from airflow", "dagId", dataActionRun.DagId, "dagRunId", dataActionRun.DagRunId)
			continue
		}

		switch response.State {
		case "success":
			err = j.repository.DataActionRunRepository.UpdateDataActionRunStatus(ctx, dataActionRun.ID, model.DataActionRunStatus_Success)
			if err != nil {
				jobLog.Error(err, "cannot update data action run status", "data action run id", dataActionRun.ID)
				break
			}
			err = j.repository.DataActionRepository.UpdateDataAction(ctx, &repository.UpdateDataActionParams{
				ID:     dataActionRun.ActionId,
				Status: model.DataActionStatus_Success,
			})
			// additional handling
			go j.SyncRelatedStatusFromDataActionRunStatus(ctx, dataActionRun.ActionId)
		case "failed":
			err = j.repository.DataActionRunRepository.UpdateDataActionRunStatus(ctx, dataActionRun.ID, model.DataActionRunStatus_Failed)
			if err != nil {
				jobLog.Error(err, "cannot update data action run status", "data action run id", dataActionRun.ID)
				break
			}
			err = j.repository.DataActionRepository.UpdateDataAction(ctx, &repository.UpdateDataActionParams{
				ID:     dataActionRun.ActionId,
				Status: model.DataActionStatus_Failed,
			})
		default:
		}
		if err != nil {
			jobLog.Error(err, "cannot get dag run from airflow", "dagId", dataActionRun.DagId, "dagRunId", dataActionRun.DagRunId)
			continue
		}

	}

	return
}

func (j *job) SyncRelatedStatusFromDataActionRunStatus(ctx context.Context, dataActionId int64) {
	dataAction, err := j.repository.DataActionRepository.GetDataAction(ctx, dataActionId)
	if err != nil {
		j.logger.WithName("job:SyncRelatedStatusFromDataActionStatus").Error(err, "cannot get data action ")
	}

	switch dataAction.ActionType {
	case model.ActionType_ImportDataFromFile, model.ActionType_ImportDataFromMySQL:
		sourceTableMap, err := j.repository.SourceTableMapRepository.GetSourceTableMapById(ctx, dataAction.ObjectId)
		if err != nil {
			j.logger.WithName("job:SyncRelatedStatusFromDataActionStatus").Error(err, "cannot get source table map", "sourceTableMapId", dataAction.ObjectId)
			return
		}
		path, err := j.repository.DataTableRepository.GetDataTableDeltaPath(ctx, sourceTableMap.TableId)
		if err != nil {
			j.logger.WithName("job:SyncRelatedStatusFromDataActionStatus").Error(err, "cannot get table delta path", "TableId", sourceTableMap.TableId)
			return
		}
		response, err := j.queryAdapter.GetSchemaTable(ctx, &query.GetSchemaDataTableRequest{TablePath: path})
		if err != nil {
			j.logger.WithName("GetSchemaTable").Info("error here")
			err = j.repository.DataTableRepository.UpdateStatusDataTable(ctx, sourceTableMap.TableId, model.DataTableStatus_NEED_TO_SYNC)
			if err != nil {
				j.logger.WithName("job:SyncRelatedStatusFromDataActionStatus").Error(err, "cannot update data table status", "TableId", sourceTableMap.TableId)
				return
			}
		}
		schema := make([]model.SchemaUnit, 0, len(response.Schema))
		for _, column := range response.Schema {
			schema = append(schema, model.SchemaUnit{
				ColumnName: column.Name,
				DataType:   column.Type,
			})
		}
		jsonSchema, err := json.Marshal(schema)
		if err != nil {
			j.logger.WithName("job:SyncRelatedStatusFromDataActionStatus").Error(err, "cannot parse schema")
			return
		}
		err = j.repository.DataTableRepository.UpdateDataTable(ctx, &repository.UpdateDataTableParams{
			ID:     sourceTableMap.TableId,
			Schema: pqtype.NullRawMessage{Valid: true, RawMessage: jsonSchema},
			Status: model.DataTableStatus_UP_TO_DATE,
		})
		if err != nil {
			j.logger.WithName("job:SyncRelatedStatusFromDataActionStatus").Error(err, "cannot update data table", "TableId", sourceTableMap.TableId)
			return
		}
		j.logger.WithName("GetSchemaTable").Info("update table here", "schema", schema)
	case model.ActionType_CreateMasterSegment:

		err = j.SyncOnCreateMasterSegment(ctx, dataAction.ObjectId)
		if err != nil {
			return
		}
	default:
	}

	return
}

func (j *job) SyncOnCreateMasterSegment(ctx context.Context, masterSegmentId int64) error {
	audienceTable, err := j.repository.SegmentRepository.GetAudienceTable(ctx, repository.GetAudienceTableParams{MasterSegmentId: masterSegmentId})
	if err != nil {
		j.logger.WithName("job:SyncOnCreateMasterSegment").Error(err, "cannot get audience table", "masterSegmentId", masterSegmentId)
		return err
	}
	behaviorTables, err := j.repository.SegmentRepository.ListBehaviorTables(ctx, repository.ListBehaviorTablesParams{MasterSegmentId: masterSegmentId})
	if err != nil {
		j.logger.WithName("job:SyncOnCreateMasterSegment").Error(err, "cannot get behavior tables", "masterSegmentId", masterSegmentId)
		return err
	}
	// Sync audience table schema
	response, err := j.queryAdapter.GetSchemaTable(ctx, &query.GetSchemaDataTableRequest{
		TablePath: utils.GenerateDeltaAudiencePath(masterSegmentId),
	})
	if err != nil {
		j.logger.WithName("job:SyncOnCreateMasterSegment").Error(err, "cannot query audience table schema")
		return err
	}
	schema := utils.Map(response.Schema, func(unit query.FieldSchema) model.SchemaUnit {
		return model.SchemaUnit{
			ColumnName: unit.Name,
			DataType:   unit.Type,
		}
	})
	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		j.logger.WithName("job:SyncOnCreateMasterSegment").Error(err, "cannot parse schema")
		return err
	}
	err = j.repository.SegmentRepository.UpdateAudienceTable(ctx, &repository.UpdateAudienceTableParams{
		Id:     audienceTable.ID,
		Schema: pqtype.NullRawMessage{RawMessage: jsonSchema, Valid: jsonSchema != nil},
	})
	if err != nil {
		j.logger.WithName("job:SyncOnCreateMasterSegment").Error(err, "cannot update audience table", "id", audienceTable.ID)
		return err
	}
	// Sync behavior schemas
	for _, behaviorTable := range behaviorTables {
		response, err := j.queryAdapter.GetSchemaTable(ctx, &query.GetSchemaDataTableRequest{
			TablePath: utils.GenerateDeltaBehaviorPath(masterSegmentId, behaviorTable.Name),
		})
		if err != nil {
			j.logger.WithName("job:SyncOnCreateMasterSegment").Error(err, "cannot query behavior table schema")
			return err
		}
		schema := utils.Map(response.Schema, func(unit query.FieldSchema) model.SchemaUnit {
			return model.SchemaUnit{
				ColumnName: unit.Name,
				DataType:   unit.Type,
			}
		})
		jsonSchema, err := json.Marshal(schema)
		if err != nil {
			j.logger.WithName("job:SyncOnCreateMasterSegment").Error(err, "cannot parse schema")
			return err
		}
		err = j.repository.SegmentRepository.UpdateBehaviorTable(ctx, &repository.UpdateBehaviorTableParams{
			Id:     behaviorTable.ID,
			Schema: pqtype.NullRawMessage{RawMessage: jsonSchema, Valid: jsonSchema != nil},
		})
		if err != nil {
			j.logger.WithName("job:SyncOnCreateMasterSegment").Error(err, "cannot update behavior table", "id", behaviorTable.ID)
			return err
		}
	}
	err = j.repository.SegmentRepository.UpdateMasterSegment(ctx, &repository.UpdateMasterSegmentParams{
		Id:     masterSegmentId,
		Status: model.MasterSegmentStatus_UP_TO_DATE,
	})
	if err != nil {
		j.logger.WithName("job:SyncOnCreateMasterSegment").Error(err, "cannot update master segment status", "masterSegmentId", masterSegmentId)
		return err
	}

	return nil
}
