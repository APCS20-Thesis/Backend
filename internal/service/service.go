package service

import (
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/service/business"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type Service struct {
	log        logr.Logger
	config     *config.Config
	jwtManager *JWTManager
	s3Manger   *S3Manager
	//// more connector here
	//store  store.StoreQuerier

	business *business.Business

	// embedded unimplemented service server
	api.UnimplementedCDPServiceServer
	api.UnimplementedCDPServiceFile
}

func NewService(logger logr.Logger, config *config.Config, gormDb *gorm.DB, jwtManager *JWTManager) (*Service, error) {
	airflowAdapter, err := airflow.NewAirflowAdapter(logger, config.AirflowAdapterConfig.Address, config.AirflowAdapterConfig.Username, config.AirflowAdapterConfig.Password)
	if err != nil {
		return nil, err
	}

	business := business.NewBusiness(logger, gormDb, airflowAdapter)

	s3Manager := NewS3Manager(
		config.S3StorageConfig.Region,
		config.S3StorageConfig.AccessKeyID,
		config.S3StorageConfig.SecretAccessKey,
	)
	return &Service{
		log:        logger,
		config:     config,
		jwtManager: jwtManager,
		business:   business,
		s3Manger:   s3Manager,
	}, nil
}
