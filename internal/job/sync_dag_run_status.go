package job

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
)

func (j *job) SyncDagRunStatus(ctx context.Context) {
	jobLog := j.logger.WithName("TriggerDagRun")
	// get data action that are in status PENDING
	dataActionRuns, err := j.repository.DataActionRunRepository.GetListDataActionRunsWithExtraInfo(ctx, &repository.GetListDataActionRunsWithExtraInfoParams{
		Statuses: []model.DataActionRunStatus{model.DataActionStatus_Pending},
	})
	if err != nil {
		jobLog.Error(err, "fail to get list data actions")
		return
	}
	jobLog.Info("dataActionRuns", "data", dataActionRuns)

	for _, dataActionRun := range dataActionRuns {
		response, err := j.airflowAdapter.GetDagRun(ctx, dataActionRun.DagId, dataActionRun.DagRunId)
		if err != nil {
			jobLog.Error(err, "cannot get dag run from airflow", "dagId", dataActionRun.DagId, "dagRunId", dataActionRun.DagRunId)
			continue
		}
		jobLog.Info("response", "resp", response)
		switch response.State {
		case "success":
			err = j.repository.DataActionRunRepository.UpdateDataActionRunStatus(ctx, dataActionRun.ID, model.DataActionRunStatus_Success)
		case "failed":
			err = j.repository.DataActionRunRepository.UpdateDataActionRunStatus(ctx, dataActionRun.ID, model.DataActionRunStatus_Failed)
		default:
		}
		if err != nil {
			jobLog.Error(err, "cannot get dag run from airflow", "dagId", dataActionRun.DagId, "dagRunId", dataActionRun.DagRunId)
		}
	}

	return
}
