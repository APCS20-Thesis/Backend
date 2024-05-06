package segment

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/repository"
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
		AccountUuid:     uuid.MustParse(accountUuid),
	})
	if err != nil {
		b.log.Error(err, "cannot create segment")
		return err
	}

	return nil
}
