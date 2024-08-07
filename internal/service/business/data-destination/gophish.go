package data_destination

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/gophish"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"gorm.io/gorm"
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

	segment, err := b.repository.SegmentRepository.GetSegment(ctx, request.SegmentId)
	if err != nil {
		return err
	}

	config, err := json.Marshal(model.GophishDestinationConfiguration{
		UserGroupName: request.Name,
		Mapping:       request.Mapping,
	})
	if err != nil {
		logger.Error(err, "cannot marshal gophish config")
		return err
	}

	mappingOptions := []*api.MappingOptionItem{
		{
			SourceFieldName:      request.Mapping.Email,
			DestinationFieldName: "email",
		},
		{
			SourceFieldName:      request.Mapping.FirstName,
			DestinationFieldName: "first_name",
		},
		{
			SourceFieldName:      request.Mapping.LastName,
			DestinationFieldName: "last_name",
		},
		{
			SourceFieldName:      request.Mapping.Position,
			DestinationFieldName: "position",
		},
	}
	jsonMappingOptions, err := json.Marshal(mappingOptions)
	if err != nil {
		logger.Error(err, "cannot marshal mapping options")
		return err
	}

	var dataActionRun *model.DataActionRun
	txErr := b.db.Transaction(func(tx *gorm.DB) error {
		destination, err := b.repository.DataDestinationRepository.CreateDataDestination(ctx, &repository.CreateDataDestinationParams{
			Name:          request.Name,
			AccountUuid:   uuid.MustParse(accountUuid),
			Type:          model.DataDestinationType_GOPHISH,
			Configuration: pqtype.NullRawMessage{RawMessage: config, Valid: config != nil},
			ConnectionId:  0,
		})
		if err != nil {
			logger.Error(err, "cannot create data destination")
			return err
		}

		destSegmentMap, err := b.repository.DestSegmentMapRepository.CreateDestinationSegmentMap(ctx, &repository.CreateDestinationSegmentMapParams{
			SegmentId:      request.SegmentId,
			DestinationId:  destination.ID,
			MappingOptions: pqtype.NullRawMessage{RawMessage: jsonMappingOptions, Valid: jsonMappingOptions != nil},
		})
		if err != nil {
			logger.Error(err, "cannot create destination segment map")
			return err
		}

		dataAction, err := b.repository.DataActionRepository.CreateDataAction(ctx, &repository.CreateDataActionParams{
			TargetTable: model.TargetTable_DestSegmentMap,
			ActionType:  model.ActionType_ExportGophish,
			AccountUuid: uuid.MustParse(accountUuid),
			Status:      model.DataActionStatus_Success,
			ObjectId:    destSegmentMap.ID,
			RunCount:    1,
		})
		if err != nil {
			logger.Error(err, " cannot create data action")
			return err
		}
		dataActionRun, err = b.repository.DataActionRunRepository.CreateDataActionRun(ctx, &repository.CreateDataActionRunParams{
			ActionId:    dataAction.ID,
			RunId:       1,
			Status:      model.DataActionRunStatus_Processing,
			AccountUuid: uuid.MustParse(accountUuid),
		})
		if err != nil {
			logger.Error(err, " cannot create data action run")
			return err
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}

	// Process export Gophish
	txErr = b.db.Transaction(func(tx *gorm.DB) error {
		s3path := fmt.Sprintf("s3a://%s/%s", b.config.S3StorageConfig.Bucket, utils.GenerateDeltaAudiencePath(segment.MasterSegmentId))

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
				profile["last_name"] = fmt.Sprintf("%s", data[request.Mapping.LastName])
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

		err = b.repository.DataActionRunRepository.UpdateDataActionRunStatus(ctx, dataActionRun.ID, model.DataActionRunStatus_Success)
		if err != nil {
			logger.Error(err, "cannot update data action run status", "dataActionStatusId", dataActionRun.ID)
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
