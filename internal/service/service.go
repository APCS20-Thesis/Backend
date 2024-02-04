package service

import (
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/service/business"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type Service struct {
	log        logr.Logger
	config     *config.Config
	jwtManager *JWTManager

	//// more connector here
	//store  store.StoreQuerier

	business *business.Business

	// embedded unimplemented service server
	api.UnimplementedCDPServiceServer
	api.UnimplementedCDPServiceFile
}

func NewService(logger logr.Logger, config *config.Config, gormDb *gorm.DB, jwtManager *JWTManager) (*Service, error) {
	business := business.NewBusiness(logger, gormDb)

	return &Service{
		log:        logger,
		config:     config,
		jwtManager: jwtManager,
		business:   business,
	}, nil
}
