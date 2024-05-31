package gophish

import (
	"context"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/go-logr/logr"
)

const (
	apiKey = "413e5f98d3d8327dc1fbbedfea8d9134732477e4ab7859886fb820d536ad97f8"

	// templates
	Endpoint_LIST_TEMPLATES   = "/api/templates"
	Endpoint_CREATE_TEMPLATES = "/api/templates/"
	Endpoint_TEMPLATE         = "/api/templates/template_id"
	// sending profiles
	Endpoint_LIST_SENDING_PROFILES  = "/api/smtp/"
	Endpoint_GET_SENDING_PROFILE    = "/api/smtp/id"
	Endpoint_CREATE_SENDING_PROFILE = "/api/smtp"
	// user_groups
	Endpoint_CREATE_USER_GROUP = "/api/groups/"
)

type GophishAdapter interface {
	// Templates
	ListTemplates(ctx context.Context) ([]Template, error)
	GetTemplate(ctx context.Context, id int) (Template, error)
	CreateTemplate(ctx context.Context, params *CreateTemplateParams) (Template, error)
	// User Groups
	CreateUserGroup(ctx context.Context, params *CreateUserGroupParams) (*CreateUserGroupResponse, error)
}

type gophish struct {
	log    logr.Logger
	client utils.HttpClient
}

func NewMailAdapter(log logr.Logger, host string) (GophishAdapter, error) {
	client := utils.HttpClient{}
	client.Init("Mail Client", log, host)
	return &gophish{
		log:    log,
		client: client,
	}, nil
}
