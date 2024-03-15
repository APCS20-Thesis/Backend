package service

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/constants"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/grpc/codes"
	"strconv"
	"time"
)

func (s *Service) ImportFile(ctx context.Context, request *api.ImportFileRequest) (*api.ImportFileResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("ImportFile").
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
		s.log.WithName("ImportFile").
			WithValues("Context", ctx).
			Error(err, "Cannot uploaded file to S3")
		return nil, err
	}

	err = s.business.DataSourceBusiness.ProcessImportFile(ctx, request, accountUuid, dateTime)
	if err != nil {
		s.log.WithName("ImportFile").
			WithValues("Context", ctx).
			Error(err, "Failed to process import file")
		return nil, err
	}
	return &api.ImportFileResponse{Message: "Import Success", Code: 0}, nil
}

func (s *Service) GetListDataSources(ctx context.Context, request *api.GetListDataSourcesRequest) (*api.GetListDataSourcesResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListDataSources").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	dataSources, err := s.business.DataSourceBusiness.GetListDataSources(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetListDataSources").
			WithValues("Context", ctx).
			Error(err, "Failed to process list data-sources")
		return nil, err
	}
	return &api.GetListDataSourcesResponse{Code: int32(codes.OK), Count: int64(len(dataSources)), Results: dataSources}, nil
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

//func GetDateTimeString() string {
//	var currentTime time.Time
//	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
//	if err != nil {
//		currentTime = time.Now()
//	}
//	currentTime = time.Now().In(location)
//
//	return currentTime.Format("02012006150405")
//}

func (s *Service) GetListSourceConnections(ctx context.Context, request *api.GetListSourceConnectionsRequest) (*api.GetListSourceConnectionsResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListSourceConnections").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	sourceConnections, err := s.business.DataSourceBusiness.GetListSourceConnections(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetListSourceConnections").
			WithValues("Context", ctx).
			Error(err, "Failed to process list source connections")
		return nil, err
	}
	return &api.GetListSourceConnectionsResponse{Code: int32(codes.OK), Count: int64(len(sourceConnections)), Results: sourceConnections}, nil
}

func (s *Service) GetSourceConnection(ctx context.Context, request *api.GetSourceConnectionRequest) (*api.GetSourceConnectionResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetSourceConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	response, err := s.business.DataSourceBusiness.GetSourceConnection(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetSourceConnection").
			WithValues("Context", ctx).
			Error(err, "Failed to process get source connection")
		return nil, err
	}
	return response, nil
}

func (s *Service) CreateSourceConnection(ctx context.Context, request *api.CreateSourceConnectionRequest) (*api.CreateSourceConnectionResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("CreateSourceConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	configurations, err := json.Marshal(request.Configurations)
	if err != nil {
		s.log.WithName("ProcessImportFile").
			WithValues("Configuration", configurations).
			Error(err, "Cannot parse mappingOptions to JSON")
		return nil, err
	}
	_, err = s.business.DataSourceBusiness.CreateSourceConnection(ctx, &repository.CreateSourceConnectionParams{
		Name:           request.Name,
		Type:           model.ConnectionType(request.Type),
		Configurations: pqtype.NullRawMessage{RawMessage: configurations, Valid: true},
		AccountUuid:    uuid.MustParse(accountUuid),
	})
	if err != nil {
		s.log.WithName("CreateSourceConnection").
			WithValues("Context", ctx).
			Error(err, "Failed to process create source connection")
		return nil, err
	}
	return &api.CreateSourceConnectionResponse{Message: "Create Success", Code: 0}, nil
}
