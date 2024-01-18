package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Service) GetInfo(ctx context.Context, request *api.GetInfoRequest) (*api.Account, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := s.jwtManager.Verify(accessToken)
	if err != nil {

		return nil, status.Errorf(codes.Unauthenticated, "Access token is invalid: %v", err)
	}
	return s.business.AuthBusiness.ProcessGetInfo(ctx, claims.Username)
}
