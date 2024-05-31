package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/gophish"
	"google.golang.org/grpc/codes"
)

func (s *Service) CheckHealth(ctx context.Context, request *api.CheckHealthRequest) (*api.CommonResponse, error) {

	templates, err := s.mailAdapter.CreateTemplate(ctx, &gophish.CreateTemplateParams{
		Name:    "Test from Backend",
		Subject: "Test from Backend",
		Text:    "Test from Backend",
		Html:    "",
	})
	if err != nil {
		return nil, err
	}
	s.log.Info("templates", "values", templates)

	return &api.CommonResponse{
		Code:    int32(codes.OK),
		Message: "Success",
	}, nil
}
