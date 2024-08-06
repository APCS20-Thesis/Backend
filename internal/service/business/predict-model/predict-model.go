package predict_model

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/internal/service/business/segment"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"golang.org/x/exp/slices"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			Label1:           0,
			Label2:           1,
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
		dataAction, err := b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
			Tx:          tx,
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
		// 4. Create data action run
		_, err = b.repository.DataActionRunRepository.CreateDataActionRun(ctx, &repository.CreateDataActionRunParams{
			Tx:          tx,
			ActionId:    dataAction.ID,
			RunId:       0,
			Status:      model.DataActionRunStatus_Creating,
			AccountUuid: uuid.MustParse(accountUuid),
		})
		if err != nil {
			logger.Error(err, "cannot create data action run")
			tx.Rollback()
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
	logger := b.log.WithName("ProcessGetListPredictModels").WithValues("request", request)

	queryResult, err := b.repository.PredictModelRepository.ListPredictModels(ctx, &repository.ListPredictModelsParams{
		Page:            int(request.Page),
		PageSize:        int(request.PageSize),
		MasterSegmentId: request.MasterSegmentId,
		Statuses: utils.Map(request.Statuses, func(status string) model.PredictModelStatus {
			return model.PredictModelStatus(status)
		}),
	})
	if err != nil {
		logger.Error(err, "cannot get list predict models")
		return nil, err
	}

	if request.NotAppliedSegmentId > 0 {
		queryDataActions, err := b.repository.GetListDataActions(ctx, &repository.GetListDataActionsParams{
			ActionTypes: []string{string(model.ActionType_ApplyPredictModel)},
			AccountUuid: uuid.MustParse(accountUuid),
			TargetTable: model.TargetTable_Segment,
			ObjectId:    request.NotAppliedSegmentId,
		})
		if err != nil {
			logger.Error(err, "cannot get list data actions")
			return nil, err
		}

		excludedIds := make([]int64, 0, len(queryDataActions.DataActions))
		for _, action := range queryDataActions.DataActions {
			var config segment.ApplyPredictModelConfig
			err = json.Unmarshal(action.Payload.RawMessage, &config)
			if err != nil {
				logger.Error(err, "cannot unmarshal data action payload")
				return nil, err
			}
			excludedIds = append(excludedIds, config.PredictModelId)
		}

		var filteredPredictModels []model.PredictModel
		for _, predictModel := range queryResult.PredictModels {
			if !slices.Contains(excludedIds, predictModel.ID) {
				filteredPredictModels = append(filteredPredictModels, predictModel)
			}
		}

		queryResult.PredictModels = filteredPredictModels
		queryResult.Count = int64(len(filteredPredictModels))
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

func (b business) ProcessGetPredictModelDetail(ctx context.Context, request *api.GetPredictModelDetailRequest, accountUuid string) (*api.GetPredictModelDetailResponse, error) {
	logger := b.log.WithName("ProcessGetPredictModelDetail").WithValues("predictModelId", request.Id)
	predictModel, err := b.repository.PredictModelRepository.GetPredictModel(ctx, request.Id)
	if err != nil {
		logger.Error(err, "cannot get predict model from db")
		return nil, err
	}

	masterSegment, err := b.repository.SegmentRepository.GetMasterSegment(ctx, predictModel.MasterSegmentId)
	if err != nil {
		logger.Error(err, "cannot get master segment")
		return nil, err
	}

	if masterSegment.AccountUuid.String() != accountUuid {
		return nil, status.Error(codes.PermissionDenied, "not have access to this predict model")
	}

	var trainConfiguration model.PredictModelTrainConfiguration
	err = json.Unmarshal(predictModel.TrainConfigurations.RawMessage, &trainConfiguration)
	if err != nil {
		logger.Error(err, "cannot unmarshal train configurations")
		return nil, err
	}

	segment1, err := b.repository.SegmentRepository.GetSegment(ctx, trainConfiguration.Segment1, accountUuid)
	if err != nil {
		logger.Error(err, "cannot get segment", "id", trainConfiguration.Segment1)
		return nil, err
	}

	segment2, err := b.repository.SegmentRepository.GetSegment(ctx, trainConfiguration.Segment2, accountUuid)
	if err != nil {
		logger.Error(err, "cannot get segment", "id", trainConfiguration.Segment2)
		return nil, err
	}

	return &api.GetPredictModelDetailResponse{
		Code:    int32(code.Code_OK),
		Message: "Success",
		Id:      predictModel.ID,
		Name:    predictModel.Name,
		MasterSegment: &api.EnrichedMasterSegment{
			Id:   masterSegment.ID,
			Name: masterSegment.Name,
		},
		CreatedAt: predictModel.CreatedAt.String(),
		UpdatedAt: predictModel.UpdatedAt.String(),
		TrainSegments: []*api.EnrichedSegment{
			{
				Id:   segment1.ID,
				Name: segment1.Name,
			},
			{
				Id:   segment2.ID,
				Name: segment2.Name,
			},
		},
		Labels:          []string{trainConfiguration.Label1, trainConfiguration.Label2},
		TrainAttributes: trainConfiguration.SelectedAttributes,
		Status:          string(predictModel.Status),
	}, nil
}

func (b business) SyncOnTrainPredictModel(ctx context.Context, dataActionId int64) error {
	logger := b.log.WithName("SyncOnCreatePredictModel").WithValues("dataActionId", dataActionId)

	dataAction, err := b.repository.GetDataAction(ctx, dataActionId)
	if err != nil {
		logger.Error(err, "cannot get data action")
		return err
	}
	switch dataAction.Status {
	case model.DataActionStatus_Success:
		err = b.repository.PredictModelRepository.UpdatePredictModel(ctx, &repository.UpdatePredictModelParams{
			Id:     dataAction.ObjectId,
			Status: model.PredictModelStatus_OK,
		})
		if err != nil {
			logger.Error(err, "cannot update predict model status")
			return err
		}
	case model.DataActionStatus_Failed:
		err = b.repository.PredictModelRepository.UpdatePredictModel(ctx, &repository.UpdatePredictModelParams{
			Id:     dataAction.ObjectId,
			Status: model.PredictModelStatus_FAILED,
		})
		if err != nil {
			logger.Error(err, "cannot update predict model status")
			return err
		}
	}

	return nil
}
