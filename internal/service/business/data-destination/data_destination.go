package data_destination

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
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
