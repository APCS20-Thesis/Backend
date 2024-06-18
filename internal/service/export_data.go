package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) ExportDataToFile(ctx context.Context, request *api.ExportDataToFileRequest) (*api.ExportDataToFileResponse, error) {
	availableFileTypes := []string{string(model.FileType_CSV), string(model.FileType_PARQUET)}
	if !slices.Contains(availableFileTypes, request.GetFileType()) {
		return nil, status.Error(codes.InvalidArgument, "Invalid file type, file type can only be csv or parquet")
	}

	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("ExportDataTableToFile").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	return s.business.DataDestinationBusiness.ProcessExportDataToFile(ctx, request, accountUuid)
}

func (s *Service) GetListFileExportRecords(ctx context.Context, request *api.GetListFileExportRecordsRequest) (*api.GetListFileExportRecordsResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListFileExportRecords").Error(err, "cannot get account uuid from context")
		return nil, err
	}

	records, err := s.business.DataTableBusiness.GetListFileExportRecords(ctx, request, accountUuid)
	if err != nil {
		return nil, err
	}

	return &api.GetListFileExportRecordsResponse{
		Code:    int32(codes.OK),
		Message: "Success",
		Results: records,
	}, nil
}
