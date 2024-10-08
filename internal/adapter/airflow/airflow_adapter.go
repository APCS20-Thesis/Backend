package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/go-logr/logr"
	"strings"
	"time"
)

const (
	Endpoint_TRIGGER_NEW_DAG_RUN                        string = "/api/v1/dags/dag_id/dagRuns"
	Endpoint_LIST_DAGS                                  string = "/api/v1/dags"
	Endpoint_UPDATE_DAG                                 string = "/api/v1/dags/dag_id"
	Endpoint_GET_DAG_RUN                                string = "/api/v1/dags/dag_id/dagRuns/dag_run_id"
	Endpoint_TRIGGER_GENERATE_DAG_IMPORT_CSV            string = "/api/v1/dags/generate_import_csv/dagRuns"
	Endpoint_TRIGGER_GENERATE_DAG_IMPORT_MYSQL          string = "/api/v1/dags/generate_import_mysql/dagRuns"
	Endpoint_TRIGGER_GENERATE_DAG_EXPORT_CSV            string = "/api/v1/dags/generate_export_csv/dagRuns"
	Endpoint_TRIGGER_GENERATE_DAG_EXPORT_MYSQL          string = "/api/v1/dags/generate_export_mysql/dagRuns"
	Endpoint_TRIGGER_GENERATE_DAG_CREATE_MASTER_SEGMENT string = "/api/v1/dags/generate_create_master_segment/dagRuns"
	Endpoint_TRIGGER_GENERATE_DAG_CREATE_SEGMENT        string = "/api/v1/dags/generate_create_segment/dagRuns"
	Endpoint_TRIGGER_GENERATE_DAG_TRAIN_PREDICT_MODEL   string = "/api/v1/dags/generate_train_predict_model/dagRuns"
	Endpoint_TRIGGER_GENERATE_DAG_APPLY_PREDICT_MODEL   string = "/api/v1/dags/generate_apply_predict_model/dagRuns"

	WriteMode_Append    DeltaWriteMode = "append"
	WriteMode_Overwrite DeltaWriteMode = "overwrite"
)

type AirflowAdapter interface {
	TriggerGenerateDagImportCsv(ctx context.Context, request *TriggerGenerateDagImportCsvRequest) (*TriggerNewDagRunResponse, error)
	TriggerGenerateDagImportMySQL(ctx context.Context, request *TriggerGenerateDagImportMySQLRequest) (*TriggerNewDagRunResponse, error)
	TriggerGenerateDagExportFile(ctx context.Context, request *TriggerGenerateDagExportFileRequest) (*TriggerNewDagRunResponse, error)
	TriggerGenerateDagExportMySQL(ctx context.Context, request *TriggerGenerateDagExportMySQLRequest) (*TriggerNewDagRunResponse, error)
	TriggerNewDagRun(ctx context.Context, dagId string, request *TriggerNewDagRunRequest) (*TriggerNewDagRunResponse, error)
	TriggerGenerateDagCreateMasterSegment(ctx context.Context, request *TriggerGenerateDagCreateMasterSegmentRequest) error
	TriggerGenerateDagCreateSegment(ctx context.Context, request *TriggerGenerateDagCreateSegmentRequest) error
	TriggerGenerateDagTrainPredictModel(ctx context.Context, request *TriggerGenerateDagTrainPredictModelRequest) (*TriggerNewDagRunResponse, error)
	TriggerGenerateDagApplyPredictModel(ctx context.Context, request *TriggerGenerateDagApplyPredictModelRequest) error

	ListDags(ctx context.Context, request *ListDagsParams) (*ListDagsResponse, error)
	UpdateDag(ctx context.Context, dagId string, request *UpdateDagRequest) (*UpdateDagResponse, error)
	GetDagRun(ctx context.Context, dagId string, dagRunId string) (*GetDagRunResponse, error)
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
	TriggerNewDagRunRequest struct{}

	TriggerNewDagRunResponse struct {
		DagId    string `json:"dag_id"`
		DagRunId string `json:"dag_run_id"`
		//DataIntervalEnd   time.Time   `json:"data_interval_end"`
		//DataIntervalStart time.Time   `json:"data_interval_start"`
		//EndDate           interface{} `json:"end_date"`
		//ExecutionDate     time.Time   `json:"execution_date"`
		//ExternalTrigger   bool        `json:"external_trigger"`
		//LastSchedulingDecision interface{} `json:"last_scheduling_decision"`
		//LogicalDate time.Time `json:"logical_date"`
		//Note        interface{} `json:"note"`
		RunType string `json:"run_type"`
		//StartDate interface{} `json:"start_date"`
		State string `json:"state"`
	}
)

func (c *airflow) TriggerNewDagRun(ctx context.Context, dagId string, request *TriggerNewDagRunRequest) (*TriggerNewDagRunResponse, error) {
	endpoint := strings.Replace(Endpoint_TRIGGER_NEW_DAG_RUN, "dag_id", dagId, 1)
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

type (
	ListDagsParams struct {
		Limit        string   `json:"limit,omitempty"`
		Offset       string   `json:"offset,omitempty"`
		OrderBy      string   `json:"order_by,omitempty"`
		Tags         []string `json:"tags,omitempty"`
		OnlyActive   string   `json:"only_active,omitempty"`
		Paused       string   `json:"paused,omitempty"`
		DagIdPattern string   `json:"dag_id_pattern,omitempty"`
	}

	ListDagsResponse struct {
		Dags         []*Dag `json:"dags"`
		TotalEntries int    `json:"total_entries"`
	}

	Dag struct {
		DagID                    string `json:"dag_id"`
		DefaultView              string `json:"default_view"`
		Description              string `json:"description"`
		FileToken                string `json:"file_token"`
		Fileloc                  string `json:"fileloc"`
		HasImportErrors          bool   `json:"has_import_errors"`
		HasTaskConcurrencyLimits bool   `json:"has_task_concurrency_limits"`
		IsActive                 bool   `json:"is_active"`
		IsPaused                 bool   `json:"is_paused"`
		IsSubdag                 bool   `json:"is_subdag"`
	}
)

func (c *airflow) ListDags(ctx context.Context, request *ListDagsParams) (*ListDagsResponse, error) {
	var response ListDagsResponse

	params, err := utils.ConvertStructRequestToParams(&request)
	if err != nil {
		return nil, err
	}

	err = c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: Endpoint_LIST_DAGS,
		Method:   utils.Method_GET,
		Params:   params,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type (
	UpdateDagRequest struct {
		IsPaused bool `json:"is_paused"`
	}

	UpdateDagResponse struct {
		DagID          string    `json:"dag_id"`
		RootDagID      string    `json:"root_dag_id"`
		IsPaused       bool      `json:"is_paused"`
		IsActive       bool      `json:"is_active"`
		IsSubdag       bool      `json:"is_subdag"`
		LastParsedTime time.Time `json:"last_parsed_time"`
		LastPickled    time.Time `json:"last_pickled"`
		LastExpired    time.Time `json:"last_expired"`
		SchedulerLock  bool      `json:"scheduler_lock"`
		PickleID       string    `json:"pickle_id"`
		DefaultView    string    `json:"default_view"`
		Fileloc        string    `json:"fileloc"`
		FileToken      string    `json:"file_token"`
		Owners         []string  `json:"owners"`
		Description    string    `json:"description"`
	}
)

func (c *airflow) UpdateDag(ctx context.Context, dagId string, request *UpdateDagRequest) (*UpdateDagResponse, error) {
	endpoint := strings.Replace(Endpoint_UPDATE_DAG, "dag_id", dagId, 1)
	c.log.WithName("Airflow-UpdateDag").Info("Endpoint", "endpoint", endpoint)
	response := UpdateDagResponse{}

	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: endpoint,
		Method:   utils.Method_PATCH,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, response)

	c.log.WithName("Airflow-UpdateDag").Info("Done calling", "response", response, "error", err)
	return &response, err
}

type (
	GetDagRunResponse struct {
		DagRunId          string    `json:"dag_run_id"`
		DagId             string    `json:"dag_id"`
		LogicalDate       time.Time `json:"logical_date"`
		ExecutionDate     time.Time `json:"execution_date"`
		StartDate         time.Time `json:"start_date"`
		EndDate           time.Time `json:"end_date"`
		DataIntervalStart time.Time `json:"data_interval_start"`
		DataIntervalEnd   time.Time `json:"data_interval_end"`
		State             string    `json:"state"`
		ExternalTrigger   bool      `json:"external_trigger"`
		Note              string    `json:"note"`
	}
)

func (c *airflow) GetDagRun(ctx context.Context, dagId string, dagRunId string) (*GetDagRunResponse, error) {
	endpoint := strings.Replace(Endpoint_GET_DAG_RUN, "dag_id", dagId, 1)
	endpoint = strings.Replace(endpoint, "dag_run_id", dagRunId, 1)
	c.log.Info("endpoint", "endpoint", endpoint)
	var response GetDagRunResponse
	err := c.client.SendHttpRequestWithBasicAuth(ctx, utils.BasicAuth{
		Username: c.username,
		Password: c.password,
	}, utils.Request{
		Endpoint: endpoint,
		Method:   utils.Method_GET,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
