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

func (b business) CreateConnection(ctx context.Context, params *repository.CreateConnectionParams) (*model.Connection, error) {
	connection, err := b.repository.ConnectionRepository.CreateConnection(ctx, params)
	if err != nil {
		b.log.WithName("CreateConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot create connection")
		return nil, err
	}
	return connection, nil
}

func (b business) UpdateConnection(ctx context.Context, params *repository.UpdateConnectionParams) error {
	err := b.repository.ConnectionRepository.UpdateConnection(ctx, params)
	if err != nil {
		b.log.WithName("UpdateConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot update connection")
		return err
	}
	return nil
}

func (b business) GetListConnections(ctx context.Context, request *api.GetListConnectionsRequest, accountUuid string) ([]*api.GetListConnectionsResponse_Connection, error) {
	Connections, err := b.repository.ConnectionRepository.ListConnections(ctx,
		&repository.FilterConnection{
			Name:        request.Name,
			AccountUuid: uuid.MustParse(accountUuid),
		})
	if err != nil {
		b.log.WithName("GetListConnections").
			WithValues("Context", ctx).
			Error(err, "Cannot get list connection")
		return nil, err
	}
	var response []*api.GetListConnectionsResponse_Connection
	for _, connection := range Connections {
		response = append(response, &api.GetListConnectionsResponse_Connection{
			Id:        connection.ID,
			Name:      connection.Name,
			UpdatedAt: connection.UpdatedAt.String(),
		})
	}
	return response, nil
}

func (b business) GetConnection(ctx context.Context, request *api.GetConnectionRequest, accountUuid string) (*api.GetConnectionResponse, error) {
	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.Id)
	if err != nil {
		b.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get data_table")
		return nil, err
	}
	if connection.AccountUuid != uuid.MustParse(accountUuid) {
		b.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Error(status.Error(codes.Code(code.Code_PERMISSION_DENIED), "Only owner can get data_table"),
				"Only owner can get connection")
		return nil, status.Error(codes.Code(code.Code_PERMISSION_DENIED), "Only owner can get data_table")
	}
	var configurations map[string]string
	err = json.Unmarshal(connection.Configurations.RawMessage, &configurations)
	if err != nil {
		return nil, err
	}

	return &api.GetConnectionResponse{
		Code:           int32(code.Code_OK),
		Id:             connection.ID,
		Name:           connection.Name,
		CreatedAt:      connection.CreatedAt.String(),
		UpdatedAt:      connection.UpdatedAt.String(),
		Configurations: configurations,
	}, nil
}
