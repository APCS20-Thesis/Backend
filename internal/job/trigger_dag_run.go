package job

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
)

func (j *job) TriggerDagRuns(ctx context.Context) error {
	// get data action that are in status PENDING
	dataActions, err := j.repository.DataActionRepository.GetListDataActions(ctx, &repository.GetListDataActionsParams{
		Statuses: []model.DataActionStatus{model.DataActionStatus_Pending},
	})
	if err != nil {
		j.logger.WithName("TriggerDagRun").Error(err, "fail to get list data actions")
		return err
	}

	// gọi xuống airflow để lấy danh sách Dags tương ứng dataActions

	for _, dataAction := range dataActions {
		// TODO: Tạo transaction chỗ này
		// nếu có Dag tương ứng thì
		// - Update DataAction status
		err := j.repository.DataActionRepository.UpdateDataAction(ctx, &repository.UpdateDataActionParams{
			ID:          dataAction.ID,
			ActionType:  dataAction.ActionType,
			Payload:     dataAction.Payload,
			Schedule:    dataAction.Schedule,
			AccountUuid: dataAction.AccountUuid,
			Status:      model.DataActionStatus_Triggered,
		})
		if err != nil {
			return err
		}
		// - Gọi airflow trigger DagRun

		// - Tạo DataActionRun với thông tin DagRun trên

	}

	// list dags with filter according to dag_id
	return nil
}

func (j *job) LogHello(ctx context.Context) error {
	j.logger.WithName("LogHello").Info("HELLO")
	return nil
}
