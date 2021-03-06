package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/squaresun/go-grpc-http-rest-microservice-tutorial/pkg/repo/memory"

	// mysql driver
	v1 "github.com/squaresun/go-grpc-http-rest-microservice-tutorial/pkg/svc/v1"

	"github.com/squaresun/go-grpc-http-rest-microservice-tutorial/pkg/logger"
	"github.com/squaresun/go-grpc-http-rest-microservice-tutorial/pkg/protocol/grpc"
	"github.com/squaresun/go-grpc-http-rest-microservice-tutorial/pkg/protocol/rest"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// GRPCPort is TCP port to listen by gRPC server
	GRPCPort string

	// HTTP/REST gateway start parameters section
	// HTTPPort is TCP port to listen by HTTP/REST gateway
	HTTPPort string

	// DB Datastore parameters section
	// DatastoreDBHost is host of database
	DatastoreDBHost string
	// DatastoreDBUser is username to connect to database
	DatastoreDBUser string
	// DatastoreDBPassword password to connect to database
	DatastoreDBPassword string
	// DatastoreDBSchema is schema of database
	DatastoreDBSchema string

	// Log parameters section
	// LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	LogLevel int
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "", "HTTP port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
	flag.IntVar(&cfg.LogLevel, "log-level", 0, "Global log level")
	flag.StringVar(&cfg.LogTimeFormat, "log-time-format", "",
		"Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfg.HTTPPort)
	}

	// initialize logger
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	todoRepo, err := memory.NewToDoDB()
	if err != nil {
		return fmt.Errorf("failed to init memory ToDoDB: %v", err)
	}

	v1API := v1.NewToDoServiceServer(todoRepo)

	// run HTTP gateway
	go func() {
		err := rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
		panic(err)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
