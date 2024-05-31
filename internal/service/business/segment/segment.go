package segment

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/genproto/googleapis/rpc/code"
)

func (b business) CreateSegment(ctx context.Context, request *api.CreateSegmentRequest, accountUuid string) error {
	jsonCondition, err := json.Marshal(request.Condition)
	if err != nil {
		b.log.Error(err, "cannot parse condition into json", "condition", request.Condition)
		return err
	}

	err = b.repository.SegmentRepository.CreateSegment(ctx, &repository.CreateSegmentParams{
		Name:            request.Name,
		Description:     request.Description,
		MasterSegmentId: request.MasterSegmentId,
		Condition:       pqtype.NullRawMessage{RawMessage: jsonCondition, Valid: true},
		SqlCondition:    request.SqlCondition,
		AccountUuid:     uuid.MustParse(accountUuid),
	})
	if err != nil {
		b.log.Error(err, "cannot create segment")
		return err
	}

	return nil
}

func (b business) ListSegments(ctx context.Context, request *api.GetListSegmentsRequest, accountUuid string) ([]*api.Segment, error) {
	segments, err := b.repository.SegmentRepository.ListSegments(ctx, accountUuid)
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
	segment, err := b.repository.SegmentRepository.GetSegment(ctx, request.Id, accountUuid)
	if err != nil {
		b.log.WithName("GetSegmentDetail").Error(err, "cannot get segment")
		return nil, err
	}

	masterSegment, err := b.repository.SegmentRepository.GetMasterSegment(ctx, segment.MasterSegmentId, accountUuid)
	if err != nil {
		b.log.WithName("GetSegmentDetail").Error(err, "cannot get segment")
		return nil, err
	}

	var condition api.GetSegmentDetailResponse_Rule
	err = json.Unmarshal(segment.Condition.RawMessage, &condition)
	if err != nil {
		b.log.WithName("GetSegmentDetail").Error(err, "cannot unmarshal condition")
		return nil, err
	}

	audienceTable, err := b.repository.SegmentRepository.GetAudienceTable(ctx, repository.GetAudienceTableParams{MasterSegmentId: masterSegment.ID})
	if err != nil {
		b.log.WithName("GetSegmentDetail").Error(err, "cannot get audience table", "masterSegmentId", masterSegment.ID)
		return nil, err
	}

	var audienceSchema []*api.SchemaColumn
	err = json.Unmarshal(audienceTable.Schema.RawMessage, &audienceSchema)
	if err != nil {
		b.log.WithName("GetSegmentDetail").Error(err, "cannot parse audience schema")
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
