package job

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/adapter/alert"
	"github.com/APCS20-Thesis/Backend/internal/adapter/mqtt"
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

			j.mqttAdapter.PublishNotification(dataActionRun.AccountUuid.String(), mqtt.Notification{
				Status:     409,
				ActionType: dataActionRun.ActionType,
				Severity:   mqtt.MqttSeverity_Error,
			})
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
			// additional handling
			go j.SyncRelatedStatusFromDataActionRunStatus(ctx, dataActionRun.ActionId)
		default:
		}
		if err != nil {
			jobLog.Error(err, "cannot get dag run from airflow", "dagId", dataActionRun.DagId, "dagRunId", dataActionRun.DagRunId)
			if j.config.AlertAdapterConfig.EnableAlert {
				_ = j.alertAdapter.AlertError(ctx, &alert.ErrorMessage{
					Title:        "Job error",
					Description:  "There is error in job sync dag run status",
					Request:      syncStatusInfo{DataActionRunId: dataActionRun.ID},
					ErrorMessage: err.Error(),
				})
			}
			continue
		}

	}

	return
}

type syncStatusInfo struct {
	DataActionRunId int64 `json:"data_action_run_id"`
}

func (j *job) SyncRelatedStatusFromDataActionRunStatus(ctx context.Context, dataActionId int64) {
	dataAction, err := j.repository.DataActionRepository.GetDataAction(ctx, dataActionId)
	if err != nil {
		j.logger.WithName("job:SyncRelatedStatusFromDataActionStatus").Error(err, "cannot get data action ")
	}

	switch dataAction.ActionType {
	case model.ActionType_ImportDataFromFile, model.ActionType_ImportDataFromMySQL:
		if dataAction.Status == model.DataActionStatus_Success {
			err = j.business.DataSourceBusiness.SyncOnImportFromSourceSuccess(ctx, dataActionId)
		} else if dataAction.Status == model.DataActionStatus_Failed {
			err = j.business.DataSourceBusiness.SynOnImportFromSourceFailed(ctx, dataActionId)
		}
		if err != nil {
			if j.config.AlertAdapterConfig.EnableAlert {
				_ = j.alertAdapter.AlertError(ctx, &alert.ErrorMessage{
					Title:        "Job error",
					Description:  "There is error in job sync dag run status",
					Request:      syncStatusInfo{DataActionRunId: dataActionId},
					ErrorMessage: err.Error(),
				})
			}
			return
		}
	case model.ActionType_CreateMasterSegment:
		err = j.business.SegmentBusiness.SyncOnCreateMasterSegment(ctx, dataAction.ObjectId, dataAction.Status)
		if err != nil {
			if j.config.AlertAdapterConfig.EnableAlert {
				_ = j.alertAdapter.AlertError(ctx, &alert.ErrorMessage{
					Title:        "Job error",
					Description:  "There is error in job sync dag run status",
					Request:      syncStatusInfo{DataActionRunId: dataActionId},
					ErrorMessage: err.Error(),
				})
			}
			return
		}
	case model.ActionType_TrainPredictModel:
		err = j.business.PredictModelBusiness.SyncOnTrainPredictModel(ctx, dataActionId)
		if err != nil {
			if j.config.AlertAdapterConfig.EnableAlert {
				_ = j.alertAdapter.AlertError(ctx, &alert.ErrorMessage{
					Title:        "Job error",
					Description:  "There is error in job sync dag run status",
					Request:      syncStatusInfo{DataActionRunId: dataActionId},
					ErrorMessage: err.Error(),
				})
			}
			return
		}
	case model.ActionType_CreateSegment:
		var status string
		switch dataAction.Status {
		case model.DataActionStatus_Success:
			status = "UP_TO_DATE"
		case model.DataActionStatus_Failed:
			status = "FAILED"
		}
		err = j.repository.SegmentRepository.UpdateSegment(ctx, &repository.UpdateSegmentParams{
			Tx:     nil,
			Id:     dataAction.ObjectId,
			Status: status,
		})
	default:
	}

	if dataAction.Status == model.DataActionStatus_Success {
		j.mqttAdapter.PublishNotification(dataAction.AccountUuid.String(), mqtt.Notification{
			Status:     200,
			ActionType: dataAction.ActionType,
			Severity:   mqtt.MqttSeverity_Success,
		})
	} else if dataAction.Status == model.DataActionStatus_Failed {
		j.mqttAdapter.PublishNotification(dataAction.AccountUuid.String(), mqtt.Notification{
			Status:     409,
			ActionType: dataAction.ActionType,
			Severity:   mqtt.MqttSeverity_Error,
		})
	}

	return
}
