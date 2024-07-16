package data_action

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
)

func (b business) ProcessGetListDataActions(ctx context.Context, request *api.GetListDataActionsRequest, accountUuid string) (*api.GetListDataActionsResponse, error) {
	queryDataActions, err := b.repository.DataActionRepository.GetListDataActions(ctx, &repository.GetListDataActionsParams{
		ActionTypes: request.Types,
		AccountUuid: uuid.MustParse(accountUuid),
		Page:        int(request.Page),
		PageSize:    int(request.PageSize),
		TargetTable: model.DataActionTargetTable(request.TargetTable),
		ObjectId:    request.ObjectId,
	})
	if err != nil {
		b.log.WithName("list data actions").Error(err, "cannot get data actions")
		return nil, err
	}

	return &api.GetListDataActionsResponse{
		Code:    0,
		Message: "Success",
		Count:   50,
		Results: utils.Map(queryDataActions.DataActions, func(modelDataAction model.DataAction) *api.DataAction {
			return &api.DataAction{
				Id:         modelDataAction.ID,
				ActionType: string(modelDataAction.ActionType),
				Status:     string(modelDataAction.Status),
				CreatedAt:  modelDataAction.CreatedAt.String(),
				UpdatedAt:  modelDataAction.UpdatedAt.String(),
			}
		}),
	}, nil

}
