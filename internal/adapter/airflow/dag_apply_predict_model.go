package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
)

type (
	TriggerGenerateDagApplyPredictModelRequest struct {
		Conf DagApplyPredictModelConfig `json:"conf"`
	}

	DagApplyPredictModelConfig struct {
		DagId            string   `json:"dag_id"`
		DataKey          string   `json:"data_key"`
		ModelPath        string   `json:"model_path"`
		ResultPath       string   `json:"result_path"`
		SelectAttributes []string `json:"select_attributes"`
	}
)

func (c *airflow) TriggerGenerateDagApplyPredictModel(ctx context.Context, request *TriggerGenerateDagApplyPredictModelRequest) error {
	response := &TriggerNewDagRunResponse{}

	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: Endpoint_TRIGGER_GENERATE_DAG_APPLY_PREDICT_MODEL,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	return err
}
