package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"google.golang.org/grpc/codes"
)

func (s *Service) CreateMasterSegment(ctx context.Context, request *api.CreateMasterSegmentRequest) (*api.CreateMasterSegmentResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("CreateMasterSegment").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	err = s.business.SegmentBusiness.CreateMasterSegment(ctx, request, accountUuid)
	if err != nil {
		return nil, err
	}

	return &api.CreateMasterSegmentResponse{
		Code:    int32(codes.OK),
		Message: "Success",
	}, nil
}
