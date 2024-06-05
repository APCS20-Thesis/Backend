package gophish

import (
	"context"
	"github.com/go-logr/logr"
)

const (
	// user_groups
	Endpoint_CREATE_USER_GROUP = "/api/groups/"
)

type GophishAdapter interface {
	// User Groups
	CreateUserGroup(ctx context.Context, params *CreateUserGroupParams) (*CreateUserGroupResponse, error)
}

type gophish struct {
	log logr.Logger
}

func NewGophishAdapter(log logr.Logger) (GophishAdapter, error) {
	return &gophish{
		log: log,
	}, nil
}
