package auth

import (
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
)

type Business interface {
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
