package service

import (
	"context"
	pb "github.com/APCS20-Thesis/Backend/api"
	"google.golang.org/grpc/codes"
)

func (s *Service) CheckHealth(ctx context.Context, request *pb.CheckHealthRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{
		Code:    int32(codes.OK),
		Message: "Hello",
	}, nil
}
