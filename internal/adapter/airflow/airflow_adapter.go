package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/go-logr/logr"
	"strings"
)

const (
	Endpoint_TRIGGER_NEW_DAG_RUN string = "/api/v1/dags/dag_id/dagRuns"
)

type AirflowAdapter interface {
	TriggerNewDagRunImportFile(ctx context.Context, request *TriggerNewDagRunImportFileRequest, dagId string) (*TriggerNewDagRunImportFileResponse, error)
}

type airflow struct {
	log      logr.Logger
	client   utils.HttpClient
	username string
	password string
}

func NewAirflowAdapter(log logr.Logger, host string, username string, password string) (AirflowAdapter, error) {
	client := utils.HttpClient{}
	client.Init("Airflow Client", log, host)
	return &airflow{
		log:      log,
		client:   client,
		username: username,
		password: password,
	}, nil
}

type (
	TriggerNewDagRunImportFileRequest struct {
		Config ImportFileRequestConfig `json:"conf"`
	}

	ImportFileRequestConfig struct {
		AccountUuid            string `json:"account_uuid"`
		DeltaTableName         string `json:"delta_table_name"`
		CsvFilePath            string `json:"csv_file_path"`
		WriteMode              string `json:"write_mode"`
		CsvReadOptionHeader    bool   `json:"csv_read_option_header"`
		CsvReadOptionMultiline bool   `json:"csv_read_option_multiline"`
		CsvReadOptionDelimiter string `json:"csv_read_option_delimiter"`
	}

	TriggerNewDagRunImportFileResponse struct {
		DagRunId string `json:"dag_run_id"`
	}
)

func (c *airflow) TriggerNewDagRunImportFile(ctx context.Context, request *TriggerNewDagRunImportFileRequest, dagId string) (*TriggerNewDagRunImportFileResponse, error) {
	endpoint := strings.Replace(Endpoint_TRIGGER_NEW_DAG_RUN, "dag_id", dagId, 1)
	c.log.Info("Endpoint", "endpoint", endpoint)
	response := &TriggerNewDagRunImportFileResponse{}

	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: endpoint,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	return response, err
}
