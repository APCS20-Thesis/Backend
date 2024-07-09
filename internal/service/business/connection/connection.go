package connection

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"golang.org/x/exp/slices"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strconv"
)

func (b business) CreateConnection(ctx context.Context, request *api.CreateConnectionRequest, accountUuid string) (*model.Connection, error) {
	if !slices.Contains(model.ConnectionTypes, model.ConnectionType(request.Type)) {
		return nil, status.Error(codes.InvalidArgument, "Invalid connection type")
	}

	configurations, err := json.Marshal(request.Configurations)
	if err != nil {
		b.log.WithName("CreateConnection").
			WithValues("Configuration", request.Configurations).
			Error(err, "Cannot parse configuration to JSON")
		return nil, err
	}

	connection, err := b.repository.ConnectionRepository.CreateConnection(ctx, &repository.CreateConnectionParams{
		Name:           request.Name,
		Type:           model.ConnectionType(request.Type),
		Configurations: pqtype.NullRawMessage{RawMessage: configurations, Valid: configurations != nil},
		AccountUuid:    uuid.MustParse(accountUuid),
	})
	if err != nil {
		b.log.WithName("CreateConnection").Error(err, "Cannot create connection")
		return nil, err
	}
	return connection, nil
}

func (b business) UpdateConnection(ctx context.Context, params *repository.UpdateConnectionParams) error {
	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, params.ID)
	if err != nil {
		b.log.WithName("UpdateConnection").
			WithValues("ConnectionID", params.ID).
			Error(err, "Can not get record with id "+strconv.FormatInt(params.ID, 10))
		return err
	}
	if connection == nil {
		b.log.WithName("UpdateConnection").
			WithValues("ConnectionID", params.ID).
			Error(err, "No record with id "+strconv.FormatInt(params.ID, 10))
		return gorm.ErrRecordNotFound
	}
	if connection.AccountUuid != params.AccountUuid {
		b.log.WithName("UpdateConnection").
			WithValues("ConnectionID", params.ID).
			Info("Only owner can get connection")
		return status.Error(codes.PermissionDenied, "Only owner can update connection")
	}
	err = b.repository.ConnectionRepository.UpdateConnection(ctx, params)
	if err != nil {
		b.log.WithName("UpdateConnection").
			WithValues("Context", ctx).
			Error(err, "Cannot update connection")
		return err
	}
	return nil
}

func (b business) GetListConnections(ctx context.Context, request *api.GetListConnectionsRequest, accountUuid string) ([]*api.GetListConnectionsResponse_Connection, int64, error) {
	connections, count, err := b.repository.ConnectionRepository.ListConnections(ctx,
		&repository.FilterConnection{
			Name:        request.Name,
			Type:        model.ConnectionType(request.Type),
			AccountUuid: uuid.MustParse(accountUuid),
		})
	if err != nil {
		b.log.WithName("GetListConnections").
			Error(err, "Cannot get list connection")
		return nil, 0, err
	}
	var response []*api.GetListConnectionsResponse_Connection
	for _, connection := range connections {
		response = append(response, &api.GetListConnectionsResponse_Connection{
			Id:               connection.ID,
			Name:             connection.Name,
			Type:             string(connection.Type),
			UpdatedAt:        connection.UpdatedAt.String(),
			DataSources:      nil,
			DataDestinations: nil,
		})
	}
	return response, count, nil
}

func (b business) GetConnection(ctx context.Context, request *api.GetConnectionRequest, accountUuid string) (*api.GetConnectionResponse, error) {
	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Not found connection with id "+strconv.FormatInt(request.Id, 10))
		}
		b.log.WithName("GetConnection").Error(err, "Cannot get connection")
		return nil, err
	}
	if connection.AccountUuid != uuid.MustParse(accountUuid) {
		b.log.WithName("GetConnection").Info("Only owner can get connection")
		return nil, status.Error(codes.PermissionDenied, "Only owner can get connection")
	}
	var configurations map[string]string
	err = json.Unmarshal(connection.Configurations.RawMessage, &configurations)
	if err != nil {
		return nil, err
	}

	// Transform secret field
	switch connection.Type {
	case model.ConnectionType_MySQL:
		configurations["password"] = utils.TransformPassword(configurations["password"])
	case model.ConnectionType_Gophish:
		configurations["api_key"] = utils.TransformPassword(configurations["api_key"])
		//case model.ConnectionType_S3:
		//	configurations["secret_access_key"] = utils.TransformPassword(configurations["secret_access_key"])
	}

	return &api.GetConnectionResponse{
		Code:           int32(code.Code_OK),
		Id:             connection.ID,
		Name:           connection.Name,
		Type:           string(connection.Type),
		Configurations: configurations,
		CreatedAt:      connection.CreatedAt.String(),
		UpdatedAt:      connection.UpdatedAt.String(),
	}, nil
}
func (b business) DeleteConnection(ctx context.Context, request *api.DeleteConnectionRequest, accountUuid string) error {
	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.Id)
	if err != nil {
		b.log.WithName("DeleteConnection").
			WithValues("ConnectionID", request.Id).
			Error(err, "Can not get record with id "+strconv.FormatInt(request.Id, 10))
		return err
	}
	if connection == nil {
		b.log.WithName("DeleteConnection").
			WithValues("ConnectionID", request.Id).
			Error(err, "No record with id "+strconv.FormatInt(request.Id, 10))
		return gorm.ErrRecordNotFound
	}
	if connection.AccountUuid != uuid.MustParse(accountUuid) {
		b.log.WithName("DeleteConnection").
			WithValues("ConnectionID", request.Id).
			Info("Only owner can get connection")
		return status.Error(codes.PermissionDenied, "Only owner can delete connection")
	}
	err = b.repository.ConnectionRepository.DeleteConnection(ctx, request.Id)
	if err != nil {
		b.log.WithName("DeleteConnection").Error(err, "Cannot delete connection")
		return err
	}
	return nil
}

func (b business) ProcessGetMySQLTableSchema(ctx context.Context, request *api.GetMySQLTableSchemaRequest, accountUuid string) ([]*api.SchemaColumn, error) {
	logger := b.log.WithName("ProcessGetMySQLTableSchema")

	connection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.ConnectionId)
	if err != nil {
		logger.Error(err, "cannot get connection")
		return nil, err
	}
	if connection.AccountUuid.String() != accountUuid {
		return nil, status.Error(codes.PermissionDenied, "Not have permission on this connection")
	}

	var config model.MySQLConfiguration
	err = json.Unmarshal(connection.Configurations.RawMessage, &config)
	if err != nil {
		logger.Error(err, "cannot unmarshal connection configuration")
		return nil, err
	}

	mysqlStringConfig := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password, config.Host, config.Port, config.Database)
	// Connect to the MySQL database
	mysqlDB, err := sql.Open("mysql", mysqlStringConfig)
	if err != nil {
		logger.Error(err, "cannot connect to database")
		return nil, err
	}
	defer mysqlDB.Close()

	// Query the information_schema to get the schema for the specified table
	rows, err := mysqlDB.Query("SELECT COLUMN_NAME, DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", request.TableName)
	if err != nil {
		logger.Error(err, "cannot query the database")
		return nil, err
	}
	defer rows.Close()

	var schema []*api.SchemaColumn
	for rows.Next() {
		var columnName, dataType string
		err := rows.Scan(&columnName, &dataType)
		if err != nil {
			logger.Error(err, "error scanning result rows")
			return nil, err
		}
		schema = append(schema, &api.SchemaColumn{
			ColumnName: columnName,
			DataType:   dataType,
		})
	}

	if len(schema) == 0 {
		return nil, status.Error(codes.NotFound, "Table is empty or wrong table name")
	}

	return schema, nil
}
