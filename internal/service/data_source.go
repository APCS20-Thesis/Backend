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

func (s *Service) ImportCsv(ctx context.Context, request *api.ImportCsvRequest) (*api.ImportCsvResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("ImportCsv").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	dateTime := strconv.FormatInt(time.Now().Unix(), 10)

	if request.ConnectionId == 0 {
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

func (s *Service) GetListConnections(ctx context.Context, request *api.GetListConnectionsRequest) (*api.GetListConnectionsResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetListConnections").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	connections, err := s.business.DataSourceBusiness.GetListConnections(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetListConnections").
			WithValues("Context", ctx).
			Error(err, "Failed to process list connections")
		return nil, err
	}
	return &api.GetListConnectionsResponse{Code: int32(codes.OK), Count: int64(len(connections)), Results: connections}, nil
}

func (s *Service) GetConnection(ctx context.Context, request *api.GetConnectionRequest) (*api.GetConnectionResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	response, err := s.business.DataSourceBusiness.GetConnection(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Error(err, "Failed to process get connection")
		return nil, err
	}
	return response, nil
}

func (s *Service) CreateConnection(ctx context.Context, request *api.CreateConnectionRequest) (*api.CreateConnectionResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("CreateConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	configurations, err := json.Marshal(request.Configurations)
	if err != nil {
		s.log.WithName("CreateConnection").
			WithValues("Configuration", request.Configurations).
			Error(err, "Cannot parse configuration to JSON")
		return nil, err
	}
	_, err = s.business.DataSourceBusiness.CreateConnection(ctx, &repository.CreateConnectionParams{
		Name:           request.Name,
		Type:           model.ConnectionType(request.Type),
		Configurations: pqtype.NullRawMessage{RawMessage: configurations, Valid: true},
		AccountUuid:    uuid.MustParse(accountUuid),
	})
	if err != nil {
		s.log.WithName("CreateConnection").
			WithValues("Context", ctx).
			Error(err, "Failed to process create connection")
		return nil, err
	}
	return &api.CreateConnectionResponse{Message: "Create Success", Code: int32(codes.OK)}, nil
}

func (s *Service) UpdateConnection(ctx context.Context, request *api.UpdateConnectionRequest) (*api.UpdateConnectionResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("UpdateConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	configurations, err := json.Marshal(request.Configurations)
	if err != nil {
		s.log.WithName("UpdateConnection").
			WithValues("Configuration", request.Configurations).
			Error(err, "Cannot parse configuration to JSON")
		return nil, err
	}
	err = s.business.DataSourceBusiness.UpdateConnection(ctx, &repository.UpdateConnectionParams{
		ID:             request.Id,
		Name:           request.Name,
		Configurations: pqtype.NullRawMessage{RawMessage: configurations, Valid: true},
		AccountUuid:    uuid.MustParse(accountUuid),
	})
	if err != nil {
		s.log.WithName("UpdateConnection").
			WithValues("Context", ctx).
			Error(err, "Failed to process update connection")
		return nil, err
	}
	return &api.UpdateConnectionResponse{Message: "Update Success", Code: int32(codes.OK)}, nil
}

func (s *Service) DeleteConnection(ctx context.Context, request *api.DeleteConnectionRequest) (*api.DeleteConnectionResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	err = s.business.DataSourceBusiness.DeleteConnection(ctx, request, accountUuid)
	if err != nil {
		s.log.WithName("GetConnection").
			WithValues("Context", ctx).
			Error(err, "Failed to process get connection")
		return nil, err
	}
	return &api.DeleteConnectionResponse{Message: "Delete Success", Code: int32(codes.OK)}, nil
}
