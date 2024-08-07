package query

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/go-logr/logr"
)

const (
	Endpoint_GET_DATA_TABLE       string = "/api/data-table"
	Endpoint_GET_SCHEMA_TABLE     string = "/api/schema-table"
	Endpoint_GET_DATA_TABLE_V2    string = "/api/v2/data-table"
	Endpoint_QUERY_SQL            string = "/api/delta-query"
	Endpoint_QUERY_SQL_V2         string = "/api/v2/delta-query"
	Endpoint_GET_COUNT_DATA_TABLE string = "/api/count-table"
)

type QueryAdapter interface {
	GetDataTableV2(ctx context.Context, request *GetQueryDataTableV2Request) (*GetQueryDataTableV2Response, error)
	GetDataTable(ctx context.Context, request *GetQueryDataTableRequest) (*GetQueryDataTableResponse, error)
	GetSchemaTable(ctx context.Context, request *GetSchemaDataTableRequest) (*GetSchemaDataTableResponse, error)
	QueryRawSQL(ctx context.Context, request *QueryRawSQLRequest) (*QueryRawSQLResponse, error)
	QueryRawSQLV2(ctx context.Context, request *QueryRawSQLV2Request) (*QueryRawSQLV2Response, error)
	GetCountDataTable(ctx context.Context, request *GetCountDataTableRequest) (*GetCountDataTableResponse, error)
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
	GetQueryDataTableV2Request struct {
		Limit     int32  `json:"limit"`
		TablePath string `json:"table_path"`
	}

	GetQueryDataTableV2Response struct {
		Count int64    `json:"count"`
		Data  []string `json:"data"`
	}
)

func (c *query) GetDataTableV2(ctx context.Context, request *GetQueryDataTableV2Request) (*GetQueryDataTableV2Response, error) {
	response := &GetQueryDataTableV2Response{}

	err := c.client.SendHttpRequest(ctx, utils.Request{
		Endpoint: Endpoint_GET_DATA_TABLE_V2,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	return response, err
}

type (
	GetQueryDataTableRequest struct {
		Limit     int32  `json:"limit"`
		TablePath string `json:"table_path"`
	}

	GetQueryDataTableResponse struct {
		Count int64               `json:"count"`
		Data  []map[string]string `json:"data"`
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
		Schema []FieldSchema `json:"schema"`
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

type (
	QueryRawSQLRequest struct {
		Query string `json:"query"`
	}
	QueryRawSQLResponse struct {
		Count int              `json:"count"`
		Data  []map[string]any `json:"data"`
	}
)

func (c *query) QueryRawSQL(ctx context.Context, request *QueryRawSQLRequest) (*QueryRawSQLResponse, error) {
	var response QueryRawSQLResponse

	err := c.client.SendHttpRequest(ctx, utils.Request{
		Endpoint: Endpoint_QUERY_SQL,
		Method:   utils.Method_POST,
		Body:     request,
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type (
	QueryRawSQLV2Request struct {
		Query string `json:"query"`
	}
	QueryRawSQLV2Response struct {
		Count int      `json:"count"`
		Data  []string `json:"data"`
	}
)

func (c *query) QueryRawSQLV2(ctx context.Context, request *QueryRawSQLV2Request) (*QueryRawSQLV2Response, error) {
	var response QueryRawSQLV2Response

	err := c.client.SendHttpRequest(ctx, utils.Request{
		Endpoint: Endpoint_QUERY_SQL_V2,
		Method:   utils.Method_POST,
		Body:     request,
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func QueryV2Paginate(page int32, pageSize int32, list []string) []string {
	if page <= 0 {
		page = 1
	}

	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	return list[offset : offset+pageSize]
}

type (
	GetCountDataTableRequest struct {
		TablePath string `json:"table_path"`
	}
	GetCountDataTableResponse struct {
		Count int `json:"count"`
	}
)

func (c *query) GetCountDataTable(ctx context.Context, request *GetCountDataTableRequest) (*GetCountDataTableResponse, error) {
	response := &GetCountDataTableResponse{}

	err := c.client.SendHttpRequest(ctx, utils.Request{
		Endpoint: Endpoint_GET_COUNT_DATA_TABLE,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	return response, err
}
