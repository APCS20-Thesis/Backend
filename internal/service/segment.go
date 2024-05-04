package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CreateMasterSegment(ctx context.Context, request *api.CreateMasterSegmentRequest) (*api.CreateMasterSegmentResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("CreateMasterSegment").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	// validate table names
	tableNames := []string{"audience"}
	for _, behaviorTable := range request.BehaviorTables {
		if slices.Contains(tableNames, behaviorTable.Name) {
			return nil, status.Error(codes.InvalidArgument, "behavior table cannot have name as audience or have same name as each other")
		}
		tableNames = append(tableNames, behaviorTable.Name)
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
