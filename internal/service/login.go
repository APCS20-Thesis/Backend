package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"google.golang.org/grpc/codes"
)

func (s *Service) Login(ctx context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {
	account, err := s.business.AuthBusiness.ProcessLogin(ctx, request)
	if err != nil {
		return nil, err
	}
	token, err := s.jwtManager.Generate(account, "user")
	if err != nil {
		s.log.WithName("Login").WithValues("request", request).Error(err, "Cannot generate access token")
		return nil, err
	}

	return &api.LoginResponse{
		Code:        int32(codes.OK),
		Message:     "Login success",
		AccessToken: token,
		Account:     account,
	}, nil
}
