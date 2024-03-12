package airflow

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/go-logr/logr"
	"strings"
	"time"
)

const (
	Endpoint_TRIGGER_NEW_DAG_RUN              string = "/api/v1/dags/dag_id/dagRuns"
	Endpoint_TRIGGER_GENERATE_DAG_IMPORT_FILE string = "/api/v1/dags/generate_import_file_type/dagRuns"
	Endpoint_LIST_DAGS                        string = "/api/v1/dags"
	Endpoint_UPDATE_DAG                       string = "/api/v1/dags/dag_id"
)

type AirflowAdapter interface {
	TriggerGenerateDagImportFile(ctx context.Context, request *TriggerGenerateDagImportFileRequest, file_type string) (*TriggerNewDagRunResponse, error)
	TriggerNewDagRun(ctx context.Context, dagId string, request *TriggerNewDagRunRequest) (*TriggerNewDagRunResponse, error)
	ListDags(ctx context.Context, request *ListDagsParams) (*ListDagsResponse, error)
	UpdateDag(ctx context.Context, dagId string, request *UpdateDagRequest) (*UpdateDagResponse, error)
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

	TriggerNewDagRunRequest struct{}

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
