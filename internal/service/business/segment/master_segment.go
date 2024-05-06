package segment

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
)

func (b business) CreateMasterSegment(ctx context.Context, request *api.CreateMasterSegmentRequest, accountUuid string) error {
	repoAttributeTables := make([]*repository.AttributeTableInfo, 0, len(request.AttributeTables))
	for _, attributeTable := range request.AttributeTables {
		repoAttributeTables = append(repoAttributeTables, &repository.AttributeTableInfo{
			TableId:         attributeTable.TableId,
			ForeignKey:      attributeTable.ForeignKey,
			JoinKey:         attributeTable.JoinKey,
			SelectedColumns: attributeTable.SelectedColumns,
		})
	}

	repoBehaviorTables := make([]*repository.CreateBehaviorTableParams, 0, len(request.BehaviorTables))
	for _, behaviorTable := range request.BehaviorTables {
		repoBehaviorTables = append(repoBehaviorTables, &repository.CreateBehaviorTableParams{
			Name:            behaviorTable.Name,
			TableId:         behaviorTable.TableId,
			ForeignKey:      behaviorTable.ForeignKey,
			JoinKey:         behaviorTable.JoinKey,
			SelectedColumns: behaviorTable.SelectedColumns,
		})
	}

	err := b.repository.TransactionRepository.CreateMasterSegmentTransaction(ctx, &repository.CreateMasterSegmentTransactionParams{
		MasterSegmentName: request.Name,
		Description:       request.Description,
		AccountUuid:       uuid.MustParse(accountUuid),
		AudienceName:      "audience",
		BuildConfiguration: repository.AudienceBuildConfiguration{
			MainTableId:     request.MainTableId,
			SelectedColumns: request.SelectedColumns,
			AttributeTables: repoAttributeTables,
		},
		BehaviorTables: repoBehaviorTables,
	}, b.airflowAdapter)
	if err != nil {
		b.log.WithName("CreateMasterSegment").Error(err, "cannot create master segment")
		return err
	}

	return nil
}

func (b business) ListMasterSegments(ctx context.Context, request *api.GetListMasterSegmentsRequest, accountUuid string) (int64, []*api.MasterSegment, error) {
	modelMasterSegments, err := b.repository.SegmentRepository.ListMasterSegments(ctx, &repository.ListMasterSegmentsParams{
		AccountUuid: uuid.MustParse(accountUuid),
	})
	if err != nil {
		b.log.WithName("ListMasterSegments").Error(err, "cannot get list master segment")
		return 0, nil, err
	}

	returnMasterSegments := make([]*api.MasterSegment, 0, len(modelMasterSegments))
	for _, masterSegment := range modelMasterSegments {
		returnMasterSegments = append(returnMasterSegments, &api.MasterSegment{
			Id:          masterSegment.ID,
			Name:        masterSegment.Name,
			Description: masterSegment.Description,
			Status:      string(masterSegment.Status),
			CreatedAt:   masterSegment.CreatedAt.String(),
			UpdatedAt:   masterSegment.UpdatedAt.String(),
		})
	}

	return int64(len(returnMasterSegments)), returnMasterSegments, nil
}
