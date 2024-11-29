package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	pb "github.com/opplieam/bb-grpc/protogen/go/product"
	"github.com/opplieam/bb-product-server/internal/product"
	"github.com/opplieam/bb-product-server/internal/store"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var build = "dev"

const (
	jaegerEndpoint = "localhost:4317"
	serviceName    = "bb-product-server"
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

func initTracerProvider() error {
	// Ensure default SDK resources and service name are set
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to merge resources: %w", err)
	}
	// Create and configure the OTel exporter for Jaeger
	otelExporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint(jaegerEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("failed to build OtlpExporter: %w", err)
	}

	// Create Tracer Provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(otelExporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}

func main() {
	// TODO: Add health check
	// TODO: Make address as environment variable
	lis, err := net.Listen("tcp", ":3031")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Setup DB
	db, err := setupDB()
	if err != nil {
		log.Fatalf("failed to setup database: %v", err)
	}
	// Add Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = logger.With("service", "bb-product-server", "build", build)

	// Setup Tracer
	err = initTracerProvider()
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}

	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
		),
	)
	productService := product.NewServer(store.NewProductStore(db))
	pb.RegisterProductServiceServer(grpcServer, productService)
	log.Printf("Starting gRPC Server at %v\n", lis.Addr())

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
