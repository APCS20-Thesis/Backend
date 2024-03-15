package data_source

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (b business) CreateSourceConnection(ctx context.Context, params *repository.CreateSourceConnectionParams) (*model.SourceConnection, error) {
	SourceConnection, err := b.repository.SourceConnectionRepository.CreateSourceConnection(ctx, params)
	if err != nil {
		b.log.WithName("CreateSourceConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot create data_table")
		return nil, err
	}
	return SourceConnection, nil
}

func (b business) UpdateSourceConnection(ctx context.Context, params *repository.UpdateSourceConnectionParams) error {
	err := b.repository.SourceConnectionRepository.UpdateSourceConnection(ctx, params)
	if err != nil {
		b.log.WithName("UpdateSourceConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot update data_table")
		return err
	}
	return nil
}

func (b business) GetListSourceConnections(ctx context.Context, request *api.GetListSourceConnectionsRequest, accountUuid string) ([]*api.SourceConnectionBase, error) {
	sourceConnections, err := b.repository.SourceConnectionRepository.ListSourceConnections(ctx,
		&repository.FilterSourceConnection{
			Name:        request.Name,
			AccountUuid: uuid.MustParse(accountUuid),
		})
	if err != nil {
		b.log.WithName("GetListSourceConnections").
			WithValues("Context", ctx).
			Error(err, "Cannot get list data_Tables")
		return nil, err
	}
	var response []*api.SourceConnectionBase
	for _, SourceConnection := range sourceConnections {
		response = append(response, &api.SourceConnectionBase{
			Id:        SourceConnection.ID,
			Name:      SourceConnection.Name,
			UpdatedAt: SourceConnection.UpdatedAt.String(),
		})
	}
	return response, nil
}

func (b business) GetSourceConnection(ctx context.Context, request *api.GetSourceConnectionRequest, accountUuid string) (*api.GetSourceConnectionResponse, error) {
	sourceConnection, err := b.repository.SourceConnectionRepository.GetSourceConnection(ctx, request.Id)
	if err != nil {
		b.log.WithName("GetSourceConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get data_table")
		return nil, err
	}
	if sourceConnection.AccountUuid != uuid.MustParse(accountUuid) {
		b.log.WithName("GetSourceConnection").
			WithValues("Context", ctx).
			Error(status.Error(codes.Code(code.Code_PERMISSION_DENIED), "Only owner can get data_table"),
				"Only owner can get data_table")
		return nil, status.Error(codes.Code(code.Code_PERMISSION_DENIED), "Only owner can get data_table")
	}
	var configurations map[string]string
	err = json.Unmarshal(sourceConnection.Configurations.RawMessage, &configurations)
	if err != nil {
		return nil, err
	}

	return &api.GetSourceConnectionResponse{
		Code:           int32(code.Code_OK),
		Id:             sourceConnection.ID,
		Name:           sourceConnection.Name,
		CreatedAt:      sourceConnection.CreatedAt.String(),
		UpdatedAt:      sourceConnection.UpdatedAt.String(),
		Configurations: configurations,
	}, nil
}
