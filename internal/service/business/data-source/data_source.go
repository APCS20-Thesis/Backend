package data_source

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (b business) CreateDataSource(ctx context.Context, params *repository.CreateDataSourceParams) (*model.DataSource, error) {
	dataSource, err := b.repository.DataSourceRepository.CreateDataSource(ctx, params)
	if err != nil {
		b.log.WithName("CreateDataSource").
			WithValues("Context", ctx).
			Error(err, "Cannot create data_source")
		return nil, err
	}
	return dataSource, nil
}

func (b business) GetListDataSources(ctx context.Context, request *api.GetListDataSourcesRequest, accountUuid string) ([]*api.GetListDataSourcesResponse_DataSource, error) {
	dataSources, err := b.repository.DataSourceRepository.ListDataSources(ctx,
		&repository.ListDataSourcesFilters{
			Name:        request.Name,
			AccountUuid: uuid.MustParse(accountUuid),
			Type:        model.DataSourceType(request.Type),
		})
	if err != nil {
		b.log.WithName("GetListDataSources").
			WithValues("Context", ctx).
			Error(err, "Cannot get list data_sources")
		return nil, err
	}
	var response []*api.GetListDataSourcesResponse_DataSource
	for _, dataSource := range dataSources {
		response = append(response, &api.GetListDataSourcesResponse_DataSource{Id: dataSource.ID, Name: dataSource.Name, Type: string(dataSource.Type), UpdatedAt: dataSource.UpdatedAt.String()})
	}
	return response, nil
}

func (b business) GetDataSource(ctx context.Context, request *api.GetDataSourceRequest, accountUuid string) (*api.GetDataSourceResponse, error) {
	dataSource, err := b.repository.DataSourceRepository.GetDataSource(ctx, request.Id)
	if err != nil {
		b.log.WithName("GetDataSource").
			WithValues("Context", ctx).
			Error(err, "Cannot get data_source")
		return nil, err
	}
	if dataSource.AccountUuid != uuid.MustParse(accountUuid) {
		b.log.WithName("GetDataSource").
			WithValues("Context", ctx).
			Info("Only owner can get data_source")
		return nil, status.Error(codes.PermissionDenied, "Only owner can get data_source")
	}
	var configurations map[string]string
	if dataSource.Configurations.RawMessage != nil {
		err = json.Unmarshal(dataSource.Configurations.RawMessage, &configurations)
		if err != nil {
			return nil, err
		}
	}
	return &api.GetDataSourceResponse{
		Code:           int32(code.Code_OK),
		Id:             dataSource.ID,
		Name:           dataSource.Name,
		Type:           string(dataSource.Type),
		Description:    dataSource.Description,
		CreatedAt:      dataSource.CreatedAt.String(),
		UpdatedAt:      dataSource.UpdatedAt.String(),
		Configurations: configurations,
	}, nil
}
