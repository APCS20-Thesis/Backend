package data_destination

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/gophish"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
)

type Business interface {
	CreateGophishUserGroupFromSegment(ctx context.Context, accountUuid string, request *api.CreateGophishUserGroupFromSegmentRequest) error
}

type business struct {
	log            logr.Logger
	repository     *repository.Repository
	gophishAdapter gophish.GophishAdapter
	queryAdapter   query.QueryAdapter
}

func NewDataDestinationBusiness(log logr.Logger, repository *repository.Repository, gophishAdapter gophish.GophishAdapter, queryAdapter query.QueryAdapter) Business {
	return &business{
		log:            log.WithName("DataDestinationBiz"),
		repository:     repository,
		gophishAdapter: gophishAdapter,
		queryAdapter:   queryAdapter,
	}
}
