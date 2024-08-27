package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/opplieam/bb-grpc/protogen/go/product"
	"github.com/opplieam/bb-product-server/internal/product"
	"github.com/opplieam/bb-product-server/internal/store"

	"google.golang.org/grpc"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func setupDB() (*sql.DB, error) {
	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DSN"))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	lis, err := net.Listen("tcp", ":3031")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	db, err := setupDB()
	if err != nil {
		log.Fatalf("failed to setup database: %v", err)
	}

	// TODO: Add logging interceptor
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, product.NewServer(store.NewProductStore(db)))
	log.Printf("Starting gRPC Server at %v\n", lis.Addr())

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
