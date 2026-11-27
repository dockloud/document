package main

import (
	"log"
	"net"

	docv1 "github.com/dockloud/contract/gen/document/v1"
	myserver "github.com/dockloud/document/internal/delivery/grpc"
	"github.com/dockloud/document/internal/delivery/grpc/auth"
	"github.com/dockloud/document/internal/infrastructure/cache/redis"
	"github.com/dockloud/document/internal/infrastructure/config"
	"github.com/dockloud/document/internal/infrastructure/database/postgres"
	"github.com/dockloud/document/internal/infrastructure/document"
	usecase "github.com/dockloud/document/internal/usecase/document"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	conf, err := config.LoadMergedConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	db, err := postgres.NewDB(conf.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	cache, err := redis.NewCache(conf.Redis.URL)
	if err != nil {
		log.Fatalf("failed to connect to cache: %s", err)
	}

	repo := document.NewRepository(db, cache)
	service := usecase.NewService(repo)

	listener, err := net.Listen("tcp", conf.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gwayConf, ok := conf.Dependencies["dockloud-gateway"]
	if !ok {
		log.Fatalf("failed to load dependencies")
	}
	linkConf, ok := conf.Dependencies["dockloud-link"]
	if !ok {
		log.Fatalf("failed to load dependencies")
	}

	keyring := auth.Keyring{
		gwayConf.ID: {
			ID:      gwayConf.ID,
			Service: gwayConf.Name,
			Secret:  conf.JWT.GatewaySecret,
		},
		linkConf.ID: {
			ID:      linkConf.ID,
			Service: linkConf.Name,
			Secret:  conf.JWT.GatewaySecret,
		},
	}
	authInterceptor, err := auth.NewKeyringInterceptor(keyring)
	if err != nil {
		log.Fatalf("failed to create interceptor: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor))
	healthServer := health.NewServer()

	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	docv1.RegisterDocumentServiceServer(grpcServer, myserver.NewDocumentServer(*service))
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
