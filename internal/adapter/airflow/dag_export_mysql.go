package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
)

type TriggerGenerateDagExportMySQLRequest struct {
	Conf DagExportMySQLConfig `json:"conf"`
}

type DagExportMySQLConfig struct {
	DagId                 string                              `json:"dag_id"`
	AccountUuid           string                              `json:"account_uuid"`
	DeltaTableName        string                              `json:"delta_table_name"`
	MasterSegmentId       int64                               `json:"master_segment_id"`
	Condition             string                              `json:"condition"`
	DatabaseConfiguration DagExportMySQLDatabaseConfiguration `json:"database_configuration"`
	DestinationTableName  string                              `json:"destination_table_name"`
}

type DagExportMySQLDatabaseConfiguration struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func (c *airflow) TriggerGenerateDagExportMySQL(ctx context.Context, request *TriggerGenerateDagExportMySQLRequest) (*TriggerNewDagRunResponse, error) {
	c.log.Info("Endpoint", "endpoint", Endpoint_TRIGGER_GENERATE_DAG_IMPORT_MYSQL)
	response := &TriggerNewDagRunResponse{}

	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: Endpoint_TRIGGER_GENERATE_DAG_EXPORT_MYSQL,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
