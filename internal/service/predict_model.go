package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"google.golang.org/genproto/googleapis/rpc/code"
)

func (s *Service) TrainPredictModel(ctx context.Context, request *api.TrainPredictModelRequest) (*api.TrainPredictModelResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("TrainPredictModel").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	err = s.business.PredictModelBusiness.ProcessTrainPredictModel(ctx, request, accountUuid)
	if err != nil {
		return nil, err
	}

	return &api.TrainPredictModelResponse{
		Code:    int32(code.Code_OK),
		Message: "Success",
	}, nil
}

func (s *Service) GetListPredictModels(ctx context.Context, request *api.GetListPredictModelsRequest) (*api.GetListPredictModelsResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListPredictModels").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	return s.business.PredictModelBusiness.ProcessGetListPredictModels(ctx, request, accountUuid)
}

func (s *Service) GetPredictModelDetail(ctx context.Context, request *api.GetPredictModelDetailRequest) (*api.GetPredictModelDetailResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetPredictModelDetail").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	return s.business.PredictModelBusiness.ProcessGetPredictModelDetail(ctx, request, accountUuid)
}
