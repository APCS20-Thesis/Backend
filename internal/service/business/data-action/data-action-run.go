package data_action

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
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
