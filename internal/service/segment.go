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

func (s *Service) CreateSegment(ctx context.Context, request *api.CreateSegmentRequest) (*api.CreateSegmentResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("CreateSegment").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	err = s.business.SegmentBusiness.CreateSegment(ctx, request, accountUuid)
	if err != nil {
		return nil, err
	}

	return &api.CreateSegmentResponse{
		Code:    int32(codes.OK),
		Message: "Success",
	}, nil
}

func (s *Service) GetListMasterSegments(ctx context.Context, request *api.GetListMasterSegmentsRequest) (*api.GetListMasterSegmentsResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListMasterSegments").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	count, masterSegments, err := s.business.SegmentBusiness.ListMasterSegments(ctx, request, accountUuid)
	if err != nil {
		return nil, err
	}

	return &api.GetListMasterSegmentsResponse{
		Code:    int32(codes.OK),
		Message: "Success",
		Count:   count,
		Results: masterSegments,
	}, nil
}

func (s *Service) GetMasterSegmentDetail(ctx context.Context, request *api.GetMasterSegmentDetailRequest) (*api.GetMasterSegmentDetailResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListMasterSegments").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	masterSegment, err := s.business.SegmentBusiness.GetMasterSegmentDetail(ctx, request, accountUuid)
	if err != nil {
		return nil, err
	}

	return &api.GetMasterSegmentDetailResponse{
		Code:    int32(codes.OK),
		Message: "Success",
		Data:    masterSegment,
	}, nil
}

func (s *Service) GetListSegments(ctx context.Context, request *api.GetListSegmentsRequest) (*api.GetListSegmentsResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListMasterSegments").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	segments, err := s.business.SegmentBusiness.ListSegments(ctx, request, accountUuid)
	if err != nil {
		return nil, err
	}

	return &api.GetListSegmentsResponse{
		Code:    int32(codes.OK),
		Message: "Success",
		Results: segments,
	}, nil
}
