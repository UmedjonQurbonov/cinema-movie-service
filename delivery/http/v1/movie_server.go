package v1

import (
	"context"

	moviev1 "github.com/UmedjonQurbonov/cinema-libs/gen/go/movie/v1"
	"github.com/UmedjonQurbonov/cinema-movie-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"go.uber.org/zap"
)

type MovieService interface {
	Create(ctx context.Context, movie *domain.Movie) (int, error)
}

type MovieServer struct {
	moviev1.UnimplementedMovieServiceServer
	service MovieService
	log     *zap.Logger
}

func NewServer(service MovieService, log *zap.Logger) *MovieServer {
	return &MovieServer{
		service: service,
		log:     log,
	}
}

func (s *MovieServer) Create(ctx context.Context, req *moviev1.CreateMovieRequest) (*moviev1.CreateMovieResponse, error) {
	// 1. Формируем доменную модель с безопасными геттерами Protobuf
	movie := &domain.Movie{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Duration:    int(req.GetDuration()),
		Age_limit:   int(req.GetAgeLimit()),
	}

	// 2. Вызываем сервис
	id, err := s.service.Create(ctx, movie)
	if err != nil {
		s.log.Error("failed to create movie",
			zap.String("title", req.GetTitle()),
			zap.Error(err),
		)

		return nil, status.Error(codes.Internal, "internal error")
	}

	s.log.Info("movie created successfully",
		zap.Int64("id", int64(id)),
		zap.String("title", req.GetTitle()),
		zap.Int("duration_min", int(req.GetDuration())),
	)

	// 3. Возвращаем правильную структуру CreateMovieResponse (преобразовав int в int64)
	return &moviev1.CreateMovieResponse{
		Id: int64(id),
	}, nil
}		