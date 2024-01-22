package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
)

func (s *Service) SignUp(ctx context.Context, request *api.SignUpRequest) (*api.CommonResponse, error) {
	return s.business.AuthBusiness.ProcessSignUp(ctx, request)
}
