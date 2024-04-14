package job

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
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
		case "failed":
			err = j.repository.DataActionRunRepository.UpdateDataActionRunStatus(ctx, dataActionRun.ID, model.DataActionRunStatus_Failed)
		default:
		}
		if err != nil {
			jobLog.Error(err, "cannot get dag run from airflow", "dagId", dataActionRun.DagId, "dagRunId", dataActionRun.DagRunId)
			continue
		}

		// additional handling
		go j.SyncRelatedStatusFromDataActionRunStatus(ctx, dataActionRun.ActionId)
	}

	return
}

func (j *job) SyncRelatedStatusFromDataActionRunStatus(ctx context.Context, dataActionId int64) {
	dataAction, err := j.repository.DataActionRepository.GetDataAction(ctx, dataActionId)
	if err != nil {
		j.logger.WithName("job:SyncRelatedStatusFromDataActionStatus").Error(err, "cannot get data action ")
	}

	switch dataAction.ActionType {
	case model.ActionType_ImportDataFromFile:
		sourceTableMap, err := j.repository.SourceTableMapRepository.GetSourceTableMapById(ctx, dataAction.ObjectId)
		if err != nil {
			j.logger.WithName("job:SyncRelatedStatusFromDataActionStatus").Error(err, "cannot get source table map", "sourceTableMapId", dataAction.ObjectId)
			return
		}
		err = j.repository.DataTableRepository.UpdateStatusDataTable(ctx, sourceTableMap.TableId, model.DataTableStatus_NEED_TO_SYNC)
		if err != nil {
			j.logger.WithName("job:SyncRelatedStatusFromDataActionStatus").Error(err, "cannot update data table status", "TableId", sourceTableMap.TableId)
			return
		}
	default:
	}

	return
}
