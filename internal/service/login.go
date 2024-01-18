package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/APCS20-Thesis/Backend/api"
)

func (s *Service) Login(ctx context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {
	account, err := s.business.AuthBusiness.ProcessLogin(ctx, request)
	if err != nil {
		s.log.WithName("Login").WithValues("request", request).Error(err, "Cannot find user")
		return nil, status.Errorf(codes.Internal, "Cannot find user")
	}
	token, err := s.jwtManager.Generate(account, "user")
	if err != nil {
		s.log.WithName("Login").WithValues("request", request).Error(err, "Cannot generate access token")
		return nil, status.Errorf(codes.Internal, "Cannot generate access token")
	}

	return &api.LoginResponse{
		Code:        int32(codes.OK),
		Message:     "Login success",
		AccessToken: token,
		Account:     account,
	}, nil
}
