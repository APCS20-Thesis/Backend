package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"google.golang.org/grpc/codes"
)

func (s *Service) GetAccountInfo(ctx context.Context, request *api.GetAccountInfoRequest) (*api.GetAccountInfoResponse, error) {
	account, err := s.business.AuthBusiness.ProcessGetAccountInfo(ctx)
	if err != nil {
		return nil, err
	}
	return &api.GetAccountInfoResponse{
		Code:    int32(codes.OK),
		Message: "Get account info success",
		Account: account,
	}, nil
}
