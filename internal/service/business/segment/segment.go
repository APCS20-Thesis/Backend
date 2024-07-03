package segment

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
	"google.golang.org/genproto/googleapis/rpc/code"
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

	var condition api.GetSegmentDetailResponse_Rule
	err = json.Unmarshal(segment.Condition.RawMessage, &condition)
	if err != nil {
		logger.Error(err, "cannot unmarshal condition")
		return nil, err
	}

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
