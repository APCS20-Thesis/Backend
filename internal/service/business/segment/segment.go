package segment

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
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
