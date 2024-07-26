package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/go-logr/logr"
)

type AlertAdapter interface {
	AlertError(ctx context.Context, message *ErrorMessage) error
}

type alertAdapter struct {
	log     logr.Logger
	client  utils.HttpClient
	webhook string
}

func NewAlertAdapter(log logr.Logger, webhook string) (AlertAdapter, error) {
	client := utils.HttpClient{}
	client.Init("Airflow Client", log, webhook)

	return &alertAdapter{
		log:     log,
		client:  client,
		webhook: webhook,
	}, nil
}

func (c *alertAdapter) AlertError(ctx context.Context, message *ErrorMessage) error {
	strReq, err := json.Marshal(message.Request)
	if err != nil {
		c.log.WithName("AlertError").Error(err, "cannot marshal request", "request", message.Request)
		return err
	}
	jsonStrReq, err := json.Marshal(string(strReq))
	if err != nil {
		c.log.WithName("AlertError").Error(err, "cannot marshal request", "request string", string(strReq))
		return err
	}

	request := &alertRequest{
		Embeds: []embeds{
			{
				Title:       message.Title,
				Description: message.Description,
				Color:       16711680,
				Fields: []embedField{
					{
						Name:   "Request",
						Value:  fmt.Sprintf("```json\n%s\n```", string(jsonStrReq)),
						Inline: false,
					},
					{
						Name:   "Error",
						Value:  message.ErrorMessage,
						Inline: false,
					},
				},
			},
		},
	}

	response := emptyStruct{}
	err = c.client.SendHttpRequest(ctx, utils.Request{
		Endpoint: "",
		Method:   utils.Method_POST,
		Body:     request,
		Headers:  map[string]string{utils.Header_CONTENT_TYPE: "application/json"},
	}, &response)
	if err != nil {
		c.log.WithName("AlertError").Error(err, "failed to send alert")
		return err
	}

	return nil
}

type emptyStruct struct{}
