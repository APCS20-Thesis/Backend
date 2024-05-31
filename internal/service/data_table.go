package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"google.golang.org/grpc/codes"
)

func (s *Service) GetListDataTables(ctx context.Context, request *api.GetListDataTablesRequest) (*api.GetListDataTablesResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListDataTables").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	dataTables, err := s.business.DataTableBusiness.GetListDataTables(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetListDataTables").
			WithValues("Context", ctx).
			Error(err, "Failed to process list data tables")
		return nil, err
	}
	return &api.GetListDataTablesResponse{Code: int32(codes.OK), Count: int64(len(dataTables)), Results: dataTables}, nil
}

func (s *Service) GetDataTable(ctx context.Context, request *api.GetDataTableRequest) (*api.GetDataTableResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetDataTable").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	response, err := s.business.DataTableBusiness.GetDataTable(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetDataTable").
			WithValues("Context", ctx).
			Error(err, "Failed to process get data table")
		return nil, err
	}
	return response, nil
}

func (s *Service) GetQueryDataTable(ctx context.Context, request *api.GetQueryDataTableRequest) (*api.GetQueryDataTableResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetQueryDataTable").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	response, err := s.business.DataTableBusiness.GetQueryDataTable(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetDataTable").
			WithValues("Context", ctx).
			Error(err, "Failed to process get data table")
		return nil, err
	}
	return response, nil
}
