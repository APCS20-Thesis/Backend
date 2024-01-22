package business

import (
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/internal/service/business/auth"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type Business struct {
	db           *gorm.DB
	repository   *repository.Repository
	AuthBusiness auth.Business
}

func NewBusiness(
	log logr.Logger,
	db *gorm.DB,
) *Business {
	repo := repository.NewRepository(db)
	return &Business{
		db:           db,
		repository:   repo,
		AuthBusiness: auth.NewAuthBusiness(log, repo),
	}
}
