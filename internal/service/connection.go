package service

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/grpc/codes"
)

func (s *Service) GetListConnections(ctx context.Context, request *api.GetListConnectionsRequest) (*api.GetListConnectionsResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListConnections").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	connections, count, err := s.business.ConnectionBusiness.GetListConnections(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetListConnections").
			WithValues("Context", ctx).
			Error(err, "Failed to process list connections")
		return nil, err
	}
	return &api.GetListConnectionsResponse{Code: int32(codes.OK), Count: count, Results: connections}, nil
}

func (s *Service) GetConnection(ctx context.Context, request *api.GetConnectionRequest) (*api.GetConnectionResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	response, err := s.business.ConnectionBusiness.GetConnection(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Error(err, "Failed to process get connection")
		return nil, err
	}
	return response, nil
}

func (s *Service) CreateConnection(ctx context.Context, request *api.CreateConnectionRequest) (*api.CreateConnectionResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("CreateConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}

	_, err = s.business.ConnectionBusiness.CreateConnection(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("CreateConnection").
			WithValues("Context", ctx).
			Error(err, "Failed to process create connection")
		return nil, err
	}
	return &api.CreateConnectionResponse{Message: "Create Success", Code: int32(codes.OK)}, nil
}

func (s *Service) UpdateConnection(ctx context.Context, request *api.UpdateConnectionRequest) (*api.UpdateConnectionResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("UpdateConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	configurations, err := json.Marshal(request.Configurations)
	if err != nil {
		s.log.WithName("UpdateConnection").
			WithValues("Configuration", request.Configurations).
			Error(err, "Cannot parse configuration to JSON")
		return nil, err
	}
	err = s.business.ConnectionBusiness.UpdateConnection(ctx, &repository.UpdateConnectionParams{
		ID:             request.Id,
		Name:           request.Name,
		Configurations: pqtype.NullRawMessage{RawMessage: configurations, Valid: true},
		AccountUuid:    uuid.MustParse(accountUuid),
	})
	if err != nil {
		s.log.WithName("UpdateConnection").Error(err, "Failed to process update connection")
		return nil, err
	}
	return &api.UpdateConnectionResponse{Message: "Update Success", Code: int32(codes.OK)}, nil
}

func (s *Service) DeleteConnection(ctx context.Context, request *api.DeleteConnectionRequest) (*api.DeleteConnectionResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	err = s.business.ConnectionBusiness.DeleteConnection(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Error(err, "Failed to process get connection")
		return nil, err
	}
	return &api.DeleteConnectionResponse{Message: "Delete Success", Code: int32(codes.OK)}, nil
}
