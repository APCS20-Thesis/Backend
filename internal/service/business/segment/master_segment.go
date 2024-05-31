package segment

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
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
			Id:        masterSegment.ID,
			Name:      masterSegment.Name,
			Status:    string(masterSegment.Status),
			CreatedAt: masterSegment.CreatedAt.String(),
			UpdatedAt: masterSegment.UpdatedAt.String(),
		})
	}

	return int64(len(returnMasterSegments)), returnMasterSegments, nil
}

func (b business) GetMasterSegmentDetail(ctx context.Context, request *api.GetMasterSegmentDetailRequest, accountUuid string) (*api.MasterSegmentDetail, error) {
	// Get master segment
	masterSegment, err := b.repository.GetMasterSegment(ctx, request.Id, accountUuid)
	if err != nil {
		b.log.WithName("GetMasterSegmentDetail").Error(err, "cannot get master segment data")
		return nil, err
	}

	// Get audience table
	audienceTable, err := b.repository.GetAudienceTable(ctx, repository.GetAudienceTableParams{
		MasterSegmentId: request.Id,
	})
	if err != nil {
		b.log.WithName("GetMasterSegmentDetail").Error(err, "cannot get audience table data")
		return nil, err
	}
	// Parse build configuration info from audience table
	var buildConfiguration repository.AudienceBuildConfiguration
	err = json.Unmarshal(audienceTable.BuildConfiguration.RawMessage, &buildConfiguration)
	if err != nil {
		b.log.WithName("GetMasterSegmentDetail").Error(err, "cannot parse build configuration info")
		return nil, err
	}
	var audienceSchema []*api.SchemaColumn
	if audienceTable.Schema.RawMessage != nil && audienceTable.Schema.Valid {
		err = json.Unmarshal(audienceTable.Schema.RawMessage, &audienceSchema)
		if err != nil {
			b.log.WithName("GetMasterSegmentDetail").Error(err, "cannot parse audience schema")
			return nil, err
		}
	}

	// Get behavior tables
	behaviorTables, err := b.repository.ListBehaviorTables(ctx, repository.ListBehaviorTablesParams{
		MasterSegmentId: request.Id,
	})
	if err != nil {
		b.log.WithName("GetMasterSegmentDetail").Error(err, "cannot get behavior tables data")
		return nil, err
	}

	// Enrich table names for all tables
	var (
		enrichDataTableIds   []int64
		mapDataTableIdToName = make(map[int64]string)
	)
	enrichDataTableIds = append(enrichDataTableIds, buildConfiguration.MainTableId)
	for _, table := range buildConfiguration.AttributeTables {
		enrichDataTableIds = append(enrichDataTableIds, table.TableId)
	}
	for _, table := range behaviorTables {
		enrichDataTableIds = append(enrichDataTableIds, table.DataTableId)
	}
	dataTables, err := b.repository.DataTableRepository.ListDataTables(ctx, &repository.ListDataTablesFilters{
		DataTableIds: enrichDataTableIds,
	})
	for _, table := range dataTables {
		mapDataTableIdToName[table.ID] = table.Name
	}

	return &api.MasterSegmentDetail{
		Id:               masterSegment.ID,
		Name:             masterSegment.Name,
		Description:      masterSegment.Description,
		Status:           string(masterSegment.Status),
		CreatedAt:        masterSegment.CreatedAt.String(),
		UpdatedAt:        masterSegment.UpdatedAt.String(),
		AudienceTableId:  audienceTable.ID,
		MainRawTableId:   buildConfiguration.MainTableId,
		MainRawTableName: mapDataTableIdToName[buildConfiguration.MainTableId],
		AudienceSchema:   audienceSchema,
		AttributeTables: utils.Map(buildConfiguration.AttributeTables, func(table *repository.AttributeTableInfo) *api.MasterSegmentDetail_AttributeTable {
			return &api.MasterSegmentDetail_AttributeTable{
				RawTableId:      table.TableId,
				RawTableName:    mapDataTableIdToName[table.TableId],
				ForeignKey:      table.ForeignKey,
				JoinKey:         table.JoinKey,
				SelectedColumns: table.SelectedColumns,
			}
		}),
		BehaviorTables: utils.Map(behaviorTables, func(table model.BehaviorTable) *api.MasterSegmentDetail_BehaviorTable {
			return &api.MasterSegmentDetail_BehaviorTable{
				Id:           table.ID,
				Name:         table.Name,
				RawTableId:   table.DataTableId,
				RawTableName: mapDataTableIdToName[table.DataTableId],
				ForeignKey:   table.ForeignKey,
				JoinKey:      table.JoinKey,
			}
		}),
	}, nil
}
