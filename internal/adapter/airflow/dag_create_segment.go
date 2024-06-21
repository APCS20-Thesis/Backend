package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
)

type (
	TriggerGenerateDagCreateSegmentRequest struct {
		Config CreateSegmentConfig `json:"conf"`
	}
	CreateSegmentConfig struct {
		DagId                       string `json:"dag_id"`
		AudienceTableKey            string `json:"audience_table_key"`
		AudienceCondition           string `json:"audience_condition"`
		SegmentTableKey             string `json:"segment_table_key"`
		SegmentTableName            string `json:"segment_table_name"`
		BehaviorTableConfigurations string `json:"behavior_table_configurations"`
	}
	SegmentBehaviorCondition struct {
		BehaviorTableKey   string   `json:"behavior_table_key"`
		JoinKey            string   `json:"join_key"`
		ForeignKey         string   `json:"foreign_key"`
		WhereClauseValue   string   `json:"where_clause_value"`
		GroupByClauseValue []string `json:"group_by_clause_value"`
		HavingClauseValue  string   `json:"having_clause_value"`
	}
)

func (c *airflow) TriggerGenerateDagCreateSegment(ctx context.Context, request *TriggerGenerateDagCreateSegmentRequest) error {
	response := &TriggerNewDagRunResponse{}

	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: Endpoint_TRIGGER_GENERATE_DAG_CREATE_SEGMENT,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	return err
}
