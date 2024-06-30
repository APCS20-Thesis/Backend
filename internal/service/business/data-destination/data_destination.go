package data_destination

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (b business) ProcessGetListDataDestinations(ctx context.Context, request *api.GetListDataDestinationsRequest, accountUuid string) (*api.GetListDataDestinationsResponse, error) {
	queryResult, err := b.repository.DataDestinationRepository.ListDataDestinations(ctx, &repository.ListDataDestinationsParams{
		Page:        int(request.Page),
		PageSize:    int(request.PageSize),
		Type:        model.DataDestinationType(request.Type),
		AccountUuid: accountUuid,
	})
	if err != nil {
		b.log.WithName("ProcessGetListDataDestinations").Error(err, "cannot get list data destinations", "request", request)
		return nil, err
	}

	return &api.GetListDataDestinationsResponse{
		Code:    0,
		Message: "Success",
		Count:   queryResult.Count,
		Results: utils.Map(queryResult.Destinations, func(modelDest model.DataDestination) *api.DataDestination {
			return &api.DataDestination{
				Id:        modelDest.ID,
				Name:      modelDest.Name,
				Type:      string(modelDest.Type),
				CreatedAt: modelDest.CreatedAt.String(),
				UpdatedAt: modelDest.UpdatedAt.String(),
			}
		}),
	}, nil
}

func (b business) ProcessGetDataDestinationDetail(ctx context.Context, request *api.GetDataDestinationDetailRequest, accountUuid string) (*api.GetDataDestinationDetailResponse, error) {
	logger := b.log.WithName("ProcessGetDataDestinationDetail").WithValues("id", request.Id)
	destination, err := b.repository.DataDestinationRepository.GetDataDestination(ctx, request.Id)
	if err != nil {
		logger.Error(err, "cannot query data destination")
		return nil, err
	}
	if destination.AccountUuid.String() != accountUuid {
		return nil, status.Error(codes.PermissionDenied, "cannot access to this data destination")
	}

	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, destination.ConnectionId)
	if err != nil {
		logger.Error(err, "cannot query connection")
		return nil, err
	}

	return &api.GetDataDestinationDetailResponse{
		Code:    int32(code.Code_OK),
		Message: code.Code_OK.String(),
		Id:      destination.ID,
		Name:    destination.Name,
		Type:    string(destination.Type),
		Connection: &api.EnrichedConnection{
			Id:   connection.ID,
			Name: connection.Name,
			Type: string(connection.Type),
		},
		Configurations: string(destination.Configurations.RawMessage),
		CreatedAt:      destination.CreatedAt.String(),
		UpdatedAt:      destination.UpdatedAt.String(),
	}, nil
}
