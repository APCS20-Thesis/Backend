package connection

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
)

type Business interface {
	CreateConnection(ctx context.Context, request *api.CreateConnectionRequest, accountUuid string) (*model.Connection, error)
	UpdateConnection(ctx context.Context, params *repository.UpdateConnectionParams) error
	GetListConnections(ctx context.Context, request *api.GetListConnectionsRequest, accountUuid string) ([]*api.GetListConnectionsResponse_Connection, int64, error)
	GetConnection(ctx context.Context, request *api.GetConnectionRequest, accountUuid string) (*api.GetConnectionResponse, error)
	DeleteConnection(ctx context.Context, request *api.DeleteConnectionRequest, accountUuid string) error
}

type business struct {
	log        logr.Logger
	repository *repository.Repository
}

func NewConnectionBusiness(log logr.Logger, repository *repository.Repository) Business {
	return &business{
		log:        log.WithName("ConnectionBiz"),
		repository: repository,
	}
}
