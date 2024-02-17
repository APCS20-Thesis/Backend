package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
)

func (s *Service) ImportFile(ctx context.Context, request *api.ImportFileRequest) (*api.ImportFileResponse, error) {
	err := s.s3Manger.S3Uploader("cdp-thesis-apcs", "datas/"+request.GetFileName(), request.GetFileContent())
	if err != nil {
		return nil, err
	}
	return &api.ImportFileResponse{Message: "Import Success", Code: 0}, nil
}
