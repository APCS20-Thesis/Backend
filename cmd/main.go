package main

import (
	"context"
	"fmt"
	pb "github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/adapter/mqtt"
	"github.com/APCS20-Thesis/Backend/internal/job"
	"github.com/APCS20-Thesis/Backend/internal/service"
	"github.com/go-logr/logr"
	migrateV4 "github.com/golang-migrate/migrate/v4"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// global variables
var logger logr.Logger

const (
	secretKey     = "secret"
	tokenDuration = 360 * time.Minute
)
const versionTimeFormat = "20060102150405"

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}

	logger = cfg.Log.MustBuildLogR()

	if cfg.Env == "dev" {
		fmt.Println(cfg)
	}

	app := cli.NewApp()
	app.Name = "service"
	// app.Usage = "tekit tool"
	// app.Version = Version
	app.Commands = []*cli.Command{
		{
			Name:   "server",
			Usage:  "start grpc/http server",
			Action: serverAction,
		},
		{
			Name:        "migrate",
			Usage:       "doing database migration",
			Subcommands: MigrateCliCommand(cfg.MigrationFolder, cfg.PostgreSQL.String()),
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
	return nil
}

func serverAction(cliCtx *cli.Context) error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln("Failed to load config:", err)
		return err
	}

	logger = cfg.Log.MustBuildLogR()

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", cfg.ServerConfig.GrpcServerAddress) // ":10443"
	if err != nil {
		log.Fatalln("Failed to listen:", err)
		return err
	}

	// Create a gRPC server object
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	interceptor := service.NewAuthInterceptor(jwtManager, config.AccessibleRoles())

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptor.Unary(), grpc_validator.UnaryServerInterceptor())),
	)
	reflection.Register(s)
	// Attach the Greeter service to the server

	gormDb, err := ConnectPostgresql(cfg.PostgreSQL.String())
	if err != nil {
		log.Fatalln("Failed to init gorm db:", err)
		return err
	}
	cdpService, err := service.NewService(logger, cfg, gormDb, jwtManager)
	if err != nil {
		log.Fatalln("Failed to create new service:", err)
		return err
	}
	pb.RegisterCDPServiceServer(s, cdpService)
	pb.RegisterCDPServiceFileServer(s, cdpService)
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0" + cfg.ServerConfig.GrpcServerAddress)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()
	// mqtt
	mqtt, err := mqtt.NewMqttAdapter(cfg, logger, gormDb)
	if err != nil {
		log.Fatalln("Failed to create new mqtt client:", err)
		return err
	}
	mqtt.Connect()

	// job
	job, err := job.NewJob(cfg, logger, gormDb, mqtt)
	if err != nil {
		log.Fatalln("Failed to create new job:", err)
		return err
	}
	job.RegisterCronJobs()
	job.StartCron()

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
		return err
	}

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		OrigName:     true,
		EmitDefaults: true,
	}))
	err = pb.RegisterCDPServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
		return err
	}
	err = pb.RegisterCDPServiceFileClient(context.Background(), gwmux, pb.NewCDPServiceFileClient(conn))
	if err != nil {
		log.Fatalln("Failed to register import csv:", err)
		return err
	}
	gwServer := &http.Server{
		Addr:    cfg.ServerConfig.HttpServerAddress,
		Handler: cors(gwmux),
	}
	log.Println("Serving gRPC-Gateway for REST on http://0.0.0.0" + cfg.ServerConfig.HttpServerAddress)
	log.Fatalln(gwServer.ListenAndServe())
	return nil
}

func MigrateCliCommand(sourceURL string, databaseURL string) []*cli.Command {
	// Migration should always run on development mode
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return []*cli.Command{
		{
			Name:  "up",
			Usage: "lift migration up to date",
			Action: func(c *cli.Context) error {
				m, err := migrateV4.New(sourceURL, databaseURL)
				if err != nil {
					logger.Fatal("Error create migration", zap.Error(err))
				}

				logger.Info("migration up")
				if err := m.Up(); err != nil && err != migrateV4.ErrNoChange {
					logger.Fatal(err.Error())
				}
				return err
			},
		},
		{
			Name:  "down",
			Usage: "step down migration by N(int)",
			Action: func(c *cli.Context) error {
				m, err := migrateV4.New(sourceURL, databaseURL)
				if err != nil {
					logger.Fatal("Error create migration", zap.Error(err))
				}

				down, err := strconv.Atoi(c.Args().Get(0))
				if err != nil {
					logger.Fatal("rev should be a number", zap.Error(err))
				}

				logger.Info("migration down", zap.Int("down", -down))
				if err := m.Steps(-down); err != nil {
					logger.Fatal(err.Error())
				}
				return err
			},
		},
		{
			Name:  "force",
			Usage: "Enforce dirty migration with verion (int)",
			Action: func(c *cli.Context) error {
				m, err := migrateV4.New(sourceURL, databaseURL)
				if err != nil {
					logger.Fatal("Error create migration", zap.Error(err))
				}

				ver, err := strconv.Atoi(c.Args().Get(0))
				if err != nil {
					logger.Fatal("rev should be a number", zap.Error(err))
				}

				logger.Info("force", zap.Int("ver", ver))

				if err := m.Force(ver); err != nil {
					logger.Fatal(err.Error())
				}
				return err
			},
		},
		{
			Name: "create",
			Action: func(c *cli.Context) error {
				folder := strings.ReplaceAll(sourceURL, "file://", "")
				now := time.Now()
				ver := now.Format(versionTimeFormat)
				name := strings.Join(c.Args().Slice(), "-")

				up := fmt.Sprintf("%s/%s_%s.up.sql", folder, ver, name)
				down := fmt.Sprintf("%s/%s_%s.down.sql", folder, ver, name)

				logger.Info("create migration", zap.String("name", name))
				logger.Info("up script", zap.String("up", up))
				logger.Info("down script", zap.String("down", up))

				if err := ioutil.WriteFile(up, []byte{}, 0600); err != nil {
					logger.Fatal("Create migration up error", zap.Error(err))
				}
				if err := ioutil.WriteFile(down, []byte{}, 0600); err != nil {
					logger.Fatal("Create migration down error", zap.Error(err))
				}
				return nil
			},
		},
	}
}

func allowedOrigin(origin string) bool {
	if viper.GetString("cors") == "*" {
		return true
	}
	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}
	return false
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin(r.Header.Get("Origin")) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType,")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
