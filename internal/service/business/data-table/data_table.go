package data_table

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (b business) CreateDataTable(ctx context.Context, params *repository.CreateDataTableParams) (*model.DataTable, error) {
	dataTable, err := b.repository.DataTableRepository.CreateDataTable(ctx, params)
	if err != nil {
		b.log.WithName("CreateDataTable").
			WithValues("Context", ctx).
			Error(err, "Cannot create data_table")
		return nil, err
	}
	return dataTable, nil
}

func (b business) UpdateDataTable(ctx context.Context, params *repository.UpdateDataTableParams) error {
	err := b.repository.DataTableRepository.UpdateDataTable(ctx, params)
	if err != nil {
		b.log.WithName("UpdateDataTable").
			WithValues("Context", ctx).
			Error(err, "Cannot update data_table")
		return err
	}
	return nil
}

func (b business) GetListDataTables(ctx context.Context, request *api.GetListDataTablesRequest, accountUuid string) ([]*api.GetListDataTablesResponse_DataTable, int64, error) {
	logger := b.log.WithName("GetListDataTables").WithValues("request", request)

	queryResult, err := b.repository.DataTableRepository.ListDataTables(ctx,
		&repository.ListDataTablesFilters{
			Name:        request.Name,
			AccountUuid: accountUuid,
			Page:        int(request.Page),
			PageSize:    int(request.PageSize),
		})
	if err != nil {
		logger.Error(err, "Cannot get list data_Tables")
		return nil, 0, err
	}

	tableIds := utils.Map(queryResult.DataTables, func(table model.DataTable) int64 { return table.ID })
	enrichedSources, err := b.repository.DataTableRepository.GetSourcesOfDataTables(ctx, tableIds)
	if err != nil {
		logger.Error(err, "cannot enrich data sources", "table ids", tableIds)
		return nil, 0, err
	}

	enrichedDestinations, err := b.repository.DataTableRepository.GetDestinationsOfDataTables(ctx, tableIds)
	if err != nil {
		logger.Error(err, "cannot enrich data destinations", "table ids", tableIds)
		return nil, 0, err
	}

	var response []*api.GetListDataTablesResponse_DataTable
	for _, dataTable := range queryResult.DataTables {
		sources := enrichedSources[dataTable.ID]
		destinations := enrichedDestinations[dataTable.ID]
		response = append(response, &api.GetListDataTablesResponse_DataTable{
			Id:        dataTable.ID,
			Name:      dataTable.Name,
			CreatedAt: dataTable.CreatedAt.String(),
			UpdatedAt: dataTable.UpdatedAt.String(),
			DataSources: utils.Map(sources, func(dataSource model.DataSource) *api.EnrichedDataSource {
				return &api.EnrichedDataSource{
					Id:   dataSource.ID,
					Name: dataSource.Name,
					Type: string(dataSource.Type),
				}
			}),
			DataDestinations: utils.Map(destinations, func(dataDestination model.DataDestination) *api.EnrichedDataDestination {
				return &api.EnrichedDataDestination{
					Id:   dataDestination.ID,
					Name: dataDestination.Name,
					Type: string(dataDestination.Type),
				}
			}),
		})
	}

	return response, queryResult.Count, nil
}

func (b business) GetDataTable(ctx context.Context, request *api.GetDataTableRequest, accountUuid string) (*api.GetDataTableResponse, error) {
	dataTable, err := b.repository.DataTableRepository.GetDataTable(ctx, request.Id)
	if err != nil {
		b.log.WithName("GetDataTable").
			WithValues("Context", ctx).
			Error(err, "Cannot get data_table")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Not exists data table")
		}
		return nil, err
	}
	if dataTable.AccountUuid != uuid.MustParse(accountUuid) {
		b.log.WithName("GetDataTable").
			WithValues("Context", ctx).
			Info("Only owner can get data_table")
		return nil, status.Error(codes.PermissionDenied, "Only owner can get data_table")
	}
	var schema []*api.SchemaColumn
	if dataTable.Schema.RawMessage != nil {
		err = json.Unmarshal(dataTable.Schema.RawMessage, &schema)
		if err != nil {
			return nil, err
		}
	}

	return &api.GetDataTableResponse{
		Code:      int32(code.Code_OK),
		Id:        dataTable.ID,
		Name:      dataTable.Name,
		CreatedAt: dataTable.CreatedAt.String(),
		UpdatedAt: dataTable.UpdatedAt.String(),
		Schema:    schema,
	}, nil
}

func (b business) GetListFileExportRecords(ctx context.Context, request *api.GetListFileExportRecordsRequest, accountUuid string) ([]*api.GetListFileExportRecordsResponse_FileExportRecord, error) {
	records, err := b.repository.FileExportRecordRepository.ListFileExportRecords(ctx, request.Id, accountUuid)
	if err != nil {
		b.log.WithName("GetListFileExportRecords").Error(err, "cannot get file export records data")
		return nil, err
	}

	returnRecords := make([]*api.GetListFileExportRecordsResponse_FileExportRecord, 0, len(records))
	for _, record := range records {
		returnRecords = append(returnRecords, &api.GetListFileExportRecordsResponse_FileExportRecord{
			Id:             record.ID,
			DataTableId:    record.DataTableId,
			Format:         string(record.Format),
			Status:         record.Status,
			DownloadUrl:    record.DownloadUrl,
			ExpirationTime: utils.ToTimeString(record.ExpirationTime),
			CreatedAt:      utils.ToTimeString(record.CreatedAt),
		})
	}

	return returnRecords, nil
}

func (b business) GetQueryDataTable(ctx context.Context, request *api.GetQueryDataTableRequest, accountUuid string) (*api.GetQueryDataTableResponse, error) {
	dataTable, err := b.repository.DataTableRepository.GetDataTable(ctx, request.GetId())
	if err != nil {
		b.log.WithName("GetQueryDataTable").Error(err, "Table not found", "id", request.Id)
		return nil, err
	}
	if dataTable.AccountUuid != uuid.MustParse(accountUuid) {
		b.log.WithName("GetQueryDataTable").
			WithValues("Context", ctx).
			Info("Only owner can access data_table")
		return nil, status.Error(codes.PermissionDenied, "Only owner can access data_table")
	}

	res, err := b.queryAdapter.GetDataTableV2(ctx, &query.GetQueryDataTableV2Request{
		Limit:     request.Limit,
		TablePath: utils.GenerateDeltaTablePath(accountUuid, dataTable.Name),
	})

	if err != nil {
		b.log.WithName("GetQueryDataTable").Error(err, "cannot get query table", "id", dataTable.ID)
		return nil, err
	}

	return &api.GetQueryDataTableResponse{Code: int32(codes.OK), Count: res.Count, Data: res.Data}, nil
}
