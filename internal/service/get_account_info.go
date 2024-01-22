package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/constants"
	"google.golang.org/grpc/codes"
)

func (s *Service) GetAccountInfo(ctx context.Context, request *api.GetAccountInfoRequest) (*api.GetAccountInfoResponse, error) {
	accountUuid, err := GetMetadata(ctx, constants.KeyAccountUuid)
	if err != nil {
		s.log.WithName("GetAccountInfo").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	account, err := s.business.AuthBusiness.ProcessGetAccountInfo(ctx, accountUuid)
	if err != nil {
		return nil, err
	}
	return &api.GetAccountInfoResponse{
		Code:    int32(codes.OK),
		Message: "Get account info success",
		Account: account,
	}, nil
}
