package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/constants"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"strconv"
	"time"
)

func (s *Service) ImportCsv(ctx context.Context, request *api.ImportCsvRequest) (*api.ImportCsvResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("ImportCsv").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	dateTime := strconv.FormatInt(time.Now().Unix(), 10)

	err = s.s3Manger.S3Uploader(
		constants.S3BucketName,
		"data/"+accountUuid+"/"+dateTime+"_"+request.GetFileName(),
		request.GetFileContent())

	if err != nil {
		s.log.WithName("ImportCsv").
			WithValues("Context", ctx).
			Error(err, "Cannot uploaded Csv to S3")
		return nil, err
	}

	err = s.business.DataSourceBusiness.ProcessImportCsv(ctx, request, accountUuid, dateTime)
	if err != nil {
		s.log.WithName("ImportCsv").
			WithValues("Context", ctx).
			Error(err, "Failed to process import Csv")
		return nil, err
	}
	return &api.ImportCsvResponse{Message: "Import Success", Code: int32(codes.OK)}, nil
}

func (s *Service) GetListDataSources(ctx context.Context, request *api.GetListDataSourcesRequest) (*api.GetListDataSourcesResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListDataSources").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}

	return s.business.DataSourceBusiness.GetListDataSources(ctx, request, accountUuid)
}

func (s *Service) GetDataSource(ctx context.Context, request *api.GetDataSourceRequest) (*api.GetDataSourceResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetDataSource").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	response, err := s.business.DataSourceBusiness.GetDataSource(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetDataSource").
			WithValues("Context", ctx).
			Error(err, "Failed to process get data-source")
		return nil, err
	}
	return response, nil
}

func (s *Service) ImportCsvFromS3(ctx context.Context, request *api.ImportCsvFromS3Request) (*api.ImportCsvFromS3Response, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("ImportCsv").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	dateTime := strconv.FormatInt(time.Now().Unix(), 10)

	err = s.business.DataSourceBusiness.ProcessImportCsvFromS3(ctx, request, accountUuid, dateTime)
	if err != nil {
		s.log.WithName("ImportCsv").
			WithValues("Context", ctx).
			Error(err, "Failed to process import Csv from s3")
		return nil, err
	}
	return &api.ImportCsvFromS3Response{Message: "Import Success", Code: int32(codes.OK)}, nil
}

func (s *Service) ImportFromMySQLSource(ctx context.Context, request *api.ImportFromMySQLSourceRequest) (*api.ImportFromMySQLSourceResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("ImportCsv").Error(err, "Cannot get account_uuid from context")
		return nil, err
	}

	err = s.business.DataSourceBusiness.ProcessImportFromMySQLSource(ctx, request, uuid.MustParse(accountUuid))
	if err != nil {
		return nil, err
	}

	return &api.ImportFromMySQLSourceResponse{
		Code:    int32(code.Code_OK),
		Message: "Success",
	}, nil
}

func (s *Service) GetListSourceTableMap(ctx context.Context, request *api.GetListSourceTableMapRequest) (*api.GetListSourceTableMapResponse, error) {
	return s.business.DataSourceBusiness.GetListSourceTableMappings(ctx, request)
}
