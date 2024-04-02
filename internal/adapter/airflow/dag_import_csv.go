package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/utils"
)

type (
	TriggerGenerateDagImportCsvRequest struct {
		Config ImportCsvRequestConfig `json:"conf"`
	}

	ImportCsvRequestConfig struct {
		DagId            string                                        `json:"dag_id"`
		AccountUuid      string                                        `json:"account_uuid"`
		DeltaTableName   string                                        `json:"delta_table_name"`
		S3Configurations *S3Configurations                             `json:"s3_configurations"`
		WriteMode        DeltaWriteMode                                `json:"write_mode"`
		CsvReadOptions   *api.ImportCsvRequest_ImportCsvConfigurations `json:"csv_read_options"`
		Headers          []string                                      `json:"headers"`
	}

	S3Configurations struct {
		AccessKeyId     string `json:"access_key_id"`
		SecretAccessKey string `json:"secret_access_key"`
		BucketName      string `json:"bucket_name"`
		Region          string `json:"region"`
		Key             string `json:"key"`
	}

	DeltaWriteMode string
)

func (c *airflow) TriggerGenerateDagImportCsv(ctx context.Context, request *TriggerGenerateDagImportCsvRequest) (*TriggerNewDagRunResponse, error) {
	c.log.Info("Endpoint", "endpoint", Endpoint_TRIGGER_GENERATE_DAG_IMPORT_CSV)
	response := &TriggerNewDagRunResponse{}

	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: Endpoint_TRIGGER_GENERATE_DAG_IMPORT_CSV,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	return response, err
}
