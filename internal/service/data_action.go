package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"google.golang.org/genproto/googleapis/rpc/code"
)

func (s *Service) GetListDataActionRuns(ctx context.Context, request *api.GetListDataActionRunsRequest) (*api.GetListDataActionRunsResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListDataActions").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	return s.business.DataActionBusiness.ProcessGetListDataActionRuns(ctx, request, accountUuid)
}

func (s *Service) TriggerDataActionRun(ctx context.Context, request *api.TriggerDataActionRunRequest) (*api.TriggerDataActionRunResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("TriggerDataActionRun").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	err = s.business.DataActionBusiness.ProcessNewDataActionRun(ctx, request, accountUuid)
	if err != nil {
		return nil, err
	}

	return &api.TriggerDataActionRunResponse{
		Code:    int32(code.Code_OK),
		Message: "Success",
	}, nil
}

func (s *Service) GetDataActionRunsPerDay(ctx context.Context, request *api.GetDataActionRunsPerDayRequest) (*api.GetDataActionRunsPerDayResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetDataActionRunsPerDay").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	return s.business.DataActionBusiness.ProcessGetTotalRunsPerDay(ctx, accountUuid)
}

func (s *Service) GetDataRunsProportion(ctx context.Context, request *api.GetDataRunsProportionRequest) (*api.GetDataRunsProportionResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetDataRunsProportion").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	return s.business.DataActionBusiness.ProcessGetDataRunsProportion(ctx, accountUuid)
}
