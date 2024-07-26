package job

import (
	"context"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/adapter/alert"
	"github.com/APCS20-Thesis/Backend/internal/adapter/mqtt"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/internal/service/business"
	"github.com/go-logr/logr"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Job interface {
	RegisterCronJobs()
	StartCron()

	// jobs
	TriggerDagRuns(ctx context.Context)
	SyncDagRunStatus(ctx context.Context)
}

type job struct {
	cronJob        *cron.Cron
	config         *config.Config
	logger         logr.Logger
	airflowAdapter airflow.AirflowAdapter
	repository     *repository.Repository
	db             *gorm.DB
	queryAdapter   query.QueryAdapter
	mqttAdapter    mqtt.MqttAdapter
	alertAdapter   alert.AlertAdapter
	business       *business.Business
}

func NewJob(config *config.Config, logger logr.Logger, db *gorm.DB, mqttAdapter mqtt.MqttAdapter) (Job, error) {
	logger.Info("Create new Job")

	airflowAdapter, err := airflow.NewAirflowAdapter(logger, config.AirflowAdapterConfig.Address, config.AirflowAdapterConfig.Username, config.AirflowAdapterConfig.Password)
	if err != nil {
		return nil, err
	}
	queryAdapter, err := query.NewQueryAdapter(logger, config.QueryAdapterConfig.Address)
	if err != nil {
		return nil, err
	}
	alertAdapter, err := alert.NewAlertAdapter(logger, config.AlertAdapterConfig.Webhook)
	if err != nil {
		return nil, err
	}

	// Repository
	repo := repository.NewRepository(db)

	// cron
	cronJob := cron.New(cron.WithLogger(logger))

	// business
	biz := business.NewBusiness(logger, db, airflowAdapter, config, queryAdapter, nil)

	return &job{
		cronJob:        cronJob,
		config:         config,
		logger:         logger.WithName("Job"),
		airflowAdapter: airflowAdapter,
		repository:     repo,
		db:             db,
		queryAdapter:   queryAdapter,
		mqttAdapter:    mqttAdapter,
		alertAdapter:   alertAdapter,
		business:       biz,
	}, nil
}

func (j *job) StartCron() {
	j.cronJob.Start()
}

func (j *job) RegisterCronJobs() {
	_, err := j.cronJob.AddFunc("* * * * *", func() { j.TriggerDagRuns(context.Background()) })
	if err != nil {
		j.logger.Error(err, "error add cronjob TriggerDagRuns")
	}
	_, err = j.cronJob.AddFunc("* * * * *", func() { j.SyncDagRunStatus(context.Background()) })
	if err != nil {
		j.logger.Error(err, "error add cronjob SyncDagRunStatus")
	}
}
