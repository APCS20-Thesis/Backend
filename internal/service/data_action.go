package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
)

func (s *Service) GetListDataActions(ctx context.Context, request *api.GetListDataActionsRequest) (*api.GetListDataActionsResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListFileExportRecords").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	return s.business.DataTableBusiness.ProcessGetListDataActions(ctx, request, accountUuid)
}
