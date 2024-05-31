package gophish

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/utils"
	"strconv"
	"time"
)

type CreateUserGroupPayload struct {
	Name    string              `json:"name"`
	Targets []map[string]string `json:"targets"`
}

type CreateUserGroupParams struct {
	GophishConfig model.GophishConfiguration
	Payload       CreateUserGroupPayload
}

type CreateUserGroupResponse struct {
	Id           int
	Name         string
	ProfileCount int
}

func (c *gophish) CreateUserGroup(ctx context.Context, params *CreateUserGroupParams) (*CreateUserGroupResponse, error) {

	var endpoint = params.GophishConfig.Host
	if params.GophishConfig.Port != 0 {
		endpoint = endpoint + ":" + strconv.FormatInt(int64(params.GophishConfig.Port), 10)
	}

	var response struct {
		Id           int       `json:"id"`
		Name         string    `json:"name"`
		ModifiedDate time.Time `json:"modified_date"`
		Targets      []struct {
			Email string `json:"email"`
		}
	}

	err := c.client.SendHttpRequest(ctx, utils.Request{
		Endpoint: endpoint,
		Method:   utils.Method_GET,
		Headers: map[string]string{
			"Authorization": apiKey,
		},
	}, &response)
	if err != nil {
		return nil, err
	}

	return &CreateUserGroupResponse{
		Id:           response.Id,
		Name:         response.Name,
		ProfileCount: len(response.Targets),
	}, nil
}
