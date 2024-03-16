package job

import (
	"context"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/go-logr/logr"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Job interface {
	RegisterCronJobs()
	StartCron()

	// jobs
	TriggerDagRuns(ctx context.Context)
}

type job struct {
	cronJob        *cron.Cron
	config         *config.Config
	logger         logr.Logger
	airflowAdapter airflow.AirflowAdapter
	repository     *repository.Repository
	db             *gorm.DB
}

func NewJob(config *config.Config, logger logr.Logger, db *gorm.DB) (Job, error) {
	logger.Info("Create new Job")

	airflowAdapter, err := airflow.NewAirflowAdapter(logger, config.AirflowAdapterConfig.Address, config.AirflowAdapterConfig.Username, config.AirflowAdapterConfig.Password)
	if err != nil {
		return nil, err
	}

	// Repository
	repo := repository.NewRepository(db)

	// cron
	cronJob := cron.New(cron.WithLogger(logger))

	return &job{
		cronJob:        cronJob,
		config:         config,
		logger:         logger.WithName("Job"),
		airflowAdapter: airflowAdapter,
		repository:     repo,
		db:             db,
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
}
