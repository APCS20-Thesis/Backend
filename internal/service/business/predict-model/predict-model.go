package predict_model

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
	"gorm.io/gorm"
)

func (b business) ProcessTrainPredictModel(ctx context.Context, request *api.TrainPredictModelRequest, accountUuid string) error {
	logger := b.log.WithName("TrainPredictModel").WithValues("request", request)
	dagId := utils.GenerateDagId(accountUuid, model.ActionType_TrainPredictModel)
	trainConfig := model.PredictModelTrainConfiguration{
		Segment1:           request.TrainSegmentIds[0],
		Segment2:           request.TrainSegmentIds[1],
		Label1:             request.Labels[0],
		Label2:             request.Labels[1],
		SelectedAttributes: request.SelectedAttributes,
	}
	trainConfigJson, err := json.Marshal(trainConfig)
	if err != nil {
		return err
	}

	errTx := b.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create predict model
		predictModel, err := b.repository.CreatePredictModel(ctx, &repository.CreatePredictModelParams{
			Tx:                  tx,
			Name:                request.Name,
			MasterSegmentId:     request.MasterSegmentId,
			TrainConfigurations: pqtype.NullRawMessage{RawMessage: trainConfigJson, Valid: trainConfigJson != nil},
		})
		if err != nil {
			logger.Error(err, "cannot create predict model")
			return err
		}
		// 2. Build airflow payload
		payload := airflow.TriggerGenerateDagTrainPredictModelRequest{Conf: airflow.DagTrainPredictModelConfig{
			DagId:            dagId,
			Segment1Key:      utils.GenerateDeltaSegmentPath(request.MasterSegmentId, request.TrainSegmentIds[0]),
			Segment2Key:      utils.GenerateDeltaSegmentPath(request.MasterSegmentId, request.TrainSegmentIds[1]),
			Label1:           "0",
			Label2:           "1",
			PredictModelKey:  utils.GenerateDeltaPredictModelFilePath(request.MasterSegmentId, predictModel.ID),
			SelectAttributes: request.SelectedAttributes,
		}}
		payloadJson, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		// 2. Call airflow
		_, err = b.airflowAdapter.TriggerGenerateDagTrainPredictModel(ctx, &payload)
		if err != nil {
			logger.Error(err, "cannot trigger dag generate train predict model airflow")
			return err
		}
		// 3. Save data action
		_, err = b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
			TargetTable: model.TargetTable_PredictModel,
			ActionType:  model.ActionType_TrainPredictModel,
			Schedule:    "",
			AccountUuid: uuid.MustParse(accountUuid),
			DagId:       dagId,
			Status:      model.DataActionStatus_Pending,
			ObjectId:    predictModel.ID,
			Payload:     pqtype.NullRawMessage{RawMessage: payloadJson, Valid: payloadJson != nil},
		})
		if err != nil {
			logger.Error(err, "cannot create data action")
			return err
		}
		return nil
	})
	if errTx != nil {
		return errTx
	}

	return nil
}

func (b business) ProcessGetListPredictModels(ctx context.Context, request *api.GetListPredictModelsRequest, accountUuid string) (*api.GetListPredictModelsResponse, error) {
	queryResult, err := b.repository.PredictModelRepository.ListPredictModels(ctx, &repository.ListPredictModelsParams{
		Page:            int(request.Page),
		PageSize:        int(request.PageSize),
		MasterSegmentId: request.MasterSegmentId,
	})
	if err != nil {
		b.log.WithName("ProcessGetListPredictModels").Error(err, "cannot get list predict models", "request", request)
		return nil, err
	}

	return &api.GetListPredictModelsResponse{
		Code:    0,
		Message: "Success",
		Count:   queryResult.Count,
		Results: utils.Map(queryResult.PredictModels, func(model model.PredictModel) *api.PredictModel {
			return &api.PredictModel{
				Id:              model.ID,
				Name:            model.Name,
				MasterSegmentId: model.MasterSegmentId,
				Status:          string(model.Status),
				Labels:          nil,
				CreatedAt:       model.CreatedAt.String(),
				UpdatedAt:       model.UpdatedAt.String(),
			}
		}),
	}, nil
}
