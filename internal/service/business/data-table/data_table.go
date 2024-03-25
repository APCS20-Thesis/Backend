package data_table

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (b business) GetListDataTables(ctx context.Context, request *api.GetListDataTablesRequest, accountUuid string) ([]*api.GetListDataTablesResponse_DataTable, error) {
	dataTables, err := b.repository.DataTableRepository.ListDataTables(ctx,
		&repository.ListDataTablesFilters{
			Name:        request.Name,
			AccountUuid: uuid.MustParse(accountUuid),
		})
	if err != nil {
		b.log.WithName("GetListDataTables").
			WithValues("Context", ctx).
			Error(err, "Cannot get list data_Tables")
		return nil, err
	}
	var response []*api.GetListDataTablesResponse_DataTable
	for _, dataTable := range dataTables {
		response = append(response, &api.GetListDataTablesResponse_DataTable{
			Id:        dataTable.ID,
			Name:      dataTable.Name,
			CreatedAt: dataTable.CreatedAt.String(),
			UpdatedAt: dataTable.UpdatedAt.String(),
		})
	}
	return response, nil
}

func (b business) GetDataTable(ctx context.Context, request *api.GetDataTableRequest, accountUuid string) (*api.GetDataTableResponse, error) {
	dataTable, err := b.repository.DataTableRepository.GetDataTable(ctx, request.Id)
	if err != nil {
		b.log.WithName("GetDataTable").
			WithValues("Context", ctx).
			Error(err, "Cannot get data_table")
		return nil, err
	}
	if dataTable.AccountUuid != uuid.MustParse(accountUuid) {
		b.log.WithName("GetDataTable").
			WithValues("Context", ctx).
			Info("Only owner can get data_table")
		return nil, status.Error(codes.PermissionDenied, "Only owner can get data_table")
	}
	var schema []*api.GetDataTableResponse_Field
	err = json.Unmarshal(dataTable.Schema.RawMessage, &schema)
	if err != nil {
		return nil, err
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
