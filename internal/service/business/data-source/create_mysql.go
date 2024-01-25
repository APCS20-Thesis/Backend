package data_source

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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (b *business) ProcessCreateDataSourceMySQL(ctx context.Context, request *api.CreateDataSourceMySQLRequest) (*api.CreateDataSourceMySQLResponse, error) {
	accountUuidString, err := getMetadata(ctx, constants.KeyAccountUuid)
	if err != nil {
		b.log.WithName("ProcessCreateDataSourceMySQL").Error(err, "cannot get account uuid from metadata")
		return nil, err
	}

	accountUuid, err := uuid.Parse(accountUuidString)
	if err != nil {
		b.log.WithName("ProcessCreateDataSourceMySQL").Error(err, "cannot parse account uuid")
		return nil, err
	}

	port := request.Port
	if port == "" {
		port = "3306"
	}

	cfg, err := json.Marshal(&MySQLConfig{
		Host:     request.Host,
		Port:     request.Port,
		Database: request.Database,
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}

	err = b.repository.DataSourceRepository.CreateDataSource(ctx, &repository.CreateDataSourceParams{
		Name:          request.Name,
		Description:   request.Description,
		Type:          model.DataSourceType_MySQL,
		Configuration: pqtype.NullRawMessage{RawMessage: cfg, Valid: cfg != nil},
		AccountUuid:   accountUuid,
	})
	if err != nil {
		b.log.WithName("ProcessCreateDataSourceMySQL").Error(err, "cannot create data source")
		return nil, err
	}

	return &api.CreateDataSourceMySQLResponse{
		Code:    0,
		Message: "success",
	}, nil
}

type MySQLConfig struct {
	// host
	Host string `json:"host,omitempty"`
	// port - 3306 by default
	Port string `json:"port,omitempty"`
	// database - database name
	Database string `json:"database,omitempty"`
	// username
	Username string `json:"username,omitempty"`
	// password
	Password string `json:"password,omitempty"`
}

func getMetadata(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Metadata is not provided")
	}
	values := md[key]
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "Key is not provided")
	}

	return values[0], nil
}
