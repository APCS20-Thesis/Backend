package segment

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (b business) ParseHavingCondition(havingConditions []*api.BehaviorCondition_HavingCondition) string {
	listTopConditions := make([]string, 0, len(havingConditions))
	for _, condition := range havingConditions {
		conditionClause := strings.Join(condition.ConditionClause, " "+condition.Combinator+" ")
		if len(condition.ConditionClause) > 1 {
			conditionClause = "(" + conditionClause + ")"
		}
		listTopConditions = append(listTopConditions, conditionClause)
	}
	return strings.Join(listTopConditions, " AND ")
}

func (b business) CreateSegment(ctx context.Context, request *api.CreateSegmentRequest, accountUuid string) error {
	logger := b.log.WithName("CreateSegment").WithValues("request", request)
	condition := model.SegmentBuildConditions{
		AudienceCondition:  request.Condition,
		BehaviorConditions: request.BehaviorConditions,
	}
	jsonCondition, err := json.Marshal(condition)
	if err != nil {
		logger.Error(err, "cannot marshal condition")
		return err
	}

	behaviorTables, err := b.repository.SegmentRepository.ListBehaviorTables(ctx, repository.ListBehaviorTablesParams{
		MasterSegmentId: request.MasterSegmentId,
	})
	behaviorTablesMap := make(map[int64]model.BehaviorTable)
	for _, table := range behaviorTables {
		behaviorTablesMap[table.ID] = table
	}

	behaviorConfig := make([]airflow.SegmentBehaviorCondition, 0, len(request.BehaviorConditions))
	for _, each := range request.BehaviorConditions {
		behaviorTable := behaviorTablesMap[each.BehaviorTableId]
		key := utils.GenerateDeltaBehaviorPath(request.MasterSegmentId, behaviorTable.Name)
		behaviorConfig = append(behaviorConfig, airflow.SegmentBehaviorCondition{
			BehaviorTableKey:   key,
			JoinKey:            behaviorTable.JoinKey,
			ForeignKey:         behaviorTable.ForeignKey,
			WhereClauseValue:   each.WhereSqlCondition,
			GroupByClauseValue: each.GroupByKeys,
			HavingClauseValue:  b.ParseHavingCondition(each.HavingConditions),
		})
	}
	behaviorsJson, err := json.Marshal(behaviorConfig)
	if err != nil {
		logger.Error(err, "cannot marshal behavior config", "behaviorConfig", behaviorConfig)
		return err
	}

	// Prepare calling airflow
	dagId := utils.GenerateDagId(accountUuid, model.ActionType_CreateSegment)

	tx := b.db.Begin()

	// 1. Save Segment
	segment, err := b.repository.SegmentRepository.CreateSegment(ctx, &repository.CreateSegmentParams{
		Tx:              tx,
		Name:            request.Name,
		Description:     request.Description,
		MasterSegmentId: request.MasterSegmentId,
		Condition:       pqtype.NullRawMessage{RawMessage: jsonCondition, Valid: true},
		SqlCondition:    request.SqlCondition,
		AccountUuid:     uuid.MustParse(accountUuid),
	})
	if err != nil {
		logger.Error(err, "cannot create segment")
		tx.Rollback()
		return err
	}

	// 2. Airflow generate segment
	payload := &airflow.TriggerGenerateDagCreateSegmentRequest{Config: airflow.CreateSegmentConfig{
		DagId:                       dagId,
		AudienceTableKey:            utils.GenerateDeltaAudiencePath(request.MasterSegmentId),
		AudienceCondition:           request.SqlCondition,
		SegmentTableKey:             utils.GenerateDeltaSegmentPath(request.MasterSegmentId, segment.ID),
		SegmentTableName:            "segment",
		BehaviorTableConfigurations: string(behaviorsJson),
	}}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err, "cannot marshal payload", "payload", payload)
		tx.Rollback()
		return err
	}
	err = b.airflowAdapter.TriggerGenerateDagCreateSegment(ctx, payload)
	if err != nil {
		logger.Error(err, "cannot trigger dag generate create segment")
		tx.Rollback()
		return err
	}

	// 3. Save data action
	_, err = b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		TargetTable: model.TargetTable_Segment,
		ActionType:  model.ActionType_CreateSegment,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    segment.ID,
		Payload:     pqtype.NullRawMessage{RawMessage: payloadJson, Valid: payloadJson != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (b business) ListSegments(ctx context.Context, request *api.GetListSegmentsRequest, accountUuid string) ([]*api.Segment, error) {
	segments, err := b.repository.SegmentRepository.ListSegments(ctx, &repository.ListSegmentFilter{
		AccountUuid:      accountUuid,
		MasterSegmentIds: request.MasterSegmentIds,
		Statuses:         request.Statuses,
	})
	if err != nil {
		b.log.WithName("ListSegments").Error(err, "cannot get list segments")
		return nil, err
	}

	return utils.Map(segments, func(segment repository.SegmentListItem) *api.Segment {
		return &api.Segment{
			Id:                segment.ID,
			Name:              segment.Name,
			MasterSegmentId:   segment.MasterSegmentId,
			MasterSegmentName: segment.MasterSegmentName,
			CreatedAt:         segment.CreatedAt.String(),
			UpdatedAt:         segment.UpdatedAt.String(),
			Status:            segment.Status,
		}
	}), nil
}

func (b business) GetSegmentDetail(ctx context.Context, request *api.GetSegmentDetailRequest, accountUuid string) (*api.GetSegmentDetailResponse, error) {
	logger := b.log.WithName("GetSegmentDetail").WithValues("id", request.Id)
	segment, err := b.repository.SegmentRepository.GetSegment(ctx, request.Id, accountUuid)
	if err != nil {
		logger.Error(err, "cannot get segment")
		return nil, err
	}

	masterSegment, err := b.repository.SegmentRepository.GetMasterSegment(ctx, segment.MasterSegmentId)
	if err != nil {
		logger.Error(err, "cannot get segment")
		return nil, err
	}

	var condition api.SegmentCondition
	err = json.Unmarshal(segment.Condition.RawMessage, &condition)
	if err != nil {
		logger.Error(err, "cannot unmarshal condition")
		return nil, err
	}
	condition.AudienceSqlCondition = segment.SqlCondition

	audienceTable, err := b.repository.SegmentRepository.GetAudienceTable(ctx, repository.GetAudienceTableParams{MasterSegmentId: masterSegment.ID})
	if err != nil {
		logger.Error(err, "cannot get audience table", "masterSegmentId", masterSegment.ID)
		return nil, err
	}

	var audienceSchema []*api.SchemaColumn
	err = json.Unmarshal(audienceTable.Schema.RawMessage, &audienceSchema)
	if err != nil {
		logger.Error(err, "cannot parse audience schema")
		return nil, err
	}

	return &api.GetSegmentDetailResponse{
		Code:              int32(code.Code_OK),
		Message:           "Success",
		Id:                segment.ID,
		Name:              segment.Name,
		Description:       segment.Description,
		MasterSegmentId:   segment.MasterSegmentId,
		MasterSegmentName: masterSegment.Name,
		CreatedAt:         segment.CreatedAt.String(),
		UpdatedAt:         segment.UpdatedAt.String(),
		Condition:         &condition,
		Schema:            audienceSchema,
	}, nil
}

type ApplyPredictModelConfig struct {
	DagConfig      airflow.DagApplyPredictModelConfig
	PredictModelId int64
	SegmentId      int64
}

func (b business) ProcessApplyPredictModel(ctx context.Context, request *api.ApplyPredictModelRequest, accountUuid string) (*api.ApplyPredictModelResponse, error) {
	logger := b.log.WithName("ProcessApplyPredictModel").WithValues("request", request)

	predictModel, err := b.repository.PredictModelRepository.GetPredictModel(ctx, request.PredictModelId)
	if err != nil {
		logger.Error(err, "cannot get predict model")
		return nil, err
	}

	segment, err := b.repository.SegmentRepository.GetSegment(ctx, request.SegmentId, accountUuid)
	if err != nil {
		logger.Error(err, "cannot get segment")
		return nil, err
	}

	if segment.MasterSegmentId != predictModel.MasterSegmentId {
		return nil, status.Error(codes.InvalidArgument, "cannot predict segment with other master segment's model")
	}

	var config model.PredictModelTrainConfiguration
	err = json.Unmarshal(predictModel.TrainConfigurations.RawMessage, &config)
	if err != nil {
		logger.Error(err, "cannot unmarshal predict model train configurations")
		return nil, err
	}

	dagId := utils.GenerateDagId(accountUuid, model.ActionType_ApplyPredictModel)
	payload := &airflow.TriggerGenerateDagApplyPredictModelRequest{
		Conf: airflow.DagApplyPredictModelConfig{
			DagId:            dagId,
			DataKey:          utils.GenerateDeltaSegmentPath(predictModel.MasterSegmentId, request.SegmentId),
			ModelPath:        utils.GenerateDeltaPredictModelFilePath(predictModel.MasterSegmentId, request.PredictModelId),
			ResultPath:       utils.GenerateDeltaPredictResult(predictModel.MasterSegmentId),
			SelectAttributes: config.SelectedAttributes,
		},
	}
	jsonPayload, err := json.Marshal(ApplyPredictModelConfig{
		DagConfig:      payload.Conf,
		PredictModelId: request.PredictModelId,
		SegmentId:      request.SegmentId,
	})
	if err != nil {
		logger.Error(err, "cannot marshal payload", "payload", payload)
		return nil, err
	}

	err = b.airflowAdapter.TriggerGenerateDagApplyPredictModel(ctx, payload)
	if err != nil {
		logger.Error(err, "cannot trigger generate dag apply predict model")
		return nil, err
	}

	_, err = b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		TargetTable: model.TargetTable_Segment,
		ActionType:  model.ActionType_ApplyPredictModel,
		Schedule:    "",
		AccountUuid: uuid.MustParse(accountUuid),
		DagId:       dagId,
		Status:      model.DataActionStatus_Pending,
		ObjectId:    request.SegmentId,
		Payload:     pqtype.NullRawMessage{RawMessage: jsonPayload, Valid: jsonPayload != nil},
	})
	if err != nil {
		logger.Error(err, "cannot create data action")
		return nil, err
	}

	return &api.ApplyPredictModelResponse{
		Code:    0,
		Message: "Success",
	}, nil
}

func (b business) ProcessGetListPredictionActions(ctx context.Context, request *api.GetListPredictionActionsRequest, accountUuid string) (*api.GetListPredictionActionsResponse, error) {
	logger := b.log.WithName("ProcessGetListPredictionActions")

	queryResult, err := b.repository.DataActionRepository.GetListDataActions(ctx, &repository.GetListDataActionsParams{
		ActionTypes: []string{string(model.ActionType_ApplyPredictModel)},
		AccountUuid: uuid.MustParse(accountUuid),
		Page:        int(request.Page),
		PageSize:    int(request.PageSize),
		TargetTable: model.TargetTable_Segment,
		ObjectId:    request.Id,
	})
	if err != nil {
		logger.Error(err, "cannot get list data actions")
		return nil, err
	}

	modelIds := make([]int64, 0, len(queryResult.DataActions))
	configMap := make(map[int64]ApplyPredictModelConfig)
	for _, action := range queryResult.DataActions {
		var config ApplyPredictModelConfig
		err := json.Unmarshal(action.Payload.RawMessage, &config)
		if err != nil {
			logger.Error(err, "cannot unmarshal data action payload", "actionId", action.ID)
			return nil, err
		}
		configMap[action.ID] = config
		modelIds = append(modelIds, config.PredictModelId)
	}

	predictModels, err := b.repository.PredictModelRepository.ListPredictModels(ctx, &repository.ListPredictModelsParams{
		Ids: modelIds,
	})
	if err != nil {
		logger.Error(err, "cannot get list predict models")
		return nil, err
	}
	predictModelMap := make(map[int64]string)
	for _, predictModel := range predictModels.PredictModels {
		predictModelMap[predictModel.ID] = predictModel.Name
	}

	predictActions := utils.Map(queryResult.DataActions, func(action model.DataAction) *api.GetListPredictionActionsResponse_PredictionAction {
		modelId := configMap[action.ID].PredictModelId
		return &api.GetListPredictionActionsResponse_PredictionAction{
			Id:        action.ID,
			ModelId:   modelId,
			ModelName: predictModelMap[modelId],
			Status:    string(action.Status),
			CreatedAt: action.CreatedAt.String(),
			UpdatedAt: action.UpdatedAt.String(),
		}
	})

	return &api.GetListPredictionActionsResponse{
		Code:    0,
		Message: "Success",
		Count:   queryResult.Count,
		Results: predictActions,
	}, nil
}

func (b business) ProcessGetResultPredictionActions(ctx context.Context, request *api.GetResultPredictionActionsRequest, accountUuid string) (*api.GetResultPredictionActionsResponse, error) {
	logger := b.log.WithName("ProcessGetResultPredictionActions")

	action, err := b.repository.DataActionRepository.GetDataAction(ctx, request.ActionId)
	if err != nil {
		logger.Error(err, "cannot get data actions")
		return nil, err
	}

	var config ApplyPredictModelConfig
	err = json.Unmarshal(action.Payload.RawMessage, &config)
	if err != nil {
		logger.Error(err, "cannot unmarshal data action payload", "actionId", action.ID)
		return nil, err
	}

	segmentPath := fmt.Sprintf("s3a://%s/%s", b.config.S3StorageConfig.Bucket, config.DagConfig.DataKey)
	resultPath := fmt.Sprintf("s3a://%s/%s", b.config.S3StorageConfig.Bucket, config.DagConfig.ResultPath)
	queryResponse, err := b.queryAdapter.QueryRawSQLV2(ctx, &query.QueryRawSQLV2Request{
		Query: fmt.Sprintf("SELECT * FROM delta.`%s` AS segment LEFT JOIN delta.`%s` AS result ON segment.cdp_system_uuid = result.cdp_system_uuid;", segmentPath, resultPath),
	})
	if err != nil {
		return nil, err
	}

	return &api.GetResultPredictionActionsResponse{
		Code:  0,
		Count: int64(queryResponse.Count),
		Data:  queryResponse.Data,
	}, nil
}
