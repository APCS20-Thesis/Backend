package auth

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
)

type Business interface {
	ProcessLogin(ctx context.Context, request *api.LoginRequest) (*api.Account, error)
	ProcessSignUp(ctx context.Context, request *api.SignUpRequest) (*api.CommonResponse, error)
	ProcessGetAccountInfo(ctx context.Context) (*api.Account, error)
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
