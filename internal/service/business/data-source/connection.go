package data_source

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"golang.org/x/exp/slices"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strconv"
)

func (b business) CreateConnection(ctx context.Context, request *api.CreateConnectionRequest, accountUuid string) (*model.Connection, error) {
	if !slices.Contains(model.ConnectionTypes, model.ConnectionType(request.Type)) {
		return nil, status.Error(codes.InvalidArgument, "Invalid connection type")
	}

	configurations, err := json.Marshal(request.Configurations)
	if err != nil {
		b.log.WithName("CreateConnection").
			WithValues("Configuration", request.Configurations).
			Error(err, "Cannot parse configuration to JSON")
		return nil, err
	}

	connection, err := b.repository.ConnectionRepository.CreateConnection(ctx, &repository.CreateConnectionParams{
		Name:           request.Name,
		Type:           model.ConnectionType(request.Type),
		Configurations: pqtype.NullRawMessage{RawMessage: configurations, Valid: configurations != nil},
		AccountUuid:    uuid.MustParse(accountUuid),
	})
	if err != nil {
		b.log.WithName("CreateConnection").Error(err, "Cannot create connection")
		return nil, err
	}
	return connection, nil
}

func (b business) UpdateConnection(ctx context.Context, params *repository.UpdateConnectionParams) error {
	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, params.ID)
	if err != nil {
		b.log.WithName("UpdateConnection").
			WithValues("ConnectionID", params.ID).
			Error(err, "Can not get record with id "+strconv.FormatInt(params.ID, 10))
		return err
	}
	if connection == nil {
		b.log.WithName("UpdateConnection").
			WithValues("ConnectionID", params.ID).
			Error(err, "No record with id "+strconv.FormatInt(params.ID, 10))
		return gorm.ErrRecordNotFound
	}
	if connection.AccountUuid != params.AccountUuid {
		b.log.WithName("UpdateConnection").
			WithValues("ConnectionID", params.ID).
			Info("Only owner can get connection")
		return status.Error(codes.PermissionDenied, "Only owner can update connection")
	}
	err = b.repository.ConnectionRepository.UpdateConnection(ctx, params)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Not found connection with id "+strconv.FormatInt(request.Id, 10))
		}
		b.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get connection")
		return nil, err
	}
	if connection.AccountUuid != uuid.MustParse(accountUuid) {
		b.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Info("Only owner can get connection")
		return nil, status.Error(codes.PermissionDenied, "Only owner can get connection")
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
func (b business) DeleteConnection(ctx context.Context, request *api.DeleteConnectionRequest, accountUuid string) error {
	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.Id)
	if err != nil {
		b.log.WithName("DeleteConnection").
			WithValues("ConnectionID", request.Id).
			Error(err, "Can not get record with id "+strconv.FormatInt(request.Id, 10))
		return err
	}
	if connection == nil {
		b.log.WithName("DeleteConnection").
			WithValues("ConnectionID", request.Id).
			Error(err, "No record with id "+strconv.FormatInt(request.Id, 10))
		return gorm.ErrRecordNotFound
	}
	if connection.AccountUuid != uuid.MustParse(accountUuid) {
		b.log.WithName("DeleteConnection").
			WithValues("ConnectionID", request.Id).
			Info("Only owner can get connection")
		return status.Error(codes.PermissionDenied, "Only owner can delte connection")
	}
	err = b.repository.ConnectionRepository.DeleteConnection(ctx, request.Id)
	if err != nil {
		b.log.WithName("DeleteConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot delete connection")
		return err
	}
	return nil
}
