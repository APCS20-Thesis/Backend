package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
	"strings"
)

type (
	TriggerGenerateDagExportFileRequest struct {
		Config ExportFileRequestConfig `json:"conf"`
	}

	ExportFileRequestConfig struct {
		DagId          string `json:"dag_id"`
		AccountUuid    string `json:"account_uuid"`
		DeltaTableName string `json:"delta_table_name"`
		SavedS3Path    string `json:"saved_s3_path"`
	}
)

func (c *airflow) TriggerGenerateDagExportFile(ctx context.Context, request *TriggerGenerateDagExportFileRequest) (*TriggerNewDagRunResponse, error) {
	response := &TriggerNewDagRunResponse{}
	endpoint := strings.Replace(Endpoint_TRIGGER_GENERATE_DAG_EXPORT_CSV, "dag_id", request.Config.DagId, 1)

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
