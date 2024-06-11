package job

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
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

	default:
	}

	return
}
