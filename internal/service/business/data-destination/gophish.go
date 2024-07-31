package data_destination

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"

	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/gophish"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/utils"
)

func (b business) CreateGophishUserGroupFromSegment(ctx context.Context, accountUuid string, request *api.CreateGophishUserGroupFromSegmentRequest) error {
	logger := b.log.WithName("CreateGophishUserGroupFromSegment")

	gophishConnection, err := b.repository.ConnectionRepository.GetConnection(ctx, request.ConnectionId)
	if err != nil {
		logger.Error(err, "cannot get connection", "connection id", request.ConnectionId)
		return err
	}

	var configuration model.GophishConfiguration
	err = json.Unmarshal(gophishConnection.Configurations.RawMessage, &configuration)
	if err != nil {
		logger.Error(err, "cannot get connection", "connection id", request.ConnectionId)
		return err
	}

	segment, err := b.repository.SegmentRepository.GetSegment(ctx, request.SegmentId, accountUuid)
	if err != nil {
		return err
	}

	dataAction, err := b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
		TargetTable: model.TargetTable_Segment,
		ActionType:  model.ActionType_ExportGophish,
		AccountUuid: uuid.MustParse(accountUuid),
		Status:      model.DataActionStatus_Success,
		ObjectId:    request.SegmentId,
		RunCount:    1,
	})
	if err != nil {
		logger.Error(err, " cannot create data action")
		return err
	}
	dataActionRun, err := b.repository.DataActionRunRepository.CreateDataActionRun(ctx, &repository.CreateDataActionRunParams{
		ActionId:    dataAction.ID,
		RunId:       1,
		DagRunId:    "",
		Status:      model.DataActionRunStatus_Processing,
		AccountUuid: uuid.MustParse(accountUuid),
	})

	s3path := fmt.Sprintf("s3a://%s/%s", b.config.S3StorageConfig.Bucket, utils.GenerateDeltaAudiencePath(segment.MasterSegmentId))

	txErr := b.db.Transaction(func(tx *gorm.DB) error {
		queryResponse, err := b.queryAdapter.QueryRawSQL(ctx, &query.QueryRawSQLRequest{
			Query: fmt.Sprintf("SELECT * FROM delta.`%s` WHERE %s;", s3path, segment.SqlCondition),
		})
		if err != nil {
			return err
		}

		targets := make([]map[string]string, 0, queryResponse.Count)
		for _, data := range queryResponse.Data {
			profile := make(map[string]string)
			profile["email"] = fmt.Sprintf("%s", data[request.Mapping.Email])
			if request.Mapping.FirstName != "" {
				profile["first_name"] = fmt.Sprintf("%s", data[request.Mapping.FirstName])
			}
			if request.Mapping.LastName != "" {
				profile["first_name"] = fmt.Sprintf("%s", data[request.Mapping.LastName])
			}
			if request.Mapping.Position != "" {
				profile["position"] = fmt.Sprintf("%s", data[request.Mapping.Position])
			}
			targets = append(targets, profile)
		}

		resp, err := b.gophishAdapter.CreateUserGroup(ctx, &gophish.CreateUserGroupParams{
			GophishConfig: configuration,
			Payload: gophish.CreateUserGroupPayload{
				Name:    request.Name,
				Targets: targets,
			},
		})
		if err != nil {
			logger.Error(err, "cannot create user group gophish")
			return err
		}
		logger.Info("create gophish user group response", "response", resp)

		config, err := json.Marshal(model.GophishDestinationConfiguration{
			UserGroupName: request.Name,
			Mapping:       request.Mapping,
		})

		_, err = b.repository.DataDestinationRepository.CreateDataDestination(ctx, &repository.CreateDataDestinationParams{
			Name:          request.Name,
			AccountUuid:   uuid.MustParse(accountUuid),
			Type:          model.DataDestinationType_GOPHISH,
			Configuration: pqtype.NullRawMessage{RawMessage: config, Valid: config != nil},
			ConnectionId:  0,
		})
		if err != nil {
			logger.Error(err, "cannot create data destination")
			tx.Rollback()
			return err
		}

		err = b.repository.DataActionRunRepository.UpdateDataActionRunStatus(ctx, dataActionRun.ID, model.DataActionRunStatus_Success)
		if err != nil {
			logger.Error(err, "cannot update data action run status", "dataActionStatusId", dataActionRun.ID)
			tx.Rollback()
			return err
		}

		return nil
	})
	if txErr != nil {
		createErr := b.repository.DataActionRunRepository.UpdateDataActionRunStatus(ctx, dataActionRun.ID, model.DataActionRunStatus_Failed)
		if createErr != nil {
			logger.Error(err, "cannot update data action run status", "dataActionStatusId", dataActionRun.ID)
			return createErr
		}
		return txErr
	}

	return nil
}
