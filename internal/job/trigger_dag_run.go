package job

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (j *job) TriggerDagRuns(ctx context.Context) {
	jobLog := j.logger.WithName("TriggerDagRun")

	// get data action that are in status PENDING
	queryDataActions, err := j.repository.DataActionRepository.GetListDataActions(ctx, &repository.GetListDataActionsParams{
		Statuses: []model.DataActionStatus{model.DataActionStatus_Pending},
	})
	if err != nil {
		jobLog.Error(err, "fail to get list data actions")
		return
	}
	jobLog.Info("dataActions", "data", queryDataActions.DataActions)

	for _, dataAction := range queryDataActions.DataActions {

		// - Gọi airflow update dag active
		_, err = j.airflowAdapter.UpdateDag(ctx, dataAction.DagId, &airflow.UpdateDagRequest{IsPaused: false})
		if err != nil && status.Code(err) == codes.Aborted && strings.Contains(err.Error(), "not found") {
			j.logger.WithName("TriggerDagRuns").Info("dag not created", "dataAction", dataAction.ID, "dagId", dataAction.DagId)
			continue
		}
		if err != nil {
			jobLog.Error(err, "cannot update Dag active in airflow")
			return
		}

		// nếu có Dag tương ứng thì
		// - Update DataAction status
		tx := j.db.Begin()
		// update data action status
		tx.WithContext(ctx).Table("data_action").Where("id = ?", dataAction.ID).Update("status", model.DataActionStatus_Processing)

		// - Gọi airflow trigger DagRun
		triggerDagRunResponse, err := j.airflowAdapter.TriggerNewDagRun(ctx, dataAction.DagId, &airflow.TriggerNewDagRunRequest{})
		if err != nil {
			jobLog.Error(err, "cannot trigger DagRun in airflow", "dagId", dataAction.ID)
			tx.Rollback()
			return
		}

		// - Tạo DataActionRun với thông tin DagRun
		tx.WithContext(ctx).Table("data_action_run").Create(&model.DataActionRun{
			ActionId:    dataAction.ID,
			RunId:       1,
			DagRunId:    triggerDagRunResponse.DagRunId,
			Status:      model.DataActionRunStatus_Processing,
			AccountUuid: dataAction.AccountUuid,
		})

		// - Xoá DataActionRun dummy
		tx.WithContext(ctx).Table("data_action_run").Where("action_id = ? and run_id = ?", dataAction.ID, 0).Delete(&model.DataActionRun{})

		tx.WithContext(ctx).Table("data_action").Where("id = ?", dataAction.ID).Update("run_count", 1)

		err = tx.Commit().Error
		if err != nil {
			j.logger.Error(err, "trigger dag run transaction err")
		}
	}

	return
}
