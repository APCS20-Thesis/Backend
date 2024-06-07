package business

import (
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/adapter/gophish"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/internal/service/business/auth"
	data_destination "github.com/APCS20-Thesis/Backend/internal/service/business/data-destination"
	data_source "github.com/APCS20-Thesis/Backend/internal/service/business/data-source"
	data_table "github.com/APCS20-Thesis/Backend/internal/service/business/data-table"
	"github.com/APCS20-Thesis/Backend/internal/service/business/segment"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type Business struct {
	db                      *gorm.DB
	repository              *repository.Repository
	AuthBusiness            auth.Business
	DataSourceBusiness      data_source.Business
	DataTableBusiness       data_table.Business
	SegmentBusiness         segment.Business
	DataDestinationBusiness data_destination.Business
}

func NewBusiness(
	log logr.Logger,
	db *gorm.DB,
	airflowAdapter airflow.AirflowAdapter,
	config *config.Config,
	queryAdapter query.QueryAdapter,
	gophishAdapter gophish.GophishAdapter,
) *Business {
	repo := repository.NewRepository(db)
	return &Business{
		db:                      db,
		repository:              repo,
		AuthBusiness:            auth.NewAuthBusiness(log, repo),
		DataSourceBusiness:      data_source.NewDataSourceBusiness(db, log, repo, airflowAdapter, config),
		DataTableBusiness:       data_table.NewDataTableBusiness(log, repo, airflowAdapter, queryAdapter),
		SegmentBusiness:         segment.NewSegmentBusiness(log, repo, airflowAdapter),
		DataDestinationBusiness: data_destination.NewDataDestinationBusiness(db, log, repo, gophishAdapter, queryAdapter, airflowAdapter),
	}
}
