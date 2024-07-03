package auth

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
)

type Business interface {
	ProcessLogin(ctx context.Context, request *api.LoginRequest) (*model.Account, error)
	ProcessSignUp(ctx context.Context, request *api.SignUpRequest) (*api.CommonResponse, error)
	ProcessGetAccountInfo(ctx context.Context, accountUuid string) (*api.Account, *api.Setting, error)
	ProcessUpdateAccountInfo(ctx context.Context, request *api.UpdateAccountInfoRequest, accountUuid string) (*api.Account, error)
	ProcessUpdateAccountSetting(ctx context.Context, request *api.UpdateAccountSettingRequest, accountUuid string) (*api.Setting, error)
}

type business struct {
	log        logr.Logger
	repository *repository.Repository
}

func NewAuthBusiness(log logr.Logger, repository *repository.Repository) Business {
	return &business{
		log:        log.WithName("AuthBiz"),
		repository: repository,
	}
}
