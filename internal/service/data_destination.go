package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
)

func (s *Service) CreateGophishUserGroupFromSegment(ctx context.Context, request *api.CreateGophishUserGroupFromSegmentRequest) (*api.CreateGophishUserGroupFromSegmentResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListDataTables").Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	err = s.business.DataDestinationBusiness.CreateGophishUserGroupFromSegment(ctx, accountUuid, request)
	if err != nil {
		return nil, err
	}

	return &api.CreateGophishUserGroupFromSegmentResponse{
		Code:    0,
		Message: "Success",
	}, nil
}

func (s *Service) ExportToMySQLDestination(ctx context.Context, request *api.ExportToMySQLDestinationRequest) (*api.ExportToMySQLDestinationResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("ExportToMySQLDestination").Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	err = s.business.DataDestinationBusiness.ProcessExportToMySQLDestination(ctx, request, accountUuid)
	if err != nil {
		return nil, err
	}

	return &api.ExportToMySQLDestinationResponse{
		Code:    0,
		Message: "Success",
	}, nil
}

func (s *Service) GetListDataDestinations(ctx context.Context, request *api.GetListDataDestinationsRequest) (*api.GetListDataDestinationsResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListDataDestinations").Error(err, "Cannot get account_uuid from context")
		return nil, err
	}

	return s.business.DataDestinationBusiness.ProcessGetListDataDestinations(ctx, request, accountUuid)
}

func (s *Service) GetDataDestinationDetail(ctx context.Context, request *api.GetDataDestinationDetailRequest) (*api.GetDataDestinationDetailResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListDataDestinations").Error(err, "Cannot get account_uuid from context")
		return nil, err
	}

	return s.business.DataDestinationBusiness.ProcessGetDataDestinationDetail(ctx, request, accountUuid)
}
