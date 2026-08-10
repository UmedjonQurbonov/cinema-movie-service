package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	moviev1 "github.com/UmedjonQurbonov/cinema-libs/gen/go/movie/v1"
	"github.com/UmedjonQurbonov/cinema-libs/pkg/postgres"
	server "github.com/UmedjonQurbonov/cinema-movie-service/delivery/http/v1"
	"github.com/UmedjonQurbonov/cinema-movie-service/internal/repository"
	"github.com/UmedjonQurbonov/cinema-movie-service/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// 1. Инициализируем Zap логгер
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Warn(".env file not found, relying on system environment variables")
	}

	dbPool, err := postgres.NewPgxPool(ctx, logger)
	if err != nil {
		logger.Fatal("failed to initialize postgres pool", zap.Error(err))
	}
	defer dbPool.Close()

	logger.Info("application initialized successfully")

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	movieRepository := repository.NewMovieRepo(dbPool)
	movieService := service.NewMovieService(movieRepository)
	movieServer := server.NewServer(movieService, logger)

	moviev1.RegisterMovieServiceServer(grpcServer, movieServer)

	reflection.Register(grpcServer)

	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	log.Println("shutting down server...")
	grpcServer.GracefulStop()
	log.Println("server stopped")
}