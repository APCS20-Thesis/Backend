package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"google.golang.org/grpc/codes"
)

func (s *Service) CheckHealth(ctx context.Context, request *api.CheckHealthRequest) (*api.CommonResponse, error) {

	return &api.CommonResponse{
		Code:    int32(codes.OK),
		Message: "Success",
	}, nil
}
