package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
)

type (
	TriggerGenerateDagTrainPredictModelRequest struct {
		Conf DagTrainPredictModelConfig `json:"conf"`
	}

	DagTrainPredictModelConfig struct {
		DagId            string   `json:"dag_id"`
		Segment1Key      string   `json:"segment1_key"`
		Segment2Key      string   `json:"segment2_key"`
		Label1           int      `json:"label1"`
		Label2           int      `json:"label2"`
		PredictModelKey  string   `json:"predict_model_key"`
		SelectAttributes []string `json:"select_attributes"`
	}
)

func (c *airflow) TriggerGenerateDagTrainPredictModel(ctx context.Context, request *TriggerGenerateDagTrainPredictModelRequest) (*TriggerNewDagRunResponse, error) {
	c.log.Info("Endpoint", "endpoint", Endpoint_TRIGGER_GENERATE_DAG_TRAIN_PREDICT_MODEL)
	response := &TriggerNewDagRunResponse{}

	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: Endpoint_TRIGGER_GENERATE_DAG_TRAIN_PREDICT_MODEL,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)
	if err != nil {
		return nil, err
	}

	return response, err
}
