package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
)

type (
	TriggerGenerateDagExportFileRequest struct {
		Config ExportFileRequestConfig `json:"conf"`
	}

	ExportFileRequestConfig struct {
		DagId            string           `json:"dag_id"`
		Key              string           `json:"key"`
		Condition        string           `json:"condition"`
		S3Configurations S3Configurations `json:"s3_configurations"`
	}
)

func (c *airflow) TriggerGenerateDagExportFile(ctx context.Context, request *TriggerGenerateDagExportFileRequest) (*TriggerNewDagRunResponse, error) {
	response := &TriggerNewDagRunResponse{}

	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: Endpoint_TRIGGER_GENERATE_DAG_EXPORT_CSV,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	return response, err
}
