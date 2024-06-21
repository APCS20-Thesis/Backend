package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/constants"
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
	}, nil
}

func (s *Service) SignUp(ctx context.Context, request *api.SignUpRequest) (*api.CommonResponse, error) {
	return s.business.AuthBusiness.ProcessSignUp(ctx, request)
}

func (s *Service) GetAccountInfo(ctx context.Context, request *api.GetAccountInfoRequest) (*api.GetAccountInfoResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetAccountInfo").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	account, setting, err := s.business.AuthBusiness.ProcessGetAccountInfo(ctx, accountUuid)
	if err != nil {
		return nil, err
	}

	return &api.GetAccountInfoResponse{
		Code:    int32(codes.OK),
		Message: "Get account info success",
		Account: account,
		Setting: setting,
	}, nil
}

func (s *Service) UpdateAccountInfo(ctx context.Context, request *api.UpdateAccountInfoRequest) (*api.UpdateAccountInfoResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetAccountInfo").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	account, err := s.business.AuthBusiness.ProcessUpdateAccountInfo(ctx, request, accountUuid)
	if err != nil {
		return nil, err
	}
	return &api.UpdateAccountInfoResponse{
		Code:    int32(codes.OK),
		Message: "Get account info success",
		Account: account,
	}, nil
}

func (s *Service) UpdateAccountSetting(ctx context.Context, request *api.UpdateAccountSettingRequest) (*api.UpdateAccountSettingResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetAccountInfo").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	setting, err := s.business.AuthBusiness.ProcessUpdateAccountSetting(ctx, request, accountUuid)
	if err != nil {
		return nil, err
	}
	return &api.UpdateAccountSettingResponse{
		Code:    int32(codes.OK),
		Message: "Get account info success",
		Setting: setting,
	}, nil
}

/*
COMMON AUTH FUNCTIONS
*/

func GetAccountUuidFromCtx(ctx context.Context) (string, error) {
	return GetMetadata(ctx, constants.KeyAccountUuid)
}
