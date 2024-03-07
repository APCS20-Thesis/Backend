package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/go-logr/logr"
	"strings"
)

const (
	Endpoint_TRIGGER_NEW_DAG_RUN              string = "/api/v1/dags/dag_id/dagRuns"
	Endpoint_TRIGGER_GENERATE_DAG_IMPORT_FILE string = "/api/v1/dags/generate_import_file_type/dagRuns"
)

type AirflowAdapter interface {
	TriggerGenerateDagImportFile(ctx context.Context, request *TriggerGenerateDagImportFileRequest, file_type string) (*TriggerNewDagRunResponse, error)
	TriggerNewDagRun(ctx context.Context, dagId string) (*TriggerNewDagRunResponse, error)
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
	TriggerGenerateDagImportFileRequest struct {
		Config ImportFileRequestConfig `json:"conf"`
	}

	ImportFileRequestConfig struct {
		DagId                  string `json:"dag_id"`
		AccountUuid            string `json:"account_uuid"`
		DeltaTableName         string `json:"delta_table_name"`
		BucketName             string `json:"bucket_name"`
		Key                    string `json:"key"`
		WriteMode              string `json:"write_mode"`
		CsvReadOptionHeader    bool   `json:"csv_read_option_header"`
		CsvReadOptionMultiline bool   `json:"csv_read_option_multiline"`
		CsvReadOptionDelimiter string `json:"csv_read_option_delimiter"`
		CsvReadOptionSkipRow   int64  `json:"csv_read_option_skip_row"`
	}

	TriggerNewDagRunResponse struct {
		DagRunId string `json:"dag_run_id"`
	}
)

func (c *airflow) TriggerNewDagRun(ctx context.Context, dagId string) (*TriggerNewDagRunResponse, error) {
	endpoint := strings.Replace(Endpoint_TRIGGER_NEW_DAG_RUN, "dag_id", dagId, 1)
	c.log.Info("Endpoint", "endpoint", endpoint)
	response := &TriggerNewDagRunResponse{}

	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: endpoint,
		Method:   utils.Method_POST,
		Body:     "",
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	return response, err
}

func (c *airflow) TriggerGenerateDagImportFile(ctx context.Context, request *TriggerGenerateDagImportFileRequest, fileType string) (*TriggerNewDagRunResponse, error) {
	endpoint := strings.Replace(Endpoint_TRIGGER_GENERATE_DAG_IMPORT_FILE, "file_type", fileType, 1)
	c.log.Info("Endpoint", "endpoint", endpoint)
	response := &TriggerNewDagRunResponse{}

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
