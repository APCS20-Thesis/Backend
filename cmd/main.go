package main

import (
	"context"
	pb "github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/service"
	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

// global variables
var logger logr.Logger

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}

	logger = cfg.Log.MustBuildLogR()

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", cfg.ServerConfig.GrpcServerAddress) // ":10443"
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	reflection.Register(s)
	// Attach the Greeter service to the server
	cdpService, err := service.NewService(logger)
	if err != nil {
		log.Fatalln("Failed to create new service:", err)
		return
	}
	pb.RegisterCDPServiceServer(s, cdpService)

	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0" + cfg.ServerConfig.GrpcServerAddress)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0"+cfg.ServerConfig.GrpcServerAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	// Register Greeter
	err = pb.RegisterCDPServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    cfg.ServerConfig.HttpServerAddress,
		Handler: gwmux,
	}
	log.Println("Serving gRPC-Gateway for REST on http://0.0.0.0" + cfg.ServerConfig.HttpServerAddress)
	log.Fatalln(gwServer.ListenAndServe())
}
