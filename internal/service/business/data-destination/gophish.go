package data_destination

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"

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

	s3path := "s3a://cdp-thesis-apcs/" + utils.GenerateDeltaAudiencePath(segment.MasterSegmentId)

	queryResponse, err := b.queryAdapter.QueryRawSQL(ctx, &query.QueryRawSQLRequest{
		Query: fmt.Sprintf("SELECT * FROM delta.`%s`;", s3path),
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
		Status:        "DONE",
		Configuration: pqtype.NullRawMessage{RawMessage: config, Valid: config != nil},
		ConnectionId:  0,
	})
	if err != nil {
		logger.Error(err, "cannot create data destination")
		return err
	}

	return nil
}
