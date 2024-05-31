package query

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/go-logr/logr"
)

const (
	Endpoint_GET_DATA_TABLE   string = "/data_table"
	Endpoint_GET_SCHEMA_TABLE string = "/schema_table"
)

type QueryAdapter interface {
	GetDataTable(ctx context.Context, request *GetQueryDataTableRequest) (*GetQueryDataTableResponse, error)
	GetSchemaTable(ctx context.Context, request *GetSchemaDataTableRequest) (*GetSchemaDataTableResponse, error)
}

type query struct {
	log    logr.Logger
	client utils.HttpClient
}

func NewQueryAdapter(log logr.Logger, host string) (QueryAdapter, error) {
	client := utils.HttpClient{}
	client.Init("Airflow Client", log, host)
	return &query{
		log:    log,
		client: client,
	}, nil
}

type (
	GetQueryDataTableRequest struct {
		Limit     int32  `json:"limit"`
		TablePath string `json:"table_path"`
	}

	GetQueryDataTableResponse struct {
		Count int64    `json:"count"`
		Data  []string `json:"data"`
	}
)

func (c *query) GetDataTable(ctx context.Context, request *GetQueryDataTableRequest) (*GetQueryDataTableResponse, error) {
	response := &GetQueryDataTableResponse{}

	err := c.client.SendHttpRequest(ctx, utils.Request{
		Endpoint: Endpoint_GET_DATA_TABLE,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	return response, err
}

type (
	GetSchemaDataTableRequest struct {
		TablePath string `json:"table_path"`
	}
	FieldSchema struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	GetSchemaDataTableResponse struct {
		Fields []FieldSchema `json:"fields"`
	}
)

func (c *query) GetSchemaTable(ctx context.Context, request *GetSchemaDataTableRequest) (*GetSchemaDataTableResponse, error) {
	response := &GetSchemaDataTableResponse{}

	err := c.client.SendHttpRequest(ctx, utils.Request{
		Endpoint: Endpoint_GET_SCHEMA_TABLE,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	return response, err
}
