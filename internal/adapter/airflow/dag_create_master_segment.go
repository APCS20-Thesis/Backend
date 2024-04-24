package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
)

type (
	TriggerGenerateDagCreateMasterSegmentRequest struct {
		Config CreateMasterSegmentConfig `json:"conf"`
	}
	CreateMasterSegmentConfig struct {
		DagId           string `json:"dag_id"`
		AccountUuid     string `json:"account_uuid"`
		MasterSegmentId int64  `json:"master_segment_id"`
		MainTableName   string `json:"main_table_name"`
		MainAttributes  string `json:"main_attributes"`
		AttributeTables string `json:"attribute_tables"`
		BehaviorTables  string `json:"behavior_tables"`
	}
	CreateMasterSegmentConfig_TableColumns struct {
		TableColumnName    string `json:"table_column_name"`
		AudienceColumnName string `json:"audience_column_name"`
	}
	CreateMasterSegmentConfig_AttributeTable struct {
		TableName  string                                   `json:"table_name"`
		JoinKey    string                                   `json:"join_key"`
		ForeignKey string                                   `json:"foreign_key"`
		Columns    []CreateMasterSegmentConfig_TableColumns `json:"columns"`
	}
	CreateMasterSegmentConfig_BehaviorTable struct {
		TableName         string                                   `json:"table_name"`
		BehaviorTableName string                                   `json:"behavior_table_name"`
		JoinKey           string                                   `json:"join_key"`
		ForeignKey        string                                   `json:"foreign_key"`
		Columns           []CreateMasterSegmentConfig_TableColumns `json:"columns"`
	}
)

func (c *airflow) TriggerGenerateDagCreateMasterSegment(ctx context.Context, request *TriggerGenerateDagCreateMasterSegmentRequest) error {
	response := &TriggerNewDagRunResponse{}

	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: Endpoint_TRIGGER_GENERATE_DAG_CREATE_MASTER_SEGMENT,
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	return err
}
