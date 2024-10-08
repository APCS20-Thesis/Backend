package data_source

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (b business) CreateDataSource(ctx context.Context, params *repository.CreateDataSourceParams) (*model.DataSource, error) {
	dataSource, err := b.repository.DataSourceRepository.CreateDataSource(ctx, params)
	if err != nil {
		b.log.WithName("CreateDataSource").
			WithValues("Context", ctx).
			Error(err, "Cannot create data_source")
		return nil, err
	}
	return dataSource, nil
}

func (b business) GetListDataSources(ctx context.Context, request *api.GetListDataSourcesRequest, accountUuid string) (*api.GetListDataSourcesResponse, error) {
	listDataSourcesResult, err := b.repository.DataSourceRepository.ListDataSources(ctx,
		&repository.ListDataSourcesFilters{
			Type:        model.DataSourceType(request.Type),
			AccountUuid: uuid.MustParse(accountUuid),
			Name:        request.Name,
			Page:        int(request.Page),
			PageSize:    int(request.PageSize),
		})
	if err != nil {
		b.log.WithName("GetListDataSources").
			WithValues("Context", ctx).
			Error(err, "Cannot get list data_sources")
		return nil, err
	}
	var dataSources []*api.GetListDataSourcesResponse_DataSource
	for _, dataSource := range listDataSourcesResult.DataSource {
		dataSources = append(dataSources, &api.GetListDataSourcesResponse_DataSource{Id: dataSource.ID, Name: dataSource.Name, Type: string(dataSource.Type), UpdatedAt: dataSource.UpdatedAt.String()})
	}
	return &api.GetListDataSourcesResponse{
		Code:    0,
		Count:   listDataSourcesResult.Count,
		Results: dataSources,
	}, nil
}

func (b business) GetDataSource(ctx context.Context, request *api.GetDataSourceRequest, accountUuid string) (*api.GetDataSourceResponse, error) {
	dataSource, err := b.repository.DataSourceRepository.GetDataSource(ctx, request.Id)
	if err != nil {
		b.log.WithName("GetDataSource").Error(err, "Cannot get data_source")
		return nil, err
	}
	if dataSource.AccountUuid != uuid.MustParse(accountUuid) {
		b.log.WithName("GetDataSource").Info("Only owner can get data_source")
		return nil, status.Error(codes.PermissionDenied, "Only owner can get data_source")
	}
	//var configurations map[string]interface{}
	//if dataSource.Configurations.RawMessage != nil {
	//	err = json.Unmarshal(dataSource.Configurations.RawMessage, &configurations)
	//	if err != nil {
	//		b.log.WithName("GetDataSource").Error(err, "Cannot get parse configuration")
	//		return nil, err
	//	}
	//}
	var enrichedConnection *api.EnrichedConnection
	if dataSource.ConnectionId != 0 {
		connection, err := b.repository.ConnectionRepository.GetConnection(ctx, dataSource.ConnectionId)
		if err != nil {
			b.log.WithName("GetDataSource").Info("cannot enrich connection info")
			return nil, err
		}
		enrichedConnection = &api.EnrichedConnection{
			Id:   connection.ID,
			Name: connection.Name,
			Type: string(connection.Type),
		}
	}

	return &api.GetDataSourceResponse{
		Code:           int32(code.Code_OK),
		Id:             dataSource.ID,
		Name:           dataSource.Name,
		Type:           string(dataSource.Type),
		Description:    dataSource.Description,
		CreatedAt:      dataSource.CreatedAt.String(),
		UpdatedAt:      dataSource.UpdatedAt.String(),
		Configurations: string(dataSource.Configurations.RawMessage),
		Connection:     enrichedConnection,
	}, nil
}

func (b business) GetListSourceTableMappings(ctx context.Context, request *api.GetListSourceTableMapRequest) (*api.GetListSourceTableMapResponse, error) {
	queryResult, err := b.repository.SourceTableMapRepository.ListSourceTableMap(ctx, &repository.ListSourceTableMapParams{
		TableId:  request.TableId,
		SourceId: request.DataSourceId,
	})
	if err != nil {
		b.log.WithName("GetListSourceTableMappings").Error(err, "cannot get list source table maps from db")
		return nil, err
	}

	returnMaps := make([]*api.SourceTableMap, 0, len(queryResult.TableSourceMaps))
	for _, modelMap := range queryResult.TableSourceMaps {
		var mappingOptions []*api.MappingOptionItem
		if modelMap.MappingOptions.Valid {
			err := json.Unmarshal(modelMap.MappingOptions.RawMessage, &mappingOptions)
			if err != nil {
				b.log.WithName("GetListSourceTableMappings").Error(err, "cannot unmarshal mapping options", "sourceTableMapId", modelMap.ID)
				return nil, err
			}
		}
		returnMaps = append(returnMaps, &api.SourceTableMap{
			Id: modelMap.ID,
			Table: &api.EnrichedTable{
				Id:   modelMap.TableId,
				Name: modelMap.TableName,
			},
			Source: &api.EnrichedDataSource{
				Id:   modelMap.SourceId,
				Name: modelMap.SourceName,
				Type: string(modelMap.SourceType),
			},
			Mappings: mappingOptions,
		})
	}

	return &api.GetListSourceTableMapResponse{
		Code:    0,
		Message: "Success",
		Count:   queryResult.Count,
		Results: returnMaps,
	}, nil
}

func (b business) SyncOnImportFromSourceSuccess(ctx context.Context, dataActionId int64) error {
	dataAction, err := b.repository.DataActionRepository.GetDataAction(ctx, dataActionId)
	if err != nil {
		b.log.WithName("SyncOnImportFromSourceSuccess").Error(err, "cannot get data action ")
	}
	sourceTableMap, err := b.repository.SourceTableMapRepository.GetSourceTableMapById(ctx, dataAction.ObjectId)
	if err != nil {
		b.log.WithName("SyncOnImportFromSourceSuccess").Error(err, "cannot get source table map", "sourceTableMapId", dataAction.ObjectId)
		return err
	}
	path, err := b.repository.DataTableRepository.GetDataTableDeltaPath(ctx, sourceTableMap.TableId)
	if err != nil {
		b.log.WithName("SyncOnImportFromSourceSuccess").Error(err, "cannot get table delta path", "TableId", sourceTableMap.TableId)
		return err
	}
	response, err := b.queryAdapter.GetSchemaTable(ctx, &query.GetSchemaDataTableRequest{TablePath: path})
	if err != nil {
		b.log.WithName("GetSchemaTable").Info("error here")
		err = b.repository.DataTableRepository.UpdateStatusDataTable(ctx, sourceTableMap.TableId, model.DataTableStatus_NEED_TO_SYNC)
		if err != nil {
			b.log.WithName("SyncOnImportFromSourceSuccess").Error(err, "cannot update data table status", "TableId", sourceTableMap.TableId)
			return err
		}
	}
	schema := make([]model.SchemaUnit, 0, len(response.Schema))
	for _, column := range response.Schema {
		schema = append(schema, model.SchemaUnit{
			ColumnName: column.Name,
			DataType:   column.Type,
		})
	}
	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		b.log.WithName("SyncOnImportFromSourceSuccess").Error(err, "cannot parse schema")
		return err
	}
	_, err = b.repository.DataTableRepository.UpdateDataTable(ctx, &repository.UpdateDataTableParams{
		ID:     sourceTableMap.TableId,
		Schema: pqtype.NullRawMessage{Valid: true, RawMessage: jsonSchema},
		Status: model.DataTableStatus_UP_TO_DATE,
	})
	if err != nil {
		b.log.WithName("SyncOnImportFromSourceSuccess").Error(err, "cannot update data table", "TableId", sourceTableMap.TableId)
		return err
	}

	return nil
}

func (b business) SynOnImportFromSourceFailed(ctx context.Context, dataActionId int64) error {
	dataAction, err := b.repository.DataActionRepository.GetDataAction(ctx, dataActionId)
	if err != nil {
		b.log.WithName("SynOnImportFromSourceFailed").Error(err, "cannot get data action ")
	}
	sourceTableMap, err := b.repository.SourceTableMapRepository.GetSourceTableMapById(ctx, dataAction.ObjectId)
	if err != nil {
		b.log.WithName("SynOnImportFromSourceFailed").Error(err, "cannot get source table map", "sourceTableMapId", dataAction.ObjectId)
		return err
	}
	err = b.repository.DataSourceRepository.UpdateDataSource(ctx, &repository.UpdateDataSourceParams{
		ID:     sourceTableMap.SourceId,
		Status: model.DataSourceStatus_Failed,
	})
	if err != nil {
		b.log.WithName("SynOnImportFromSourceFailed").Error(err, "failed to update source status", "sourceTableMapId", dataAction.ObjectId)
		return err
	}
	return nil
}
