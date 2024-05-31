package gophish

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
	"strconv"
	"strings"
	"time"
)

type Template struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	EnvelopeSender string    `json:"envelope_sender"`
	Subject        string    `json:"subject"`
	Text           string    `json:"text"`
	Html           string    `json:"html"`
	ModifiedDate   time.Time `json:"modified_date"`
}

type CreateTemplateParams struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Subject      string    `json:"subject"`
	Text         string    `json:"text"`
	Html         string    `json:"html"`
	ModifiedDate time.Time `json:"modified_date"`
}

func (c *gophish) CreateTemplate(ctx context.Context, params *CreateTemplateParams) (Template, error) {
	var response Template
	err := c.client.SendHttpRequest(ctx, utils.Request{
		Endpoint: Endpoint_CREATE_TEMPLATES,
		Method:   utils.Method_POST,
		Body:     params,
		Headers: map[string]string{
			"Authorization": apiKey,
		},
	}, &response)
	if err != nil {
		return Template{}, err
	}

	return response, nil
}

func (c *gophish) ListTemplates(ctx context.Context) ([]Template, error) {
	var response []Template
	err := c.client.SendHttpRequest(ctx, utils.Request{
		Endpoint: Endpoint_LIST_TEMPLATES,
		Method:   utils.Method_POST,
		Headers: map[string]string{
			"Authorization": apiKey,
		},
	}, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *gophish) GetTemplate(ctx context.Context, id int) (Template, error) {
	var response Template
	endpoint := strings.Replace(Endpoint_TEMPLATE, "template_id", strconv.FormatInt(int64(id), 10), 1)

	err := c.client.SendHttpRequest(ctx, utils.Request{
		Endpoint: endpoint,
		Method:   utils.Method_GET,
		Headers: map[string]string{
			"Authorization": apiKey,
		},
	}, &response)
	if err != nil {
		return Template{}, err
	}

	return response, err
}
