package main

import (
	"log"
	"net"
	"os"

	docv1 "github.com/dockloud/contract/gen/document/v1"
	myserver "github.com/dockloud/document/internal/delivery/grpc"
	"github.com/dockloud/document/internal/delivery/grpc/auth"
	"github.com/dockloud/document/internal/infrastructure/cache/redis"
	"github.com/dockloud/document/internal/infrastructure/database/postgres"
	"github.com/dockloud/document/internal/infrastructure/document"
	usecase "github.com/dockloud/document/internal/usecase/document"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	db, err := postgres.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	cache, err := redis.NewCache(os.Getenv("CACHE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to cache: %s", err)
	}

	repo := document.NewRepository(db, cache)
	service := usecase.NewService(repo)

	listener, err := net.Listen("tcp", ":50150")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	keyering := auth.Keyring{
		"gateway-key": {
			ID:      "dck-gateway-v1",
			Service: "dockloud-gateway",
			Secret:  []byte(os.Getenv("GATEWAY_SERVICE_KEY")),
		},
		"link-key": {
			ID:      "dck-link-v1",
			Service: "dockloud-link",
			Secret:  []byte(os.Getenv("LINK_SERVICE_KEY")),
		},
	}
	authInterceptor, err := auth.NewKeyringInterceptor(keyering)
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
