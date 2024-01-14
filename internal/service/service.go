package service

import (
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/go-logr/logr"
)

type Service struct {
	log logr.Logger
	//// more connector here
	//store  store.StoreQuerier
	//gormDb *gorm.DB

	//biz *business.Business

	// embedded unimplemented service server
	api.UnimplementedCDPServiceServer
}

func NewService(logger logr.Logger) (*Service, error) {
	return &Service{
		log: logger,
	}, nil
}
