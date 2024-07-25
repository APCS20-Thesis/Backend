package segment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (b business) ListMasterSegments(ctx context.Context, request *api.GetListMasterSegmentsRequest, accountUuid string) (*api.GetListMasterSegmentsResponse, error) {
	modelMasterSegments, err := b.repository.SegmentRepository.ListMasterSegments(ctx, &repository.ListMasterSegmentsFilter{
		AccountUuid: uuid.MustParse(accountUuid),
		Name:        request.Name,
		Status:      model.MasterSegmentStatus(request.Status),
		PageSize:    int(request.PageSize),
		Page:        int(request.Page),
	})
	if err != nil {
		b.log.WithName("ListMasterSegments").Error(err, "cannot get list master segment")
		return nil, err
	}

	var returnMasterSegments []*api.MasterSegment
	for _, masterSegment := range modelMasterSegments.MasterSegments {
		returnMasterSegments = append(returnMasterSegments, &api.MasterSegment{
			Id:        masterSegment.ID,
			Name:      masterSegment.Name,
			Status:    string(masterSegment.Status),
			UpdatedAt: masterSegment.UpdatedAt.String(),
			CreatedAt: masterSegment.CreatedAt.String(),
		})
	}
	return &api.GetListMasterSegmentsResponse{
		Code:    0,
		Count:   modelMasterSegments.Count,
		Results: returnMasterSegments,
	}, nil
}

func (b business) GetMasterSegmentDetail(ctx context.Context, request *api.GetMasterSegmentDetailRequest, accountUuid string) (*api.MasterSegmentDetail, error) {
	// Get master segment
	masterSegment, err := b.repository.GetMasterSegment(ctx, request.Id)
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
	listTablesResult, err := b.repository.DataTableRepository.ListDataTables(ctx, &repository.ListDataTablesFilters{
		DataTableIds: enrichDataTableIds,
	})
	for _, table := range listTablesResult.DataTables {
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
			var behaviorSchema []*api.SchemaColumn
			if table.Schema.RawMessage != nil && table.Schema.Valid {
				err = json.Unmarshal(table.Schema.RawMessage, &behaviorSchema)
				if err != nil {
					b.log.WithName("GetMasterSegmentDetail").Error(err, "cannot parse audience schema")
					return nil
				}
			}
			return &api.MasterSegmentDetail_BehaviorTable{
				Id:           table.ID,
				Name:         table.Name,
				RawTableId:   table.DataTableId,
				RawTableName: mapDataTableIdToName[table.DataTableId],
				ForeignKey:   table.ForeignKey,
				JoinKey:      table.JoinKey,
				Schema:       behaviorSchema,
			}
		}),
	}, nil
}

func (b business) ListMasterSegmentProfiles(ctx context.Context, request *api.GetListMasterSegmentProfilesRequest, accountUuid string) (int64, []string, error) {
	masterSegment, err := b.repository.GetMasterSegment(ctx, request.Id)
	if err != nil {
		b.log.WithName("GetMasterSegmentDetail").Error(err, "cannot get master segment data")
		return 0, nil, err
	}

	if masterSegment.AccountUuid.String() != accountUuid {
		b.log.WithName("ListMasterSegmentProfiles").WithValues("masterSegmentId", request.Id).Error(err, "No have permission with dataTable")
		return 0, nil, status.Error(codes.PermissionDenied, "No have permission with master segment")
	}
	path := "s3a://cdp-thesis-apcs/" + utils.GenerateDeltaAudiencePath(request.Id)
	queryResponse, err := b.queryAdapter.QueryRawSQLV2(ctx, &query.QueryRawSQLV2Request{
		Query: fmt.Sprintf("SELECT * FROM delta.`%s`;", path),
	})
	if err != nil {
		return 0, nil, err
	}

	res := query.QueryV2Paginate(request.Page, request.PageSize, queryResponse.Data)
	return int64(queryResponse.Count), res, nil
}

func (b business) GetMasterSegmentProfile(ctx context.Context, request *api.GetMasterSegmentProfileRequest, accountUuid string) (string, error) {
	masterSegment, err := b.repository.GetMasterSegment(ctx, request.Id)
	if err != nil {
		b.log.WithName("GetMasterSegmentDetail").Error(err, "cannot get master segment data")
		return "", err
	}

	if masterSegment.AccountUuid.String() != accountUuid {
		b.log.WithName("ListMasterSegmentProfiles").WithValues("masterSegmentId", request.Id).Error(err, "No have permission with dataTable")
		return "", status.Error(codes.PermissionDenied, "No have permission with master segment")
	}
	path := "s3a://cdp-thesis-apcs/" + utils.GenerateDeltaAudiencePath(request.Id)
	queryResponse, err := b.queryAdapter.QueryRawSQLV2(ctx, &query.QueryRawSQLV2Request{
		Query: fmt.Sprintf("SELECT * FROM delta.`%s`;", path),
	})
	if err != nil {
		return "", err
	}

	return queryResponse.Data[0], nil
}

func (b business) SyncOnCreateMasterSegment(ctx context.Context, masterSegmentId int64, actionStatus model.DataActionStatus) error {
	if actionStatus == model.DataActionStatus_Failed {
		err := b.repository.SegmentRepository.UpdateMasterSegment(ctx, &repository.UpdateMasterSegmentParams{
			Id:     masterSegmentId,
			Status: model.MasterSegmentStatus_FAILED,
		})
		if err != nil {
			b.log.WithName("SyncOnCreateMasterSegment").Error(err, "cannot update master segment status", "masterSegmentId", masterSegmentId)
			return err
		}
		return nil
	} else if actionStatus != model.DataActionStatus_Success {
		return nil
	}
	audienceTable, err := b.repository.SegmentRepository.GetAudienceTable(ctx, repository.GetAudienceTableParams{MasterSegmentId: masterSegmentId})
	if err != nil {
		b.log.WithName("SyncOnCreateMasterSegment").Error(err, "cannot get audience table", "masterSegmentId", masterSegmentId)
		return err
	}
	behaviorTables, err := b.repository.SegmentRepository.ListBehaviorTables(ctx, repository.ListBehaviorTablesParams{MasterSegmentId: masterSegmentId})
	if err != nil {
		b.log.WithName("SyncOnCreateMasterSegment").Error(err, "cannot get behavior tables", "masterSegmentId", masterSegmentId)
		return err
	}
	// Sync audience table schema
	response, err := b.queryAdapter.GetSchemaTable(ctx, &query.GetSchemaDataTableRequest{
		TablePath: utils.GenerateDeltaAudiencePath(masterSegmentId),
	})
	if err != nil {
		b.log.WithName("SyncOnCreateMasterSegment").Error(err, "cannot query audience table schema")
		return err
	}
	schema := utils.Map(response.Schema, func(unit query.FieldSchema) model.SchemaUnit {
		return model.SchemaUnit{
			ColumnName: unit.Name,
			DataType:   unit.Type,
		}
	})
	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		b.log.WithName("SyncOnCreateMasterSegment").Error(err, "cannot parse schema")
		return err
	}

	tx := b.db.Begin()

	err = b.repository.SegmentRepository.UpdateAudienceTable(ctx, &repository.UpdateAudienceTableParams{
		Tx:     tx,
		Id:     audienceTable.ID,
		Schema: pqtype.NullRawMessage{RawMessage: jsonSchema, Valid: jsonSchema != nil},
	})
	if err != nil {
		b.log.WithName("SyncOnCreateMasterSegment").Error(err, "cannot update audience table", "id", audienceTable.ID)
		tx.Rollback()
		return err
	}
	// Sync behavior schemas
	for _, behaviorTable := range behaviorTables {
		response, err := b.queryAdapter.GetSchemaTable(ctx, &query.GetSchemaDataTableRequest{
			TablePath: utils.GenerateDeltaBehaviorPath(masterSegmentId, behaviorTable.Name),
		})
		if err != nil {
			b.log.WithName("job:SyncOnCreateMasterSegment").Error(err, "cannot query behavior table schema")
			return err
		}
		schema := utils.Map(response.Schema, func(unit query.FieldSchema) model.SchemaUnit {
			return model.SchemaUnit{
				ColumnName: unit.Name,
				DataType:   unit.Type,
			}
		})
		jsonSchema, err := json.Marshal(schema)
		if err != nil {
			b.log.WithName("SyncOnCreateMasterSegment").Error(err, "cannot parse schema")
			return err
		}
		err = b.repository.SegmentRepository.UpdateBehaviorTable(ctx, &repository.UpdateBehaviorTableParams{
			Tx:     tx,
			Id:     behaviorTable.ID,
			Schema: pqtype.NullRawMessage{RawMessage: jsonSchema, Valid: jsonSchema != nil},
		})
		if err != nil {
			b.log.WithName("SyncOnCreateMasterSegment").Error(err, "cannot update behavior table", "id", behaviorTable.ID)
			tx.Rollback()
			return err
		}
	}
	err = b.repository.SegmentRepository.UpdateMasterSegment(ctx, &repository.UpdateMasterSegmentParams{
		Tx:     tx,
		Id:     masterSegmentId,
		Status: model.MasterSegmentStatus_UP_TO_DATE,
	})
	if err != nil {
		b.log.WithName("SyncOnCreateMasterSegment").Error(err, "cannot update master segment status", "masterSegmentId", masterSegmentId)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
