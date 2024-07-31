package data_action

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (b business) ProcessNewDataActionRun(ctx context.Context, request *api.TriggerDataActionRunRequest, accountUuid string) error {
	logger := b.log.WithName("ProcessNewDataActionRun").WithValues("request", request)
	dataAction, err := b.repository.DataActionRepository.GetDataAction(ctx, request.Id)
	if err != nil {
		logger.Error(err, "cannot get data action")
		return err
	}
	if dataAction.AccountUuid.String() != accountUuid {
		return status.Error(codes.PermissionDenied, "user not have access to this data action")
	}

	err = b.db.Transaction(func(tx *gorm.DB) error {
		// trigger airflow
		dagRun, txErr := b.airflowAdapter.TriggerNewDagRun(ctx, dataAction.DagId, &airflow.TriggerNewDagRunRequest{})
		if txErr != nil {
			logger.Error(err, "cannot trigger new dag run")
			return txErr
		}

		// create new data action run
		_, txErr = b.repository.DataActionRunRepository.CreateDataActionRun(ctx, &repository.CreateDataActionRunParams{
			Tx:          tx,
			ActionId:    dataAction.ID,
			RunId:       dataAction.RunCount + 1,
			DagRunId:    dagRun.DagRunId,
			Status:      model.DataActionRunStatus_Processing,
			AccountUuid: uuid.MustParse(accountUuid),
		})
		if txErr != nil {
			logger.Error(err, "cannot create new data action run")
			return txErr
		}

		// increase data action run count by 1
		txErr = b.repository.DataActionRepository.UpdateDataAction(ctx, &repository.UpdateDataActionParams{
			Tx:       tx,
			ID:       dataAction.ID,
			RunCount: dataAction.RunCount + 1,
		})
		if txErr != nil {
			logger.Error(err, "cannot update data action run count")
			return txErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (b business) ProcessGetListDataActionRuns(ctx context.Context, request *api.GetListDataActionRunsRequest, accountUuid string) (*api.GetListDataActionRunsResponse, error) {
	queryDataActionRuns, err := b.repository.DataActionRunRepository.GetListDataActionRuns(ctx, &repository.GetListDataActionRunsParams{
		ActionTypes: request.Types,
		AccountUuid: uuid.MustParse(accountUuid),
		Page:        int(request.Page),
		PageSize:    int(request.PageSize),
	})
	if err != nil {
		b.log.WithName("list data actions").Error(err, "cannot get data actions")
		return nil, err
	}
	var returnDataActionRuns []*api.DataActionRun
	for _, dataActionRun := range queryDataActionRuns.DataActionRuns {
		returnDataActionRuns = append(returnDataActionRuns, &api.DataActionRun{
			Id:         dataActionRun.ID,
			ActionId:   dataActionRun.ActionId,
			ActionType: string(dataActionRun.ActionType),
			Status:     string(dataActionRun.Status),
			CreatedAt:  dataActionRun.CreatedAt.String(),
			UpdatedAt:  dataActionRun.UpdatedAt.String(),
		})
	}

	return &api.GetListDataActionRunsResponse{
		Code:    0,
		Message: "Success",
		Count:   queryDataActionRuns.Count,
		Results: returnDataActionRuns,
	}, nil

}

func (b business) ProcessGetTotalRunsPerDay(ctx context.Context, accountUuid string) (*api.GetDataActionRunsPerDayResponse, error) {
	results, err := b.repository.DataActionRunRepository.GetTotalRunsPerDay(ctx, &repository.GetTotalRunsPerDayParams{AccountUuid: accountUuid})
	if err != nil {
		b.log.WithName("ProcessGetTotalRunsPerDay").Error(err, "cannot get data action runs per day", "uuid", accountUuid)
		return nil, err
	}

	runsPerDay := utils.Map(results, func(data repository.TotalRunsPerDay) *api.GetDataActionRunsPerDayResponse_TotalActionRunsPerDay {
		return &api.GetDataActionRunsPerDayResponse_TotalActionRunsPerDay{
			Date:  data.Date.String(),
			Total: int32(data.Total),
		}
	})

	return &api.GetDataActionRunsPerDayResponse{
		Code:    0,
		Message: "Success",
		Result:  runsPerDay,
	}, nil
}

func (b business) ProcessGetDataRunsProportion(ctx context.Context, accountUuid string) (*api.GetDataRunsProportionResponse, error) {
	results, err := b.repository.DataActionRunRepository.GetTotalRunsPerType(ctx, &repository.GetTotalRunsPerTypeParams{AccountUuid: accountUuid})
	if err != nil {
		b.log.WithName("ProcessGetTotalRunsPerDay").Error(err, "cannot get data action runs per day", "uuid", accountUuid)
		return nil, err
	}

	typeCount := map[string]int32{
		"segmentation": 0,
		"source":       0,
		"destination":  0,
	}
	for _, each := range results {
		switch each.Type {
		case model.ActionType_ImportDataFromMySQL, model.ActionType_ImportDataFromFile, model.ActionType_ImportDataFromS3:
			typeCount["source"] = typeCount["source"] + int32(each.Total)
		case model.ActionType_ExportDataToS3CSV, model.ActionType_ExportToMySQL, model.ActionType_ExportGophish:
			typeCount["destination"] = typeCount["destination"] + int32(each.Total)
		case model.ActionType_CreateSegment, model.ActionType_CreateMasterSegment, model.ActionType_TrainPredictModel, model.ActionType_ApplyPredictModel:
			typeCount["segmentation"] = typeCount["segmentation"] + int32(each.Total)
		default:
		}
	}

	return &api.GetDataRunsProportionResponse{
		Code:    0,
		Message: "Success",
		Result: []*api.GetDataRunsProportionResponse_CategoryCount{
			{
				Category: "Source",
				Count:    typeCount["source"],
			},
			{
				Category: "Destination",
				Count:    typeCount["destination"],
			},
			{
				Category: "Segmentation",
				Count:    typeCount["segmentation"],
			},
		},
	}, nil
}
