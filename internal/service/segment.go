package service

import (
	"context"

	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/utils"
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
	
	return s.business.SegmentBusiness.ListMasterSegments(ctx, request, accountUuid)
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
		Code:             int32(codes.OK),
		Message:          "Success",
		Id:               masterSegment.Id,
		Name:             masterSegment.Name,
		Description:      masterSegment.Description,
		Status:           masterSegment.Status,
		CreatedAt:        masterSegment.CreatedAt,
		UpdatedAt:        masterSegment.UpdatedAt,
		AudienceTableId:  masterSegment.AudienceTableId,
		MainRawTableId:   masterSegment.MainRawTableId,
		MainRawTableName: masterSegment.MainRawTableName,
		AttributeTables: utils.Map(masterSegment.AttributeTables, func(table *api.MasterSegmentDetail_AttributeTable) *api.GetMasterSegmentDetailResponse_AttributeTable {
			return &api.GetMasterSegmentDetailResponse_AttributeTable{
				RawTableId:      table.RawTableId,
				RawTableName:    table.RawTableName,
				ForeignKey:      table.ForeignKey,
				JoinKey:         table.JoinKey,
				SelectedColumns: table.SelectedColumns,
			}
		}),
		BehaviorTables: utils.Map(masterSegment.BehaviorTables, func(table *api.MasterSegmentDetail_BehaviorTable) *api.GetMasterSegmentDetailResponse_BehaviorTable {
			return &api.GetMasterSegmentDetailResponse_BehaviorTable{
				Id:           table.Id,
				Name:         table.Name,
				RawTableId:   table.RawTableId,
				RawTableName: table.RawTableName,
				ForeignKey:   table.ForeignKey,
				JoinKey:      table.JoinKey,
				Schema:       table.Schema,
			}
		}),
		AudienceSchema: masterSegment.AudienceSchema,
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

func (s *Service) GetSegmentDetail(ctx context.Context, request *api.GetSegmentDetailRequest) (*api.GetSegmentDetailResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetSegmentDetail").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	return s.business.SegmentBusiness.GetSegmentDetail(ctx, request, accountUuid)
}

func (s *Service) GetMasterSegmentProfiles(ctx context.Context, request *api.GetListMasterSegmentProfilesRequest) (*api.GetListMasterSegmentProfilesResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetMasterSegmentProfiles").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	count, res, err := s.business.SegmentBusiness.ListMasterSegmentProfiles(ctx, request, accountUuid)

	return &api.GetListMasterSegmentProfilesResponse{
		Code:    int32(codes.OK),
		Message: "Success",
		Count:   count,
		Results: res,
	}, nil
}

func (s *Service) ApplyPredictModel(ctx context.Context, request *api.ApplyPredictModelRequest) (*api.ApplyPredictModelResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("ApplyPredictModel").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	return s.business.SegmentBusiness.ProcessApplyPredictModel(ctx, request, accountUuid)
}

func (s *Service) GetListPredictionActions(ctx context.Context, request *api.GetListPredictionActionsRequest) (*api.GetListPredictionActionsResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListPredictionActions").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	return s.business.SegmentBusiness.ProcessGetListPredictionActions(ctx, request, accountUuid)
}

func (s *Service) GetMasterSegmentProfile(ctx context.Context, request *api.GetMasterSegmentProfileRequest) (*api.GetMasterSegmentProfileResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetMasterSegmentProfiles").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	profile, err := s.business.SegmentBusiness.GetMasterSegmentProfile(ctx, request, accountUuid)

	return &api.GetMasterSegmentProfileResponse{
		Code:    int32(codes.OK),
		Message: "Success",
		Profile: profile,
	}, nil
}
