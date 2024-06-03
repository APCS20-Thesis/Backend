package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
)

func (s *Service) CreateGophishUserGroupFromSegment(ctx context.Context, request *api.CreateGophishUserGroupFromSegmentRequest) (*api.CreateGophishUserGroupFromSegmentResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListDataTables").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
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
