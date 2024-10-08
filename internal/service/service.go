package service

import (
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/adapter/alert"
	"github.com/APCS20-Thesis/Backend/internal/adapter/gophish"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/service/business"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type Service struct {
	log          logr.Logger
	config       *config.Config
	jwtManager   *JWTManager
	s3Manger     *S3Manager
	alertAdapter alert.AlertAdapter
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
	queryAdapter, err := query.NewQueryAdapter(logger, config.QueryAdapterConfig.Address)
	if err != nil {
		return nil, err
	}
	gophishAdapter, err := gophish.NewGophishAdapter(logger)
	if err != nil {
		return nil, err
	}
	alertAdapter, err := alert.NewAlertAdapter(logger, config.AlertAdapterConfig.Webhook)
	if err != nil {
		return nil, err
	}

	business := business.NewBusiness(logger, gormDb, airflowAdapter, config, queryAdapter, gophishAdapter)

	s3Manager := NewS3Manager(
		config.S3StorageConfig.Region,
		config.S3StorageConfig.AccessKeyID,
		config.S3StorageConfig.SecretAccessKey,
	)
	return &Service{
		log:          logger,
		config:       config,
		jwtManager:   jwtManager,
		s3Manger:     s3Manager,
		business:     business,
		alertAdapter: alertAdapter,
	}, nil
}
