package business

import (
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/internal/service/business/auth"
	data_source "github.com/APCS20-Thesis/Backend/internal/service/business/data-source"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type Business struct {
	db                 *gorm.DB
	repository         *repository.Repository
	AuthBusiness       auth.Business
	DataSourceBusiness data_source.Business
}

func NewBusiness(
	log logr.Logger,
	db *gorm.DB,
	airflowAdapter airflow.AirflowAdapter,
) *Business {
	repo := repository.NewRepository(db)
	return &Business{
		db:                 db,
		repository:         repo,
		AuthBusiness:       auth.NewAuthBusiness(log, repo),
		DataSourceBusiness: data_source.NewDataSourceBusiness(log, repo, airflowAdapter),
	}
}
