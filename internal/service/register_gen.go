package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// RegisterWithServer implementing service server interface
//func (s *Service) RegisterWithServer(server *grpc.Server) {
//	api.RegisterCDPServiceServer(server, s)
//}

// RegisterWithHandler implementing service server interface
func (s *Service) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := api.RegisterCDPServiceHandler(ctx, mux, conn)
	//if err != nil {
	//	s.log.Error(err, "Error register servers")
	//}

	return err
}
