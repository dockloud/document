package main

import (
	"log"
	"net"
	"os"

	docv1 "github.com/dockloud/contract/gen/document/v1"
	myserver "github.com/dockloud/document/internal/delivery/grpc"
	"github.com/dockloud/document/internal/infrastructure/cache/redis"
	"github.com/dockloud/document/internal/infrastructure/database/postgres"
	"github.com/dockloud/document/internal/infrastructure/document"
	usecase "github.com/dockloud/document/internal/usecase/document"
	"google.golang.org/grpc"
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

	listener, err := net.Listen("tcp", ":50250")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	docv1.RegisterDocumentServiceServer(grpcServer, myserver.NewDocumentServer(*service))

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
