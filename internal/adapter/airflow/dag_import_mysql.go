package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
)

type (
	TriggerGenerateDagImportMySQLRequest struct {
		Conf DagImportMySQLConfig `json:"conf"`
	}

	DagImportMySQLConfig struct {
		DagId                 string                              `json:"dag_id"`
		AccountUuid           string                              `json:"account_uuid"`
		DeltaTableName        string                              `json:"delta_table_name"`
		Headers               []string                            `json:"headers"`
		DatabaseConfiguration DagImportMySQLDatabaseConfiguration `json:"database_configuration"`
	}

	DagImportMySQLDatabaseConfiguration struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
		User     string `json:"user"`
		Password string `json:"password"`
		Table    string `json:"table"`
	}
)

func (c *airflow) TriggerGenerateDagImportMySQL(ctx context.Context, request *TriggerGenerateDagImportMySQLRequest) (*TriggerNewDagRunResponse, error) {
	c.log.Info("Endpoint", "endpoint", Endpoint_TRIGGER_GENERATE_DAG_IMPORT_MYSQL)
	response := &TriggerNewDagRunResponse{}

	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: Endpoint_TRIGGER_GENERATE_DAG_IMPORT_MYSQL,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)
	if err != nil {
		return nil, err
	}

	return response, err
}
