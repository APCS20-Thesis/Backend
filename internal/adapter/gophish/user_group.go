package gophish

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/utils"
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
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	ModifiedDate time.Time `json:"modified_date"`
	Targets      []struct {
		Email string `json:"email"`
	}
}

func (c *gophish) CreateUserGroup(ctx context.Context, params *CreateUserGroupParams) (*CreateUserGroupResponse, error) {

	var endpoint = params.GophishConfig.Host
	if params.GophishConfig.Port != "" {
		endpoint = endpoint + ":" + params.GophishConfig.Port
	}

	client := utils.HttpClient{}
	client.Init("Gophish Client", c.log, endpoint)

	var response CreateUserGroupResponse

	err := client.SendHttpRequest(ctx, utils.Request{
		Endpoint: Endpoint_CREATE_USER_GROUP,
		Method:   utils.Method_POST,
		Headers: map[string]string{
			"Authorization": params.GophishConfig.ApiKey,
		},
		Body: params.Payload,
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
